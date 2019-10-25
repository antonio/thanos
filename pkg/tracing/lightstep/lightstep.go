package lightstep

import (
	"context"
	"errors"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
	"gopkg.in/yaml.v2"
)

// Config - YAML configuration.
type Config struct {
	// AccessToken is the unique API key for your LightStep project. It is
	// available on your account page at https://app.lightstep.com/account.
	AccessToken string `yaml:"access_token"`

	// Collector is the host, port, and plaintext option to use
	// for the collector.
	Collector lightstep.Endpoint `yaml:"collector"`
}

// Tracer wraps the Lightstep tracer and the context.
type Tracer struct {
	lightstep.Tracer
	ctx context.Context
}

// Close synchronously flushes the Lightstep tracer, then terminates it.
func (t *Tracer) Close() error {
	t.Tracer.Close(t.ctx)

	return nil
}

type LoggingRecorder struct {
	r      lightstep.SpanRecorder
	Logger log.Logger
}

func (lr *LoggingRecorder) RecordSpan(s lightstep.RawSpan) {
	if s.Operation == "promql_instant_query" {
		level.Info(lr.Logger).Log("operation", s.Operation, "query", s.Tags["query"], "duration", s.Duration.Seconds())
	} else if s.Operation == "promql_range_query" {
		level.Info(lr.Logger).Log("operation", s.Operation, "query", s.Tags["query"],
			"range", s.Tags["range"], "step", s.Tags["step"], "duration", s.Duration.Seconds())
	}
}

// NewTracer creates a Tracer with the options present in the YAML config.
func NewTracer(ctx context.Context, logger log.Logger, yamlConfig []byte) (opentracing.Tracer, io.Closer, error) {
	config := Config{}
	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return nil, nil, err
	}

	recorder := &LoggingRecorder{
		Logger: logger,
	}
	options := lightstep.Options{
		AccessToken: config.AccessToken,
		Collector:   config.Collector,
		Recorder:    recorder,
	}
	lighstepTracer := lightstep.NewTracer(options)
	if lighstepTracer == nil { // lightstep.NewTracer returns nil when there is an error
		return nil, nil, errors.New("error creating Lightstep tracer")
	}

	t := &Tracer{
		lighstepTracer,
		ctx,
	}
	return t, t, nil
}
