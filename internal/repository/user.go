//go:generate mockgen -source=$GOFILE -destination=./mock/mock_$GOFILE -package=mock -typed
package repository

import (
	"context"
)

type UserStore interface {
	Get(ctx context.Context, id string) (string, error)
}

type userStore struct{}

func NewUserStore() UserStore {
	return &userStore{}
}

func (s *userStore) Get(ctx context.Context, id string) (string, error) {
	return "user", nil
}
