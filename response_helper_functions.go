package main

import (
	"encoding/json"
	"net/http"
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
