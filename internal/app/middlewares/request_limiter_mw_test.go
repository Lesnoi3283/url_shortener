package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
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

func TestRequestLimiterMW_OkAndTooMany(t *testing.T) {
	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()
	timeToWait := time.Second * 3
	//prepare manager
	manager := NewRequestManager(1, timeToWait)

	//prepare handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Some response"))
	})

	//prepare response recorder
	recorderOk := httptest.NewRecorder()
	recorderTooMany := httptest.NewRecorder()
	recorderOk2 := httptest.NewRecorder()

	//test MW
	mw := RequestLimiterMW(*sugar, manager)
	handler := mw(next)

	handler.ServeHTTP(recorderOk, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusOK, recorderOk.Code)
	handler.ServeHTTP(recorderTooMany, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusTooManyRequests, recorderTooMany.Code)
	time.Sleep(timeToWait + time.Second)
	handler.ServeHTTP(recorderOk2, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusOK, recorderOk2.Code)
}

func BenchmarkRequestLimiterMW(b *testing.B) {
	//prepare logger
	logger := zap.NewNop()
	sugar := logger.Sugar()

	//prepare manager
	manager := NewRequestManager(2, time.Second*1)

	//prepare handler
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Some response"))
	})

	//test MW
	mw := RequestLimiterMW(*sugar, manager)
	handler := mw(next)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		//prepare response recorder
		recorder := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)

		//test
		b.StartTimer()
		handler.ServeHTTP(recorder, req)
	}
}
