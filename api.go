package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func HandlerHealth(w http.ResponseWriter, _ *http.Request) {
	body := []byte("OK")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(body)
}

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
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
