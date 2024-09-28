package middlewares

import (
	"bytes"
	"compress/gzip"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCompressionMW_WithGzip(t *testing.T) {
	//prepare data
	correctData := "some correct data"
	var dataBuf bytes.Buffer
	writer := gzip.NewWriter(&dataBuf)
	writer.Write([]byte(correctData))
	writer.Close()
	data := bytes.NewReader(dataBuf.Bytes())

	//prepare request and response recorder
	r := httptest.NewRequest(http.MethodPost, "/", data)
	r.Header.Set("Content-Encoding", "gzip")
	r.Header.Set("Accept-Encoding", "gzip")

	w := httptest.NewRecorder()

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		receivedData, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		assert.Equal(t, correctData, string(receivedData))

		_, err = w.Write([]byte(correctData))
		require.NoError(t, err, "Error while writing response")
	})

	//testing MW
	mw := CompressionMW(*sugar)
	mw(nextHandler).ServeHTTP(w, r)

	//check if result is encoded
	require.Equal(t, "gzip", w.Header().Get("Content-Encoding"))
	compressedData := w.Body.Bytes()

	//decoding
	reader, err := gzip.NewReader(bytes.NewReader(compressedData))
	require.NoError(t, err, "Error while creating gzip reader")
	decompressedData, err := io.ReadAll(reader)
	require.NoError(t, err, "Error while reading from gzip reader")
	err = reader.Close()
	require.NoError(t, err, "Error while closing gzip reader")

	//result check
	assert.Equal(t, correctData, string(decompressedData))
}

func TestCompressionMW_UnsupportedEncoding(t *testing.T) {
	//prepare request and response recorder
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	r.Header.Set("Content-Encoding", "SomeUnsupportedEncoding")

	w := httptest.NewRecorder()

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		t.Error("Handler should not be called")
	})

	//testing MW
	mw := CompressionMW(*sugar)
	mw(nextHandler).ServeHTTP(w, r)

	//check result
	assert.Equal(t, http.StatusUnsupportedMediaType, w.Code)
}

func TestCompressionMW_NoEncoding(t *testing.T) {
	//prepare data
	correctData := "some correct data"

	//prepare request and response recorder
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	//prepare logger
	logger := zaptest.NewLogger(t)
	sugar := logger.Sugar()

	//prepare handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(correctData))
		require.NoError(t, err, "Error while writing response")
	})

	//testing MW
	mw := CompressionMW(*sugar)
	mw(nextHandler).ServeHTTP(w, r)

	//check if result is encoded (should be NOT)
	require.Equal(t, "", w.Header().Get("Content-Encoding"))
	dataGot := w.Body.Bytes()

	//result check
	assert.Equal(t, correctData, string(dataGot))
}

func BenchmarkCompressionMW(b *testing.B) {
	//prepare data

	//prepare logger
	logger := zaptest.NewLogger(b)
	sugar := logger.Sugar()

	//prepare handler
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	//testing MW
	mw := CompressionMW(*sugar)
	testable := mw(nextHandler)

	b.ResetTimer()

	b.Run("With GZIP", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			b.StopTimer()
			//prepare request and response recorder
			correctData := "some correct data"
			var buf bytes.Buffer
			writer := gzip.NewWriter(&buf)
			writer.Write([]byte(correctData))
			writer.Close()
			data := bytes.NewReader(buf.Bytes())
			requestWithGZIP := httptest.NewRequest(http.MethodPost, "/", data)
			requestWithGZIP.Header.Set("Accept-Encoding", "gzip")
			requestWithGZIP.Header.Set("Content-Encoding", "gzip")
			w := httptest.NewRecorder()

			//test
			b.StartTimer()
			testable.ServeHTTP(w, requestWithGZIP)
		}
	})

	b.Run("No encoding", func(b *testing.B) {
		for i := 0; i < 100; i++ {
			b.StopTimer()
			//prepare request and response recorder
			requestNoEncoding := httptest.NewRequest(http.MethodPost, "/", nil)
			w := httptest.NewRecorder()

			//test
			b.StartTimer()
			testable.ServeHTTP(w, requestNoEncoding)
		}
	})
}
