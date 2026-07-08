package constant

type Env string

const (
	EnvLocal   Env = "local"
	EnvDev     Env = "dev"
	EnvStaging Env = "staging"
	EnvPreprod Env = "preprod"
	EnvProd    Env = "prod"
)

var NonProdEnvs = []Env{
	EnvLocal,
	EnvDev,
	EnvStaging,
	EnvPreprod,
}
