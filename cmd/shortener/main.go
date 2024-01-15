package main

import (
	"github.com/Lesnoi3283/url_shortener/internal/app/handlers"
	"log"
	"net/http"
)

func main() {
	r := handlers.MyRouter()
	log.Fatal(http.ListenAndServe("127.0.0.1:8080", r))
}
