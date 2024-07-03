package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

func getMessageFromReq(req *http.Request) (string, error) {
	type chirpDef struct {
		ChirpMessage string `json:"body"`
	}
	decoder := json.NewDecoder(req.Body)
	chirp := chirpDef{}
	err := decoder.Decode(&chirp)
	if err != nil {
		return "", fmt.Errorf("cannot decode json")
	}
	if len(chirp.ChirpMessage) > 140 {
		return "", fmt.Errorf("chirp is too long")
	}

	return checkProfane(chirp.ChirpMessage), nil
}

func (dbCfg *dbConig) createChiprs(res http.ResponseWriter, req *http.Request) {

	ChirpMessage, err := getMessageFromReq(req)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	payload, err := dbCfg.dbClient.AddChirp(ChirpMessage)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	respondWithJSON(res, 201, payload)
}

func (dbCfg *dbConig) getChirps(res http.ResponseWriter, req *http.Request) {
	sliceOfChirps, err := dbCfg.dbClient.GetALlChirps()
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	respondWithJSON(res, 200, sliceOfChirps)
}

func (dbCfg *dbConig) getAChirp(res http.ResponseWriter, req *http.Request) {

	id, err := strconv.Atoi(req.PathValue("id"))
	if err != nil {
		respondWithError(res, 400, "id should be integer")
		return
	}
	mychirp, statCode, err := dbCfg.dbClient.GetsingleChirp(id)
	if err != nil {
		respondWithError(res, statCode, err.Error())
		return
	}

	respondWithJSON(res, 200, mychirp)
}
