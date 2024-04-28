package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const TOKEN_EXP = time.Hour * 3
const SECRET_KEY = "supersecretkey"
const JWT_COOCKIE_NAME = "JWT"

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func BuildNewJWTString(userID int) (string, error) {
	claims := Claims{RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
	},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	stringToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return "", fmt.Errorf("building new jwt: %w", err)
	}

	return stringToken, nil
}

func GetUserId(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SECRET_KEY), nil
		})
	if err != nil {
		return -1
	}

	if !token.Valid {
		fmt.Println("Token is not valid")
		return -1
	}

	fmt.Println("Token is valid")
	return claims.UserID
}

type UserCreater interface {
	CreateUser(ctx context.Context) (int, error)
}

func AuthMW(h http.Handler, store UserCreater, logger zap.SugaredLogger) http.HandlerFunc {
	authFn := func(w http.ResponseWriter, r *http.Request) {
		coockie, err := r.Cookie(JWT_COOCKIE_NAME)
		if err == nil {
			userIDd := GetUserId(coockie.Value)
			if userIDd != -1 {
				h.ServeHTTP(w, r)
			}
		}

		userID, err := store.CreateUser(r.Context())
		if err != nil {
			logger.Errorf("err while creating new user in auth mw: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		jwt, err := BuildNewJWTString(userID)
		if err != nil {
			logger.Errorf("err while building new jwt in auth mw: %v", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:  JWT_COOCKIE_NAME,
			Value: jwt,
		})

		h.ServeHTTP(w, r)

	}

	return authFn
}
