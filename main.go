package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("can't start miniserver: ", err.Error())
	}
}
