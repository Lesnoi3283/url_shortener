package middlewares

import (
	"net/http"
	"time"

	"go.uber.org/zap"
)

type responceData struct {
	size   int
	status int
}

type loggingResponceWriter struct {
	data responceData
	rw   http.ResponseWriter
}

func (l *loggingResponceWriter) Write(b []byte) (int, error) {
	size, err := l.rw.Write(b)
	l.data.size += size
	return size, err
}

func (l *loggingResponceWriter) WriteHeader(statusCode int) {
	l.rw.WriteHeader(statusCode)
	l.data.status = statusCode
}

func (l *loggingResponceWriter) Header() http.Header {
	return l.rw.Header()
}

// LoggerMW logs request`s params: URL, method, duration.
// And response`s params: status code and size.
func LoggerMW(logger zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			lw := loggingResponceWriter{
				data: responceData{},
				rw:   w,
			}

			next.ServeHTTP(&lw, r)

			duration := time.Since(start)

			logger.Info("request", zap.String("url", r.URL.String()), zap.String("method", r.Method), zap.Duration("duration", duration))
			logger.Info("response", zap.Int("status code", lw.data.status), zap.Int("size", lw.data.size))

		})
	}
}
