package main

import (
	"encoding/json"
	"net/http"
	"strings"
)

func respondWithJSON(w http.ResponseWriter, code int, payload any) error {
	response, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	w.Write(response)
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}

func cleanInput(s string) string {
	forbiddenWords := []string{"kerfuffle", "sharbert", "fornax"}
	sArray := strings.Split(s, " ")

	for i, word := range sArray {
		for _, forbWord := range forbiddenWords {
			if strings.EqualFold(word, forbWord) {
				sArray[i] = "****"
			}
		}
	}

	return strings.Join(sArray, " ")
}
