package reusablecode

import (
	"promotion/pkg/databases"
	"promotion/pkg/logger"
)

type Module struct {
	Repo    *Repo
	Service *Service
}

func NewModule(log *logger.Logger, db databases.MySQLDB) *Module {
	repo := newRepo(db)
	s := newService(log, repo)
	return &Module{repo, s}
}
