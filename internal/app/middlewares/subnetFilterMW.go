package middlewares

import (
	"go.uber.org/zap"
	"net"
	"net/http"
	"strings"
)

func SubnetFilterMW(subnetMask *net.IPNet, logger zap.SugaredLogger) func(handlerFunc http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			if strings.HasSuffix(r.URL.Path, "/api/internal/stats") {
				if subnetMask == nil {
					w.WriteHeader(http.StatusForbidden)
					logger.Debugf("request forbidden (all requests to this endpoint are forbidden, because trusted subnet is not set)")
				}

				ip := r.Header.Get("X-Real-IP")
				if !subnetMask.Contains(net.ParseIP(ip)) {
					logger.Debugf("request forbidden (IP is not in the allowed subnet)")
					w.WriteHeader(http.StatusForbidden)
					return
				}
				next.ServeHTTP(w, r)
			} else {
				next.ServeHTTP(w, r)
			}
		})
	}
}
