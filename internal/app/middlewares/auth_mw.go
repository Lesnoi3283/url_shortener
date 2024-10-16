// Package middlewares contains middlewares (functions which will run before handlers).
package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
)

// JWT params.
const (
	TokenExp      = time.Hour * 3
	SecretKey     = "supersecretkey"
	JwtCookieName = "JWT"
)

type contextKey string

// UserIDContextKey is a key to get a userID from context values.
const UserIDContextKey contextKey = "userID"

// Claims is a jwt.RegisteredClaims struct with custom field "Claims.UserID".
type Claims struct {
	jwt.RegisteredClaims
	UserID int
}

// BuildNewJWTString returns new JWT string with userID inside.
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

// GetUserID parses JWT and returns a userID from it.
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
		return -1
	}

	return claims.UserID
}

//go:generate mockgen -source=auth_mw.go -destination=mocks/mocks_AuthMW.go -package=mocks_MW github.com/Lesnoi3283/url_shortener/internal/app/middlewares UserCreater

// UserCreater can create a new user.
type UserCreater interface {
	CreateUser(ctx context.Context) (int, error)
}

// AuthMW parses AuthJWT from cookie and puts UserID to http.Request.Context values.
func AuthMW(store UserCreater, logger zap.SugaredLogger) func(handlerFunc http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			coockie, err := r.Cookie(JwtCookieName)
			if err == nil {
				userID := GetUserID(coockie.Value)
				if userID != -1 {
					ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
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

			ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
