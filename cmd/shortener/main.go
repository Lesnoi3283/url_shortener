package main

import (
	"fmt"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/Lesnoi3283/url_shortener/pkg/databases"
	"go.uber.org/zap"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"net/http"
)

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

// naOrValue returns "N/A" if v contains a default value. Returns v if not.
func naOrValue(v string) string {
	if v == "" {
		return "N/A"
	} else {
		return v
	}
}

func main() {

	fmt.Printf("Build version: %s\n", naOrValue(buildVersion))
	fmt.Printf("Build date: %s\n", naOrValue(buildDate))
	fmt.Printf("Build commit: %s\n", naOrValue(buildCommit))

	//conf
	conf := config.Config{}
	conf.Configure()

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
	var listener net.Listener

	//HTTPS configuration
	if conf.EnableHTTPS {
		sugar.Info("Starting HTTPS server")
		listener = autocert.NewListener("urlshortener.ru")
	} else {
		sugar.Info("Starting HTTP server")
		listener, err = net.Listen("tcp", conf.ServerAddress)
		if err != nil {
			sugar.Fatalf("cant create a default listener, err: %v", err)
		}
	}

	//server starting
	r := handlers.NewRouter(conf, URLStore, *sugar)
	sugar.Fatal(http.Serve(listener, r))
}
