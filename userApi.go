package main

import (
	"encoding/json"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

type returnUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func getemailAndPassFromReq(req *http.Request) (string, string, error) {
	type emaiDef struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(req.Body)
	email := emaiDef{}
	err := decoder.Decode(&email)
	if err != nil {
		return "", "", err
		// return "", []byte{}, fmt.Errorf("cannot decode json")
	}

	return email.Email, email.Password, nil
}

func (dbCfg *userdbConig) createUser(res http.ResponseWriter, req *http.Request) {

	UserEmail, UserPass, err := getemailAndPassFromReq(req)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	UserPassbyte, err := bcrypt.GenerateFromPassword([]byte(UserPass), 12)
	if err != nil {
		respondWithError(res, 400, "error in generating hash")
		return
	}

	payload, err := dbCfg.dbClient.AddUser(UserEmail, UserPassbyte)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	respondWithJSON(res, 201, returnUser{
		payload.Id,
		payload.Email,
	})
}

func (dbCfg *userdbConig) userLogin(res http.ResponseWriter, req *http.Request) {
	UserEmail, UserPass, err := getemailAndPassFromReq(req)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	code, err := dbCfg.dbClient.VerifyUser(UserEmail, []byte(UserPass))
	if err != nil {
		if code == 401 {
			respondWithError(res, 401, "Unauthorized")
			return
		}
		respondWithError(res, 400, err.Error())
		return
	}
	respondWithJSON(res, 200, returnUser{
		code,
		UserEmail,
	})
}
