package main

import (
	"crypto/rand"
	"encoding/hex"
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
	ID          int    `json:"id"`
	Email       string `json:"email"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type returnLoggedinUser struct {
	ID           int    `json:"id"`
	Email        string `json:"email"`
	JwtToken     string `json:"token"`
	RefreshToken string `json:"refresh_token"`
	IsChirpyRed  bool   `json:"is_chirpy_red"`
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
		payload.IsChirpyRed,
	})
}

type Myclaim struct {
	jwt.RegisteredClaims
}

func (dbCfg *userdbConig) userLogin(res http.ResponseWriter, req *http.Request) {
	UserEmail, UserPass, _, err := getemailAndPassFromReq(req)
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

	expiresInS := time.Second * (60 * 60)
	// if expiry_from_req != -1 {
	// 	if expiry_from_req < 24*60*60 {
	// 		expiresInS = time.Second * time.Duration((expiry_from_req))
	// 	}
	// }

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

	b := make([]byte, 32)
	rand.Read(b)
	refreshToken := hex.EncodeToString(b)

	refDbStrct, err := dbCfg.dbClient.GetRefreshTokenStruct()
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}
	refDbStrct.Token[refreshToken] = struct {
		Expiry time.Time
		Id     int
	}{time.Now().Add(time.Hour * 24 * 60), code}
	err = dbCfg.dbClient.WriteRefreshTokenStruct(refDbStrct)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	chirpyRedStatus, err := dbCfg.dbClient.GetChirpyRedStatus(code)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	respondWithJSON(res, 200, returnLoggedinUser{
		code,
		UserEmail,
		ss,
		refreshToken,
		chirpyRedStatus,
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

	chirpyRedStatus, err := dbCfg.dbClient.GetChirpyRedStatus(code)
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	respondWithJSON(res, 200, returnUser{
		id,
		UserEmail,
		chirpyRedStatus,
	})
}

func (dbCfg *userdbConig) RefreshTokens(res http.ResponseWriter, req *http.Request) {
	tokenString := req.Header.Get("Authorization")
	tokenString = strings.Split(tokenString, " ")[1]

	refDbStrct, err := dbCfg.dbClient.GetRefreshTokenStruct()
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	v, ok := refDbStrct.Token[tokenString]

	if !ok || (v.Expiry.Before(time.Now())) {
		respondWithError(res, 401, "does not exist or has expired")
	}

	now := time.Now()
	expirationTime := now.Add(time.Hour)

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Myclaim{jwt.RegisteredClaims{
		Issuer:    "chirpy",
		IssuedAt:  jwt.NewNumericDate(now),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Subject:   strconv.Itoa(v.Id),
	}})

	ss, err := jwtToken.SignedString(dbCfg.dbClient.JwtSecret)
	if err != nil {
		respondWithError(res, 400, "error in creating/signing jwt tokens")
		return
	}

	respondWithJSON(res, 200, struct {
		Token string `json:"token"`
	}{
		ss,
	})
}

func (dbCfg *userdbConig) RevokeTokens(res http.ResponseWriter, req *http.Request) {
	tokenString := req.Header.Get("Authorization")
	tokenString = strings.Split(tokenString, " ")[1]

	refDbStrct, err := dbCfg.dbClient.GetRefreshTokenStruct()
	if err != nil {
		respondWithError(res, 400, err.Error())
		return
	}

	_, ok := refDbStrct.Token[tokenString]

	if ok {
		delete(refDbStrct.Token, tokenString)
		err = dbCfg.dbClient.WriteRefreshTokenStruct(refDbStrct)
		if err != nil {
			respondWithError(res, 400, err.Error())
			return
		}
	}

	res.WriteHeader(204)

}

type ChirpyRedHookDefn struct {
	Event string `json:"event"`
	Data  struct {
		UserID int `json:"user_id"`
	} `json:"data"`
}

func (dbCfg *userdbConig) ChirpyRedWebhook(res http.ResponseWriter, req *http.Request) {
	givenSecret := req.Header.Get("Authorization")
	if givenSecret == "" {
		respondWithError(res, 401, "no authorization given")
		return
	}
	givenSecret = strings.Split(givenSecret, " ")[1]
	if givenSecret != dbCfg.dbClient.PolkaSecret {
		respondWithError(res, 401, "wrong authorization key")
		return
	}

	deocder := json.NewDecoder(req.Body)
	ChirpyRedHookRes := ChirpyRedHookDefn{}
	err := deocder.Decode(&ChirpyRedHookRes)
	if err != nil {
		respondWithError(res, 400, "cannot decode the input")
		return
	}
	if ChirpyRedHookRes.Event != "user.upgraded" {
		res.WriteHeader(204)
		return
	}
	code, err := dbCfg.dbClient.UpdateChirpyRed(ChirpyRedHookRes.Data.UserID)
	if err != nil {
		respondWithError(res, code, err.Error())
		return
	}
	res.WriteHeader(204)
}
