package main

import (
	"flag"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"log"
	"net/http"
)

func init() {
	flag.StringVar(&config.ServerAddress, "a", "localhost:8080", "Address where server will work. Example: \"localhost:8080\".")
	flag.StringVar(&config.DefaultShortAddress, "b", "http://localhost:8080", "Base address before a shorted URL")
}

func main() {
	flag.Parse()

	r := handlers.MyRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, r))
}
