package lg

import "fmt"

type PrometheusErrorLogger interface {
	Println(v ...any)
}

type prometheusErrorLogger struct {
	Logger
}

func NewPrometheusErrorLogger(cfg *Config, l Logger) PrometheusErrorLogger {
	promLogger := &prometheusErrorLogger{}
	if l == nil {
		promLogger.Logger = New(cfg)
	} else {
		promLogger.Logger = l
	}
	return promLogger
}

// prometheus/client_golang/prometheus/promhttp.HandlerOpts.ErrorLog
func (l *prometheusErrorLogger) Println(v ...any) {
	if l.Level() > LogLevelError {
		return
	}

	l.Error("failed to collect and serve metrics",
		Str("kind", "prometheus"),
		Str("error", fmt.Sprintln(v...)),
	)
}
