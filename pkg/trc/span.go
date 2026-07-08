package trc

import (
	"context"
	"runtime"
	"runtime/debug"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

var modulePrefix string

func init() {
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Path != "" {
		modulePrefix = info.Main.Path + "/"
	}
}

// Start opens a span auto-named from the caller's function. Use exactly like:
//
//	ctx, span := trc.Start(ctx)
//	defer span.End()
//
// Span name and tracer name are derived from runtime.Caller — no strings to choose,
// no attribute boilerplate. Add span.SetAttributes(...) after Start when extra
// fields are useful.
func Start(ctx context.Context) (context.Context, trace.Span) {
	pc, _, _, ok := runtime.Caller(1)
	if !ok {
		return otel.Tracer("").Start(ctx, "unknown")
	}
	pkg, name := splitCallerName(runtime.FuncForPC(pc).Name())
	return otel.Tracer(pkg).Start(ctx, name)
}

func splitCallerName(full string) (pkg, name string) {
	full = strings.TrimPrefix(full, modulePrefix)
	slash := strings.LastIndex(full, "/")
	var leaf string
	if slash >= 0 {
		pkg = full[:slash]
		leaf = full[slash+1:]
	} else {
		leaf = full
	}
	leaf = strings.ReplaceAll(leaf, "(*", "")
	leaf = strings.ReplaceAll(leaf, ")", "")
	if pkg == "" {
		if before, _, ok := strings.Cut(leaf, "."); ok {
			pkg = before
		}
	}
	return pkg, leaf
}
