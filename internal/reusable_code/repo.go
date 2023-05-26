package reusablecode

import (
	"context"
	"promotion/pkg/databases"
	"promotion/pkg/failure"

	"gorm.io/gorm"
)

type Repo struct {
	db databases.MySQLDB
}

func newRepo(db databases.MySQLDB) *Repo {
	return &Repo{db}
}

func (r *Repo) dbWithContext(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

func (r *Repo) GetByCode(ctx context.Context, code string) (*ReusableCode, error) {
	var rc ReusableCode
	if err := r.dbWithContext(ctx).Where(&ReusableCode{Code: code}).Take(&rc).
		Error; err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	return &rc, nil
}
