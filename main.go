package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func handlerImage(w http.ResponseWriter, r *http.Request) {
	buf, err := os.ReadFile("./assets/logo.png")
	if err != nil {
		log.Fatal("can't read file")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	w.Write(buf)
}

func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(".")))
	mux.HandleFunc("GET /assets/logo.png", handlerImage)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
