package auth

import (
	"context"
	"crypto/subtle"

	"github.com/golang-jwt/jwt/v5"
)

type claimsContextKey struct{}

// Claims holds JWT registered claims after successful validation.
type Claims struct {
	jwt.RegisteredClaims
}

// ClaimsFromContext returns JWT claims stored by middleware on success.
func ClaimsFromContext(ctx context.Context) (*Claims, bool) {
	claims, ok := ctx.Value(claimsContextKey{}).(*Claims)
	return claims, ok
}

// ContextWithClaims stores validated claims on the context.
func ContextWithClaims(ctx context.Context, claims *Claims) context.Context {
	return context.WithValue(ctx, claimsContextKey{}, claims)
}

// PrepareAPIKeys converts configured keys to byte slices for constant-time comparison.
// Empty keys are skipped.
func PrepareAPIKeys(keys []string) [][]byte {
	prepared := make([][]byte, 0, len(keys))
	for _, k := range keys {
		if k == "" {
			continue
		}
		prepared = append(prepared, []byte(k))
	}
	return prepared
}

// ValidateAPIKey checks provided against expected keys
// using constant-time comparison to mitigate timing attacks.
func ValidateAPIKey(provided string, expectedKeys [][]byte) bool {
	if provided == "" || len(expectedKeys) == 0 {
		return false
	}

	candidate := []byte(provided)
	for _, expected := range expectedKeys {
		if subtle.ConstantTimeCompare(candidate, expected) == 1 {
			return true
		}
	}
	return false
}
