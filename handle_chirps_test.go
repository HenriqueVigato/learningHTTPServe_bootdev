package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHandlerChirps(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("nao foi possivel se conectar ao bando de dados %err", err)
	}

	input := `{"body":"Hello, testing the handler chirps","user_id":"123e4567-e89b-12d3-a456-426614174000"}`
	req := httptest.NewRequest("POST", "/api/chirps", strings.NewReader(input))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiConfig.addChirps(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status code retornado '%v' nao condiz com o esperado '%v'", w.Code, http.StatusCreated)
	}

	if !strings.Contains(w.Body.String(), "testing the handler chirps") {
		t.Errorf("nao foi retornado o body com o chirp criado, foi retornado: %v", w.Body)
	}
}
