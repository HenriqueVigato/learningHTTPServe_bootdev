package main

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
)

type Tk struct {
	Token string `json:"token"`
}

func (api *apiConfig) validateRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao obter o token do header %v", err)
		return
	}

	refreshToken, err := api.db.GetRefreshTokenByToken(r.Context(), bearerToken)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "")
		log.Printf("erro ao buscar o token no banco de dados %v", err)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	if refreshToken.RevokedAt.Valid {
		respondWithError(w, http.StatusUnauthorized, "")
		return
	}

	user, err := api.db.GetUserFromRefresToken(r.Context(), refreshToken.Token)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao buscar o usuario %v", err)
		return
	}

	tokenJWT, err := auth.MakeJWT(user.ID, api.secrete, 0)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao gerar o JWT token %v", err)
		return
	}

	response := Tk{
		Token: tokenJWT,
	}
	respondWithJSON(w, http.StatusOK, response)
}

func (api *apiConfig) revokeRefreshToken(w http.ResponseWriter, r *http.Request) {
	bearerToken, err := auth.GetBearerToken(r.Header)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao obter o token do header %v", err)
		return
	}

	revokeTokens := database.RevokeRefreshTokenParams{
		RevokedAt: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		UpdatedAt: time.Now(),
		Token:     bearerToken,
	}

	_, err = api.db.RevokeRefreshToken(r.Context(), revokeTokens)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "")
		log.Printf("erro ao atualizar o refreshToken no db %v", err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
