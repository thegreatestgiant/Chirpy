package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/thegreatestgiant/go-server/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Bad Auth Header")
		return
	}

	issuer, idString, err := auth.ValidateJWT(token, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, "Invalid Token")
		return
	}
	if issuer != "chirpy-refresh" {
		respondWithError(w, 401, "Not a Refresh Token")
		return
	}

	isRevoked, err := cfg.DB.IsRevoked(token)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	if isRevoked {
		respondWithError(w, 401, "Token Revoked")
		return
	}

	id, err := strconv.Atoi(idString)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error converting id to int: %v", err))
		return
	}

	newToken, err := auth.MakeJWT(id, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error creating new token: %v", err))
		return
	}

	respondWithJSON(w, 200, response{
		Token: newToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Bad Auth Header")
		return
	}

	err = cfg.DB.RevokeToken(token)
	if err != nil {
		respondWithError(w, 401, "Unable to revoke token: " + err.Error())
	}

	respondWithJSON(w, 200, response{})
}
