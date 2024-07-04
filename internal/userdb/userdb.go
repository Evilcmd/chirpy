package userdb

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type UserDB struct {
	Path        string
	id          int
	Mutex       *sync.RWMutex
	JwtSecret   []byte
	TokenDbPath string
	PolkaSecret string
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    []byte `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type UserDBStructure struct {
	Users map[int]User `json:"chirps"`
}

func NewDB() UserDB {
	filepath := "UserDB.json"
	return UserDB{filepath, 0, &sync.RWMutex{}, []byte{}, "TokenDBPath.json", ""}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (db *UserDB) loadDatabase() (UserDBStructure, error) {
	dbStrct := UserDBStructure{Users: make(map[int]User)}
	if !fileExists(db.Path) {
		return dbStrct, nil
	}
	dbread, err := os.ReadFile(db.Path)
	if err != nil {
		return dbStrct, fmt.Errorf("error in opening/reading file: %v", err)
	}
	dbreader := bytes.NewReader(dbread)
	decoder := json.NewDecoder(dbreader)
	err = decoder.Decode(&dbStrct)
	if err != nil {
		return UserDBStructure{make(map[int]User)}, fmt.Errorf("error in decoding: %v", err)
	}

	return dbStrct, nil
}

func (db *UserDB) writeDatabase(dbStrct UserDBStructure) error {
	dat, _ := json.Marshal(dbStrct)

	err := os.WriteFile(db.Path, dat, 0666)
	if err != nil {
		return fmt.Errorf("error in opening/writing file")
	}

	return nil
}
