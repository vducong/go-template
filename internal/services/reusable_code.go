package services

import (
	"promotion/internal/entity"
	"promotion/internal/repos"
	"promotion/pkg/failure"
	"promotion/pkg/logger"

	"github.com/gin-gonic/gin"
)

type ReusableCodeService struct {
	log  *logger.Logger
	repo *repos.ReusableCodeRepo
}

func NewReusableCodeService(
	log *logger.Logger, repo *repos.ReusableCodeRepo,
) *ReusableCodeService {
	return &ReusableCodeService{log, repo}
}

func (s *ReusableCodeService) GetByCode(ctx *gin.Context, code string) (*entity.ReusableCode, error) {
	rc, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, failure.ErrorWithTrace(err)
	}
	return rc, nil
}
