package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/BellOriba/Chirpy/internal/auth"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerPolkaWebhook(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetAPIKey(r.Header)
	if err != nil {
		log.Printf("Error getting api key: %s", err)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	if apiKey != cfg.polka_key {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	type RedEvent struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	decoder := json.NewDecoder(r.Body)
	newEvent := RedEvent{}
	err = decoder.Decode(&newEvent)
	if err != nil {
		log.Printf("Error decoding parameters: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	if newEvent.Event != "user.upgraded" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	_, err = cfg.dbQueries.UpgradeToRed(r.Context(), newEvent.Data.UserID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, 404, "User not found")
			return
		}

		log.Printf("Error updating user in DB: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
