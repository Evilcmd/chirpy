package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
)

type DB struct {
	Path      string
	id        int
	Mutex     *sync.RWMutex
	JwtSecret []byte
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

func NewDB() DB {
	filepath := "database.json"
	return DB{filepath, 0, &sync.RWMutex{}, []byte{}}
}

func fileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

func (db *DB) loadDatabase() (DBStructure, error) {
	dbStrct := DBStructure{Chirps: make(map[int]Chirp)}
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
		return DBStructure{make(map[int]Chirp)}, fmt.Errorf("error in decoding: %v", err)
	}

	return dbStrct, nil
}
