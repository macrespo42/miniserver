package main

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

func HandlerImage(w http.ResponseWriter, r *http.Request) {
	buf, err := os.ReadFile("./assets/logo.png")
	if err != nil {
		log.Fatal("can't read file")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	w.Write(buf)
}
