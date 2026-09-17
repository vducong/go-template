package mtr

import (
	"context"
	"strings"
	"testing"

	"go.opentelemetry.io/otel"
)

func TestSetupOtel_PrometheusReaderWritesToServedRegistry(t *testing.T) {
	ctx := context.Background()
	mp, err := Setup(ctx, &Config{
		ServiceName: "test",
		Kind:        KindOtel,
		Readers:     []*ReaderConfig{{ExporterKind: ExporterKindOtelPrometheus}},
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mp.StopFn(context.Background()) })

	counter, err := otel.Meter("test").Int64Counter("probe_requests")
	if err != nil {
		t.Fatal(err)
	}
	counter.Add(ctx, 7)

	if !servesFamily(t, mp, "probe_requests") {
		t.Fatal("an OpenTelemetry counter must appear in the registry that /metrics serves")
	}
}

func TestSetupOtel_ServesRuntimeMetricsWithoutPrometheusReader(t *testing.T) {
	mp, err := Setup(context.Background(), &Config{ServiceName: "test", Kind: KindOtel})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = mp.StopFn(context.Background()) })

	if mp.PromClient == nil {
		t.Fatal("/metrics is only mounted when PromClient is set")
	}
	if !servesFamily(t, mp, "go_goroutines") {
		t.Fatal("runtime metrics must be served")
	}
}

func servesFamily(t *testing.T, mp *MeterProvider, prefix string) bool {
	t.Helper()
	families, err := mp.PromClient.Registry.Gather()
	if err != nil {
		t.Fatal(err)
	}
	for _, f := range families {
		if strings.HasPrefix(f.GetName(), prefix) {
			return true
		}
	}
	return false
}
