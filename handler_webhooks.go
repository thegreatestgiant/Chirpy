package main

import (
	"errors"
	"net/http"

	"github.com/thegreatestgiant/go-server/internal/auth"
	"github.com/thegreatestgiant/go-server/internal/database"
)

type data struct {
	UserID int `json:"user_id"`
}

func (cfg *apiConfig) handlerPolkaWebhooks(w http.ResponseWriter, r *http.Request) {
	apiKey, err := auth.GetApiKey(r.Header)
	if err != nil {
		respondWithError(w, 401, "Bad api key: "+err.Error())
		return
	}

	if apiKey != cfg.polkaApi {
		respondWithError(w, 401, "Bad API key"+apiKey)
		return
	}

	params := decodeJSON(w, r)

	if params.Event != "user.upgraded" {
		respondWithJSON(w, 200, struct{}{})
		return
	}

	err = cfg.DB.UpgradeUser(params.Data.UserID)
	if err != nil {
		if errors.Is(err, database.ErrNotExist) {
			respondWithError(w, http.StatusNotFound, "Couldn't find user")
			return
		}
		respondWithError(w, 404, "could not upgrade user: "+err.Error())
		return
	}

	respondWithJSON(w, 200, response{})
}
