package trc

import (
	"context"
	"testing"

	"go.opentelemetry.io/otel"
	otelres "go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	semconv "go.opentelemetry.io/otel/semconv/v1.38.0"
	"go.opentelemetry.io/otel/trace"
)

func TestInitTracerProvider_Sampling(t *testing.T) {
	ratioZero := ProviderConfig{SamplerKind: SamplerKindRatio, SampleRate: 0}
	tests := []struct {
		name          string
		provider      ProviderConfig
		parentSampled bool
		wantRecorded  bool
	}{
		{"sampled caller keeps the trace whole under ratio 0", ratioZero, true, true},
		{"not-sampled caller falls back to always", ProviderConfig{SamplerKind: SamplerKindAlways}, false, true},
		{"not-sampled caller falls back to ratio 0", ratioZero, false, false},
		{"never stays off under a sampled caller", ProviderConfig{SamplerKind: SamplerKindNever}, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			stop, err := initTracerProvider(
				&Config{Provider: tt.provider},
				otelres.Empty(),
				tracetest.NewInMemoryExporter(),
			)
			if err != nil {
				t.Fatal(err)
			}
			t.Cleanup(func() { _ = stop(context.Background()) })

			_, span := otel.Tracer("test").Start(remoteParent(tt.parentSampled), "child")
			defer span.End()
			if got := span.SpanContext().IsSampled(); got != tt.wantRecorded {
				t.Fatalf("sampled = %v, want %v", got, tt.wantRecorded)
			}
		})
	}
}

func TestNewResource_EnvironmentTag(t *testing.T) {
	res, err := newResource(&Config{ServiceName: "svc", Environment: "staging"})
	if err != nil {
		t.Fatal(err)
	}
	if got, ok := res.Set().Value(semconv.DeploymentEnvironmentNameKey); !ok || got.AsString() != "staging" {
		t.Fatalf("deployment.environment.name = %q (set %v), want staging", got.AsString(), ok)
	}

	res, err = newResource(&Config{ServiceName: "svc"})
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := res.Set().Value(semconv.DeploymentEnvironmentNameKey); ok {
		t.Fatal("an empty environment must not produce an empty tag")
	}
}

func remoteParent(sampled bool) context.Context {
	var flags trace.TraceFlags
	if sampled {
		flags = trace.FlagsSampled
	}
	return trace.ContextWithRemoteSpanContext(context.Background(), trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{0x4b, 0xf9},
		SpanID:     trace.SpanID{0x00, 0xf0},
		TraceFlags: flags,
		Remote:     true,
	}))
}
