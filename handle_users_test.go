package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
)

func TestAddUser(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Logf("erro ao se conenctar ao bando de dados %v", err)
	}

	userInput := `{"email":"email@email.com","password":"passwordnew"}`

	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(userInput))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiConfig.addUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("era esperado o codigo de sucesso mas recebeu o codigo %v", w.Code)
		t.Logf("%v", w.Body)
	}

	user, err := apiConfig.db.GetUserByEmail(context.Background(), "email@email.com")
	if err != nil {
		t.Fatalf("erro ao buscar o usuario no banco de dados %v", err)
	}

	if !strings.Contains(user.HashedPassword, "$argon2id$v=19$m=65536,t=1") {
		t.Errorf("deveria conter o hash da senha inves de %v", user.HashedPassword)
	}
}

func TestLoginUser(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados %v", err)
	}

	testsCase := []struct {
		name       string
		body       string
		wantedCode int
		wantedBody string
	}{
		{
			name:       "Correct test",
			body:       `{"email":"user@test.test","password":"password"}`,
			wantedCode: 200,
			wantedBody: "user@test.test",
		},
		{
			name:       "Correct test - verifica token e refresh_token",
			body:       `{"email":"user@test.test","password":"password"}`,
			wantedCode: 200,
			wantedBody: "token",
		},
		{
			name:       "Correct test - verifica refresh_token",
			body:       `{"email":"user@test.test","password":"password"}`,
			wantedCode: 200,
			wantedBody: "refresh_token",
		},
		{
			name:       "Incorrect password",
			body:       `{"email":"user@test.test","password":"wrongpassword"}`,
			wantedCode: 401,
			wantedBody: "",
		},
	}

	for _, tc := range testsCase {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/login", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			apiConfig.loginUser(w, req)

			if w.Code != tc.wantedCode {
				t.Errorf("deveria retornar %d instead de: %v", tc.wantedCode, w.Code)
			}

			if !strings.Contains(w.Body.String(), tc.wantedBody) {
				t.Errorf("deveria retornar '%v' inves de: \n%v", tc.wantedBody, w.Body.String())
			}
		})
	}
}

func TestUpdateUser(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados %v", err)
	}

	user, err := apiConfig.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro ao buscar o usuario %v", err)
	}

	jwtToken, err := auth.MakeJWT(user.ID, apiConfig.secrete, time.Hour)
	if err != nil {
		t.Fatalf("erro ao ober um jwtToken %v", err)
	}

	testsCase := []struct {
		name       string
		body       string
		header     http.Header
		wantedCode int
		wantedBody string
	}{
		{
			name:       "without token",
			body:       `{"email":"user@test.test", "password":"souUmaSenhaNova"}`,
			header:     http.Header{"Authorization": []string{"Bearer"}},
			wantedCode: 401,
			wantedBody: "error",
		},
		{
			name:       "invalide token",
			body:       `{"email":"user@test.test", "password":"souUmaSenhaNova"}`,
			header:     http.Header{"Authorization": []string{"Bearer invalid token"}},
			wantedCode: 401,
			wantedBody: "error",
		},
		{
			name:       "valid token",
			body:       `{"email":"user@test.test", "password":"souUmaSenhaNova"}`,
			header:     http.Header{"Authorization": []string{"Bearer " + jwtToken}},
			wantedCode: 200,
			wantedBody: `"email":"user@test.test"`,
		},
	}

	for _, tc := range testsCase {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("PUT", "/api/users", strings.NewReader(tc.body))
			req.Header = tc.header
			w := httptest.NewRecorder()

			apiConfig.updateUser(w, req)

			if w.Code != tc.wantedCode {
				t.Errorf("codigo obtido(%v) diferente do esperado(%d)", w.Code, tc.wantedCode)
			}

			if !strings.Contains(w.Body.String(), tc.wantedBody) {
				t.Errorf("body obtido diferente do esperado %v", w.Body.String())
			}
		})
	}
}
