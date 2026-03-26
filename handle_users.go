package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAT time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func (api *apiConfig) addUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Email string `json:"email"`
	}
	param := params{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}

	user, err := api.db.CreateUser(r.Context(), param.Email)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create user")
	}
	mapedUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAT: user.UpdatedAt,
		Email:     user.Email,
	}
	if err = respondWithJSON(w, http.StatusCreated, mapedUser); err != nil {
		log.Printf("erro ao responder a request %v", err)
	}
}
