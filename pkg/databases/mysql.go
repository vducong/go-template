package databases

import (
	"fmt"
	"promotion/configs"
	"promotion/pkg/logger"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MySQLDB = *gorm.DB

func NewMySQLDB(cfg *configs.Config) (MySQLDB, error) {
	db, err := gorm.Open(
		mysql.Open(getDBURI(cfg)),
		getGORMConfig(),
	)
	if err != nil {
		return nil, err
	}

	instance, err := db.DB()
	if err != nil {
		return nil, err
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	instance.SetMaxIdleConns(cfg.MySQL.DBMaxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	instance.SetMaxOpenConns(cfg.MySQL.DBMaxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	instance.SetConnMaxLifetime(time.Duration(cfg.MySQL.DBConnMaxLifetime))

	return db, nil
}

func getDBURI(cfg *configs.Config) string {
	// e.g. 'project:region:instance'
	if cfg.Server.InstanceConnectionName != "" {
		return fmt.Sprintf(
			"%s:%s@unix(/%s/%s)/%s?parseTime=true",
			cfg.MySQL.DBUsername, cfg.MySQL.DBPassword,
			cfg.MySQL.DBSocketDir, cfg.Server.InstanceConnectionName,
			cfg.MySQL.DBName,
		)
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		cfg.MySQL.DBUsername, cfg.MySQL.DBPassword,
		cfg.MySQL.DBHost, cfg.MySQL.DBPort,
		cfg.MySQL.DBName,
	)
}

func getGORMConfig() *gorm.Config {
	return &gorm.Config{
		Logger: logger.InitGORMLogger(),
	}
}
