package auth

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokensSecret string, expiresIn time.Duration) (string, error) {
	claims := jwt.RegisteredClaims{}

	claims.Issuer = "chirpy-access"
	claims.IssuedAt = jwt.NewNumericDate(time.Now())
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(expiresIn))
	claims.Subject = userID.String()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(tokensSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ValidateJWT(tokenString, tokenSecrete string) (uuid.UUID, error) {
	myClaims := jwt.RegisteredClaims{}

	_, err := jwt.ParseWithClaims(tokenString, &myClaims, func(token *jwt.Token) (any, error) {
		return []byte(tokenSecrete), nil
	})
	if err != nil {
		return uuid.Nil, err
	}

	if myClaims.Issuer != "chirpy-access" {
		return uuid.Nil, fmt.Errorf("erro Issuer diferente do esperado")
	}

	return uuid.Parse(myClaims.Subject)
}

func GetBearerToken(headers http.Header) (string, error) {
	auth := headers.Get("Authorization")
	parts := strings.SplitN(auth, " ", 2)
	if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
		return "", fmt.Errorf("authorization header format must be: Bearer <token>")
	}

	return parts[1], nil
}
