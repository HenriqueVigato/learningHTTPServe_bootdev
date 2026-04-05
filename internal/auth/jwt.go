package auth

import (
	"time"

	"github.com/google/uuid"
)

func MakeJWT(userID uuid.UUID, tokensSecret string, expiresIn time.Duration) (string, error) {
	return "", nil
}
