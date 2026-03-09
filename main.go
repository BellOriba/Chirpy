package main

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"

	"github.com/BellOriba/Chirpy/internal/database"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
	jwt_secret     string
	polka_key      string
}

func main() {
	const filepathRoot = "."
	const port = "8080"

	godotenv.Load()

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return
	}

	apiCfg := apiConfig{
		fileserverHits: atomic.Int32{},
		dbQueries:      database.New(db),
		platform:       os.Getenv("PLATFORM"),
		jwt_secret:     os.Getenv("JWT_SECRET"),
		polka_key:      os.Getenv("POLKA_KEY"),
	}

	srvMux := http.NewServeMux()
	srvMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot)))))

	srvMux.HandleFunc("GET /api/healthz", handlerReadiness)
	srvMux.HandleFunc("POST /api/login", apiCfg.handlerLogin)
	srvMux.HandleFunc("POST /api/refresh", apiCfg.handlerRefreshToken)
	srvMux.HandleFunc("POST /api/revoke", apiCfg.handlerRevokeToken)
	srvMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	srvMux.HandleFunc("PUT /api/users", apiCfg.handlerUpdateUser)
	srvMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	srvMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	srvMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)
	srvMux.HandleFunc("DELETE /api/chirps/{chirpID}", apiCfg.handlerDeleteChirp)

	srvMux.HandleFunc("POST /api/polka/webhooks", apiCfg.handlerPolkaWebhook)

	srvMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	srvMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Handler: srvMux,
		Addr:    ":" + port,
	}

	log.Printf("Listening at port: %v", port)
	log.Fatal(srv.ListenAndServe())
}
