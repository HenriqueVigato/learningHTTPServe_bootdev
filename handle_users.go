package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAT    time.Time `json:"updated_at"`
	Email        string    `json:"email"`
	Token        string    `json:"token"`
	RefreshToken string    `json:"refresh_token"`
	IsChirpyRed  bool      `json:"is_chirpy_red"`
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
		ID:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAT:   user.UpdatedAt,
		Email:       user.Email,
		IsChirpyRed: user.IsChirpyRed,
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
		return
	}

	user, err := api.db.GetUserByEmail(r.Context(), param.Email)
	if err != nil {
		log.Printf("erro ao buscar o usuario no banco de dados: \n %v", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}

	correctHash, err := auth.CheckPasswordHash(param.Password, user.HashedPassword)
	if err != nil {
		log.Printf("erro o comparar a senha com o hash: \n%v", err)
		respondWithError(w, http.StatusInternalServerError, "")
		return
	}
	if !correctHash {
		respondWithError(w, http.StatusUnauthorized, "something is wrong whith your input")
		return
	}

	token, err := auth.MakeJWT(user.ID, api.secrete, 0)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "nao foi possivel gerar o token")
		return
	}

	// Registra o refresh_token no banco de dados
	refreshToken := auth.MakeRefreshToken()
	refreshParams := database.CreateRefreshTokenParams{
		Token:     refreshToken,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().AddDate(0, 0, 60),
	}
	refreshTokenDB, err := api.db.CreateRefreshToken(r.Context(), refreshParams)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "erro ao registrar o refresh token")
		return
	}

	mapedUser := User{
		ID:           user.ID,
		CreatedAt:    user.CreatedAt,
		UpdatedAT:    user.UpdatedAt,
		Email:        user.Email,
		Token:        token,
		RefreshToken: refreshTokenDB.Token,
		IsChirpyRed:  user.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, mapedUser)
}

func (api *apiConfig) updateUser(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Password string `json:"password"`
		Email    string `json:"email"`
		Token    string `json:"token"`
	}

	param := requestParams{}
	var err error

	param.Token, err = auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "nao foi possivel localizar o token")
		return
	}

	userByToken, err := auth.ValidateJWT(param.Token, api.secrete)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "token Invalido")
		return
	}

	decoder := json.NewDecoder(r.Body)
	if err = decoder.Decode(&param); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid reqcuest body")
		return
	}

	hashedPassword, err := auth.HashPassword(param.Password)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao hashear a senha %v", err)
		return
	}

	updateArgs := database.UpdateUserParams{
		Email:          param.Email,
		HashedPassword: hashedPassword,
		ID:             userByToken,
	}

	updatedUser, err := api.db.UpdateUser(r.Context(), updateArgs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao atualizar o usuario %v", err)
		return
	}

	mapedUser := User{
		ID:          updatedUser.ID,
		CreatedAt:   updatedUser.CreatedAt,
		UpdatedAT:   updatedUser.UpdatedAt,
		Email:       updatedUser.Email,
		Token:       param.Token,
		IsChirpyRed: updatedUser.IsChirpyRed,
	}

	respondWithJSON(w, http.StatusOK, mapedUser)
}

func (api *apiConfig) updateUserRed(w http.ResponseWriter, r *http.Request) {
	type requestParams struct {
		Event string `json:"event"`
		Data  struct {
			UserID uuid.UUID `json:"user_id"`
		} `json:"data"`
	}

	params := requestParams{}

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&params); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if params.Event != "user.upgraded" {
		respondWithJSON(w, http.StatusNoContent, "")
		return
	}

	userRedParams := database.UpdateUserRedParams{
		IsChirpyRed: true,
		ID:          params.Data.UserID,
	}

	_, err := api.db.UpdateUserRed(r.Context(), userRedParams)
	if err == sql.ErrNoRows {
		respondWithError(w, http.StatusNotFound, "User not found")
		return
	}
	if err != nil {
		respondWithError(w, http.StatusBadRequest, fmt.Sprintf("%v", err))
		return
	}

	respondWithJSON(w, http.StatusNoContent, "")
}
