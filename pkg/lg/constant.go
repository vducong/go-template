package lg

const (
	messageFieldName = "message"
	callerFieldName  = "caller"
	ErrorFieldName   = "error"
	StackFieldName   = "stack"
)

var (
	logLevelValueMapping = map[string]LogLevel{
		"debug": LogLevelDebug,
		"info":  LogLevelInfo,
		"warn":  LogLevelWarn,
		"error": LogLevelError,
		"fatal": LogLevelFatal,
		"panic": LogLevelPanic,
	}

	logLevelStrMapping = map[LogLevel]string{
		LogLevelDebug: LogLevelDebugStr,
		LogLevelInfo:  LogLevelInfoStr,
		LogLevelWarn:  LogLevelWarnStr,
		LogLevelError: LogLevelErrorStr,
		LogLevelFatal: LogLevelFatalStr,
		LogLevelPanic: LogLevelPanicStr,
	}

	logModeValueMapping = map[string]LogMode{
		"console": LogModeConsole,
		"json":    LogModeJSON,
	}
)

type LogLevel int8

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
	LogLevelFatal
	LogLevelPanic
)

type LogLevelStr string

const (
	LogLevelDebugStr = "DEBUG"
	LogLevelInfoStr  = "INFO"
	LogLevelWarnStr  = "WARN"
	LogLevelErrorStr = "ERROR"
	LogLevelFatalStr = "FATAL"
	LogLevelPanicStr = "PANIC"
)

func (l LogLevel) String() string {
	return logLevelStrMapping[l]
}

type LogMode string

const (
	LogModeConsole LogMode = "console"
	LogModeJSON    LogMode = "json"
)
