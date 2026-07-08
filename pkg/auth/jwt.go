package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrMissingJWTSecret = errors.New("jwt secret is required")
	ErrMissingExpiry    = errors.New("jwt exp claim is required")
)

const defaultJWTLeeway = 30 * time.Second

// JWTConfig holds JWT validation settings.
type JWTConfig struct {
	Secret   []byte
	Leeway   *time.Duration // nil uses defaultJWTLeeway
	Issuer   string
	Audience []string
}

func (c JWTConfig) leeway() time.Duration {
	if c.Leeway == nil {
		return defaultJWTLeeway
	}
	return *c.Leeway
}

// ValidateJWT parses and verifies an HMAC-signed JWT using cfg.
// exp is required; exp, nbf, and iat are validated when present.
func ValidateJWT(tokenString string, cfg JWTConfig) (*Claims, error) {
	if len(cfg.Secret) == 0 {
		return nil, ErrMissingJWTSecret
	}

	claims := &Claims{}
	parserOpts := []jwt.ParserOption{
		jwt.WithLeeway(cfg.leeway()),
		jwt.WithIssuedAt(),
	}
	if cfg.Issuer != "" {
		parserOpts = append(parserOpts, jwt.WithIssuer(cfg.Issuer))
	}
	if len(cfg.Audience) > 0 {
		parserOpts = append(parserOpts, jwt.WithAudience(cfg.Audience...))
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrTokenSignatureInvalid
		}
		return cfg.Secret, nil
	}, parserOpts...)
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, jwt.ErrTokenInvalidClaims
	}
	if claims.ExpiresAt == nil {
		return nil, ErrMissingExpiry
	}
	return claims, nil
}

// JWTConfigFromAuth builds JWT validation settings from auth config fields.
func JWTConfigFromAuth(secret string, leeway time.Duration, issuer string, audience []string) (JWTConfig, error) {
	if secret == "" {
		return JWTConfig{}, ErrMissingJWTSecret
	}
	if leeway <= 0 {
		leeway = defaultJWTLeeway
	}
	return JWTConfig{
		Secret:   []byte(secret),
		Leeway:   &leeway,
		Issuer:   issuer,
		Audience: audience,
	}, nil
}
