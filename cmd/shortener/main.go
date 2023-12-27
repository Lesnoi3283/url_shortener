package main

import "net/http"

import "github.com/Lesnoi3283/url_shortener/internal/app/handlers"

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlers.URLShortenerHandler)

	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		panic(err)
	}
}
