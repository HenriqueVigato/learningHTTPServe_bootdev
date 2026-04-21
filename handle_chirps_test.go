package main

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
	"github.com/google/uuid"
)

func TestAddChirps(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("nao foi possivel se conectar ao bando de dados %err", err)
	}
	user, err := apiConfig.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro a buscar o usuario no bando de dados %v", err)
	}

	token, err := auth.MakeJWT(user.ID, apiConfig.secrete, time.Hour)
	if err != nil {
		t.Fatalf("erro ao gerar o JWTToken %v", err)
	}

	tests := []struct {
		name     string
		body     string
		userID   uuid.UUID
		token    string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid chirp",
			body:     "Hello, testing the handler chirps",
			userID:   user.ID,
			token:    token, // token válido
			wantCode: 201,
			wantBody: "Hello, testing the handler chirps",
		},
		{
			name:     "Too Long chirp",
			body:     "i'm too long.........................................................................................................................................................................................................................0.",
			userID:   user.ID,
			token:    token, // token válido, falha pelo tamanho
			wantCode: 406,
			wantBody: "error",
		},
		{
			name:     "Missing token",
			body:     "Hello",
			userID:   user.ID,
			token:    "", // sem token
			wantCode: 401,
			wantBody: "error",
		},
		{
			name:     "Invalid token",
			body:     "Hello",
			userID:   user.ID,
			token:    "Bearer invalido.token.aqui",
			wantCode: 401,
			wantBody: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf(`{"body":"%v","user_id":"%v"}`, tt.body, tt.userID)
			req := httptest.NewRequest("POST", "/api/chirps", strings.NewReader(input))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", tt.token))
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
	api.getChirps(w, req)

	if !strings.Contains(w.Body.String(), "testing the handler chirps") {
		t.Errorf("Nao foi encontrado o chirp esperado %v", w.Body.String())
	}
}

func TestGetChirpsByUserID(t *testing.T) {
	api, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados para os testes %v", err)
	}

	user, err := api.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro ao buscar o useario no banco de dados %v", err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/chirps", api.getChirps)

	req := httptest.NewRequest("GET", fmt.Sprintf("/api/chirps?author_id=%s", user.ID), nil)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	mux.ServeHTTP(w, req)

	var chirps []Chirp
	if err := json.NewDecoder(w.Body).Decode(&chirps); err != nil {
		t.Fatalf("erro ao converter o body para um array de struct %v", err)
	}

	if len(chirps) == 0 {
		t.Errorf("a resposta nao deveria estar vazia %v", chirps)
	}

	for _, chirp := range chirps {
		if chirp.UserID != user.ID {
			t.Errorf("deveria apenas retornar os chips do usuario informada (%v) e foi retornado chirp do usuario %v", user.ID, chirp.UserID)
		}
	}
}

func TestSortChirp(t *testing.T) {
	api, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("error preparing the database: %v", err)
	}

	chirps, err := getAllChirps(api)
	if err != nil {
		t.Fatalf("error fetching chirps: %v", err)
	}

	if len(chirps) == 0 {
		t.Fatal("expected chirps in the database but got none")
	}

	testCases := []struct {
		name        string
		order       string
		expectError bool
		validate    func(t *testing.T, chirps []database.Chirp)
	}{
		{
			name:        "sort ascending",
			order:       "asc",
			expectError: false,
			validate: func(t *testing.T, chirps []database.Chirp) {
				for i := 1; i < len(chirps); i++ {
					if chirps[i].CreatedAt.Before(chirps[i-1].CreatedAt) {
						t.Errorf("expected ascending order but got wrong order at index %d", i)
					}
				}
			},
		},
		{
			name:        "sort descending",
			order:       "desc",
			expectError: false,
			validate: func(t *testing.T, chirps []database.Chirp) {
				for i := 1; i < len(chirps); i++ {
					if chirps[i].CreatedAt.After(chirps[i-1].CreatedAt) {
						t.Errorf("expected descending order but got wrong order at index %d", i)
					}
				}
			},
		},
		{
			name:        "invalid order",
			order:       "random",
			expectError: true,
			validate:    nil,
		},
		{
			name:        "empty order",
			order:       "",
			expectError: true,
			validate:    nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			sorted, err := sortChirp(chirps, tc.order)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(sorted) == 0 {
				t.Fatal("expected chirps but got empty slice")
			}

			if tc.validate != nil {
				tc.validate(t, sorted)
			}
		})
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

	if !strings.Contains(w.Body.String(), "Ola eu sou um chirp pra test") {
		t.Errorf("nao foi encontrado o chirp desejado, ao inves foi retornado:\n %v", w.Body.String())
	}
}

func TestDeleteChirp(t *testing.T) {
	api, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados para teste %v", err)
	}

	user, err := api.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro ao buscar o usuario %v", err)
	}
	userToken, err := auth.MakeJWT(user.ID, api.secrete, time.Hour)
	if err != nil {
		t.Fatalf("erro ao geraro o JWTToken %err", err)
	}

	user1, err := api.db.GetUserByEmail(context.Background(), "user1@test.test")
	if err != nil {
		t.Fatalf("erro ao buscar o usuario %v", err)
	}

	user1Token, err := auth.MakeJWT(user1.ID, api.secrete, time.Hour)
	if err != nil {
		t.Fatalf("erro ao geraro o JWTToken %err", err)
	}

	testCases := []struct {
		name       string
		header     http.Header
		wantedCode int
	}{
		{
			name:       "wrongAuthor",
			header:     http.Header{"Authorization": []string{"Bearer " + user1Token}},
			wantedCode: 403,
		},
		{
			name:       "Chirp not found",
			header:     http.Header{"Authorization": []string{"Bearer " + userToken}},
			wantedCode: 404,
		},
		{
			name:       "Succes delete",
			header:     http.Header{"Authorization": []string{"Bearer " + userToken}},
			wantedCode: 204,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var chirp database.Chirp
			if tc.name != "Chirp not found" {
				chirp, _ = api.db.CreateChirp(context.Background(), database.CreateChirpParams{
					Body:   fmt.Sprintf("ola eu sou um chirp Aleatorio %d", rand.IntN(234234)),
					UserID: user.ID,
				})
			} else {
				chirp = database.Chirp{
					ID: uuid.New(),
				}
			}

			req := httptest.NewRequest("DELETE", fmt.Sprintf("/api/chirps/%v", chirp.ID), strings.NewReader(""))
			req.SetPathValue("chirpID", chirp.ID.String())
			req.Header = tc.header
			w := httptest.NewRecorder()

			api.deleteChirp(w, req)

			if w.Code != tc.wantedCode {
				t.Errorf("status code recived (%v) != spected (%v).", w.Code, tc.wantedCode)
				t.Logf("Erro: \n%v", w.Body.String())
			}
		})
	}
}
