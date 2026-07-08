package lg

import "context"

type RespWritLogger interface {
	Error(err error)
	CtxError(ctx context.Context, err error)
}

type respWritLogger struct {
	Logger
}

func NewRespWritLogger(cfg *Config, l Logger) RespWritLogger {
	if l == nil {
		l = New(cfg)
	}
	return &respWritLogger{
		Logger: l,
	}
}

func (l *respWritLogger) Error(err error) {
	l.Logger.Error("failed to handle request", Err(err))
}

func (l *respWritLogger) CtxError(ctx context.Context, err error) {
	l.Logger.CtxError(ctx, "failed to handle request", Err(err))
}
