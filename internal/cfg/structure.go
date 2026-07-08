//nolint:revive
package cfg

import (
	"encoding/json"
	"fmt"
	"gotemplate/internal/constant"
	"gotemplate/pkg/mtr"
	"time"
)

type Config struct {
	App     AppConfig     `yaml:"app"`
	HTTP    HTTPConfig    `yaml:"http"`
	GRPC    GRPCConfig    `yaml:"grpc"`
	Auth    AuthConfig    `yaml:"auth"`
	Log     LogConfig     `yaml:"log"`
	Tracing TracingConfig `yaml:"tracing"`
	Metric  MetricConfig  `yaml:"metric"`
}

type AppConfig struct {
	Name             string        `yaml:"name" env:"APP_NAME" env-default:"gotemplate"`
	Version          string        `yaml:"version" env:"APP_VERSION" env-default:"1.0.0"`
	Env              constant.Env  `yaml:"env" env:"GO_ENV" env-default:"development"`
	ServiceCode      int           `yaml:"service_code" env:"APP_SERVICE_CODE" env-default:"12"`
	MessageKeyPrefix string        `yaml:"message_key_prefix" env:"APP_MESSAGE_KEY_PREFIX" env-default:""`
	MaxProcs         int           `yaml:"max_procs" env:"MAX_PROCS" env-default:"300"`
	ShutdownTimeout  time.Duration `yaml:"shutdown_timeout" env:"APP_SHUTDOWN_TIMEOUT" env-default:"5s"`
}

type HTTPConfig struct {
	Port              string        `yaml:"port" env:"HTTP_PORT" env-default:"3000"`
	ReadTimeout       time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"5s"`
	ReadHeaderTimeout time.Duration `yaml:"read_header_timeout" env:"HTTP_READ_HEADER_TIMEOUT" env-default:"5s"`
	WriteTimeout      time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	IdleTimeout       time.Duration `yaml:"idle_timeout" env:"HTTP_IDLE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout   time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"5s"`
	CORS              CORSConfig    `yaml:"cors"`
}

type CORSConfig struct {
	Enabled          bool     `yaml:"enabled" env:"HTTP_CORS_ENABLED" env-default:"false"`
	AllowedOrigins   []string `yaml:"allowed_origins" env:"HTTP_CORS_ALLOWED_ORIGINS" env-default:""`
	AllowedMethods   []string `yaml:"allowed_methods" env:"HTTP_CORS_ALLOWED_METHODS" env-default:""`
	AllowedHeaders   []string `yaml:"allowed_headers" env:"HTTP_CORS_ALLOWED_HEADERS" env-default:""`
	ExposedHeaders   []string `yaml:"exposed_headers" env:"HTTP_CORS_EXPOSED_HEADERS" env-default:""`
	AllowCredentials bool     `yaml:"allow_credentials" env:"HTTP_CORS_ALLOW_CREDENTIALS" env-default:"false"`
	MaxAge           int      `yaml:"max_age" env:"HTTP_CORS_MAX_AGE" env-default:"300"`
}

type AuthConfig struct {
	JWTSecret   string        `yaml:"jwt_secret" env:"AUTH_JWT_SECRET" env-default:""`
	JWTLeeway   time.Duration `yaml:"jwt_leeway" env:"AUTH_JWT_LEEWAY" env-default:"30s"`
	JWTIssuer   string        `yaml:"jwt_issuer" env:"AUTH_JWT_ISSUER" env-default:""`
	JWTAudience []string      `yaml:"jwt_audience" env:"AUTH_JWT_AUDIENCE" env-default:""`
	APIKeys     []string      `yaml:"api_keys" env:"AUTH_API_KEYS" env-default:""`
}

type GRPCConfig struct {
	Port string `yaml:"port" env:"GRPC_PORT" env-default:"2000"`
}

type LogConfig struct {
	Level string `yaml:"level" env:"LOG_LEVEL" env-default:"info"`
	Mode  string `yaml:"mode" env:"LOG_MODE" env-default:"console"`
}

type TracingConfig struct {
	ServiceName string                `yaml:"service_name" env:"TRACING_SERVICE_NAME" env-default:"gotemplate"`
	Kind        string                `yaml:"kind" env:"TRACING_KIND" env-default:"otel"`
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
	Kind        string        `yaml:"kind" env:"METRIC_KIND" env-default:"otel"`
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
