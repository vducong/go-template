package lg

type Config struct {
	Level string
	Mode  string
}

func (c *Config) ToCreateLoggerDTO() *CreateLoggerDTO {
	level, ok := logLevelValueMapping[c.Level]
	if !ok {
		level = LogLevelInfo
	}

	mode, ok := logModeValueMapping[c.Mode]
	if !ok {
		mode = LogModeConsole
	}

	return &CreateLoggerDTO{
		Level: level,
		Mode:  mode,
	}
}

type CreateLoggerDTO struct {
	Level LogLevel
	Mode  LogMode
}

func New(cfg *Config) Logger {
	return newZerologLogger(cfg.ToCreateLoggerDTO())
}
