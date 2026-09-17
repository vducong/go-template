package httpif

import (
	"context"
	"gotemplate/internal/cfg"
	"gotemplate/internal/infra"
	"gotemplate/internal/interface/httpif/handler"
	"gotemplate/internal/interface/httpif/middleware"
	"gotemplate/pkg/auth"
	"gotemplate/pkg/httprespwrit"
	"gotemplate/pkg/lg"
	"gotemplate/pkg/mtr"
	"gotemplate/pkg/trc"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"go.opentelemetry.io/otel"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

func TestRouter_TracesAndLogsOnlyAPIRoutes(t *testing.T) {
	recorder := tracetest.NewSpanRecorder()
	prev := otel.GetTracerProvider()
	otel.SetTracerProvider(sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder)))
	t.Cleanup(func() { otel.SetTracerProvider(prev) })

	metric, err := mtr.Setup(context.Background(), &mtr.Config{Kind: mtr.KindPrometheus})
	if err != nil {
		t.Fatal(err)
	}
	logger := &recordingLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	writer := httprespwrit.NewWriter(lg.NewRespWritLogger(nil, logger))
	infrastructure := &infra.Infrastructure{
		Logger:         logger,
		ResponseWriter: writer,
		Tracing:        &trc.TracerProvider{},
		Metric:         metric,
	}
	configs := &cfg.Config{}
	jwtConfig, err := auth.JWTConfigFromAuth("test-secret", 0, "", nil)
	if err != nil {
		t.Fatal(err)
	}
	authenticator, err := middleware.NewAuthenticator(logger, writer, jwtConfig, nil)
	if err != nil {
		t.Fatal(err)
	}
	router := newRouter(configs, infrastructure, handler.Setup(configs, infrastructure, nil), authenticator)

	tests := []struct {
		path      string
		wantSpans int
		wantLogs  int
	}{
		{path: "/metrics", wantSpans: 0, wantLogs: 0},
		{path: "/health", wantSpans: 0, wantLogs: 0},
		{path: "/api/v1/nope", wantSpans: 1, wantLogs: 1}, // no such route, still traced and logged
	}
	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			spansBefore, logsBefore := len(recorder.Ended()), logger.count()
			router.ServeHTTP(httptest.NewRecorder(), httptest.NewRequest(http.MethodGet, tt.path, nil))
			if got := len(recorder.Ended()) - spansBefore; got != tt.wantSpans {
				t.Errorf("spans = %d, want %d", got, tt.wantSpans)
			}
			if got := logger.count() - logsBefore; got != tt.wantLogs {
				t.Errorf("request log lines = %d, want %d", got, tt.wantLogs)
			}
		})
	}
}

// recordingLogger counts the request log lines RequestLogger writes through CtxLog.
type recordingLogger struct {
	lg.Logger
	mu sync.Mutex
	n  int
}

func (l *recordingLogger) CtxLog(context.Context, lg.LogLevel, string, ...lg.Field) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.n++
}

func (l *recordingLogger) count() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.n
}
