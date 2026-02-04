package lg

import (
	"context"
)

func (l *zerologLogger) CtxLog(ctx context.Context, level LogLevel, msg string, fields ...Field) {
	if l.Level() > level {
		return
	}

	e := l.logger.WithLevel(zerologLevelMapping[level]).Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) CtxDebug(ctx context.Context, msg string, fields ...Field) {
	if l.Level() > LogLevelDebug {
		return
	}

	e := l.logger.Debug().Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) CtxInfo(ctx context.Context, msg string, fields ...Field) {
	if l.Level() > LogLevelInfo {
		return
	}

	e := l.logger.Info().Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) CtxWarn(ctx context.Context, msg string, fields ...Field) {
	if l.Level() > LogLevelWarn {
		return
	}

	e := l.logger.Warn().Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) CtxError(ctx context.Context, msg string, fields ...Field) {
	if l.Level() > LogLevelError {
		return
	}

	e := l.logger.Error().Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) CtxFatal(ctx context.Context, msg string, fields ...Field) {
	if l.Level() > LogLevelFatal {
		return
	}

	e := l.logger.Fatal().Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) CtxPanic(ctx context.Context, msg string, fields ...Field) {
	if l.Level() > LogLevelPanic {
		return
	}

	e := l.logger.Panic().Ctx(ctx)
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}
