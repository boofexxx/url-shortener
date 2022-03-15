package main

import (
	"log"
	"net/http"
)

func main() {
	mux, err := NewServerMux()
	if err != nil {
		log.Fatal(err)
	}

	mux.HandleFunc("/a/", mux.CreateShortURLHandler).
		Methods(http.MethodGet)
	mux.HandleFunc("/s/{shortURL}", mux.GetShortURLHandler).
		Methods(http.MethodGet)

	mux.Use(mux.LogMiddleware)

	mux.log.Print("started listening :8080")
	http.ListenAndServe(":8080", mux)
}
