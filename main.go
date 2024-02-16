package main

import (
	"log"
	"net/http"
)

func getMultiplexer() *http.ServeMux {
	mux := http.NewServeMux()

	// Mux handlers

	return mux
}

func main() {
	log.Printf("VEEEKTOR_api is starting...")
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	mux := getMultiplexer()
	server := &http.Server{Addr: ":8080", Handler: mux}
	defer server.Close()

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}
