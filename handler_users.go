package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/thegreatestgiant/go-server/internal/auth"
	"github.com/thegreatestgiant/go-server/internal/database"
)

type User struct {
	ID       int    `json:"id,omitempty"`
	Email    string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
	Token    string `json:"token,omitempty"`
}

type response struct {
	User
	Token string `json:"token,omitempty"`
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	params := decodeJSON(w, r)

	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 400, fmt.Sprintf("Error Parsing password: %v", err))
	}

	user, err := cfg.DB.CreateUser(params.Email, hashedPassword)
	if err != nil {
		if errors.Is(err, database.ErrAlreadyExists) {
			respondWithError(w, http.StatusConflict, "User already exists")
			return
		}

		respondWithError(w, 500, fmt.Sprintf("Couldn't create user: %v", err))
	}

	respondWithJSON(w, 201, response{
		User: User{
			Email: user.Email,
			ID:    user.ID,
		},
	})
}

func (cfg *apiConfig) handlerLogin(w http.ResponseWriter, r *http.Request) {
	params := decodeJSON(w, r)

	user, err := cfg.DB.GetUserByEmail(params.Email)
	if err != nil {
		respondWithError(w, 500, fmt.Sprintf("Couldn't Get email: %v", err))
		return
	}

	err = auth.CheckPasswordHash(params.Password, user.HashedPassword)
	if err != nil {
		respondWithError(w, 401, "Invalid Password:")
	}

	defaultExpiration := 60 * 60 * 24
	if params.Expires_In_Seconds == 0 {
		params.Expires_In_Seconds = defaultExpiration
	} else if params.Expires_In_Seconds > defaultExpiration {
		params.Expires_In_Seconds = defaultExpiration
	}

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret, time.Duration(params.Expires_In_Seconds)*time.Second)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Couldn't Create JWT: %v", err))
	}

	respondWithJSON(w, 200, response{
		User: User{
			Email: user.Email,
			ID:    user.ID,
			Token: token,
		},
		Token: token,
	})
}

func (cfg *apiConfig) handlerPutUser(w http.ResponseWriter, r *http.Request) {
	type response struct {
		User
	}

	params := decodeJSON(w, r)

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("%v", err))
		return
	}

	idString, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error parsing token: %v", err))
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error converting id to int: %v", err))
		return
	}

	fmt.Println(params.Password)
	hashedPassword, err := auth.HashPassword(params.Password)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error hashing password: %v", err))
		return
	}
	cfg.DB.UpdateUser(id, params.Email, hashedPassword)

	respondWithJSON(w, 200, response{
		User: User{
			Email: params.Email,
			ID:    id,
		},
	})
}
