package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func errorResp(w http.ResponseWriter, code int, msg string) {
	if code > 499 {
		log.Printf("Responding with 5XX error: %s", msg)
	}
	type errorResp struct {
		Error string `json:"error"`
	}
	JSONResp(w, code, errorResp{
		Error: msg,
	})
}

func JSONResp(w http.ResponseWriter, code int, payload interface{}) {
	dat, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Error marshaling response body: %s", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(dat)
}
