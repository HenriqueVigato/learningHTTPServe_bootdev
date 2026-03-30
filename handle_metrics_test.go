package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestMetrics(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao se conectar ao banco de dados %v", err)
	}

	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})
	handler := apiConfig.middlewareMetricsInt(dummyHandler)

	for range 10 {
		req := httptest.NewRequest("GET", "/app", nil)
		w := httptest.NewRecorder()
		handler.ServeHTTP(w, req)
	}

	req := httptest.NewRequest("GET", "/admin/metrics", strings.NewReader(`"body":""`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	apiConfig.metrics(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status code nao corresponde ao esperado %v", w.Code)
	}

	if !strings.Contains(w.Body.String(), "10 times!") {
		t.Errorf("nao contem a quantidade de visitas esperadas %v", w.Body)
	}
}
