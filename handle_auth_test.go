package main

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/auth"
	"github.com/HenriqueVigato/learningHTTPServe_bootdev/internal/database"
)

func TestValidateRefreshToken(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados %v", err)
	}
	token := auth.MakeRefreshToken()
	user, err := apiConfig.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro ao fazer consulta no banco de dados: %v", err)
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Second),
	}
	regRP, err := apiConfig.db.CreateRefreshToken(context.Background(), refreshParams)
	if err != nil {
		t.Fatalf("erro ao salvar o refresh token no banco %v", err)
	}

	testsCase := []struct {
		name       string
		header     string
		wantedCode int
		wantedBody string
	}{
		{
			name:       "Correct Test",
			header:     fmt.Sprintf("Bearer %s", regRP.Token),
			wantedCode: 200,
			wantedBody: "token",
		},
		{
			name:       "Incorrect Token",
			header:     fmt.Sprintf("Bearer %s", "souUmTokenInvalido"),
			wantedCode: 401,
			wantedBody: "error",
		},
		{
			name:       "expired token",
			header:     fmt.Sprintf("Bearer %s", regRP.Token),
			wantedCode: 401,
			wantedBody: "error",
		},
	}

	for _, tc := range testsCase {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "expired token" {
				time.Sleep(time.Second)
			}
			req := httptest.NewRequest("POST", "/api/refresh", strings.NewReader(""))
			req.Header.Set("Authorization", tc.header)

			w := httptest.NewRecorder()

			apiConfig.validateRefreshToken(w, req)

			if w.Code != tc.wantedCode {
				t.Errorf("wantedCode: %v got: %v", tc.wantedCode, w.Code)
			}

			if !strings.Contains(w.Body.String(), tc.wantedBody) {
				t.Errorf("Got body: %v Wanted body: %v", w.Body.String(), tc.wantedBody)
			}
		})
	}
}

func TestRevokeRefreshToken(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("erro ao preparar o banco de dados %v", err)
	}
	token := auth.MakeRefreshToken()
	user, err := apiConfig.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro ao fazer consulta no banco de dados: %v", err)
	}

	refreshParams := database.CreateRefreshTokenParams{
		Token:     token,
		UserID:    user.ID,
		ExpiresAt: time.Now().UTC().Add(time.Second),
	}
	regRP, err := apiConfig.db.CreateRefreshToken(context.Background(), refreshParams)
	if err != nil {
		t.Fatalf("erro ao salvar o refresh token no banco %v", err)
	}

	t.Run("Revoked Token", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/revoke", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", regRP.Token))
		w := httptest.NewRecorder()

		apiConfig.revokeRefreshToken(w, req)

		if w.Code != 204 {
			t.Errorf("era esperado o codigo 204 mas recebeu: %v", w.Code)
		}
	})

	t.Run("Is really revoked?", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/api/revoke", strings.NewReader(""))
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", regRP.Token))
		w := httptest.NewRecorder()

		apiConfig.validateRefreshToken(w, req)

		if w.Code != 401 {
			t.Errorf("era esperado o codigo 204 mas recebeu: %v", w.Code)
		}

		if !strings.Contains(w.Body.String(), "error") {
			t.Errorf("era esperado um erro, mas recebeu %v", w.Body.String())
		}
	})
}
