package main

import (
	"net/http"
	"sort"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"
)

type Chirp struct {
	ID   int    `json:"id"`
	Body string `json:"body"`
}

func (cfg *apiConfig) handlerGetChirps(w http.ResponseWriter, r *http.Request) {
	dbChirps, err := cfg.DB.GetChirps()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get Chirps")
	}

	chirps := []Chirp{}
	for _, dbChirp := range dbChirps {
		chirps = append(chirps, Chirp{
			ID:   dbChirp.ID,
			Body: dbChirp.Body,
		})
	}

	sort.Slice(chirps, func(i, j int) bool {
		return chirps[i].ID < chirps[j].ID
	})
	respondWithJSON(w, 200, chirps)
}

func (cfg *apiConfig) handlerCreateChirp(w http.ResponseWriter, r *http.Request) {
	params := decodeJSON(w, r)

	if len(params.Body) > 140 {
		respondWithError(w, 400, "chirp is too long")
		return
	}

	text := params.Body
	badWords := []string{
		"kerfuffle",
		"sharbert",
		"fornax",
	}
	arr := strings.Split(text, " ")
	for i, word := range arr {
		for _, badWord := range badWords {
			if strings.ToLower(word) == badWord {
				arr[i] = "****"
				break
			}
		}
	}
	body := strings.Join(arr, " ")

	chirp, err := cfg.DB.CreateChirp(body)
	if err != nil {
		respondWithError(w, 400, err.Error())
	}
	respondWithJSON(w, 201, chirp)
}

func (cfg *apiConfig) handlerResetDB(w http.ResponseWriter, r *http.Request) {
	cfg.DB.ResetDB()
	respondWithJSON(w, 205, struct{ Success string }{Success: "Deleted DB, Regenerated db file"})
}

func (cfg *apiConfig) handlerGetChirpsByID(w http.ResponseWriter, r *http.Request) {
	strChirpID := chi.URLParam(r, "chirpID")
	chirpID, err := strconv.Atoi(strChirpID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid Chirp ID")
	}
	
	chirps, err := cfg.DB.GetChirp(chirpID)
	if err != nil {
		respondWithError(w, 404, "chirpID not found")
	}
	respondWithJSON(w, 200, chirps)
}
