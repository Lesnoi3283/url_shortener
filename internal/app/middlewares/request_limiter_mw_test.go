package middlewares

import (
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestRequestLimiterMW_Ok(t *testing.T) {
	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare manager
	manager := NewRequestManager(1, time.Second*10)

	//prepare handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Some response"))
	})

	//prepare response recorder
	recorder := httptest.NewRecorder()

	//test MW
	mw := RequestLimiterMW(*sugar, manager)
	handler := mw(next)
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	//check result
	assert.Equal(t, http.StatusOK, recorder.Code)
}

func TestRequestLimiterMW_TooManyRequests(t *testing.T) {
	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare manager
	manager := NewRequestManager(0, time.Second*10)

	//prepare handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("handler must`nt be called")
	})

	//prepare response recorder
	recorder := httptest.NewRecorder()

	//test MW
	mw := RequestLimiterMW(*sugar, manager)
	handler := mw(next)
	handler.ServeHTTP(recorder, httptest.NewRequest("GET", "/", nil))

	//check result
	assert.Equal(t, http.StatusTooManyRequests, recorder.Code)
}

//TODO: Добавить бенчмарки на все мидлвари
