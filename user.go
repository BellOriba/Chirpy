package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type newUsers struct {
		Email string `json:"email"`
	}

	decoder := json.NewDecoder(r.Body)
	newUser := newUsers{}
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), newUser.Email)
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	respondWithJSON(w, 201, user)
}
