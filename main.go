package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

func respondWithJSON(res http.ResponseWriter, code int, payload interface{}) {
	res.WriteHeader(code)
	dat, _ := json.Marshal(payload)
	res.Write(dat)
}

func respondWithError(res http.ResponseWriter, code int, msg string) {
	payload := struct {
		Error string `json:"error"`
	}{
		msg,
	}
	respondWithJSON(res, code, &payload)
}

func healthz(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte("200 OK"))
}

type apiConfig struct {
	hits int
}

func (apiCfg *apiConfig) middleWareToIncreaseHits(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		apiCfg.hits++
		next.ServeHTTP(res, req)
	})
}

func (apiCfg *apiConfig) metrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte(fmt.Sprintf("Hits: %v", apiCfg.hits)))
}

func (apiCfg *apiConfig) reset(res http.ResponseWriter, req *http.Request) {
	apiCfg.hits = 0
}

func (apiCfg *apiConfig) metricsHTML(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.hits)))
}

func checkProfane(message string) string {
	splitMessage := strings.Split(message, " ")
	for i, msg := range splitMessage {
		if strings.ToLower(msg) == "kerfuffle" || strings.ToLower(msg) == "sharbert" || strings.ToLower(msg) == "fornax" {
			splitMessage[i] = "****"
		}
	}
	return strings.Join(splitMessage, " ")
}

func validateChirp(res http.ResponseWriter, req *http.Request) {
	type chirpDef struct {
		ChirpMessage string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	chirp := chirpDef{}
	err := decoder.Decode(&chirp)
	if err != nil {
		respondWithError(res, 400, "Cannot decode JSON")
		return
	}
	// fmt.Println(len(chirp.ChirpMessage))
	if len(chirp.ChirpMessage) > 140 {
		respondWithError(res, 400, "Chirp is too long")
		return
	}

	chirp.ChirpMessage = checkProfane(chirp.ChirpMessage)

	payload := struct {
		CleanedBody string `json:"cleaned_body"`
	}{
		chirp.ChirpMessage,
	}
	respondWithJSON(res, 200, payload)
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
