// Package middlewares contains middlewares (functions which will run before handlers).
package middlewares

import (
	"context"
	"github.com/Lesnoi3283/url_shortener/pkg/secure"
	"go.uber.org/zap"
	"net/http"
)

// JWT params.
const (
	JwtCookieName = "JWT"
)

type contextKey string

// UserIDContextKey is a key to get a userID from context values.
const UserIDContextKey contextKey = "userID"

//go:generate mockgen -source=auth_mw.go -destination=mocks/mocks_AuthMW.go -package=mocks_MW github.com/Lesnoi3283/url_shortener/internal/app/middlewares UserCreater

// UserCreater can create a new user.
type UserCreater interface {
	CreateUser(ctx context.Context) (int, error)
}

// AuthMW parses AuthJWT from cookie and puts UserID to http.Request.Context values.
func AuthMW(store UserCreater, logger zap.SugaredLogger, jh *secure.JWTHelper) func(handlerFunc http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie(JwtCookieName)
			if err == nil {
				userID, err := jh.GetUserID(cookie.Value)
				if err == nil {
					ctx := context.WithValue(r.Context(), UserIDContextKey, userID)
					next.ServeHTTP(w, r.WithContext(ctx))
					return
				} else {
					logger.Debugf("Error while getting userID from JWT: %v", err)
				}
			}

			userID, err := store.CreateUser(r.Context())
			if err != nil {
				logger.Errorf("err while creating new user in auth mw: %v", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			jwt, err := jh.BuildNewJWTString(userID)
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
