package configs

type ServerEnv string

const (
	ServerEnvLocalhost   ServerEnv = "local"
	ServerEnvDevelopment ServerEnv = "dev"
	ServerEnvStaging     ServerEnv = "staging"
	ServerEnvProduction  ServerEnv = "prod"
)

type Config struct {
	Server      *ServerConfig      `mapstructure:"SERVER"`
	MySQL       *MySQLConfig       `mapstructure:"MYSQL"`
	GCP         *GCPConfig         `mapstructure:"GCP"`
	PubSubTopic *PubSubTopicConfig `mapstructure:"PUB_SUB_TOPIC"`
	APIKey      *APIKeyConfig      `mapstructure:"API_KEY"`
}

type ServerConfig struct {
	Name       string    `mapstructure:"NAME"`
	Env        ServerEnv `mapstructure:"ENV"`
	ConfigFile string    `mapstructure:"CONFIG_FILE"`
	Host       string    `mapstructure:"HOST"`
	Port       int       `mapstructure:"PORT"`
	Timezone   string    `mapstructure:"TZ"`
}

type MySQLConfig struct {
	DBUsername        string `mapstructure:"DB_USERNAME"`
	DBPassword        string `mapstructure:"DB_PASSWORD"`
	DBHost            string `mapstructure:"DB_HOST"`
	DBPort            string `mapstructure:"DB_PORT"`
	DBName            string `mapstructure:"DB_NAME"`
	DBSocketDir       string `mapstructure:"DB_SOCKET_DIR"`
	DBConnMaxLifetime int64  `mapstructure:"DB_CONN_MAX_LIFETIME"`
	DBMaxIdleConns    int    `mapstructure:"DB_MAX_IDLE_CONNS"`
	DBMaxOpenConns    int    `mapstructure:"DB_MAX_OPEN_CONNS"`
}

type GCPConfig struct {
	ProjectID string `mapstructure:"PROJECT_ID"`
}

type PubSubTopicConfig struct {
	NotificationTopic string `mapstructure:"NOTIFICATION_TOPIC"`
}

type APIKeyConfig struct {
	NotificationAPIKey string `mapstructure:"NOTIFICATION_API_KEY"`
	PromotionAPIKey    string `mapstructure:"PROMOTION_API_KEY"`
	PortfolioAPIKEY    string `mapstructure:"PORTFOLIO_API_KEY"`
	StockOrderAPIKEY   string `mapstructure:"STOCK_ORDER_API_KEY"`
}

func (cfg *Config) GetServerEnv() ServerEnv {
	if cfg.Server.Env == "" {
		return DefaultServerEnv
	}
	return cfg.Server.Env
}
