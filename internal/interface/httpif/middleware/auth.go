package middleware

import (
	"gotemplate/internal/apperr"
	"gotemplate/pkg/auth"
	"gotemplate/pkg/httprespwrit"
	"gotemplate/pkg/lg"
	"net/http"
	"strings"
)

type Authenticator struct {
	logger    lg.Logger
	writer    httprespwrit.Writer
	jwtConfig auth.JWTConfig
	apiKeys   [][]byte
}

func NewAuthenticator(
	logger lg.Logger,
	writer httprespwrit.Writer,
	jwtConfig auth.JWTConfig,
	apiKeys []string,
) (*Authenticator, error) {
	if len(jwtConfig.Secret) == 0 {
		return nil, auth.ErrMissingJWTSecret
	}
	return &Authenticator{
		logger:    logger,
		writer:    writer,
		jwtConfig: jwtConfig,
		apiKeys:   auth.PrepareAPIKeys(apiKeys),
	}, nil
}

func (a *Authenticator) RequireJwt() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			tokenString, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || tokenString == "" {
				a.unauthorized(w, r)
				return
			}

			claims, err := auth.ValidateJWT(tokenString, a.jwtConfig)
			if err != nil {
				a.unauthorized(w, r)
				return
			}

			next.ServeHTTP(w, r.WithContext(auth.ContextWithClaims(r.Context(), claims)))
		})
	}
}

func (a *Authenticator) RequireAPIKey() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !auth.ValidateAPIKey(r.Header.Get("X-API-Key"), a.apiKeys) {
				a.unauthorized(w, r)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func (a *Authenticator) unauthorized(w http.ResponseWriter, r *http.Request) {
	a.writer.WriteError(w, r, &httprespwrit.ErrorResponse{
		Err: apperr.Unauthorized,
	})
}
