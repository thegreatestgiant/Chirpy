package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body  string `json:"body"`
	Email string `json:"email"`
	Password string `json:"password"`
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func validateLength(w http.ResponseWriter, r *http.Request) string {
	params := decodeJSON(w, r)

	if len(params.Body) > 140 {
		errorResp(w, 400, "chirp is too long")
		return ""
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
	text = strings.Join(arr, " ")

	return text

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
