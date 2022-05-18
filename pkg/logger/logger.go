package logger

import (
	"os"
	"promotion/configs"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	defaultLogLevel         = zapcore.InfoLevel
	minHighPriorityLogLevel = zapcore.ErrorLevel
)

type (
	Logger      = zap.SugaredLogger
	PlainLogger = zap.Logger
)

func New(cfg *configs.Config) *Logger {
	encoder := getEncoder()
	core := getCore(encoder)
	plain := zap.New(
		core,
		zap.AddCaller(),
		zap.ErrorOutput(zapcore.Lock(os.Stderr)),
	)

	sugared := plain.Sugar()
	return sugared
}

func getEncoder() zapcore.Encoder {
	encoderConfig := getEncoderConfig()
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getEncoderConfig() zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoderConfig.EncodeTime = zapcore.RFC3339TimeEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return encoderConfig
}

func getCore(encoder zapcore.Encoder) zapcore.Core {
	return zapcore.NewTee(
		coreWritingLowPriorityLog(encoder),
		coreWritingHighPriorityLog(encoder),
	)
}

func coreWritingLowPriorityLog(encoder zapcore.Encoder) zapcore.Core {
	stdout := zapcore.Lock(os.Stdout) // lock for concurrent safe
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(stdout),
		isLowPriorityLogLevel(),
	)
}

func isLowPriorityLogLevel() zapcore.LevelEnabler {
	// lowPriority used by info\debug\warn
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= defaultLogLevel && lvl < minHighPriorityLogLevel
	})
}

func coreWritingHighPriorityLog(encoder zapcore.Encoder) zapcore.Core {
	stderr := zapcore.Lock(os.Stderr)
	return zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(stderr),
		isHighPriorityLogLevel(),
	)
}

func isHighPriorityLogLevel() zapcore.LevelEnabler {
	// highPriority used by error\panic\fatal
	return zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= defaultLogLevel && lvl >= minHighPriorityLogLevel
	})
}
