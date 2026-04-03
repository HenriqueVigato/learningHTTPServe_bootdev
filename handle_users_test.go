package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
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

	if !strings.Contains(user.HashedPassword, "$argon2id$v=19$m=65536,t=1,p=32$") {
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
			name:       "Incorrect password",
			body:       `{"email":"user@test.test","password":"wrongpassword"}`,
			wantedCode: 401,
			wantedBody: "",
		},
	}

	for _, tc := range testsCase {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/users", strings.NewReader(tc.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			apiConfig.loginUser(w, req)

			if w.Code != tc.wantedCode {
				t.Errorf("deveria retorar o %d inves de: %v", tc.wantedCode, w.Code)
			}

			if !strings.Contains(w.Body.String(), tc.wantedBody) {
				t.Errorf("deveria retornar o usuario inves de: \n%v", w.Body.String())
			}
		})
	}
}
