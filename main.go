package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Evilcmd/chirpy/internal/database"
	"github.com/Evilcmd/chirpy/internal/userdb"
	"github.com/joho/godotenv"
)

type apiConfig struct {
	hits int
}

type dbConig struct {
	dbClient database.DB
}

type userdbConig struct {
	dbClient userdb.UserDB
}

func main() {

	apiCfg := apiConfig{0}

	dbConig := dbConig{
		database.NewDB(),
	}

	userdbconfig := userdbConig{
		userdb.NewDB(),
	}

	godotenv.Load()
	jwtSecret := os.Getenv("JWT_SECRET")
	polkaSecret := os.Getenv("POLKA_SECRET")
	userdbconfig.dbClient.JwtSecret = []byte(jwtSecret)
	userdbconfig.dbClient.PolkaSecret = polkaSecret
	dbConig.dbClient.JwtSecret = []byte(jwtSecret)

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

	mux.HandleFunc("DELETE /api/chirps/{chirpId}", dbConig.deleteAChirp)

	mux.HandleFunc("POST /api/users", userdbconfig.createUser)
	mux.HandleFunc("POST /api/login", userdbconfig.userLogin)
	mux.HandleFunc("PUT /api/users", userdbconfig.updateUser)

	mux.HandleFunc("POST /api/refresh", userdbconfig.RefreshTokens)
	mux.HandleFunc("POST /api/revoke", userdbconfig.RevokeTokens)

	mux.HandleFunc("POST /api/polka/webhooks", userdbconfig.ChirpyRedWebhook)

	server := http.Server{
		Addr:    ":8080",
		Handler: mux,
	}
	fmt.Println("Starting Server on :8080")
	if err := server.ListenAndServe(); err != nil {
		fmt.Println("Error starting server")
	}
}
