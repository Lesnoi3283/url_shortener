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

func TestShortenBatchHandler_ServeHTTP(t *testing.T) {
	tests := []struct {
		name          string
		query         string
		method        string
		reqBody       string
		statusWant    int
		wantEmptyBody bool
	}{
		{
			name:       "Normal POST batch (should work)",
			query:      "/api/shorten/batch",
			method:     http.MethodPost,
			statusWant: http.StatusCreated,
			reqBody: `[
                {"correlation_id": "1", "original_url": "https://example.com"},
                {"correlation_id": "2", "original_url": "https://example.org"}
            ]`,
			wantEmptyBody: false,
		},
		{
			name:          "Bad request (empty body)",
			query:         "/api/shorten/batch",
			method:        http.MethodPost,
			reqBody:       "",
			statusWant:    http.StatusInternalServerError,
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

	zapTestLogger := zaptest.NewLogger(t)
	defer zapTestLogger.Sync()
	sugar := zapTestLogger.Sugar()
	ts := httptest.NewServer(NewRouter(conf, URLStore, *sugar))

	// Запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.query, strings.NewReader(tt.reqBody))
			require.NoError(t, err, "Error while creating a request")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "Error while making a request")

			defer resp.Body.Close()
			assert.Equal(t, tt.statusWant, resp.StatusCode, "Wrong status code")

			if !tt.wantEmptyBody {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err, "Reading response body should not error")
				require.NotEmpty(t, body, "Response body should not be empty")

				type URLShorten struct {
					CorrelationID string `json:"correlation_id"`
					ShortURL      string `json:"short_url"`
				}
				var URLsToReturn []URLShorten
				err = json.Unmarshal(body, &URLsToReturn)
				require.NoError(t, err, "Unmarshalling response error")

				for _, urlShort := range URLsToReturn {
					splittedURL := strings.Split(string(urlShort.ShortURL), "/")
					urlToAsk := ts.URL + "/" + splittedURL[len(splittedURL)-1]

					req2, err := http.NewRequest(http.MethodGet, urlToAsk, nil)
					require.NoError(t, err, "Error while making a request")

					//to catch redirect
					ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
						return http.ErrUseLastResponse
					}

					resp2, err := ts.Client().Do(req2)
					require.NoError(t, err)

					assert.Equal(t, http.StatusTemporaryRedirect, resp2.StatusCode)

					resp2.Body.Close()
				}
			}

		})
	}
}

func BenchmarkShortenBatchHandler_ServeHTTP(b *testing.B) {

	c := gomock.NewController(b)
	defer c.Finish()
	storage := mocks.NewMockURLStorageInterface(c)
	storage.EXPECT().SaveBatch(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	logger := zaptest.NewLogger(b)
	sugar := logger.Sugar()

	handler := ShortenBatchHandler{
		URLStorage: storage,
		Conf: config.Config{
			BaseAddress:   "http://localhost:8080",
			ServerAddress: "localhost:8080",
			LogLevel:      "info",
		},
		Log: *sugar,
	}

	reqBody := `[
                {"correlation_id": "1", "original_url": "https://example.com"},
                {"correlation_id": "2", "original_url": "https://example.org"}
            ]`

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/api/shorten/batch", strings.NewReader(reqBody)))
	}
}
