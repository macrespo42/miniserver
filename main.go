package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
)

type apiConfig struct {
	fileServerHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cfg.fileServerHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleServerHits(w http.ResponseWriter, r *http.Request) {
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

func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileServerHits.Store(0)
	hits := cfg.fileServerHits.Load()
	body := []byte(fmt.Sprintf("Hit: %d", hits))
	w.WriteHeader(200)
	w.Write(body)
}

func handlerImage(w http.ResponseWriter, r *http.Request) {
	buf, err := os.ReadFile("./assets/logo.png")
	if err != nil {
		log.Fatal("can't read file")
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buf)))
	w.Write(buf)
}

func handlerHealth(w http.ResponseWriter, _ *http.Request) {
	body := []byte("OK")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(body)
}

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	type responseError struct {
		Error string `json:"error"`
	}

	if err != nil {
		errBody := responseError{
			Error: "Something went wrong",
		}
		dat, err := json.Marshal(errBody)
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Type", "encoding/json")
		w.WriteHeader(500)
		w.Write(dat)
	}

	if len(params.Body) > 140 {
		errBody := responseError{
			Error: "Chirp is too long",
		}
		dat, err := json.Marshal(errBody)
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}
		w.Header().Add("Content-Type", "encoding/json")
		w.WriteHeader(400)
		w.Write(dat)
	} else {
		type reponseValid struct {
			Valid bool `json:"valid"`
		}
		responseBody := reponseValid{
			Valid: true,
		}

		dat, err := json.Marshal(responseBody)
		if err != nil {
			log.Printf("Error marshalling json: %s", err)
			w.WriteHeader(500)
			return
		}

		w.Header().Add("Content-Type", "encoding/json")
		w.WriteHeader(200)
		w.Write(dat)
	}
}

func main() {
	apiCfg := apiConfig{}
	mux := http.NewServeMux()
	rootHandler := http.StripPrefix("/app/", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(rootHandler))

	mux.HandleFunc("GET /app/assets/logo.png", handlerImage)

	mux.HandleFunc("GET /api/healthz", handlerHealth)
	mux.HandleFunc("POST /api/validate_chirp", handlerValidateChirp)

	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)
	mux.HandleFunc("GET /admin/metrics", apiCfg.handleServerHits)

	server := http.Server{
		Handler: mux,
		Addr:    ":8080",
	}

	log.Fatal(server.ListenAndServe())
}
