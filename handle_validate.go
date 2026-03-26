package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func validateChirp(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Body string `json:"body"`
	}

	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Something went wrong")
		return
	}
	if len(params.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	forbiddenWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"cleaned_body": cleanInput(params.Body, forbiddenWords)})
}

func cleanInput(s string, forbiddenWords map[string]struct{}) string {
	sArray := strings.Split(s, " ")

	for i, word := range sArray {
		for forbWord := range forbiddenWords {
			if strings.EqualFold(word, forbWord) {
				sArray[i] = "****"
			}
		}
	}

	return strings.Join(sArray, " ")
}
