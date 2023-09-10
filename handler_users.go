package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/thegreatestgiant/Chirpy/internal/auth"
	"github.com/thegreatestgiant/Chirpy/internal/database"
)

type User struct {
	ID          int    `json:"id,omitempty"`
	Email       string `json:"email,omitempty"`
	Password    string `json:"password,omitempty"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
	Token       string `json:"token,omitempty"`
}

type response struct {
	User
	Token        string `json:"token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
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
			ID:          user.ID,
			Email:       user.Email,
			IsChirpyRed: user.IsChirpyRed,
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

	token, err := auth.MakeJWT(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Couldn't Create JWT: %v", err))
	}

	refreshToken, err := auth.MakeRefresh(user.ID, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("Couldn't Create Refresh Token: %v", err))
	}

	respondWithJSON(w, 200, response{
		User: User{
			Email:       user.Email,
			ID:          user.ID,
			IsChirpyRed: user.IsChirpyRed,
		},
		Token:        token,
		RefreshToken: refreshToken,
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

	issuer, id, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error parsing token: %v", err))
		return
	}
	if issuer == "chirpy-refresh" {
		respondWithError(w, 401, "Is a refresh token")
		return
	}

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
