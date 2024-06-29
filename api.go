package main

import (
	"encoding/json"
	"net/http"
)

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
