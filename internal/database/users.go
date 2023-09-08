package database

import (
	"os"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type em struct {
	ID    int    `json:"id"`
	Email string `json:"email"`
}

func (db *DB) CreateUser(pass, email string) (em, error) {

	dbStructure, err := db.loadDB()
	if err != nil {
		return em{}, err
	}
	id := len(dbStructure.Users) + 1
	dbStructure.Users[id] = User{
		ID:       id,
		Email:    email,
		Password: pass,
	}
	err = db.writeDB(dbStructure)
	if err != nil {
		return em{}, err
	}
	user := em{
		ID:    id,
		Email: email,
	}
	return user, nil
}

func (db *DB) GetUsers() ([]User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return nil, err
	}
	users := make([]User, 0, len(dbStructure.Users))
	for _, user := range dbStructure.Users {
		users = append(users, user)
	}
	return users, nil
}

func (db *DB) GetUser(id int) (User, error) {
	dbStructure, err := db.loadDB()
	if err != nil {
		return User{}, err
	}
	users, ok := dbStructure.Users[id]
	if !ok {
		return User{}, os.ErrNotExist
	}
	return users, nil
}
