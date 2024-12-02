package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func filterWords(words []string, text string) string {
	for _, word := range words {
		splitted := strings.Split(text, " ")
		for index, element := range splitted {
			if strings.ToLower(element) == word {
				splitted[index] = "****"
			}
		}
		text = strings.Join(splitted, " ")
	}
	return text
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	type responseError struct {
		Error string `json:"error"`
	}

	errBody := responseError{
		Error: msg,
	}
	dat, err := json.Marshal(errBody)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}
	w.Header().Add("Content-Type", "encoding/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshalling json: %s", err)
		w.WriteHeader(500)
		return
	}

	w.Header().Add("Content-Type", "encoding/json")
	w.WriteHeader(code)
	w.Write(dat)
}

func HandlerHealth(w http.ResponseWriter, _ *http.Request) {
	body := []byte("OK")
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write(body)
}

func HandlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	params := parameters{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		type response struct {
			CleanedBody string `json:"cleaned_body"`
		}
		forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
		filteredBody := filterWords(forbiddenWords, params.Body)

		responseBody := response{
			CleanedBody: filteredBody,
		}

		respondWithJSON(w, 200, responseBody)
	}
}

func (cfg *ApiConfig) HandleCreateUser(w http.ResponseWriter, r *http.Request) {
	type createUserParams struct {
		Email string `json:"email"`
	}
	decoder := json.NewDecoder(r.Body)
	params := createUserParams{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}

	user, err := cfg.db.CreateUser(context.Background(), params.Email)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	userJson := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email:     user.Email,
	}
	respondWithJSON(w, 201, userJson)
}
