package cfld

const (
	envPath        = "CONFIG_FILE_PATH"
	defaultEnvPath = "configs/config.yaml"
)

type LoaderType string

const (
	LoaderTypeCleanenv LoaderType = "cleanenv"
	LoaderTypeJSON     LoaderType = "json"
)

type Loader interface {
	Load(path string, structure any) error
}

func New(loaderType LoaderType) Loader {
	if loaderType == LoaderTypeJSON {
		return &JSONConfigLoader{}
	}

	return &CleanenvConfigLoader{}
}
