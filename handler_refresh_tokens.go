package main

import (
	"fmt"
	"net/http"

	"github.com/thegreatestgiant/Chirpy/internal/auth"
)

func (cfg *apiConfig) handlerRefresh(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Bad Auth Header")
		return
	}

	isRevoked, err := cfg.DB.IsRevoked(refreshToken)
	if err != nil {
		respondWithError(w, 401, err.Error())
		return
	}
	if isRevoked {
		respondWithError(w, 401, "Token Revoked")
		return
	}

	accessToken, err := auth.RefreshToken(refreshToken, cfg.jwtSecret)
	if err != nil {
		respondWithError(w, 401, fmt.Sprintf("error creating new token: %v", err))
		return
	}

	respondWithJSON(w, 200, response{
		Token: accessToken,
	})
}

func (cfg *apiConfig) handlerRevoke(w http.ResponseWriter, r *http.Request) {
	refreshToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, 401, "Bad Auth Header")
		return
	}

	err = cfg.DB.RevokeToken(refreshToken)
	if err != nil {
		respondWithError(w, 401, "Unable to revoke token: "+err.Error())
	}

	respondWithJSON(w, 200, struct{}{})
}
