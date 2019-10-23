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

type Config struct {
	// AccessToken is the unique API key for your LightStep project.  It is
	// available on your account page at https://app.lightstep.com/account
	AccessToken string `yaml:"access_token"`

	// Collector is the host, port, and plaintext option to use
	// for the collector.
	Collector lightstep.Endpoint `yaml:"collector"`
}

type Tracer struct {
	opentracing.Tracer
	ctx context.Context
}

func (t *Tracer) Close() error {
	lightstepTracer := t.Tracer.(lightstep.Tracer)
	lightstepTracer.Close(t.ctx)

	return nil
}

func NewTracer(ctx context.Context, logger log.Logger, yamlConfig []byte) (opentracing.Tracer, io.Closer, error) {
	config := Config{}
	if err := yaml.Unmarshal(yamlConfig, &config); err != nil {
		return nil, nil, err
	}

	options := lightstep.Options{
		AccessToken: config.AccessToken,
		Collector:   config.Collector,
	}
	lighstepTracer := lightstep.NewTracer(options)
	if lighstepTracer == nil { // lightstep.NewTracer returns nil when there is an error
		return nil, nil, errors.New("error creating Lightstep tracer")
	}

	logHandler := func(event lightstep.Event) {
		switch event := event.(type) {
		case lightstep.EventStatusReport:
			level.Info(logger).Log("msg", event, "duration", event.Duration)
		case lightstep.ErrorEvent:
			level.Error(logger).Log("msg", event)
		default:
			level.Info(logger).Log("msg", event)
		}
	}

	lightstep.SetGlobalEventHandler(logHandler)

	t := &Tracer{
		lighstepTracer,
		ctx,
	}
	return t, t, nil
}
