package reusablecode

import (
	"context"
	"promotion/pkg/failure"
	"promotion/pkg/logger"
)

type Service struct {
	log  *logger.Logger
	repo *Repo
}

func newService(log *logger.Logger, repo *Repo) *Service {
	return &Service{log, repo}
}

func (s *Service) GetByCode(ctx context.Context, code string) (*ReusableCode, error) {
	rc, err := s.repo.GetByCode(ctx, code)
	if err != nil {
		return nil, failure.ErrWithTrace(err)
	}
	return rc, nil
}
