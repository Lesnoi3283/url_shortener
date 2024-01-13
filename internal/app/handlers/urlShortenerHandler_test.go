package handlers

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
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
			statusWant: 400,
		},
		{
			name:       "Bad request",
			query:      "/",
			statusWant: 400,
		},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, tt.query, nil)
		recorder := httptest.NewRecorder()
		URLShortenerHandler(recorder, req)
		res := recorder.Result()
		require.Equal(t, tt.statusWant, res.StatusCode, tt.name)
	}
}
