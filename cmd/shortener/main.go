package main

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/Lesnoi3283/url_shortener/internal/storages"
	"github.com/Lesnoi3283/url_shortener/pkg/databases/justamap"
	"log"
	"net/http"
)

func main() {
	conf := config.Config{}
	conf.Configurate()

	URLStore := &storages.URLStorage{
		DB: justamap.NewJustAMap(),
	}

	r := handlers.BuildRouter(conf, URLStore)
	log.Fatal(http.ListenAndServe(conf.ServerAddress, r))
}
