package main

import (
	"log"
	"net/http"

	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"go.uber.org/zap"
)

func main() {
	//conf
	conf := config.Config{}
	conf.Configurate()

	//storages set
	var URLStore handlers.URLStorageInterface
	if conf.DBConnString != "" {
		var err error
		URLStore, err = databases.NewPostgresql(conf.DBConnString)
		if err != nil {
			log.Fatalf("Problem with starting postgresql: %v", err.Error())
		}
	} else if conf.FileStoragePath != "" {
		URLStore = databases.NewJSONFileStorage(conf.FileStoragePath)
	} else {
		URLStore = databases.NewJustAMap()
	}

	//logger set
	logLevel, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	zCfg := zap.NewProductionConfig()
	zCfg.Level = logLevel
	zCfg.DisableStacktrace = true
	zapLogger, err := zCfg.Build()
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	sugar := zapLogger.Sugar()

	//server building
	r := handlers.NewRouter(conf, URLStore, *sugar)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, r))
}
