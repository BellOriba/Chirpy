package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
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

	if utf8.RuneCountInString(chirp.Body) > 140 {
		respondWithError(w, 400, "Chirp is too long")
	}

	type validChirp struct {
		CleanedBody string `json:"cleaned_body"`
	}

	cleanChirp := checkProfane(chirp.Body)
	respondWithJSON(w, 200, validChirp{CleanedBody: cleanChirp})
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
