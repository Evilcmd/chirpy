package userdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"
)

func checkDuplicate(mp map[int]User, email string) bool {
	for _, v := range mp {
		if v.Email == email {
			return true
		}
	}
	return false
}

func (db *UserDB) AddUser(email string, password []byte) (User, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	dbStrct, err := db.loadDatabase()
	if err != nil {
		return User{}, err
	}

	if checkDuplicate(dbStrct.Users, email) {
		return User{}, fmt.Errorf("duplicate user")
	}

	db.id++
	usr := User{db.id, email, password, false}
	dbStrct.Users[db.id] = usr

	err = db.writeDatabase(dbStrct)
	if err != nil {
		return User{}, err
	}

	return usr, nil
}

func (db *UserDB) VerifyUser(email string, password []byte) (int, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	dbStrct, err := db.loadDatabase()
	if err != nil {
		return 400, err
	}

	mp := dbStrct.Users
	for id, v := range mp {
		if v.Email == email {
			if err = bcrypt.CompareHashAndPassword(v.Password, password); err != nil {
				return 401, fmt.Errorf("passwords dont match")
			} else {
				return id, nil
			}
		}
	}
	return 400, fmt.Errorf("user not found")
}

func (db *UserDB) UpdateUser(id int, email string, password []byte) (int, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	dbStrct, err := db.loadDatabase()
	if err != nil {
		return 400, err
	}

	_, ok := dbStrct.Users[id]
	if !ok {
		return 400, fmt.Errorf("id not found")
	}

	// if v.Email != email {
	// 	return 400, fmt.Errorf("email does not match")
	// }

	dbStrct.Users[id] = User{id, email, password, dbStrct.Users[id].IsChirpyRed}

	err = db.writeDatabase(dbStrct)
	if err != nil {
		return 400, err
	}

	return 0, nil
}

type RefreshTokenDefn struct {
	Token map[string]struct {
		Expiry time.Time
		Id     int
	} `json:"token"`
}

func (db *UserDB) GetRefreshTokenStruct() (RefreshTokenDefn, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()

	dbStrct := RefreshTokenDefn{Token: make(map[string]struct {
		Expiry time.Time
		Id     int
	})}
	if !fileExists(db.TokenDbPath) {
		return dbStrct, nil
	}
	dbread, err := os.ReadFile(db.TokenDbPath)
	if err != nil {
		return dbStrct, fmt.Errorf("error in opening/reading file: %v", err)
	}
	dbreader := bytes.NewReader(dbread)
	decoder := json.NewDecoder(dbreader)
	err = decoder.Decode(&dbStrct)
	if err != nil {
		return RefreshTokenDefn{Token: make(map[string]struct {
			Expiry time.Time
			Id     int
		})}, fmt.Errorf("error in decoding: %v", err)
	}
	return dbStrct, nil
}

func (db *UserDB) WriteRefreshTokenStruct(refDbStrct RefreshTokenDefn) error {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	dat, _ := json.Marshal(refDbStrct)

	err := os.WriteFile(db.TokenDbPath, dat, 0666)
	if err != nil {
		return fmt.Errorf("error in opening/writing file")
	}

	return nil
}

func (db *UserDB) UpdateChirpyRed(UserId int) (int, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	userDbStct, err := db.loadDatabase()
	if err != nil {
		return 400, err
	}
	_, ok := userDbStct.Users[UserId]
	if !ok {
		return 404, fmt.Errorf("user not found")
	}
	userDbStct.Users[UserId] = User{UserId, userDbStct.Users[UserId].Email, userDbStct.Users[UserId].Password, true}
	err = db.writeDatabase(userDbStct)
	if err != nil {
		return 400, fmt.Errorf("error in writing to database")
	}
	return 204, nil
}

func (db *UserDB) GetChirpyRedStatus(UserId int) (bool, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	userDbStct, err := db.loadDatabase()
	if err != nil {
		return false, err
	}
	v, ok := userDbStct.Users[UserId]
	if !ok {
		return false, fmt.Errorf("user not found")
	}

	return v.IsChirpyRed, nil
}
