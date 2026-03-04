package main

import (
	"log"
	"net/http"
)

func (cfg *apiConfig) handlerReset(w http.ResponseWriter, r *http.Request) {
	if cfg.platform != "dev" {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err := cfg.dbQueries.ResetUsers(r.Context())
	if err != nil {
		log.Printf("Error deleting users: %s", err)
		respondWithError(w, 500, "Something went wrong")
		return
	}

	cfg.fileserverHits.Store(0)
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Metrics reseted!"))
}
