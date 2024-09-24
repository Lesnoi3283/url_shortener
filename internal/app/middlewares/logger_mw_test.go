package middlewares

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLoggerMW(t *testing.T) {
	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Some response"))
	})

	//prepare response recorder
	recorder := httptest.NewRecorder()

	//test MW
	mw := LoggerMW(*sugar)
	handler := mw(next)
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	//check result
	assert.Equal(t, http.StatusOK, recorder.Code)

}
