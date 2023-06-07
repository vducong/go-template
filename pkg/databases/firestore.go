package databases

import (
	"context"
	"promotion/configs"
	"promotion/pkg/failure"

	"cloud.google.com/go/firestore"
)

type FirestoreDB = *firestore.Client

func NewFirestoreDB(cfg *configs.Config) (FirestoreDB, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, cfg.GCP.ProjectID)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	return client, nil
}
