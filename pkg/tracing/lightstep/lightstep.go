package lightstep

import (
	"context"
	"errors"
	"io"

	"github.com/go-kit/kit/log"
	"github.com/lightstep/lightstep-tracer-go"
	"github.com/opentracing/opentracing-go"
)

type Config struct {
	AccessToken string `yaml:"access_token"`
}

type Tracer struct {
	opentracing.Tracer
}

func (t *Tracer) Close() {
}

func NewTracer(ctx context.Context, logger log.Logger, conContentYaml []byte) (opentracing.Tracer, io.Closer, error) {
	options := lightstep.Options{
		AccessToken: "PLACEHOLDER", // TODO: Fix this
	}
	tracer := lightstep.NewTracer(options)
	if tracer == nil { // lightstep.NewTracer returns nil when there is an error
		return nil, nil, errors.New("error creating Lightstep tracer")
	}
	return tracer, tracer, nil
}
