package auth

import "testing"

func TestMakeRefreshToken(t *testing.T) {
	token := MakeRefreshToken()
	if len(token) != 64 {
		t.Errorf("era esperado um token de 32 bytes")
	}
}
