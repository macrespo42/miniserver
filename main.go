package main

import (
	"log"
	"net/http"
)

func main() {
	apiCfg := ApiConfig{}
	mux := http.NewServeMux()
	rootHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.MiddlewareMetricsInc(rootHandler))

	mux.HandleFunc("GET /app/assets/logo.png", HandlerImage)

	mux.HandleFunc("GET /api/healthz", HandlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", HandlerValidateChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.HandleReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.HandleServerHits)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
