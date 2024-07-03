package userdb

import (
	"encoding/json"
	"fmt"
	"os"

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
	usr := User{db.id, email, password}
	dbStrct.Users[db.id] = usr

	dat, _ := json.Marshal(dbStrct)

	err = os.WriteFile(db.Path, dat, 0666)
	if err != nil {
		return User{}, fmt.Errorf("error in opening/writing file")
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
