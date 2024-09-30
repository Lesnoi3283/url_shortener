package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/entities"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/Lesnoi3283/url_shortener/internal/app/middlewares"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

type testURLData struct {
	ShortURL    string `json:"short_url"`
	OriginalURL string `json:"original_url"`
}

func TestUserURLsHandler_ServeHTTP(t *testing.T) {

	//prepare data
	correctUserID := 1
	correctURLs := make([]entities.URL, 0)
	correctURLs = append(correctURLs, entities.URL{
		Short: "someUrlShort",
		Long:  "http://someUrlLong",
	}, entities.URL{
		Short: "someUrlShort2",
		Long:  "http://someUrlLong2",
	})

	conf := config.Config{
		BaseAddress: "http://baseAddress",
	}

	correctJWTToken, err := middlewares.BuildNewJWTString(correctUserID)
	require.NoErrorf(t, err, "error while building JWT in test: %v", err)
	buildARequestWithJWT := func(token string) *http.Request {
		//build a request and add a cookie to it
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req.AddCookie(&http.Cookie{
			Name:  middlewares.JwtCookieName,
			Value: token,
		})
		return req
	}

	//prepare mocks
	c := gomock.NewController(t)
	defer c.Finish()

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	type fields struct {
		URLStorage URLStorageInterface
		Conf       config.Config
		Logger     zap.SugaredLogger
		StatusWant int
		BodyWant   string
	}
	type args struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "ok",
			fields: fields{
				URLStorage: func() URLStorageInterface {
					//mock storage
					storage := mocks.NewMockURLStorageInterface(c)
					storage.EXPECT().GetUserUrls(gomock.Any(), correctUserID).Return(correctURLs, nil)
					return storage
				}(),
				Conf:       conf,
				Logger:     *sugar,
				StatusWant: http.StatusOK,
				BodyWant: func() string {
					//marshal URLs (using local structure because it has required JSON tags)
					URLs := make([]testURLData, 0)
					URLs = append(URLs, testURLData{
						ShortURL:    conf.BaseAddress + "/" + correctURLs[0].Short,
						OriginalURL: correctURLs[0].Long,
					}, testURLData{
						ShortURL:    conf.BaseAddress + "/" + correctURLs[1].Short,
						OriginalURL: correctURLs[1].Long,
					})
					JSON, err := json.Marshal(URLs)
					require.NoError(t, err, "error while marshalling test URLs for response")
					return string(JSON)
				}(),
			},
			args: args{
				res: httptest.NewRecorder(),
				req: buildARequestWithJWT(correctJWTToken),
			},
		},
		{
			name: "NotAuth",
			fields: fields{
				URLStorage: nil,
				Conf:       conf,
				Logger:     *sugar,
				StatusWant: http.StatusUnauthorized,
				BodyWant:   "",
			},
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "/api/user/urls", nil),
			},
		},
		{
			name: "Not valid JWT",
			fields: fields{
				URLStorage: nil,
				Conf:       conf,
				Logger:     *sugar,
				StatusWant: http.StatusUnauthorized,
				BodyWant:   "",
			},
			args: args{
				res: httptest.NewRecorder(),
				req: buildARequestWithJWT("not a correct JWT token"),
			},
		},
		{
			name: "No URLs",
			fields: fields{
				URLStorage: func() URLStorageInterface {
					storage := mocks.NewMockURLStorageInterface(c)
					storage.EXPECT().GetUserUrls(gomock.Any(), correctUserID).Return(make([]entities.URL, 0), nil)
					return storage
				}(),
				Conf:       conf,
				Logger:     *sugar,
				StatusWant: http.StatusNoContent,
				BodyWant:   "",
			},
			args: args{
				res: httptest.NewRecorder(),
				req: buildARequestWithJWT(correctJWTToken),
			},
		},
		{
			name: "DB error",
			fields: fields{
				URLStorage: func() URLStorageInterface {
					storage := mocks.NewMockURLStorageInterface(c)
					storage.EXPECT().GetUserUrls(gomock.Any(), correctUserID).Return(make([]entities.URL, 0), errors.New("db error"))
					return storage
				}(),
				Conf:       conf,
				Logger:     *sugar,
				StatusWant: http.StatusInternalServerError,
				BodyWant:   "",
			},
			args: args{
				res: httptest.NewRecorder(),
				req: buildARequestWithJWT(correctJWTToken),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &UserURLsHandler{
				URLStorage: tt.fields.URLStorage,
				Conf:       tt.fields.Conf,
				Logger:     tt.fields.Logger,
			}
			h.ServeHTTP(tt.args.res, tt.args.req)
			assert.Equal(t, tt.fields.StatusWant, tt.args.res.Code)
			assert.Equal(t, tt.fields.BodyWant, string(tt.args.res.Body.Bytes()))
		})
	}
}

func BenchmarkUserURLsHandler_ServeHTTP(b *testing.B) {

	//prepare data
	correctUserID := 1
	correctURLs := make([]entities.URL, 0)
	correctURLs = append(correctURLs, entities.URL{
		Short: "someUrlShort",
		Long:  "http://someUrlLong",
	}, entities.URL{
		Short: "someUrlShort2",
		Long:  "http://someUrlLong2",
	})

	conf := config.Config{
		BaseAddress: "http://baseAddress",
	}

	correctJWTToken, err := middlewares.BuildNewJWTString(correctUserID)
	require.NoErrorf(b, err, "error while building JWT in test: %v", err)
	buildARequestWithJWT := func(token string) *http.Request {
		//build a request and add a cookie to it
		req := httptest.NewRequest(http.MethodGet, "/api/user/urls", nil)
		req.AddCookie(&http.Cookie{
			Name:  middlewares.JwtCookieName,
			Value: token,
		})
		return req
	}

	//prepare mocks
	c := gomock.NewController(b)
	defer c.Finish()
	storage := mocks.NewMockURLStorageInterface(c)
	storage.EXPECT().GetUserUrls(gomock.Any(), correctUserID).Return(correctURLs, nil).AnyTimes()

	//prepare logger
	logger := zaptest.NewLogger(b)
	sugar := logger.Sugar()

	handler := UserURLsHandler{
		URLStorage: storage,
		Conf:       conf,
		Logger:     *sugar,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		//build a request
		b.StopTimer()
		req := buildARequestWithJWT(correctJWTToken)
		b.StartTimer()

		//test
		handler.ServeHTTP(httptest.NewRecorder(), req)
	}
}
