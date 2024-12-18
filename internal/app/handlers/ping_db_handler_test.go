package handlers

import (
	"database/sql"
	"github.com/Lesnoi3283/url_shortener/internal/app/logic/mocks"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_pingDBHandler_ServeHTTP(t *testing.T) {

	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	tests := []struct {
		name          string
		mockSetupFunc func(dbInterface *mocks.MockURLStorageInterface)
		statusWant    int
	}{
		{
			name: "DB works",
			mockSetupFunc: func(dbInterface *mocks.MockURLStorageInterface) {
				dbInterface.EXPECT().Ping().Return(nil)
			},
			statusWant: http.StatusOK,
		},
		{
			name: "DB doesnt work",
			mockSetupFunc: func(dbInterface *mocks.MockURLStorageInterface) {
				dbInterface.EXPECT().Ping().Return(sql.ErrConnDone)
			},
			statusWant: http.StatusInternalServerError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockController := gomock.NewController(t)
			db := mocks.NewMockURLStorageInterface(mockController)

			tt.mockSetupFunc(db)

			p := &PingDBHandler{
				DB:  db,
				log: *sugar,
			}

			req, err := http.NewRequest("GET", "/ping", nil)
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()

			p.ServeHTTP(rr, req)
			assert.Equal(t, tt.statusWant, rr.Code, "Status code is not equal.")
		})
	}
}

func BenchmarkPingDBHandler_ServeHTTP(b *testing.B) {
	mockController := gomock.NewController(b)
	defer mockController.Finish()

	db := mocks.NewMockURLStorageInterface(mockController)
	db.EXPECT().Ping().Return(nil).AnyTimes()

	logger := zaptest.NewLogger(b)
	sugar := logger.Sugar()

	handler := PingDBHandler{
		DB:  db,
		log: *sugar,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		req := httptest.NewRequest("GET", "/ping", nil)
		b.StartTimer()

		handler.ServeHTTP(httptest.NewRecorder(), req)
	}
}
