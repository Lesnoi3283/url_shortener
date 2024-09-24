package handlers

import (
	"encoding/json"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLShortenHandler_ServeHTTP(t *testing.T) {

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
			query:         "/api/shorten",
			method:        http.MethodPost,
			statusWant:    http.StatusCreated,
			reqBody:       "{\"url\":\"https://practicum.yandex.ru\"}",
			wantEmptyBody: false,
		},
	}

	//test server building
	conf := config.Config{
		BaseAddress:   "http://localhost:8080",
		ServerAddress: "localhost:8080",
		LogLevel:      "info",
	}

	URLStore := databases.NewJustAMap()

	zapTestLogger := zaptest.NewLogger(t)
	defer zapTestLogger.Sync()
	sugar := zapTestLogger.Sugar()

	ts := httptest.NewServer(NewRouter(conf, URLStore, *sugar))

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
			require.NotEmpty(t, resp.Body, tt.name)
			response, err := io.ReadAll(resp.Body)
			require.NoError(t, err, tt.name)

			result := struct {
				Result string `json:"result"`
			}{}
			require.NoError(t, json.Unmarshal(response, &result), "Error while unmarshalling json response")
		}
	}
}

func BenchmarkShortenHandler_ServeHTTP(b *testing.B) {
	c := gomock.NewController(b)
	defer c.Finish()
	storage := mocks.NewMockURLStorageInterface(c)
	storage.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	logger := zaptest.NewLogger(b)
	sugar := logger.Sugar()

	handler := ShortenHandler{
		URLStorage: storage,
		Conf: config.Config{
			BaseAddress:   "http://localhost:8080",
			ServerAddress: "localhost:8080",
			LogLevel:      "info",
		},
		Log: *sugar,
	}

	reqBody := "{\"url\":\"https://practicum.yandex.ru\"}"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/api/shorten", strings.NewReader(reqBody)))
	}
}
