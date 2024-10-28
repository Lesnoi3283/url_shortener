package handlers

import (
	"encoding/json"
	"errors"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testStatsData struct {
	URLs  int `json:"urls"`
	Users int `json:"users"`
}

func TestStatsHandler_ServeHTTP(t *testing.T) {

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare mocks
	c := gomock.NewController(t)

	//prepare data
	correctData := testStatsData{
		URLs:  1500,
		Users: 500,
	}
	correctJSONData, err := json.Marshal(correctData)
	require.NoError(t, err, "error while marshalling test data")

	type fields struct {
		log     *zap.SugaredLogger
		storage URLStorageInterface
	}
	type args struct {
		w http.ResponseWriter
		r *http.Request
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		statusWant int
		dataWant   []byte
	}{
		{
			name: "ok",
			fields: fields{
				log: sugar,
				storage: func() URLStorageInterface {
					storage := mocks.NewMockURLStorageInterface(c)
					storage.EXPECT().GetShortURLCount(gomock.Any()).Return(correctData.URLs, nil)
					storage.EXPECT().GetUserCount(gomock.Any()).Return(correctData.Users, nil)
					return storage
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/api/internal/stats", nil),
			},
			statusWant: http.StatusOK,
			dataWant:   correctJSONData,
		},
		{
			name: "db error",
			fields: fields{
				log: sugar,
				storage: func() URLStorageInterface {
					storage := mocks.NewMockURLStorageInterface(c)
					storage.EXPECT().GetShortURLCount(gomock.Any()).Return(0, errors.New("test db error"))
					return storage
				}(),
			},
			args: args{
				w: httptest.NewRecorder(),
				r: httptest.NewRequest("GET", "/api/internal/stats", nil),
			},
			statusWant: http.StatusInternalServerError,
			dataWant:   nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &StatsHandler{
				log:     *tt.fields.log,
				storage: tt.fields.storage,
			}
			h.ServeHTTP(tt.args.w, tt.args.r)
		})
	}
}
