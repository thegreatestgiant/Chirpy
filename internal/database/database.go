package database

import (
	"encoding/json"
	"errors"
	"os"
	"sync"
)

type DB struct {
	Path string
	Mux  *sync.RWMutex
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
	Users map[int]User `json:"users"`
}

type Chirp struct {
	Body string `json:"body"`
	ID   int    `json:"id"`
}



func NewDB(path string) (*DB, error) {
	db := &DB{
		Path: path,
		Mux:  &sync.RWMutex{},
	}
	err := db.ensureDB()
	return db, err
}

func (db *DB) CreateChirp(body string) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	id := len(dbStructure.Chirps) + 1
	NewChirp := Chirp{
		ID:   id,
		Body: body,
	}
	dbStructure.Chirps[NewChirp.ID] = NewChirp
	err = db.writeDB(dbStructure)
	if err != nil {
		return Chirp{}, err
	}
	return NewChirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	chirps := make([]Chirp, 0, len(dbStructure.Chirps))
	for _, chirp := range dbStructure.Chirps {
		chirps = append(chirps, chirp)
	}
	return chirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}
	chirps, ok := dbStructure.Chirps[id]
	if !ok {
		return Chirp{}, os.ErrNotExist
	}
	return chirps, nil
}

func (db *DB) createDB() error {
	dbStructure := DBStructure{
		Chirps: map[int]Chirp{},
		Users: map[int]User{},
	}
	return db.writeDB(dbStructure)
}

func (db *DB) ensureDB() error {
	_, err := os.ReadFile(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		return db.createDB()
	}
	return err
}

func (db *DB) loadDB() (DBStructure, error) {
	db.Mux.RLock()
	defer db.Mux.RUnlock()

	dbStructure := DBStructure{}
	dat, err := os.ReadFile(db.Path)
	if errors.Is(err, os.ErrNotExist) {
		return dbStructure, err
	}
	err = json.Unmarshal(dat, &dbStructure)
	if err != nil {
		return dbStructure, err
	}
	return dbStructure, nil
}

func (db *DB) writeDB(dbStructure DBStructure) error {
	db.Mux.Lock()
	defer db.Mux.Unlock()

	dat, err := json.Marshal(dbStructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.Path, dat, 0600)
	if err != nil {
		return err
	}
	return nil
}

func (db *DB) ResetDB() {
	os.Remove(db.Path)
	db.createDB()
}
