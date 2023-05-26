package middleware

import (
	"net/http"
	"promotion/constants"
	"promotion/pkg/databases"
	"promotion/pkg/response"
	"strings"

	"firebase.google.com/go/v4/auth"
	"github.com/gin-gonic/gin"
)

const (
	JWTAuthTokenPrefix = "Bearer"
	JWTAuthHeader      = "authorization"
)

type JWTAuthMiddleware struct {
	firebaseAuth databases.FirebaseAuth
}

func NewJWTAuthMiddleware(
	firebaseAuth databases.FirebaseAuth,
) *JWTAuthMiddleware {
	return &JWTAuthMiddleware{firebaseAuth}
}

func (m *JWTAuthMiddleware) Handler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader(JWTAuthHeader)
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPResponse{
				Status: http.StatusUnauthorized,
				Data:   "Missing auth header",
			})
			return
		}

		idToken := m.getIDToken(authHeader)
		if idToken == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPResponse{
				Status: http.StatusUnauthorized,
				Data:   "Missing JWT",
			})
			return
		}

		token, err := m.verifyIDToken(ctx, idToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response.HTTPResponse{
				Status: http.StatusUnauthorized,
				Data:   "Malformed JWT",
			})
			return
		}

		ctx.Set(string(constants.ContextKeyUserID), token.UID)
		ctx.Next()
	}
}

func (m *JWTAuthMiddleware) getIDToken(header string) string {
	return strings.TrimSpace(strings.Replace(header, JWTAuthTokenPrefix, "", 1))
}

func (m *JWTAuthMiddleware) verifyIDToken(ctx *gin.Context, idToken string) (token *auth.Token, err error) {
	return m.firebaseAuth.VerifyIDToken(ctx, idToken)
}
