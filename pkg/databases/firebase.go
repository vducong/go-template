package databases

import (
	"context"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

type FirebaseAuth = *auth.Client

func NewFirebaseClient() (FirebaseAuth, error) {
	ctx := context.Background()
	app, err := firebase.NewApp(ctx, nil)
	if err != nil {
		return nil, err
	}

	client, err := app.Auth(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}
