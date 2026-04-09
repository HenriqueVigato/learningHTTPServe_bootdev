package auth

import (
	"net/http"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func TestMakeJWT(t *testing.T) {
	userID := uuid.New()
	secret := "myUltraSecret"
	duration := time.Hour

	t.Run("valid token is created", func(t *testing.T) {
		tokenString, err := MakeJWT(userID, secret, duration)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if tokenString == "" {
			t.Fatalf("expected a token string, got empty string")
		}
	})

	t.Run("token has correct claims", func(t *testing.T) {
		tokenString, _ := MakeJWT(userID, secret, duration)

		token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err != nil {
			t.Fatalf("failed to parse token: %v", err)
		}

		claims, ok := token.Claims.(*jwt.RegisteredClaims)
		if !ok {
			t.Fatalf("failed to extract claims")
		}

		if claims.Subject != userID.String() {
			t.Errorf("expected Subject %s, got %s", userID.String(), claims.Subject)
		}

		if claims.Issuer != "chirpy-access" {
			t.Errorf("expected Issuer 'chirpy-access', got %s", claims.Issuer)
		}
	})

	t.Run("token expires correctly", func(t *testing.T) {
		shortDuration := time.Millisecond * 100
		tokenString, _ := MakeJWT(userID, secret, shortDuration)

		time.Sleep(time.Millisecond * 200)

		_, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			return []byte(secret), nil
		})
		if err == nil {
			t.Fatalf("expected token to be expired, but it was valid")
		}
	})

	t.Run("wrong secret fails validation", func(t *testing.T) {
		tokenString, _ := MakeJWT(userID, secret, duration)

		_, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(t *jwt.Token) (any, error) {
			return []byte("wrongSecret"), nil
		})
		if err == nil {
			t.Fatalf("expected error with wrong secret, but got none")
		}
	})
}

func TestParseJWT(t *testing.T) {
	userID := uuid.New()
	secret := "myUltraSecret"
	duration := time.Millisecond * 1000
	claims := jwt.RegisteredClaims{}

	claims.Issuer = "chirpy-access"
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(duration))
	claims.Subject = userID.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		t.Fatalf("erro ao criar o JWT para testes: %v", err)
	}

	t.Run("Valid token", func(t *testing.T) {
		claimsUserID, err := ValidateJWT(signedToken, secret)
		if err != nil {
			t.Errorf("Nao era esperado nenhum erro aqui: %v", err)
		}

		if claimsUserID != userID {
			t.Errorf("userID retornado do token diferento do esperado")
		}
	})

	t.Run("Expired token", func(t *testing.T) {
		time.Sleep(time.Millisecond * 1100)
		claimsUserID, err := ValidateJWT(signedToken, secret)

		if err == nil {
			t.Errorf("era esperado um erro informando que o token esta espirado %v", err)
		}

		if claimsUserID != uuid.Nil {
			t.Errorf("era esperado que o id retornado fosse nil")
		}
	})
}

func TestGetBearerToken(t *testing.T) {
	tests := []struct {
		name        string
		headers     http.Header
		expected    string
		expectError bool
	}{
		{
			name:        "valid bearer token",
			headers:     http.Header{"Authorization": []string{"Bearer mytoken123"}},
			expected:    "mytoken123",
			expectError: false,
		},
		{
			name:        "missing authorization header",
			headers:     http.Header{},
			expected:    "",
			expectError: true,
		},
		{
			name:        "empty authorization value",
			headers:     http.Header{"Authorization": []string{""}},
			expected:    "",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GetBearerToken(tt.headers)

			if tt.expectError && err == nil {
				t.Error("expected error but got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if token != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, token)
			}
		})
	}
}
