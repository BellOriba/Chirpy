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

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	type LoginParams struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}

	decoder := json.NewDecoder(r.Body)
	userLogin := LoginParams{}
	err := decoder.Decode(&userLogin)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	user, err := cfg.dbQueries.GetUserByEmail(r.Context(), userLogin.Email)
	if err != nil {
		log.Printf("Error querying user: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	checkPass, err := auth.CheckPasswordHash(userLogin.Password, user.HashedPassword)
	if err != nil {
		log.Printf("Error checking password: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if !checkPass {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwt_secret, time.Hour * 1)
	if err != nil {
		log.Printf("Error generation token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	refreshToken := auth.MakeRefreshToken()
	cfg.dbQueries.InsertRToken(r.Context(), database.InsertRTokenParams{
		Token: refreshToken,
		UserID: user.ID,
		ExpiresAt: time.Now().AddDate(0, 0, 60),
	})

	respondWithJSON(w, 200, struct{
		ID uuid.UUID `json:"id"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
		Email string `json:"email"`
		Token string `json:"token"`
		RefreshToken string `json:"refresh_token"`
	}{
		ID: user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Email: user.Email,
		Token: token,
		RefreshToken: refreshToken,
	})
}

