package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/BellOriba/Chirpy/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body string `json:"body"`
	UserID uuid.UUID `json:"user_id"`
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	type newChirps struct {
		Body string `json:"body"`
		UserID uuid.UUID `json:"user_id"`
	}

	decoder := json.NewDecoder(r.Body)
	newChirp := newChirps{}
	err := decoder.Decode(&newChirp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if utf8.RuneCountInString(newChirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	}

	newChirp.Body = checkProfane(newChirp.Body)

	chirp, err := cfg.dbQueries.CreateChirp(r.Context(), database.CreateChirpParams{Body: newChirp.Body, UserID: newChirp.UserID})
	if err != nil {
		log.Printf("Error creating new chirp: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	respondWithJSON(w, 201, chirp)
}

func checkProfane(s string) string {
	profaneWords := [...]string{"kerfuffle", "sharbert", "fornax"}
	sus := strings.Split(s, " ")
	ls := strings.ToLower(s)
	sls := strings.Split(ls, " ")

	for idx, word := range sls {
		for _, profaneW := range profaneWords {
			if word == profaneW {
				sus[idx] = "****"
			}
		}
	}

	fs := strings.Join(sus, " ")
	return fs
}

