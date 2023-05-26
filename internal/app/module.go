package app

import reusablecode "promotion/internal/reusable_code"

type Modules struct {
	ReusableCode *reusablecode.Module
}

func initModules(infra *infrastructure) *Modules {
	return &Modules{
		ReusableCode: reusablecode.NewModule(infra.log, infra.db.MySQL),
	}
}
