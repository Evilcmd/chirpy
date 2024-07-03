package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/golang-jwt/jwt/v5"
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

func (dbCfg *dbConig) authenticate(req *http.Request) (int, error) {
	ss := req.Header.Get("Authorization")
	ss = strings.Split(ss, " ")[1]

	jwtToken, err := jwt.ParseWithClaims(ss, &Myclaim{}, func(t *jwt.Token) (interface{}, error) {
		return dbCfg.dbClient.JwtSecret, nil
	})
	if err != nil {
		return 0, fmt.Errorf("unauthorized")
	}
	x, _ := jwtToken.Claims.GetSubject()
	y, _ := strconv.Atoi(x)
	return y, nil
}

func (dbCfg *dbConig) createChiprs(res http.ResponseWriter, req *http.Request) {

	authorID, err := dbCfg.authenticate(req)
	if err != nil {
		respondWithError(res, 401, err.Error())
		return
	}

	ChirpMessage, err := getMessageFromReq(req)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	payload, err := dbCfg.dbClient.AddChirp(ChirpMessage, authorID)
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
