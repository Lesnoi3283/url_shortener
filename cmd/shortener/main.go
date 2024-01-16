package main

import (
	"flag"
	"github.com/Lesnoi3283/url_shortener/config"
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"log"
	"net/http"
	"os"
)

func init() {
	flag.StringVar(&config.ServerAddress, "a", config.DefaultServerAddress, "Address where server will work. Example: \"localhost:8080\".")
	flag.StringVar(&config.BaseAddress, "b", config.DefaultBaseAddress, "Base address before a shorted URL")
}

func main() {

	flag.Parse()
	envServerAddress, wasFoundServerAddress := os.LookupEnv("SERVER_ADDRESS")
	envBaseAddress, wasFoundBaseAddress := os.LookupEnv("BASE_URL")

	if config.ServerAddress == config.DefaultServerAddress && wasFoundServerAddress {
		config.ServerAddress = envServerAddress
	}
	if config.BaseAddress == config.DefaultBaseAddress && wasFoundBaseAddress {
		config.BaseAddress = envBaseAddress
	}

	r := handlers.MyRouter()
	log.Fatal(http.ListenAndServe(config.ServerAddress, r))
}
