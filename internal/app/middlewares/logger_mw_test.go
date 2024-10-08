package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
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

func BenchmarkLoggerMW(b *testing.B) {
	//prepare logger
	logger := zap.NewNop()
	sugar := logger.Sugar()

	//prepare handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Some response"))
	})

	//prepare mw
	mw := LoggerMW(*sugar)
	handler := mw(next)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		//prepare request and response
		b.StopTimer()
		w := httptest.NewRecorder()
		r := httptest.NewRequest("GET", "/", nil)
		b.StartTimer()

		//test
		handler.ServeHTTP(w, r)
	}
}
