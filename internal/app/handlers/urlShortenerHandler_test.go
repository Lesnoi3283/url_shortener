package handlers

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestURLShortenerGETHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		statusWant int
	}{
		{
			name:       "URL doesnt exist",
			query:      "/veryLongUrlWichShouldntExistIhopeForIt",
			statusWant: http.StatusBadRequest,
		},
		{
			name:       "Bad request",
			query:      "/",
			statusWant: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.query, nil)
		recorder := httptest.NewRecorder()
		URLShortenerHandler(recorder, req)
		res := recorder.Result()
		defer res.Body.Close()

		require.Equal(t, tt.statusWant, res.StatusCode, tt.name)
	}
}

func TestURLShortenerPOSTHandler(t *testing.T) {
	tests := []struct {
		name       string
		query      string
		body       string
		statusWant int
	}{
		{
			name:       "Normal one",
			query:      "/",
			body:       "https://practicum.yandex.ru/",
			statusWant: http.StatusCreated,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodPost, tt.query, strings.NewReader(tt.body))
		recorder := httptest.NewRecorder()
		URLShortenerHandler(recorder, req)
		res := recorder.Result()

		defer res.Body.Close()

		require.Equal(t, tt.statusWant, res.StatusCode, tt.name)
	}
}
