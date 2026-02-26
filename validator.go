package main

import (
	"encoding/json"
	"log"
	"net/http"
	"unicode/utf8"
)

func handlerValidateChirp(w http.ResponseWriter, r *http.Request) {
	type chirps struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	chirp := chirps{}
	err := decoder.Decode(&chirp)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	type validChirp struct {
		Valid bool `json:"valid"`
	}

	if utf8.RuneCountInString(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	} else {
		type validChirp struct {
			Valid bool `json:"valid"`
		}
		respondWithJSON(w, 200, validChirp{Valid: true})
	}
}
