package main

import (
	"fmt"
	"net/http"

	"github.com/Evilcmd/chirpy/internal/database"
)

type apiConfig struct {
	hits int
}

type dbConig struct {
	dbClient database.DB
}

func main() {

	apiCfg := apiConfig{0}

	dbConig := dbConig{
		database.NewDB(),
	}

	mux := http.NewServeMux()

	fileServerRoute := http.StripPrefix("/app", http.FileServer(http.Dir(".")))
	mux.Handle("/app/", apiCfg.middleWareToIncreaseHits(fileServerRoute))

	mux.HandleFunc("GET /api/healthz", healthz)

	mux.HandleFunc("GET /api/metrics", apiCfg.metrics)

	mux.HandleFunc("/api/reset", apiCfg.reset)

	mux.HandleFunc("GET /admin/metrics", apiCfg.metricsHTML)

	mux.HandleFunc("POST /api/chirps", dbConig.createChiprs)
	mux.HandleFunc("GET /api/chirps", dbConig.getChirps)
	mux.HandleFunc("GET /api/chirps/{id}", dbConig.getAChirp)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Starting Server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server")
	}
}
