package main

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {

	chirps, err := cfg.DB.GetChirps()
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Couldn't get Chirps")
	}
	sort.Slice(chirps, func(i, j int) bool { return chirps[i].ID < chirps[j].ID })
	JSONResp(w, 200, chirps)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	body := validateLength(w, r)
	chirp, err := cfg.DB.CreateChirp(body)
	if err != nil {
		errorResp(w, 400, err.Error())
	}
	JSONResp(w, 201, chirp)
}

func (cfg *apiConfig) handlerResetDB(w http.ResponseWriter, r *http.Request) {
	cfg.DB.ResetDB()
	JSONResp(w, 205, struct{ Success string }{Success: "Deleted DB, Regenerated db file"})
}

func (cfg *apiConfig) handlerGetChirpsByID(w http.ResponseWriter, r *http.Request) {
	strChirpID := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(strChirpID)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Invalid Chirp ID")
	}
	chirps, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		errorResp(w, 404, "chirpID not found")
	}
	JSONResp(w, 200, chirps)
}
