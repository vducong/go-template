package databases

import (
	"context"
	"fmt"
	"promotion/pkg/failure"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
)

type FirebaseAuth = *auth.Client

func NewFirebaseClient() (FirebaseAuth, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to connect Firebase: %w", err))
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to get Firebase instance: %w", err))
	}
	return client, nil
}
