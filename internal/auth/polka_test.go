package auth

import (
	"net/http"
	"testing"
)

func TestGetApiKey(t *testing.T) {
	tests := []struct {
		name          string
		authorization string
		expected      string
		expectError   bool
	}{
		{
			name:          "valid api key",
			authorization: "ApiKey my-secret-key",
			expected:      "my-secret-key",
			expectError:   false,
		},
		{
			name:          "valid api key case insensitive",
			authorization: "APIKEY my-secret-key",
			expected:      "my-secret-key",
			expectError:   false,
		},
		{
			name:          "missing authorization header",
			authorization: "",
			expected:      "",
			expectError:   true,
		},
		{
			name:          "wrong scheme - Bearer instead of ApiKey",
			authorization: "Bearer my-secret-key",
			expected:      "",
			expectError:   true,
		},
		{
			name:          "missing token value",
			authorization: "ApiKey",
			expected:      "",
			expectError:   true,
		},
		{
			name:          "malformed header - no space",
			authorization: "ApiKeymy-secret-key",
			expected:      "",
			expectError:   true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			headers := make(http.Header)
			if tc.authorization != "" {
				headers.Set("Authorization", tc.authorization)
			}

			result, err := GetApiKey(headers)

			if tc.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}

				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if result != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, result)
			}
		})
	}
}
