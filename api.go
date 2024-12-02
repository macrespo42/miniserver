package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/macrespo42/miniserver/internal/database"
)

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

func (cfg *ApiConfig) HandleCreateChirp(w http.ResponseWriter, r *http.Request) {
	type createChirpArgs struct {
		Body   string    `json:"body"`
		UserId uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	params := createChirpArgs{}

	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, 500, "Something went wrong")
	}

	if len(params.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	}

	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	params.Body = filterWords(forbiddenWords, params.Body)

	createChirpParams := database.CreateChirpParams{
		Body:   params.Body,
		UserID: params.UserId,
	}
	chirp, err := cfg.db.CreateChirp(context.Background(), createChirpParams)
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	chirpJson := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}

	respondWithJSON(w, 201, chirpJson)
}

func (cfg *ApiConfig) HandleGetAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := cfg.db.GetAllChirp(context.Background())
	if err != nil {
		respondWithError(w, 500, err.Error())
	}

	chirpsJson := []Chirp{}
	for _, chirp := range chirps {
		newChirp := Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserId:    chirp.UserID,
		}
		chirpsJson = append(chirpsJson, newChirp)
	}

	respondWithJSON(w, 200, chirpsJson)
}

func (cfg *ApiConfig) HandleGetChirp(w http.ResponseWriter, r *http.Request) {
	rawId := r.PathValue("id")

	id, err := uuid.Parse(rawId)
	if err != nil {
		respondWithError(w, 404, "Invalid chirp id")
	}

	chirp, err := cfg.db.GetChirp(context.Background(), id)
	if err != nil {
		respondWithError(w, 404, "Chirp not found")
	}

	chirpJson := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserId:    chirp.UserID,
	}
	respondWithJSON(w, 200, chirpJson)
}
