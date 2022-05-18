package repos

import "promotion/pkg/databases"

type Repos struct {
	ReusableCode *ReusableCodeRepo
}

func New(db *databases.Databases) *Repos {
	return &Repos{
		ReusableCode: NewReusableCodeRepo(db.MySQL),
	}
}
