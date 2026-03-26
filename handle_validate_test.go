package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestValidateChirp(t *testing.T) {
	tests := []struct {
		name     string
		body     string
		wantCode int
		wantBody string
	}{
		{
			name:     "valid chirp",
			body:     `{"body":"I had something interesting for breakfast"}`,
			wantCode: http.StatusOK,
			wantBody: `{"cleaned_body":"I had something interesting for breakfast"}`,
		},
		{
			name:     "too long",
			body:     `{"body":"this is a very long chirp that exceeds the maximum allowed character limit of one hundred and forty characters so it should fail so there are some more characters"}`,
			wantCode: http.StatusBadRequest,
			wantBody: `{"error":"Chirp is too long"}`,
		},
		{
			name:     "forbidden word",
			body:     `{"body":"I heard that kerfuffle was quite bad"}`,
			wantCode: http.StatusOK,
			wantBody: `{"cleaned_body":"I heard that **** was quite bad"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("POST", "/api/validate_chirp", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			validateChirp(w, req)

			if w.Code != tt.wantCode {
				t.Errorf("got status %d, want %d", w.Code, tt.wantCode)
			}
			if strings.TrimSpace(w.Body.String()) != tt.wantBody {
				t.Errorf("got body %s, want %s", w.Body.String(), tt.wantBody)
			}
		})
	}
}

func TestCleanInput(t *testing.T) {
	forbiddenWords := map[string]struct{}{
		"kerfuffle": {},
		"sharbert":  {},
		"fornax":    {},
	}
	cases := []string{
		"I heard that kerfuffle was quite bad",
		"What a KERFUFFLE that was yesterday",
		"He is such a Sharbert sometimes",
		"The stars are made of FORNAX and sharbert",
		"I had something interesting for breakfast",
		"The weather is nice today",
		"I love coding in Go",
		"kerfuffle",
		"KERFUFFLE SHARBERT FORNAX",
		"kerfuffles",
		"I like kerfuffle, sharbert, fornax",
	}

	for i, str := range cases {
		cases[i] = cleanInput(str, forbiddenWords)
	}

	if !strings.Contains(cases[0], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[0])
	}
	if !strings.Contains(cases[1], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[1])
	}
	if !strings.Contains(cases[2], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[2])
	}
	if !strings.Contains(cases[3], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[3])
	}
	if !strings.Contains(cases[7], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[7])
	}
	if !strings.Contains(cases[8], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[8])
	}
	if strings.Contains(cases[9], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[9])
	}
	if !strings.Contains(cases[10], "****") {
		t.Errorf("era esperado a palavra censurada mas veio: %v", cases[10])
	}
}
