package databases

import (
	"promotion/configs"
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
		log.Errorf("Failed to connect MySQL: %v", err)
		return db, err
	}
	log.Info("MySQL connection established")

	fb, err := NewFirebaseClient()
	if err != nil {
		log.Errorf("Failed to connect Firebase: %v", err)
		return db, err
	}
	log.Info("Firebase connection established")

	fs, err := NewFirestoreDB(cfg)
	if err != nil {
		log.Errorf("Failed to connect Firestore: %v", err)
		return db, err
	}
	log.Info("Firestore connection established")

	return Databases{
		MySQL:     mysql,
		Firebase:  fb,
		Firestore: fs,
	}, nil
}

func (d *Databases) Close() {
	d.CloseMySQL()
	d.CloseFirestore()
}

func (d *Databases) CloseMySQL() {
	mysql, _ := d.MySQL.DB()
	mysql.Close()
}

func (d *Databases) CloseFirestore() {
	d.Firestore.Close()
}
