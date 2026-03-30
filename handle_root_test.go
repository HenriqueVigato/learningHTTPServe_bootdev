package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoot(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("nao foi possivel se conectar a o bando de dados %v", err)
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := apiConfig.middlewareMetricsInt(dummyHandler)

	req := httptest.NewRequest("GET", "/app", nil)
	w := httptest.NewRecorder()
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code retornado: %v diverge do esperado 200", w.Code)
	}
}
