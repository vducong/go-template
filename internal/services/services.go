package services

import (
	"promotion/configs"
	"promotion/internal/repos"
	"promotion/pkg/databases"
	"promotion/pkg/logger"
)

type CreateServicesDTO struct {
	Cfg   *configs.Config
	Log   *logger.Logger
	DB    *databases.Databases
	Repos *repos.Repos
}

type Services struct {
	JWTAuth      *JWTAuthService
	ReusableCode *ReusableCodeService
}

func New(dto *CreateServicesDTO) *Services {
	return &Services{
		JWTAuth:      NewJWTAuthService(dto.DB.Firebase),
		ReusableCode: NewReusableCodeService(dto.Log, dto.Repos.ReusableCode),
	}
}
