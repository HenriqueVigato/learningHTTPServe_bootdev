package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
	"github.com/google/uuid"
)

type Chirp struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Body      string    `json:"body"`
	UserID    uuid.UUID `json:"user_id"`
}

func (api *apiConfig) addChirps(w http.ResponseWriter, r *http.Request) {
	type params struct {
		Body string `json:"body"`
	}

	token, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "erro ao obter o token")
		return
	}

	userID, err := auth.ValidateJWT(token, api.secrete)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token invalido")
		return
	}

	var requestBody params

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&requestBody); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	requestBody.Body, err = validateChirp(requestBody.Body)
	if err != nil {
		respondWithError(w, http.StatusNotAcceptable, err.Error())
		return
	}

	param := database.CreateChirpParams{
		Body:   requestBody.Body,
		UserID: userID,
	}

	chirp, err := api.db.CreateChirp(r.Context(), param)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Could not create chirp")
		return
	}

	mapedChirp := Chirp{
		ID:        chirp.ID,
		CreatedAt: chirp.CreatedAt,
		UpdatedAt: chirp.UpdatedAt,
		Body:      chirp.Body,
		UserID:    chirp.UserID,
	}

	if err = respondWithJSON(w, http.StatusCreated, mapedChirp); err != nil {
		log.Printf("erro ao responder a request %v", err)
	}
}

func (api *apiConfig) getAllChirps(w http.ResponseWriter, r *http.Request) {
	chirps, err := api.db.GetAllChirps(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Erro ao buscar os chirps no banco de dados: %v", err))
		return
	}

	chirpsJSON := []Chirp{}

	for _, chirp := range chirps {
		chirpsJSON = append(chirpsJSON, Chirp{
			ID:        chirp.ID,
			CreatedAt: chirp.CreatedAt,
			UpdatedAt: chirp.UpdatedAt,
			Body:      chirp.Body,
			UserID:    chirp.UserID,
		})
	}

	if err = respondWithJSON(w, http.StatusOK, chirpsJSON); err != nil {
		log.Printf("erro ao responder a request: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Erro ao responder com os chirps")
		return
	}
}

func (api *apiConfig) getChirpsByID(w http.ResponseWriter, r *http.Request) {
	chirpReqID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "erro ao obter o id")
		return
	}
	chirpRequested, err := api.db.GetChirpByID(r.Context(), chirpReqID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, "")
		return
	}

	respondWithJSON(w, http.StatusOK, Chirp{
		ID:        chirpRequested.ID,
		CreatedAt: chirpRequested.CreatedAt,
		UpdatedAt: chirpRequested.UpdatedAt,
		Body:      chirpRequested.Body,
		UserID:    chirpRequested.UserID,
	})
}
