package main

import (
	"context"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestHandlerChirps(t *testing.T) {
	apiConfig, err := getConnectionTestDB(t)
	if err != nil {
		t.Fatalf("nao foi possivel se conectar ao bando de dados %err", err)
	}
	user, err := apiConfig.db.GetUserByEmail(context.Background(), "user@test.test")
	if err != nil {
		t.Fatalf("erro a buscar o usuario no bando de dados %v", err)
	}

	tests := []struct {
		name     string
		body     string
		userID   uuid.UUID
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid chirp",
			body:     "Hello, testing the handler chirps",
			userID:   user.ID,
			wantCode: 201,
			wantBody: "Hello, testing the handler chirps",
		},
		{
			name:     "Too Long chirp",
			body:     "i'm too long to be added a chirpaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaasadfsafdaslkfdjaslkdfjaslkjflkasjflksajfklsajflksajflkasjdflkjsalkfjlasfjlskfjlskajflksajaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
			userID:   user.ID,
			wantCode: 406,
			wantBody: "error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input := fmt.Sprintf(`{"body":"%v","user_id":"%v"}`, tt.body, tt.userID)
			req := httptest.NewRequest("POST", "/api/chirps", strings.NewReader(input))
			req.Header.Set("Content-Type", "application/json")
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
