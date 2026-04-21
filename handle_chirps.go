package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"sort"
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

func (api *apiConfig) getChirps(w http.ResponseWriter, r *http.Request) {
	authorID := r.URL.Query().Get("author_id")
	order := r.URL.Query().Get("sort")
	var chirps []database.Chirp
	var err error

	if authorID == "" {
		chirps, err = getAllChirps(api)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "erro ao obter os chirps")
			return
		}
	} else {
		chirps, err = getChirpByUserID(api, authorID)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, "erro ao obter os chirps")
			return
		}
	}

	if order != "" {
		sortChirp(chirps, order)
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

func getAllChirps(api *apiConfig) ([]database.Chirp, error) {
	chirps, err := api.db.GetAllChirps(context.Background())
	if err != nil {
		return nil, err
	}
	return chirps, nil
}

func getChirpByUserID(api *apiConfig, userID string) ([]database.Chirp, error) {
	parsedID, err := uuid.Parse(userID)
	if err != nil {
		return nil, err
	}

	chirps, err := api.db.GetChirpByUserID(context.Background(), parsedID)
	if err != nil {
		return nil, err
	}
	return chirps, nil
}

func sortChirp(arrayChirp []database.Chirp, order string) ([]database.Chirp, error) {
	switch order {
	case "asc":
		sort.Slice(arrayChirp, func(i, j int) bool {
			return arrayChirp[i].CreatedAt.Before(arrayChirp[j].CreatedAt)
		})
	case "desc":
		sort.Slice(arrayChirp, func(i, j int) bool {
			return arrayChirp[i].CreatedAt.After(arrayChirp[j].CreatedAt)
		})
	default:
		return nil, fmt.Errorf("invalid order value: %s, must be 'asc' or 'desc'", order)
	}
	return arrayChirp, nil
}

func (api *apiConfig) getChirpsByID(w http.ResponseWriter, r *http.Request) {
	chirpReqID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "erro ao obter o id")
		return
	}
	chirpRequested, err := api.db.GetChirpByID(r.Context(), chirpReqID)
	if err != nil {
		respondWithError(w, http.StatusNotFound, err.Error())
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

func (api *apiConfig) deleteChirp(w http.ResponseWriter, r *http.Request) {
	chirpReqID, err := uuid.Parse(r.PathValue("chirpID"))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "erro ao obter o id")
		return
	}

	userToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	userID, err := auth.ValidateJWT(userToken, api.secrete)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, err.Error())
		return
	}

	databaseChirp, err := api.db.GetChirpByID(r.Context(), chirpReqID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			respondWithError(w, http.StatusNotFound, "")
			return
		}
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if databaseChirp.UserID != userID {
		respondWithError(w, http.StatusForbidden, "Nao pode deletar chirps dos coleguinhas")
		return
	}

	err = api.db.DeleteChirpByID(r.Context(), chirpReqID)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
