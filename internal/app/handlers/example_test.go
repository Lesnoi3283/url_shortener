package handlers

import (
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
)

func ExampleShortenHandler_ServeHTTP() {

	//Server:

	r := chi.NewRouter()

	//prepare a storage
	storage := databases.NewJustAMap()
	conf := config.Config{
		BaseAddress: "baseAddress",
		LogLevel:    "info", // different fields are not needed in this example
	}

	//prepare a logger
	loggerConf := zap.NewDevelopmentConfig()
	level, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("failed to parse log level: %s", conf.LogLevel)
	}
	loggerConf.Level = level
	logger, err := loggerConf.Build()
	if err != nil {
		log.Fatalf("failed to build logger: %v", err)
	}
	sugar := logger.Sugar()

	//prepare some handler
	handler := ShortenHandler{
		URLStorage: storage,
		Conf:       conf,
		Log:        *sugar,
	}

	//add this handler to router
	r.Post("/shorten", handler.ServeHTTP)

	//run server
	server := httptest.NewServer(r)
	defer server.Close()

	//Client:

	//prepare request
	req, err := http.NewRequest(http.MethodPost, server.URL+"/shorten", strings.NewReader("{\"url\": \"www.somelongurl.ru/loooooong\"}"))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	//send request
	resp, err := server.Client().Do(req)
	if err != nil {
		sugar.Errorf("Error while sending a request: %v", err)
	}
	defer resp.Body.Close()

	//see the answer
	fmt.Println("Response sent, status:")
	fmt.Println(resp.StatusCode)
	// Output:
	// Response sent, status:
	// 201
	body, _ := io.ReadAll(resp.Body)
	sugar.Infof("Response body: %s", string(body))
}
