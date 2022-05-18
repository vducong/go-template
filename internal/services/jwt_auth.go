package services

import (
	"context"
	"promotion/pkg/databases"

	"firebase.google.com/go/auth"
)

type JWTAuthService struct {
	firebaseAuth databases.FirebaseAuth
}

func NewJWTAuthService(client databases.FirebaseAuth) *JWTAuthService {
	return &JWTAuthService{
		firebaseAuth: client,
	}
}

func (s *JWTAuthService) VerifyIDToken(idToken string) (token *auth.Token, err error) {
	return s.firebaseAuth.VerifyIDToken(context.Background(), idToken)
}
