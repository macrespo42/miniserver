package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	router := http.NewServeMux()
	server := http.Server{
		Handler: router,
		Addr:    ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("can't start miniserver: ", err.Error())
	}

	fmt.Printf("listening on http://localhost%s\n", server.Addr)
}
