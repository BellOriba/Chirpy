package main

import (
	"log"
	"net/http"
)

func main() {
	const port = "8080"

	srvMux := http.NewServeMux()
	srv := &http.Server{
		Handler: srvMux,
		Addr: ":" + port,
	}

	log.Printf("Listening at port: %v", port)
	log.Fatal(srv.ListenAndServe())
}

