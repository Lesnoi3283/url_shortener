package main

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/justamap"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	//conf
	conf := config.Config{}
	conf.Configurate()

	//storages set
	URLStore := justamap.NewJustAMap()

	//logger set
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
	sugar := zapLogger.Sugar()

	//server building
	r := handlers.BuildRouter(conf, URLStore, *sugar)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, r))
}
