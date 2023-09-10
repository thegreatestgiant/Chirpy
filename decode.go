package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type parameters struct {
	Body     string `json:"body"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Event    string `json:"event"`
	Data     data   `json:"data"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte(http.StatusText(http.StatusOK)))
}

func decodeJSON(w http.ResponseWriter, r *http.Request) parameters {
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		log.Printf("Error decoding body: %s", err)
		w.WriteHeader(500)
		return parameters{}
	}
	return params
}
