package main

import (
	"fmt"
	"strings"
)

func validateChirp(chirp string) (string, error) {
	if len(chirp) > 140 {
		return "", fmt.Errorf("Chirp is too long")
	}

	forbiddenWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	return cleanInput(chirp, forbiddenWords), nil
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
