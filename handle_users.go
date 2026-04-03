package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
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
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	param := params{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}

	passwordHash, err := auth.HashPassword(param.Password)
	if err != nil {
		log.Printf("erro ao hashear a senha %v", err)
		respondWithError(w, http.StatusInternalServerError, "")
	}

	userParams := database.CreateUserParams{
		Email:          param.Email,
		HashedPassword: passwordHash,
	}

	user, err := api.db.CreateUser(r.Context(), userParams)
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

func (api *apiConfig) loginUser(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	param := params{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&param); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}

	user, err := api.db.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		log.Printf("erro ao buscar o usuario no banco de dados: \n %v", err)
		respondWithError(w, http.StatusInternalServerError, "")
	}

	correctHash, err := auth.CheckPasswordHash(param.Password, user.HashedPassword)
	if err != nil {
		log.Printf("erro o comparar a senha com o hash: \n%v", err)
		respondWithError(w, http.StatusInternalServerError, "")
	}
	mapedUser := User{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAT: user.UpdatedAt,
		Email:     user.Email,
	}

	if correctHash {
		respondWithJSON(w, http.StatusOK, mapedUser)
	} else {
		respondWithError(w, http.StatusUnauthorized, "")
	}
}
