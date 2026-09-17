package mtr

import (
	"context"
	"testing"
)

// infra.Close calls StopFn for every metric kind, so a nil StopFn panics at shutdown.
func TestSetup_PrometheusKindHasStopFn(t *testing.T) {
	mp, err := Setup(context.Background(), &Config{Kind: KindPrometheus})
	if err != nil {
		t.Fatal(err)
	}
	if mp.StopFn == nil {
		t.Fatal("StopFn must be set for the prometheus kind")
	}
	if err := mp.StopFn(context.Background()); err != nil {
		t.Fatalf("StopFn returned %v", err)
	}
}
