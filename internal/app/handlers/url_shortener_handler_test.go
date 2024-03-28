package handlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLShortenerHandler(t *testing.T) {

	tests := []struct {
		name          string
		query         string
		method        string
		reqBody       string
		statusWant    int
		wantEmptyBody bool
	}{
		{
			name:          "Normal POST (should work)",
			query:         "/",
			method:        http.MethodPost,
			statusWant:    http.StatusCreated,
			reqBody:       "https://practicum.yandex.ru/",
			wantEmptyBody: false,
		},
		{
			name:          "url doesnt exist",
			query:         "/veryLongUrlWichShouldntExistIhopeForIt",
			method:        http.MethodGet,
			statusWant:    http.StatusBadRequest,
			wantEmptyBody: true,
		},
		{
			name:          "Method not allowed request (chi error)",
			query:         "/",
			method:        http.MethodGet,
			statusWant:    http.StatusMethodNotAllowed,
			wantEmptyBody: true,
		},
	}

	//test server building
	conf := config.Config{
		BaseAddress:   "http://localhost:8080",
		ServerAddress: "localhost:8080",
		LogLevel:      "info",
	}
	URLStore := databases.NewJustAMap()
	logLevel, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}

	zCfg := zap.NewProductionConfig()
	zCfg.Level = logLevel
	zapLogger, err := zCfg.Build()
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	defer zapLogger.Sync()
	sugar := zapLogger.Sugar()
	mockControlelr := gomock.NewController(t)
	db := mocks.NewMockDBInterface(mockControlelr)
	ts := httptest.NewServer(NewRouter(conf, URLStore, *sugar, db))

	//tests run
	for _, tt := range tests {
		req, err := http.NewRequest(tt.method, ts.URL+tt.query, strings.NewReader(tt.reqBody))
		require.NoError(t, err, tt.name)

		resp, err := ts.Client().Do(req)
		require.NoError(t, err, tt.name)
		assert.Equal(t, tt.statusWant, resp.StatusCode, tt.name)

		defer resp.Body.Close()

		//redirect check
		if resp.StatusCode == http.StatusCreated {
			//here we try to get a full url back
			require.NotEmpty(t, resp.Body, tt.name)
			shortedURL, err := io.ReadAll(resp.Body)
			require.NoError(t, err, tt.name)

			//we need to split it because server returns us smg like  "127.0.0.1:8080/qqqq", but our port can be different.
			//so we need to get just `shorted url part` (for example "qqqq" from "127.0.0.1:8080/qqqq") from a full address
			splittedURL := strings.Split(string(shortedURL), "/")
			urlToAsk := ts.URL + "/" + splittedURL[len(splittedURL)-1]

			//to catch redirect
			ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
				assert.Equal(t, tt.reqBody, req.URL.String(), tt.name+" (in check redirect)")
				return http.ErrUseLastResponse
			}

			req2, err := http.NewRequest(http.MethodGet, urlToAsk, nil)
			require.NoError(t, err, tt.name)

			resp2, err := ts.Client().Do(req2)
			require.NoError(t, err, tt.name)
			defer resp2.Body.Close()

			assert.Equal(t, http.StatusTemporaryRedirect, resp2.StatusCode, tt.name)
		}
	}
}
