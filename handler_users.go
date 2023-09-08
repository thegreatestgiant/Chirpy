package main

import (
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"sort"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func (cfg *apiConfig) handlerGetusers(w http.ResponseWriter, r *http.Request) {

	users, err := cfg.DB.GetUsers()
	if err != nil {
		errorResp(w, http.StatusInternalServerError, "Couldn't get Users")
	}
	sort.Slice(users, func(i, j int) bool { return users[i].ID < users[j].ID })
	JSONResp(w, 200, users)
}

func (cfg *apiConfig) handlerCreateUser(w http.ResponseWriter, r *http.Request) {
	email := decodeJSON(w, r)
	pass, err := bcrypt.GenerateFromPassword([]byte(email.Password), 10)
	if err != nil {
		errorResp(w, 400, "Error Parsing password")
	}
	user, err := cfg.DB.CreateUser(string(pass), email.Email)
	if err != nil {
		errorResp(w, 400, err.Error())
	}
	JSONResp(w, 201, user)
}

func (cfg *apiConfig) handlerGetUserByID(w http.ResponseWriter, r *http.Request) {
	strUserID := chi.URLParam(r, "UserID")
	UserID, err := strconv.Atoi(strUserID)
	if err != nil {
		errorResp(w, http.StatusBadRequest, "Invalid UserID")
	}
	users, err := cfg.DB.GetUser(UserID)
	if err != nil {
		errorResp(w, 404, "UserID not found")
	}
	JSONResp(w, 200, users)
}
