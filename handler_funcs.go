package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type parameters struct {
	Body string `json:"body"`
}

func (cfg *apiConfig) resetMetrics(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits = 0
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hits reset to 0"))
}

func (cfg *apiConfig) getMetrics(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8 ")
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`<html>

	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>
	
	</html>`, cfg.fileserverHits)))
}

func handlerReadiness(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(200)
	w.Write([]byte("OK"))
}

func validateLength(w http.ResponseWriter, r *http.Request) {
	type returnVals struct {
		Cleaned_Body string `json:"cleaned_body"`
		// Error        string `json:"error"`
		// Valid        bool   `json:"valid"`
	}
	params := decodeJSON(w, r)

	if len(params.Body) > 140 {
		errorResp(w, 400, "chirp is too long")
		return
	}

	cleaned := cleanBody(params.Body)

	JSONResp(w, 200, returnVals{
		Cleaned_Body: cleaned,
		// Error:        "",
		// Valid:        true,
	})

}

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

func cleanBody(text string) string {
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
