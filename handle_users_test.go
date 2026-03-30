package main

import (
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

	userEmail := `{"email":"email@email.com"}`
	req := httptest.NewRequest("POST", "/api/users", strings.NewReader(userEmail))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiConfig.addUser(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("era esperado o codigo de sucesso mas recebeu o codigo %v", w.Code)
		t.Logf("%v", w.Body)
	}
}
