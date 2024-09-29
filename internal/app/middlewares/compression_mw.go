package middlewares

import (
	"bytes"
	"compress/gzip"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type gzipWriter struct {
	http.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressionMW(logger zap.SugaredLogger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			encoding := r.Header.Get("Content-Encoding")
			if encoding == "gzip" {
				//reading compressed data
				reader, err := gzip.NewReader(r.Body)
				if err != nil {
					logger.Error("Error while creating gzip reader", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}

				//decompressing
				decompressed, err := io.ReadAll(reader)
				if err != nil {
					logger.Error("Error while reading from gzip reader", zap.Error(err))
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				//logger.Debug(zap.Int("Decomressed data size", len(decompressed)))

				// data replacement
				r.Body.Close()
				r.Body = io.NopCloser(bytes.NewReader(decompressed))

			} else if encoding != "" {
				logger.Infof("Unsupported compression type `%s`", encoding)
				w.WriteHeader(http.StatusUnsupportedMediaType)
				return
			}

			//compressing
			for _, el := range r.Header.Values("Accept-Encoding") {
				if el == "gzip" {
					writer, err := gzip.NewWriterLevel(w, gzip.BestCompression)
					if err != nil {
						logger.Error("Error while creating new gzip writer", zap.Error(err))
						w.WriteHeader(http.StatusInternalServerError)
						return
					}
					defer writer.Close()

					w = gzipWriter{
						ResponseWriter: w,
						Writer:         writer,
					}

					w.Header().Set("Content-Encoding", "gzip")
					break
				}
			}

			//serving
			next.ServeHTTP(w, r)
		})
	}

}
