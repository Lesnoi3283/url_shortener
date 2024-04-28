package middlewares

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	"time"
)

const TokenExp = time.Hour * 3
const SecretKey = "supersecretkey"
const JwtCookieName = "JWT"
const UserIDContextName = "userID"

type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

func BuildNewJWTString(userID int) (string, error) {
	claims := Claims{RegisteredClaims: jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
	},
		UserID: userID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	stringToken, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", fmt.Errorf("building new jwt: %w", err)
	}

	return stringToken, nil
}

func GetUserID(tokenString string) int {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(SecretKey), nil
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
		coockie, err := r.Cookie(JwtCookieName)
		if err == nil {
			userID := GetUserID(coockie.Value)
			if userID != -1 {
				ctx := context.WithValue(r.Context(), UserIDContextName, userID)
				h.ServeHTTP(w, r.WithContext(ctx))
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
			Name:  JwtCookieName,
			Value: jwt,
		})

		ctx := context.WithValue(r.Context(), UserIDContextName, userID)
		h.ServeHTTP(w, r.WithContext(ctx))

	}

	return authFn
}
