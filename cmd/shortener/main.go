package main

import (
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"github.com/spf13/pflag"
	"log"
	"net/http"
)

func init() {
	pflag.StringVarP(&config.ServerAddress, "address", "a", "localhost:8080", "Address where server will work. Example: \"localhost:8080\".")
	pflag.StringVarP(&config.DefaultShortAddress, "base_address", "b", "http://localhost:8080/", "Base address before a shorted URL")
}

func main() {
	pflag.Parse()

	r := handlers.MyRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, r))
}
