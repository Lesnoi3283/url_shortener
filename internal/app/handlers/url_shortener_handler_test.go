package handlers

import (
	"github.com/Lesnoi3283/url_shortener/internal/app/logic/mocks"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
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

	zapTestLogger := zaptest.NewLogger(t)
	defer zapTestLogger.Sync()
	sugar := zapTestLogger.Sugar()

	r, err := NewRouter(conf, URLStore, *sugar)
	require.NoError(t, err, "error while creating a router in test")
	ts := httptest.NewServer(r)

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

func BenchmarkURLShortenerHandler_ServeHTTP(b *testing.B) {

	c := gomock.NewController(b)
	defer c.Finish()
	storage := mocks.NewMockURLStorageInterface(c)
	storage.EXPECT().Save(gomock.Any(), gomock.Any()).Return(nil).AnyTimes()

	conf := config.Config{
		BaseAddress:   "http://localhost:8080",
		ServerAddress: "localhost:8080",
		LogLevel:      "info",
	}

	handler := URLShortenerHandler{
		Conf:       conf,
		URLStorage: storage,
	}

	reqBody := "https://practicum.yandex.ru/"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		handler.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest("GET", "/", strings.NewReader(reqBody)))
	}
}
