package middlewares

import (
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSubnetFilterMW_allowedAddress(t *testing.T) {
	//prepare data
	_, allowedSubnet, err := net.ParseCIDR("192.168.1.0/24")
	require.NoError(t, err, "error while parsing CIDR (in test prepare)")

	allowedIP := "192.168.1.1"
	notAllowedIP := "110.110.1.1"
	badIP := "127.badip0.1"
	protectedTarget := "/api/internal/stats"

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//prepare tests
	tests := []struct {
		name       string
		req        *http.Request
		resp       *httptest.ResponseRecorder
		statusWant int
	}{
		{
			name: "ok",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, protectedTarget, nil)
				r.Header.Set("X-Real-IP", allowedIP)
				return r
			}(),
			resp:       httptest.NewRecorder(),
			statusWant: http.StatusOK,
		},
		{
			name: "not allowed IP",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, protectedTarget, nil)
				r.Header.Set("X-Real-IP", notAllowedIP)
				return r
			}(),
			resp:       httptest.NewRecorder(),
			statusWant: http.StatusForbidden,
		},
		{
			name: "bad ip",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, protectedTarget, nil)
				r.Header.Set("X-Real-IP", badIP)
				return r
			}(),
			resp:       httptest.NewRecorder(),
			statusWant: http.StatusForbidden,
		},
		{
			name: "not allowed ip, but target is not protected",
			req: func() *http.Request {
				r := httptest.NewRequest(http.MethodGet, "/", nil)
				r.Header.Set("X-Real-IP", notAllowedIP)
				return r
			}(),
			resp:       httptest.NewRecorder(),
			statusWant: http.StatusOK,
		},
	}

	//test
	mw := SubnetFilterMW(allowedSubnet, *sugar)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mw(handler).ServeHTTP(tt.resp, tt.req)
			require.Equal(t, tt.statusWant, tt.resp.Code)
		})
	}
}
