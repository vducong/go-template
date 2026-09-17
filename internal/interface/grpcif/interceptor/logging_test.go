package interceptor

import (
	"context"
	"gotemplate/pkg/lg"
	"sync"
	"testing"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
)

// An unsampled trace is never exported, so its trace_id would point at a trace that does not exist.
func TestLoggingInterceptor_TraceIDOnlyForSampledTraces(t *testing.T) {
	for _, sampled := range []bool{true, false} {
		logger := &fieldsLogger{Logger: lg.New(&lg.Config{Level: "panic", Mode: "json"})}
		interceptor := NewLoggingInterceptor(logger)
		info := &grpc.UnaryServerInfo{FullMethod: "/svc/Method"}
		handler := func(context.Context, any) (any, error) { return nil, nil }
		if _, err := interceptor(spanContext(sampled), nil, info, handler); err != nil {
			t.Fatal(err)
		}

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

// fieldsLogger records the fields of each line the interceptor writes through CtxInfo.
type fieldsLogger struct {
	lg.Logger
	mu    sync.Mutex
	lines []map[string]any
}

func (l *fieldsLogger) CtxInfo(_ context.Context, _ string, fields ...lg.Field) {
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
		t.Fatalf("got %d log lines, want 1", len(l.lines))
	}
	return l.lines[0]
}
