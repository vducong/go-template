package lg

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
)

var (
	zerologLevelMapping = map[LogLevel]zerolog.Level{
		LogLevelDebug: zerolog.DebugLevel,
		LogLevelInfo:  zerolog.InfoLevel,
		LogLevelWarn:  zerolog.WarnLevel,
		LogLevelError: zerolog.ErrorLevel,
		LogLevelFatal: zerolog.FatalLevel,
		LogLevelPanic: zerolog.PanicLevel,
	}

	zerologLevelMappingReverse = map[zerolog.Level]LogLevel{
		zerolog.DebugLevel: LogLevelDebug,
		zerolog.InfoLevel:  LogLevelInfo,
		zerolog.WarnLevel:  LogLevelWarn,
		zerolog.ErrorLevel: LogLevelError,
		zerolog.FatalLevel: LogLevelFatal,
		zerolog.PanicLevel: LogLevelPanic,
	}

	zerologLevelMappingReverseStr = map[zerolog.Level]string{
		zerolog.DebugLevel: LogLevelDebug.String(),
		zerolog.InfoLevel:  LogLevelInfo.String(),
		zerolog.WarnLevel:  LogLevelWarn.String(),
		zerolog.ErrorLevel: LogLevelError.String(),
		zerolog.FatalLevel: LogLevelFatal.String(),
		zerolog.PanicLevel: LogLevelPanic.String(),
	}
)

type zerologLogger struct {
	logger *zerolog.Logger
}

func newZerologLogger(cfg *CreateLoggerDTO) Logger {
	zerolog.LevelFieldName = "level"
	zerolog.LevelFieldMarshalFunc = func(level zerolog.Level) string {
		return strings.ToUpper(level.String())
	}
	zerolog.FormattedLevels = zerologLevelMappingReverseStr

	zerolog.MessageFieldName = messageFieldName

	zerolog.ErrorFieldName = ErrorFieldName
	zerolog.ErrorStackFieldName = StackFieldName
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	zerolog.TimestampFieldName = "time"
	zerolog.TimeFieldFormat = time.RFC3339
	zerolog.TimestampFunc = func() time.Time {
		return time.Now()
	}

	zerolog.CallerFieldName = callerFieldName
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		dir := path.Dir(file)
		parent := path.Base(dir)
		filename := path.Base(file)

		funcName := runtime.FuncForPC(pc).Name()
		if i := strings.LastIndex(funcName, "/"); i >= 0 {
			funcName = funcName[i+1:]
		}
		return fmt.Sprintf("%s:%d > %s", path.Join(parent, filename), line, funcName)
	}

	var l zerolog.Logger
	if cfg.Mode == "console" {
		writer := zerolog.NewConsoleWriter()
		writer.TimeFormat = zerolog.TimeFieldFormat
		l = zerolog.New(writer).With().CallerWithSkipFrameCount(3).Timestamp().Logger()
	} else {
		l = zerolog.New(os.Stderr).With().CallerWithSkipFrameCount(3).Timestamp().Logger()
	}

	logLevel, ok := zerologLevelMapping[cfg.Level]
	if !ok {
		l.Error().Msgf("invalid log level: %s", cfg.Level)
		logLevel = zerolog.InfoLevel
	}
	l = l.Level(logLevel)

	return &zerologLogger{logger: &l}
}

func (l *zerologLogger) Level() LogLevel {
	return zerologLevelMappingReverse[l.logger.GetLevel()]
}

func applyField(e *zerolog.Event, f Field) {
	switch f.Value.(type) {
	case string:
		e.Str(f.Key, f.Value.(string))
	case int:
		e.Int(f.Key, f.Value.(int))
	case float64:
		e.Float64(f.Key, f.Value.(float64))
	case bool:
		e.Bool(f.Key, f.Value.(bool))
	case time.Duration:
		e.Dur(f.Key, f.Value.(time.Duration))
	case error:
		if f.Key == ErrorFieldName {
			e.Stack().Err(f.Value.(error))
		} else {
			e.AnErr(f.Key, f.Value.(error))
		}
	default:
		e.Any(f.Key, f.Value)
	}
}
