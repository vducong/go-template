package auth

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func signJWT(t *testing.T, secret []byte, claims jwt.RegisteredClaims) string {
	t.Helper()
	token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(secret)
	if err != nil {
		t.Fatalf("sign token: %v", err)
	}
	return token
}

func validClaims() jwt.RegisteredClaims {
	now := time.Now()
	return jwt.RegisteredClaims{
		Subject:   "user-1",
		Issuer:    "gotemplate",
		Audience:  jwt.ClaimStrings{"gotemplate-api"},
		ExpiresAt: jwt.NewNumericDate(now.Add(time.Hour)),
		IssuedAt:  jwt.NewNumericDate(now),
	}
}

func jwtCfg(secret []byte, leeway time.Duration) JWTConfig {
	return JWTConfig{Secret: secret, Leeway: &leeway}
}

func TestValidateJWT(t *testing.T) {
	t.Parallel()

	secret := []byte("jwt-secret")
	leeway := 30 * time.Second
	cfg := JWTConfig{
		Secret:   secret,
		Leeway:   &leeway,
		Issuer:   "gotemplate",
		Audience: []string{"gotemplate-api"},
	}

	tokenString := signJWT(t, secret, validClaims())

	parsed, err := ValidateJWT(tokenString, cfg)
	if err != nil {
		t.Fatalf("ValidateJWT() error = %v", err)
	}
	if parsed.Subject != "user-1" {
		t.Fatalf("subject = %q, want user-1", parsed.Subject)
	}

	wrongSecret := cfg
	wrongSecret.Secret = []byte("wrong-secret")
	if _, err := ValidateJWT(tokenString, wrongSecret); err == nil {
		t.Fatal("expected error for invalid secret")
	}

	if _, err := ValidateJWT("not-a-jwt", cfg); err == nil {
		t.Fatal("expected error for malformed token")
	}

	if _, err := ValidateJWT(tokenString, JWTConfig{}); err != ErrMissingJWTSecret {
		t.Fatalf("error = %v, want %v", err, ErrMissingJWTSecret)
	}
}

func TestValidateJWT_requiresExp(t *testing.T) {
	t.Parallel()

	secret := []byte("jwt-secret")
	cfg := jwtCfg(secret, 30*time.Second)

	claims := validClaims()
	claims.ExpiresAt = nil
	tokenString := signJWT(t, secret, claims)

	_, err := ValidateJWT(tokenString, cfg)
	if err != ErrMissingExpiry {
		t.Fatalf("error = %v, want %v", err, ErrMissingExpiry)
	}
}

func TestValidateJWT_rejectsExpired(t *testing.T) {
	t.Parallel()

	secret := []byte("jwt-secret")
	cfg := jwtCfg(secret, 0)

	claims := validClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-time.Minute))
	tokenString := signJWT(t, secret, claims)

	if _, err := ValidateJWT(tokenString, cfg); err == nil {
		t.Fatal("expected error for expired token")
	}
}

func TestValidateJWT_leeway(t *testing.T) {
	t.Parallel()

	secret := []byte("jwt-secret")
	claims := validClaims()
	claims.ExpiresAt = jwt.NewNumericDate(time.Now().Add(-10 * time.Second))
	tokenString := signJWT(t, secret, claims)

	if _, err := ValidateJWT(tokenString, jwtCfg(secret, 0)); err == nil {
		t.Fatal("expected expired token to fail with zero leeway")
	}

	if _, err := ValidateJWT(tokenString, jwtCfg(secret, 30*time.Second)); err != nil {
		t.Fatalf("expected token within leeway to pass, got %v", err)
	}
}

func TestValidateJWT_issuerAndAudience(t *testing.T) {
	t.Parallel()

	secret := []byte("jwt-secret")
	leeway := 30 * time.Second
	cfg := JWTConfig{
		Secret:   secret,
		Leeway:   &leeway,
		Issuer:   "gotemplate",
		Audience: []string{"gotemplate-api"},
	}
	tokenString := signJWT(t, secret, validClaims())

	if _, err := ValidateJWT(tokenString, cfg); err != nil {
		t.Fatalf("valid token rejected: %v", err)
	}

	wrongIssuer := cfg
	wrongIssuer.Issuer = "other"
	if _, err := ValidateJWT(tokenString, wrongIssuer); err == nil {
		t.Fatal("expected error for wrong issuer")
	}

	wrongAudience := cfg
	wrongAudience.Audience = []string{"other-api"}
	if _, err := ValidateJWT(tokenString, wrongAudience); err == nil {
		t.Fatal("expected error for wrong audience")
	}
}

func TestValidateJWT_rejectsFutureNbfAndIat(t *testing.T) {
	t.Parallel()

	secret := []byte("jwt-secret")
	cfg := jwtCfg(secret, 0)
	future := time.Now().Add(time.Hour)

	nbfClaims := validClaims()
	nbfClaims.NotBefore = jwt.NewNumericDate(future)
	if _, err := ValidateJWT(signJWT(t, secret, nbfClaims), cfg); err == nil {
		t.Fatal("expected error for future nbf")
	}

	iatClaims := validClaims()
	iatClaims.IssuedAt = jwt.NewNumericDate(future)
	if _, err := ValidateJWT(signJWT(t, secret, iatClaims), cfg); err == nil {
		t.Fatal("expected error for future iat")
	}
}

func TestJWTConfigFromAuth(t *testing.T) {
	t.Parallel()

	if _, err := JWTConfigFromAuth("", 0, "", nil); err != ErrMissingJWTSecret {
		t.Fatalf("error = %v, want %v", err, ErrMissingJWTSecret)
	}

	cfg, err := JWTConfigFromAuth("secret", 45*time.Second, "issuer", []string{"aud"})
	if err != nil {
		t.Fatalf("JWTConfigFromAuth() error = %v", err)
	}
	if string(cfg.Secret) != "secret" || cfg.leeway() != 45*time.Second || cfg.Issuer != "issuer" {
		t.Fatalf("unexpected config: %+v", cfg)
	}

	defaultCfg, err := JWTConfigFromAuth("secret", 0, "", nil)
	if err != nil {
		t.Fatalf("JWTConfigFromAuth() error = %v", err)
	}
	if defaultCfg.leeway() != defaultJWTLeeway {
		t.Fatalf("leeway = %v, want %v", defaultCfg.leeway(), defaultJWTLeeway)
	}
}
