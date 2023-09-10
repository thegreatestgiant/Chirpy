package database

import (
	"errors"
	"time"
)

func (db *DB) RevokeToken(token string) error {
	dbStructure, err := db.loadDB()
	if err != nil {
		return err
	}
	if _, ok := dbStructure.Revoked[token]; ok {
		return errors.New("token already revoked")
	}

	dbStructure.Revoked[token] = time.Now().UTC()
	db.writeDB(dbStructure)

	return nil
}

func (db *DB) IsRevoked(token string) (bool, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return false, err
	}
	
	if time, ok := dbStructure.Revoked[token]; !ok || time.IsZero() {
		return false, nil
	}

	return true, nil
}