package middlewares

import (
	mocks_MW "github.com/Lesnoi3283/url_shortener/internal/app/middlewares/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthMW_WithValidJWT(t *testing.T) {

	//prepare data
	correctUserID := 1
	correctJWTString, err := BuildNewJWTString(correctUserID)
	require.NoError(t, err, "Err while preparing test")

	//prepare mocks
	c := gomock.NewController(t)
	store := mocks_MW.NewMockUserCreater(c)

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler witch will check our MW
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDContextKey)
		userIDInt, ok := userID.(int)
		assert.True(t, ok)
		assert.Equal(t, userIDInt, correctUserID)
		w.WriteHeader(http.StatusOK)
	})

	//prepare request and recorder
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.AddCookie(&http.Cookie{
		Name:  JwtCookieName,
		Value: correctJWTString,
	})

	w := httptest.NewRecorder()

	//test MW
	mw := AuthMW(store, *sugar)
	mw(nextHandler).ServeHTTP(w, r)

	//check result
	assert.Equal(t, http.StatusOK, w.Code)
}

func TestAuthMW_NoJWT(t *testing.T) {

	//prepare data
	correctUserID := 1

	//prepare mocks
	c := gomock.NewController(t)
	store := mocks_MW.NewMockUserCreater(c)
	store.EXPECT().CreateUser(gomock.Any()).Return(correctUserID, nil)

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler witch will check our MW
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(UserIDContextKey)
		userIDInt, ok := userID.(int)
		assert.True(t, ok)
		assert.Equal(t, userIDInt, correctUserID)
		w.WriteHeader(http.StatusOK)
	})

	//prepare request and recorder
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	//test MW
	mw := AuthMW(store, *sugar)
	mw(nextHandler).ServeHTTP(w, r)

	//check result
	assert.Equal(t, http.StatusOK, w.Code)
}
