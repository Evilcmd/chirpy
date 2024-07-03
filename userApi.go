package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type returnUser struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

type returnLoggedinUser struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	JwtToken string `json:"token"`
}

func getemailAndPassFromReq(req *http.Request) (string, string, int, error) {
	type emaiDef struct {
		Email              string `json:"email"`
		Password           string `json:"password"`
		Expires_in_seconds *int   `json:"expires_in_seconds"`
	}
	decoder := json.NewDecoder(req.Body)
	email := emaiDef{}
	err := decoder.Decode(&email)
	if err != nil {
		return "", "", -1, fmt.Errorf("cannot decode json: " + err.Error())
	}
	x := -1
	if email.Expires_in_seconds != nil {
		x = *email.Expires_in_seconds
	}
	return email.Email, email.Password, x, nil
}

func (dbCfg *userdbConig) createUser(res http.ResponseWriter, req *http.Request) {

	UserEmail, UserPass, _, err := getemailAndPassFromReq(req)
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

type Myclaim struct {
	jwt.RegisteredClaims
}

func (dbCfg *userdbConig) userLogin(res http.ResponseWriter, req *http.Request) {
	UserEmail, UserPass, expiry_from_req, err := getemailAndPassFromReq(req)
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

	expiresInS := time.Second * (24 * 60 * 60)
	if expiry_from_req != -1 {
		if expiry_from_req < 24*60*60 {
			expiresInS = time.Second * time.Duration((expiry_from_req))
		}
	}

	now := time.Now()
	expirationTime := now.Add(expiresInS)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Myclaim{jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.Itoa(code),
	}})

	ss, err := jwtToken.SignedString(dbCfg.dbClient.JwtSecret)
	if err != nil {
		respondWithError(res, 400, "error in creating/signing jwt tokens")
		return
	}

	respondWithJSON(res, 200, returnLoggedinUser{
		code,
		UserEmail,
		ss,
	})
}

func (dbCfg *userdbConig) updateUser(res http.ResponseWriter, req *http.Request) {
	ss := req.Header.Get("Authorization")
	ss = strings.Split(ss, " ")[1]

	jwtToken, err := jwt.ParseWithClaims(ss, &Myclaim{}, func(t *jwt.Token) (interface{}, error) {
		return dbCfg.dbClient.JwtSecret, nil
	})
	if err != nil {
		respondWithError(res, 401, "Unauthorized")
		return
	}

	strId, err := jwtToken.Claims.GetSubject()
	if err != nil {
		respondWithError(res, 400, "cannot get id")
		return
	}
	id, _ := strconv.Atoi(strId)
	UserEmail, UserPass, _, err := getemailAndPassFromReq(req)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	UserPassbyte, err := bcrypt.GenerateFromPassword([]byte(UserPass), 12)
	if err != nil {
		respondWithError(res, 400, "error in generating hash")
		return
	}

	code, err := dbCfg.dbClient.UpdateUser(id, UserEmail, UserPassbyte)
	if err != nil {
		respondWithError(res, code, err.Error())
		return
	}

	respondWithJSON(res, 200, returnUser{
		id,
		UserEmail,
	})
}
