package repos

import (
	"promotion/internal/entity"
	"promotion/pkg/databases"
	"promotion/pkg/failure"

	"github.com/gin-gonic/gin"
)

type ReusableCodeRepo struct {
	db databases.MySQLDB
}

func NewReusableCodeRepo(db databases.MySQLDB) *ReusableCodeRepo {
	return &ReusableCodeRepo{db}
}

func (r *ReusableCodeRepo) GetByCode(ctx *gin.Context, code string) (*entity.ReusableCode, error) {
	var rc entity.ReusableCode
	if err := r.db.Where(&entity.ReusableCode{Code: code}).Take(&rc).
		Error; err != nil {
		return nil, failure.ErrorWithTrace(err)
	}
	return &rc, nil
}
