package main

import (
	"fmt"
	"net/http"
)

func healthz(res http.ResponseWriter, req *http.Request) {
	res.WriteHeader(200)
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte("200 OK"))
}

func (apiCfg *apiConfig) metrics(res http.ResponseWriter, req *http.Request) {
	res.Header().Set("Content-Type", "text/plain; charset=utf-8")
	res.Write([]byte(fmt.Sprintf("Hits: %v", apiCfg.hits)))
}

func (apiCfg *apiConfig) reset(res http.ResponseWriter, req *http.Request) {
	apiCfg.hits = 0
}
