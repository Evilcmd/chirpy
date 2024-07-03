package database

import (
	"encoding/json"
	"fmt"
	"os"
)

func (db *DB) AddChirp(message string, authorId int) (Chirp, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()
	dbStrct, err := db.loadDatabase()
	if err != nil {
		return Chirp{}, err
	}

	db.id++
	chrp := Chirp{db.id, message, authorId}
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

func (db *DB) DeleteSingleChirp(authorId int, chirpId int) (int, error) {
	db.Mutex.Lock()
	defer db.Mutex.Unlock()

	dbStrct, err := db.loadDatabase()
	if err != nil {
		return 400, err
	}

	mychirp, ok := dbStrct.Chirps[chirpId]
	if !ok {
		return 400, fmt.Errorf("chirp does not exist")
	}

	if mychirp.AuthorId != authorId {
		return 403, fmt.Errorf("unauthorized")
	}

	delete(dbStrct.Chirps, chirpId)

	dat, _ := json.Marshal(dbStrct)

	err = os.WriteFile(db.Path, dat, 0666)
	if err != nil {
		return 400, fmt.Errorf("error in opening/writing file")
	}

	return 204, nil

}
