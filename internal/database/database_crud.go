package database

import (
	"encoding/json"
	"fmt"
	"os"
)

func (db *DB) AddChirp(message string) (Chirp, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	dbStrct, err := db.loadDatabase()
	if err != nil {
		return Chirp{}, err
	}

	db.id++
	chrp := Chirp{db.id, message}
	dbStrct.Chirps[db.id] = chrp

	dat, _ := json.Marshal(dbStrct)

	err = os.WriteFile(db.Path, dat, 0666)
	if err != nil {
		return Chirp{}, fmt.Errorf("error in opening/writing file")
	}

	return chrp, nil
}

func (db *DB) GetALlChirps() ([]Chirp, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()
	dbStruct, err := db.loadDatabase()
	if err != nil {
		return []Chirp{}, err
	}
	sliceOfChirps := make([]Chirp, 0, len(dbStruct.Chirps))

	for _, chrp := range dbStruct.Chirps {
		sliceOfChirps = append(sliceOfChirps, chrp)
	}

	return sliceOfChirps, nil
}

func (db *DB) GetsingleChirp(id int) (Chirp, int, error) {
	db.Mutex.RLock()
	defer db.Mutex.RUnlock()
	dbStruct, err := db.loadDatabase()
	if err != nil {
		return Chirp{}, 400, err
	}

	mychirp, ok := dbStruct.Chirps[id]
	if !ok {
		return Chirp{}, 404, fmt.Errorf("chirp does not exist")
	}

	return mychirp, 0, nil
}
