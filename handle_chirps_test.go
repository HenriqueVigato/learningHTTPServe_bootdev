package main

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
	"github.com/google/uuid"
)

func TestHandlerChirps(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("nao foi possivel se conectar ao bando de dados %err", err)
	}
	user, err := apiConfig.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro a buscar o usuario no bando de dados %v", err)
	}

	tests := []struct {
		name     string
		body     string
		userID   uuid.UUID
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid chirp",
			body:     "Hello, testing the handler chirps",
			userID:   user.ID,
			wantCode: 201,
			wantBody: "Hello, testing the handler chirps",
		},
		{
			name:     "Too Long chirp",
			body:     "i'm too long to be added a chirpaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasadfsafdaslkfdjaslkdfjaslkjflkasjflksajfklsajflksajflkasjdflkjsalkfjlasfjlskfjlskajflksajaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			userID:   user.ID,
			wantCode: 406,
			wantBody: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf(`{"body":"%v","user_id":"%v"}`, tt.body, tt.userID)
			req := httptest.NewRequest("POST", "/api/chirps", strings.NewReader(input))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			apiConfig.addChirps(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("status code retornado '%v' nao condiz com o esperado '%v'", w.Code, tt.wantCode)
			}

			if !strings.Contains(w.Body.String(), tt.wantBody) {
				t.Errorf("nao foi retornado o body com o chirp criado, foi retornado: %v", w.Body)
			}
		})
	}
}

func TestGetAllChirps(t *testing.T) {
	api, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados para os testes %v", err)
	}

	req := httptest.NewRequest("GET", "/api/chirps", nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	api.getAllChirps(w, req)

	if !strings.Contains(w.Body.String(), "testing the handler chirps") {
		t.Errorf("Nao foi encontrado o chirp esperado %v", w.Body.String())
	}
}

func TestGetChirpsByID(t *testing.T) {
	api, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados para os testes %v", err)
	}

	user, err := api.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro ao buscar o usuario no bando %v", err)
	}

	chirp, err := api.db.CreateChirp(context.Background(), database.CreateChirpParams{
		Body:   "Ola eu sou um chirp pra test",
		UserID: user.ID,
	})
	if err != nil {
		t.Fatalf("erro ao criar um chirp para o teste %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/chirps/{chirpID}", api.getChirpsByID)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/chirps/%s", chirp.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	if !strings.Contains(w.Body.String(), "Ola eu sou um chrip pra test") {
		t.Errorf("nao foi encontrado o chirp desejado, ao inves foi retornado:\n %v", w.Body.String())
	}
}
