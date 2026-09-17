package middleware

import (
	"context"
	"gotemplate/pkg/lg"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/trace"
)

func TestRequestLogger_OmitsBodiesOnSuccess(t *testing.T) {
	logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	serve(t, logger, http.StatusOK, `{"password":"hunter2"}`, `{"balance":12345}`)

	fields := logger.only(t)
	for _, key := range []string{"http.request.body", "http.response.body"} {
		if _, ok := fields[key]; ok {
			t.Errorf("%s must not be logged for a 200: %v", key, fields[key])
		}
	}
	if fields["http.request.body.size"] != 22 || fields["http.response.body.size"] != 17 {
		t.Errorf("body sizes must still be logged: %v", fields)
	}
}

func TestRequestLogger_LogsBodiesOnFailure(t *testing.T) {
	logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	serve(t, logger, http.StatusBadRequest, `{"amount":-1}`, `{"error":"bad amount"}`)

	fields := logger.only(t)
	if fields["http.request.body"] != `{"amount":-1}` || fields["http.response.body"] != `{"error":"bad amount"}` {
		t.Errorf("bodies must be logged for a 400: %v", fields)
	}
}

func TestRequestLogger_CapsLoggedBodiesButNotTheRequest(t *testing.T) {
	logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
	big := strings.Repeat("x", 10<<10)
	got, rec := serve(t, logger, http.StatusInternalServerError, big, big)

	if len(got) != len(big) || rec.Body.Len() != len(big) {
		t.Fatalf("handler read %d and client got %d bytes, want %d each", len(got), rec.Body.Len(), len(big))
	}
	fields := logger.only(t)
	for _, key := range []string{"http.request.body", "http.response.body"} {
		if body, _ := fields[key].(string); len(body) != requestLogBodyLimit {
			t.Errorf("%s logged %d bytes, want %d", key, len(body), requestLogBodyLimit)
		}
	}
}

// An unsampled trace is never exported, so its trace_id would point at a trace that does not exist.
func TestRequestLogger_TraceIDOnlyForSampledTraces(t *testing.T) {
	for _, sampled := range []bool{true, false} {
		logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
		h := RequestLogger(logger)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
		req := httptest.NewRequest(http.MethodGet, "/api/v1/wallet", nil)
		h.ServeHTTP(httptest.NewRecorder(), req.WithContext(spanContext(sampled)))

		fields := logger.only(t)
		_, hasTraceID := fields["trace_id"]
		_, hasSpanID := fields["span_id"]
		if hasTraceID != sampled || hasSpanID != sampled {
			t.Errorf("sampled=%v: trace_id logged %v, span_id logged %v", sampled, hasTraceID, hasSpanID)
		}
	}
}

func spanContext(sampled bool) context.Context {
	var flags trace.TraceFlags
	if sampled {
		flags = trace.FlagsSampled
	}
	return trace.ContextWithSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x01},
		SpanID:     trace.SpanID{0x02},
		TraceFlags: flags,
	}))
}

// serve runs one request through RequestLogger and returns what the handler read from the body.
func serve(t *testing.T, logger lg.Logger, status int, reqBody, respBody string) (string, *httptest.ResponseRecorder) {
	t.Helper()
	var read []byte
	h := RequestLogger(logger)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		if read, err = io.ReadAll(r.Body); err != nil {
			t.Error(err)
		}
		w.WriteHeader(status)
		_, _ = io.WriteString(w, respBody)
	}))
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodPost, "/api/v1/wallet", strings.NewReader(reqBody)))
	return string(read), rec
}

// fieldsLogger records the fields of each request log line RequestLogger writes through CtxLog.
type fieldsLogger struct {
	lg.Logger
	mu    sync.Mutex
	lines []map[string]any
}

func (l *fieldsLogger) CtxLog(_ context.Context, _ lg.LogLevel, _ string, fields ...lg.Field) {
	line := make(map[string]any, len(fields))
	for _, f := range fields {
		line[f.Key] = f.Value
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	l.lines = append(l.lines, line)
}

func (l *fieldsLogger) only(t *testing.T) map[string]any {
	t.Helper()
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.lines) != 1 {
		t.Fatalf("got %d request log lines, want 1", len(l.lines))
	}
	return l.lines[0]
}
