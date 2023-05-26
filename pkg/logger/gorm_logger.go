package logger

import (
	standardLog "log"
	"os"
	"time"

	gormLogger "gorm.io/gorm/logger"
)

func InitGORMLogger() gormLogger.Interface {
	return gormLogger.New(
		standardLog.New(os.Stderr, "\r\n", standardLog.LstdFlags),
		gormLogger.Config{
			SlowThreshold:             time.Second,
			Colorful:                  true,
			IgnoreRecordNotFoundError: true,
			LogLevel:                  gormLogger.Warn,
		},
	)
}
