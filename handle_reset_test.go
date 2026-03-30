package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestReset(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro a configurar o banco de dados de teste %v", err)
	}

	req := httptest.NewRequest("POST", "/admin/reset", strings.NewReader(`{"body":""}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiConfig.reset(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("response status code got %v wanted %v", w.Code, http.StatusOK)
	}
}
