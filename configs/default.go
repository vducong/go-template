package configs

import (
	"time"
)

func DefaultConfig() *Config {
	return &Config{
		Server: &ServerConfig{
			Name:       DefaultServerName,
			Env:        DefaultServerEnv,
			Host:       DefaultServerHost,
			Port:       DefaultServerPort,
			ConfigFile: DefaultServerConfigFile,
			Timezone:   DefaultServerTimezone,
		},
		MySQL: &MySQLConfig{
			DBUsername:        DefaultDBUsername,
			DBPassword:        DefaultDBPassword,
			DBHost:            DefaultDBHost,
			DBPort:            DefaultDBPort,
			DBName:            DefaultDBName,
			DBSocketDir:       DefaultDBSocketDir,
			DBMaxIdleConns:    DefaultDBMaxIdleConns,
			DBMaxOpenConns:    DefaultDBMaxOpenConns,
			DBConnMaxLifetime: DefaultDBConnMaxLifetime,
		},
		GCP: &GCPConfig{
			ProjectID: DefaultGoogleCloudProjectID,
		},
		PubSubTopic: &PubSubTopicConfig{
			NotificationTopic: DefaultNotificationTopic,
		},
		APIKey: &APIKeyConfig{
			NotificationAPIKey: DefaultNotificationAPIKey,
			PromotionAPIKey:    DefaultPromotionAPIKey,
			PortfolioAPIKEY:    DefaultPortfolioAPIKEY,
			StockOrderAPIKEY:   DefaultStockOrderAPIKEY,
		},
	}
}

const (
	DefaultServerName                             = ""
	DefaultServerInstanceConnectionName           = ""
	DefaultServerEnv                    ServerEnv = ServerEnvDevelopment
	DefaultServerHost                             = ""
	DefaultServerPort                             = 5005
	DefaultServerConfigFile                       = "configs/dev.config.yaml"
	DefaultServerTimezone                         = "Asia/Ho_Chi_Minh"
	DefaultDBUsername                             = ""
	DefaultDBPassword                             = ""
	DefaultDBHost                                 = ""
	DefaultDBPort                                 = ""
	DefaultDBName                                 = ""
	DefaultDBSocketDir                            = "/cloudsql"
	DefaultDBMaxIdleConns                         = 10
	DefaultDBMaxOpenConns                         = 100
	DefaultDBConnMaxLifetime                      = int64(time.Hour)
	DefaultGoogleCloudProjectID                   = ""
	DefaultNotificationTopic                      = ""
	DefaultNotificationAPIKey                     = ""
	DefaultPromotionAPIKey                        = ""
	DefaultPortfolioAPIKEY                        = ""
	DefaultStockOrderAPIKEY                       = ""
)
