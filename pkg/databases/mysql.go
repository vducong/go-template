package databases

import (
	"database/sql"
	"fmt"
	"os"
	"promotion/configs"
	"promotion/pkg/failure"
	"promotion/pkg/logger"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
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
		return nil, fmt.Errorf("Failed to init MySQL session: %w", err)
	}

	if err := db.Use(
		otelgorm.NewPlugin(otelgorm.WithDBName(cfg.MySQL.DBName)),
	); err != nil {
		return nil, failure.ErrWithTrace(fmt.Errorf("Failed to instrument trace: %w", err))
	}

	instance, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("Failed to get MySQL instance: %w", err)
	}
	setupDBInstance(instance, cfg)

	return db, nil
}

func getDBURI(cfg *configs.Config) string {
	// e.g. 'project:region:instance'
	instanceConnectionName := os.Getenv("INSTANCE_CONNECTION_NAME")
	if instanceConnectionName != "" {
		return fmt.Sprintf(
			"%s:%s@unix(/%s/%s)/%s?parseTime=true",
			cfg.MySQL.DBUsername, cfg.MySQL.DBPassword,
			cfg.MySQL.DBSocketDir, instanceConnectionName,
			cfg.MySQL.DBName,
		)
	}
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
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

func setupDBInstance(instance *sql.DB, cfg *configs.Config) {
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	instance.SetMaxIdleConns(cfg.MySQL.DBMaxIdleConns)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	instance.SetMaxOpenConns(cfg.MySQL.DBMaxOpenConns)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	instance.SetConnMaxLifetime(time.Duration(cfg.MySQL.DBConnMaxLifetime))
}
