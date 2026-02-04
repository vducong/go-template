package mtr

type Kind string

const (
	KindPrometheus Kind = "prometheus"
	KindOtel       Kind = "otel"
)

type ExporterKind string

const (
	ExporterKindPrometheus     ExporterKind = "prometheus"
	ExporterKindOtelPrometheus ExporterKind = "otel_prometheus"
	ExporterKindOtelGRPC       ExporterKind = "otel_grpc"
	ExporterKindOtelHTTP       ExporterKind = "otel_http"
)
