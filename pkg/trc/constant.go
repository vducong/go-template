package trc

type Kind string

const (
	KindOtel Kind = "otel"
)

type ExporterKind string

const (
	ExporterKindStdout ExporterKind = "stdout"
	ExporterKindZipkin ExporterKind = "zipkin"
	ExporterKindGRPC   ExporterKind = "grpc"
	ExporterKindHTTP   ExporterKind = "http"
)

type SamplerKind string

const (
	SamplerKindAlways SamplerKind = "always"
	SamplerKindNever  SamplerKind = "never"
	SamplerKindRatio  SamplerKind = "ratio"
)
