package main

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type ApiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *ApiConfig) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *ApiConfig) HandleServerHits(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileServerHits.Load()
	body := []byte(fmt.Sprintf(`<html>
    <body>
      <h1>Welcome, Chirpy Admin</h1>
      <p>Chirpy has been visited %d times!</p>
    </body>
  </html>`, hits))
	w.Header().Add("Content-Type", "text/html")
	w.WriteHeader(200)
	w.Write(body)
}

func (cfg *ApiConfig) HandleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	hits := cfg.fileServerHits.Load()
	body := []byte(fmt.Sprintf("Hit: %d", hits))
	w.WriteHeader(200)
	w.Write(body)
}
