package lg

func (l *zerologLogger) Log(level LogLevel, msg string, fields ...Field) {
	if l.Level() > level {
		return
	}

	e := l.logger.WithLevel(zerologLevelMapping[level])
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) Debug(msg string, fields ...Field) {
	if l.Level() > LogLevelDebug {
		return
	}

	e := l.logger.Debug()
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) Info(msg string, fields ...Field) {
	if l.Level() > LogLevelInfo {
		return
	}

	e := l.logger.Info()
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...Field) {
	if l.Level() > LogLevelWarn {
		return
	}

	e := l.logger.Warn()
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...Field) {
	if l.Level() > LogLevelError {
		return
	}

	e := l.logger.Error()
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) Fatal(msg string, fields ...Field) {
	if l.Level() > LogLevelFatal {
		return
	}

	e := l.logger.Fatal()
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}

func (l *zerologLogger) Panic(msg string, fields ...Field) {
	if l.Level() > LogLevelPanic {
		return
	}

	e := l.logger.Panic()
	for _, f := range fields {
		applyField(e, f)
	}
	e.Msg(msg)
}
