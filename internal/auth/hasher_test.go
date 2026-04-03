package auth

import "testing"

func TestHashPassword(t *testing.T) {
	password := "ola eu sou uma senha besta"

	_, err := HashPassword(password)
	if err != nil {
		t.Errorf("erro ao hashear a senha %v", err)
	}
}

func TestCheckPasswordHash(t *testing.T) {
	testCase := []struct {
		name           string
		password       string
		hash           string
		expectedResult bool
	}{
		{
			name:           "Valid hash",
			password:       "ola eu sou uma senha",
			hash:           "$argon2id$v=19$m=65536,t=1,p=8$fwgBWqPcpJZSSnIj5ojmBw$/OZxALIq46xpj92a8DydeydVtqO1EA0bO2zFg7nguC0",
			expectedResult: true,
		},
		{
			name:           "Unvalid hash",
			password:       "ola eu sou uma senha",
			hash:           "$argon2id$v=19$m=65536,t=1,p=8$fwgBWqPcpJZSSnIj5ojmBw$/asdfsadfsadf92a8DydeydVtqO1EA0bO2zFg7nguC0",
			expectedResult: false,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			value, err := CheckPasswordHash(tc.password, tc.hash)
			if err != nil {
				t.Fatalf("erro ao comparar a senha %v", err)
			}

			if value != tc.expectedResult {
				t.Errorf("resultado (%v) diferente do esperado", value)
			}
		})
	}
}
