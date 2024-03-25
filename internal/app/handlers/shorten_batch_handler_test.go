package handlers

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers/mocks"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/justamap"
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
	URLStore := justamap.NewJustAMap()
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
	mockController := gomock.NewController(t)
	db := mocks.NewMockDBInterface(mockController)
	ts := httptest.NewServer(BuildRouter(conf, URLStore, *sugar, db))

	// Запуск тестов
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req, err := http.NewRequest(tt.method, ts.URL+tt.query, strings.NewReader(tt.reqBody))
			require.NoError(t, err, "Creating request should not error")

			resp, err := http.DefaultClient.Do(req)
			require.NoError(t, err, "Making request should not error")

			defer resp.Body.Close()
			assert.Equal(t, tt.statusWant, resp.StatusCode, "Status code should match expected")

			if !tt.wantEmptyBody {
				body, err := io.ReadAll(resp.Body)
				require.NoError(t, err, "Reading response body should not error")
				require.NotEmpty(t, body, "Response body should not be empty")
			}
		})
	}
}
