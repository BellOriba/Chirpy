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
	dbQueries *database.Queries
	platform string
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
		dbQueries: database.New(db),
		platform: os.Getenv("PLATFORM"),
	}

	srvMux := http.NewServeMux()
	srvMux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir(filepathRoot))))) 

	srvMux.HandleFunc("GET /api/healthz", handlerReadiness)
	srvMux.HandleFunc("POST /api/users", apiCfg.handlerCreateUser)
	srvMux.HandleFunc("GET /api/chirps", apiCfg.handlerGetAllChirps)
	srvMux.HandleFunc("GET /api/chirps/{chirpID}", apiCfg.handlerGetChirp)
	srvMux.HandleFunc("POST /api/chirps", apiCfg.handlerCreateChirp)

	srvMux.HandleFunc("GET /admin/metrics", apiCfg.handlerMetrics)
	srvMux.HandleFunc("POST /admin/reset", apiCfg.handlerReset)

	srv := &http.Server{
		Handler: srvMux,
		Addr: ":" + port,
	}

	log.Printf("Listening at port: %v", port)
	log.Fatal(srv.ListenAndServe())
}

