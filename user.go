package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/BellOriba/Chirpy/internal/auth"
	"github.com/BellOriba/Chirpy/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	CreatedAt      time.Time `json:"create_at"`
	UpdatedAt      time.Time `json:"updated_at"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"hashed_password"`
	IsChirpyRed    bool      `json:"is_chirpy_red"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"create_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Email       string    `json:"email"`
	IsChirpyRed bool      `json:"is_chirpy_red"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	type newUsers struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	newUser := newUsers{}
	err := decoder.Decode(&newUser)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	hashedPassword, err := auth.HashPassword(newUser.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{Email: newUser.Email, HashedPassword: hashedPassword})
	if err != nil {
		log.Printf("Error creating user: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	respondWithJSON(w, 201, UserResponse{
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
	})
}

func (cfg *apiConfig) handlerUpdateUser(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	validBearer, err := auth.ValidateJWT(bearer, cfg.jwt_secret)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if validBearer == uuid.Nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type newCredentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	newCred := newCredentials{}
	err = decoder.Decode(&newCred)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	hashedPass, err := auth.HashPassword(newCred.Password)
	if err != nil {
		log.Printf("Error hashing password: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newUser, err := cfg.dbQueries.UpdateUser(r.Context(), database.UpdateUserParams{
		ID:             validBearer,
		Email:          newCred.Email,
		HashedPassword: hashedPass,
	})
	if err != nil {
		log.Printf("Error updating user: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	respondWithJSON(w, 200, UserResponse{
		ID:          newUser.ID,
		CreatedAt:   newUser.CreatedAt,
		UpdatedAt:   newUser.UpdatedAt,
		Email:       newUser.Email,
		IsChirpyRed: newUser.IsChirpyRed,
	})
}
