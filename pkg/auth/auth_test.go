package auth

import (
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

func TestValidateAPIKey(t *testing.T) {
	t.Parallel()

	keys := PrepareAPIKeys([]string{"secret-key", "other-key"})

	tests := []struct {
		name     string
		provided string
		expected bool
	}{
		{name: "match first key", provided: "secret-key", expected: true},
		{name: "match second key", provided: "other-key", expected: true},
		{name: "no match", provided: "wrong-key", expected: false},
		{name: "empty provided", provided: "", expected: false},
		{name: "partial prefix", provided: "secret", expected: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := ValidateAPIKey(tt.provided, keys); got != tt.expected {
				t.Fatalf("ValidateAPIKey(%q) = %v, want %v", tt.provided, got, tt.expected)
			}
		})
	}

	if ValidateAPIKey("secret-key", nil) {
		t.Fatal("expected false when no expected keys configured")
	}
}

func TestClaimsFromContext(t *testing.T) {
	t.Parallel()

	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{Subject: "user-42"},
	}
	ctx := ContextWithClaims(t.Context(), claims)

	got, ok := ClaimsFromContext(ctx)
	if !ok {
		t.Fatal("expected claims in context")
	}
	if got.Subject != "user-42" {
		t.Fatalf("subject = %q, want user-42", got.Subject)
	}
}
