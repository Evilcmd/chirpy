package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	hits int
}

func main() {

	apiCfg := apiConfig{0}

	mux := http.NewServeMux()

	fileServerRoute := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middleWareToIncreaseHits(fileServerRoute))

	mux.HandleFunc("GET /api/healthz", healthz)

	mux.HandleFunc("GET /api/metrics", apiCfg.metrics)

	mux.HandleFunc("/api/reset", apiCfg.reset)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHTML)

	mux.HandleFunc("POST /api/validate_chirp", validateChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Starting Server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server")
	}
}
