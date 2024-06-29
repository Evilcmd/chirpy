package main

import (
	"fmt"
	"net/http"
)

func (apiCfg *apiConfig) metricsHTML(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/html")
	res.Write([]byte(fmt.Sprintf("<html><body><h1>Welcome, Chirpy Admin</h1><p>Chirpy has been visited %d times!</p></body></html>", apiCfg.hits)))
}
