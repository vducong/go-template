//nolint:revive
package cfg

import (
	"encoding/json"
	"fmt"
	"gotemplate/pkg/mtr"
	"time"
)

type Config struct {
	App        AppConfig        `yaml:"app"`
	HTTP       HTTPConfig       `yaml:"http"`
	GRPC       GRPCConfig       `yaml:"grpc"`
	Log        LogConfig        `yaml:"log"`
	HTTPClient HTTPClientConfig `yaml:"http_client"`
	Tracing    TracingConfig    `yaml:"tracing"`
	Metric     MetricConfig     `yaml:"metric"`
}

type AppConfig struct {
	Name            string        `yaml:"name" env:"APP_NAME" env-default:"gotemplate"`
	Version         string        `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	Env             string        `yaml:"env" env:"GO_ENV" env-default:"development"`
	MaxProcs        int           `yaml:"max_procs" env:"MAX_PROCS" env-default:"300"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"APP_SHUTDOWN_TIMEOUT" env-default:"5s"`
}

type HTTPConfig struct {
	Port              string        `yaml:"port" env:"HTTP_PORT" env-default:"8080"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT" env-default:"5s"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
}

type GRPCConfig struct {
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"9090"`
}

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Mode  string `yaml:"mode" env:"LOG_MODE" env-default:"console"`
}

type HTTPClientConfig struct {
	MaxIdleConns        int           `yaml:"max_idle_conns" env:"HTTP_CLIENT_MAX_IDLE_CONNS" env-default:"100"`
	MaxIdleConnsPerHost int           `yaml:"max_idle_conns_per_host" env:"HTTP_CLIENT_MAX_IDLE_CONNS_PER_HOST" env-default:"10"`
	IdleConnTimeout     time.Duration `yaml:"idle_conn_timeout" env:"HTTP_CLIENT_IDLE_CONN_TIMEOUT" env-default:"20s"`
	Timeout             time.Duration `yaml:"timeout" env:"HTTP_CLIENT_TIMEOUT" env-default:"10s"`
}

type TracingConfig struct {
	ServiceName string                `yaml:"service_name" env:"TRACING_SERVICE_NAME" env-default:"gotemplate"`
	Kind        string                `yaml:"kind" env:"TRACING_KIND"`
	Exporter    TracingExporterConfig `yaml:"exporter" env:"TRACING_EXPORTER"`
	Provider    TracingProviderConfig `yaml:"provider" env:"TRACING_PROVIDER"`
}

type TracingExporterConfig struct {
	Kind     string `yaml:"kind" env:"TRACING_EXPORTER_KIND" env-default:"zipkin"`
	Endpoint string `yaml:"endpoint" env:"TRACING_EXPORTER_ENDPOINT" env-default:"localhost:4317"`
}

type TracingProviderConfig struct {
	SamplerKind string  `yaml:"sampler_kind" env:"TRACING_PROVIDER_SAMPLER_KIND" env-default:"always"`
	SampleRate  float64 `yaml:"sample_rate" env:"TRACING_PROVIDER_SAMPLE_RATE" env-default:"1.0"`
}

type MetricConfig struct {
	ServiceName string        `yaml:"service_name" env:"METRIC_SERVICE_NAME" env-default:"gotemplate"`
	Kind        string        `yaml:"kind" env:"METRIC_KIND"`
	Readers     ReaderConfigs `yaml:"readers" env:"METRIC_READERS"`
}

type ReaderConfigs []*ReaderConfig

type ReaderConfig struct {
	ExporterKind     string `yaml:"exporter_kind" json:"exporter_kind"`
	ExporterEndpoint string `yaml:"exporter_endpoint" json:"exporter_endpoint"`
	IntervalSec      int    `yaml:"interval" json:"interval"`
}

// SetValue implements cleanenv.Setter for ReaderConfigs.
// It expects a JSON array of reader objects, for example:
//
//	[
//	  {
//	    "exporter_kind": "otel_grpc",
//	    "exporter_endpoint": "otel-collector:4317",
//	    "interval": "10s"
//	  }
//	]
func (c *ReaderConfigs) SetValue(s string) error {
	if len(s) == 0 {
		return nil
	}

	var readers []*ReaderConfig
	if err := json.Unmarshal([]byte(s), &readers); err != nil {
		return fmt.Errorf("parse METRIC_READERS as JSON: %w", err)
	}

	*c = readers
	return nil
}

func (c *ReaderConfigs) ToMtrReaderConfigs() []*mtr.ReaderConfig {
	readers := make([]*mtr.ReaderConfig, len(*c))
	for i, reader := range *c {
		readers[i] = reader.ToMtrReaderConfig()
	}
	return readers
}

func (c *ReaderConfig) ToMtrReaderConfig() *mtr.ReaderConfig {
	return &mtr.ReaderConfig{
		ExporterKind:     mtr.ExporterKind(c.ExporterKind),
		ExporterEndpoint: c.ExporterEndpoint,
		Interval:         time.Duration(c.IntervalSec) * time.Second,
	}
}
