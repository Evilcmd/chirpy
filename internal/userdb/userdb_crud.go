package userdb

import (
	"fmt"

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

	dbStrct.Users[id] = User{id, email, password}

	err = db.writeDatabase(dbStrct)
	if err != nil {
		return 400, err
	}

	return 0, nil
}
