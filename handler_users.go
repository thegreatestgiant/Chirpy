package main

import (
	"errors"
	"net/http"

	"github.com/thegreatestgiant/go-server/internal/auth"
	"github.com/thegreatestgiant/go-server/internal/database"
)

type User struct {
	ID       int    `json:"id"`
	Email    string `json:"email"`
	Password string `json:"-"`
}

type response struct {
	User
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	params := decodeJSON(w, r)

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, "Error Parsing password")
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}

		respondWithError(w, 500, "Couldn't create user")
	}

	respondWithJSON(w, 201, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params := decodeJSON(w, r)

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, 500, "Couldn't Get email")
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Invalid Password")
	}

	respondWithJSON(w, 200, response{
		User: User{
			ID:    user.ID,
			Email: user.Email,
		},
	})
}
