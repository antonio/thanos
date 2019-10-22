package lightstep

import (
	"context"
	"errors"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
	"gopkg.in/yaml.v2"
)

type Config struct {
	AccessToken string `yaml:"access_token"`
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
