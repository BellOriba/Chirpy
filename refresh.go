package main

import (
	"log"
	"net/http"
	"time"

	"github.com/BellOriba/Chirpy/internal/auth"
	"github.com/google/uuid"
)

type RefreshToken struct {
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
	RevokedAt time.Time `json:"revoked_at"`
}

func (cfg *apiConfig) handlerRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := cfg.dbQueries.GetRTokenByToken(r.Context(), bearer)
	if err != nil {
		log.Printf("Error querying refresh token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token.ExpiresAt.Before(time.Now()) {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if token.RevokedAt.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	newJWT, err := auth.MakeJWT(token.UserID, cfg.jwt_secret, time.Hour*1)
	if err != nil {
		log.Printf("Error generating new JWT: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	respondWithJSON(w, 200, struct {
		Token string `json:"token"`
	}{
		Token: newJWT,
	})
}

func (cfg *apiConfig) handlerRevokeToken(w http.ResponseWriter, r *http.Request) {
	bearer, err := auth.GetBearerToken(r.Header)
	if err != nil {
		log.Printf("Error getting bearer token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	err = cfg.dbQueries.SetRTokenRevoke(r.Context(), bearer)
	if err != nil {
		log.Printf("Error revoking refresh token: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
