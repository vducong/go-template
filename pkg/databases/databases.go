package databases

import (
	"fmt"
	"promotion/configs"
	"promotion/pkg/failure"
	"promotion/pkg/logger"
)

type Databases struct {
	MySQL     MySQLDB
	Firebase  FirebaseAuth
	Firestore FirestoreDB
}

func New(cfg *configs.Config, log *logger.Logger) (Databases, error) {
	db := Databases{}
	mysql, err := NewMySQLDB(cfg)
	if err != nil {
		return db, failure.ErrWithTrace(fmt.Errorf("Failed to connect MySQL: %w", err))
	}
	log.Info("MySQL connection established")

	fb, err := NewFirebaseClient()
	if err != nil {
		return db, failure.ErrWithTrace(fmt.Errorf("Failed to connect Firebase: %w", err))
	}
	log.Info("Firebase connection established")

	fs, err := NewFirestoreDB(cfg)
	if err != nil {
		return db, failure.ErrWithTrace(fmt.Errorf("Failed to connect Firestore: %w", err))
	}
	log.Info("Firestore connection established")

	return Databases{
		MySQL:     mysql,
		Firebase:  fb,
		Firestore: fs,
	}, nil
}
