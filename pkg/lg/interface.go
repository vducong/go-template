package lg

import (
	"context"
	"time"
)

type Field struct {
	Key   string
	Value any
}

type Logger interface {
	Level() LogLevel
	LogPrinter
	LogCtxPrinter
}

type LogPrinter interface {
	Log(level LogLevel, msg string, fields ...Field)
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
	Panic(msg string, fields ...Field)
}

type LogCtxPrinter interface {
	CtxLog(ctx context.Context, level LogLevel, msg string, fields ...Field)
	CtxDebug(ctx context.Context, msg string, fields ...Field)
	CtxInfo(ctx context.Context, msg string, fields ...Field)
	CtxWarn(ctx context.Context, msg string, fields ...Field)
	CtxError(ctx context.Context, msg string, fields ...Field)
	CtxFatal(ctx context.Context, msg string, fields ...Field)
	CtxPanic(ctx context.Context, msg string, fields ...Field)
}

func Str(key, val string) Field {
	return Field{Key: key, Value: val}
}

func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

func Float64(key string, val float64) Field {
	return Field{Key: key, Value: val}
}

func Bool(key string, val bool) Field {
	return Field{Key: key, Value: val}
}

func Dur(key string, val time.Duration) Field {
	return Field{Key: key, Value: val}
}

func Any(key string, val any) Field {
	return Field{Key: key, Value: val}
}

func Err(err error) Field {
	return Field{Key: ErrorFieldName, Value: err}
}
