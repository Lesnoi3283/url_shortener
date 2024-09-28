package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDeleteURLsHandler_ServeHTTP(t *testing.T) {

	//prepare data
	coorectUserID := 1
	URLsToDelete := make([]string, 0)
	URLsToDelete = append(URLsToDelete, "http://someurl")
	URLsToDelete = append(URLsToDelete, "http://someurl2")
	correrctData, err := json.Marshal(URLsToDelete)
	require.NoError(t, err, "Error while marshalling json (data preparation in test)")

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare mocks
	c := gomock.NewController(t)
	defer c.Finish()

	type fields struct {
		URLStorage URLStorageInterface
		//Conf       config.Config
		Log zap.SugaredLogger
	}
	type args struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		statusWant int
	}{
		{
			name: "Ok",
			fields: fields{
				URLStorage: func() URLStorageInterface {
					URLsToDeleteChan := make(chan string)
					storage := mocks.NewMockURLStorageInterface(c)
					storage.EXPECT().DeleteBatchWithUserID(coorectUserID).Return(URLsToDeleteChan, nil)
					return storage
				}(),
				Log: *sugar,
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(correrctData)).WithContext(context.WithValue(context.Background(), middlewares.UserIDContextKey, coorectUserID)),
			},
			statusWant: http.StatusAccepted,
		},
		{
			name: "No userID",
			fields: fields{
				URLStorage: nil,
				Log:        *sugar,
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(correrctData)),
			},
			statusWant: http.StatusUnauthorized,
		},
		{
			name: "Bad request",
			fields: fields{
				URLStorage: nil,
				Log:        *sugar,
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodDelete, "/api/user/urls", strings.NewReader("{basJSON:")).WithContext(context.WithValue(context.Background(), middlewares.UserIDContextKey, coorectUserID)),
			},
			statusWant: http.StatusBadRequest,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &DeleteURLsHandler{
				URLStorage: tt.fields.URLStorage,
				Log:        tt.fields.Log,
			}
			h.ServeHTTP(tt.args.res, tt.args.req)

			assert.Equal(t, tt.statusWant, tt.args.res.Code)
		})
	}
}

func BenchmarkDeleteURLsHandler_ServeHTTP(b *testing.B) {

	//prepare data
	coorectUserID := 1
	URLsToDelete := make([]string, 0)
	URLsToDelete = append(URLsToDelete, "http://someurl")
	URLsToDelete = append(URLsToDelete, "http://someurl2")
	correrctData, err := json.Marshal(URLsToDelete)
	require.NoError(b, err, "Error while marshalling json (data preparation in test)")

	//prepare logger
	logger := zaptest.NewLogger(b)
	sugar := logger.Sugar()

	//prepare mocks
	c := gomock.NewController(b)
	defer c.Finish()

	storage := mocks.NewMockURLStorageInterface(c)
	storage.EXPECT().DeleteBatchWithUserID(coorectUserID).Return(make(chan string), nil).AnyTimes()

	//prepare handler
	h := DeleteURLsHandler{
		URLStorage: storage,
		Log:        *sugar,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		req := httptest.NewRequest(http.MethodDelete, "/api/user/urls", bytes.NewReader(correrctData)).WithContext(context.WithValue(context.Background(), middlewares.UserIDContextKey, coorectUserID))
		res := httptest.NewRecorder()

		b.StartTimer()
		h.ServeHTTP(res, req)
	}
}
