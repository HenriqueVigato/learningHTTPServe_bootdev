package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func addUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	param := params{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		log.Fatalf("erro ao ler o request body: %v", err)
	}

	// user, err :=
}
