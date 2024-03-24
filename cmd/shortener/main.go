package main

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/jsonfilestorage"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/justamap"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/postgresql"
	"go.uber.org/zap"
	"log"
	"net/http"
)

func main() {
	//conf
	conf := config.Config{}
	conf.Configurate()

	//storages set
	var URLStore handlers.URLStorageInterface
	if conf.FileStoragePath == "" {
		URLStore = justamap.NewJustAMap()
	} else {
		URLStore = jsonfilestorage.NewJSONFileStorage(conf.FileStoragePath)
	}

	//logger set
	logLevel, err := zap.ParseAtomicLevel(conf.LogLevel)
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}

	//config set
	zCfg := zap.NewProductionConfig()
	zCfg.Level = logLevel
	zapLogger, err := zCfg.Build()
	if err != nil {
		log.Fatalf("logger was not started, err: %v", err)
	}
	sugar := zapLogger.Sugar()

	//db set
	db, err := postgresql.NewPostgresql(conf.DBConnString)
	if err != nil {
		log.Printf("db was not started, err: %v", err)
	}

	//server building
	r := handlers.BuildRouter(conf, URLStore, *sugar, &db)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, r))
}
