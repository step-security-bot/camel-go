package endpoint

import (
	"github.com/pkg/errors"
	"net/url"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/core/processors"
)

const TAG = "endpoint"

func init() {
	processors.Types[TAG] = func() interface{} {
		return &Endpoint{}
	}
}

type Endpoint struct {
	api.Identifiable
	api.WithOutputs

	Identity   string                 `yaml:"id"`
	URI        string                 `yaml:"uri"`
	Parameters map[string]interface{} `yaml:"parameters,omitempty"`
}

func (e *Endpoint) ID() string {
	return e.Identity
}

func (e *Endpoint) Consumer(ctx api.Context) (api.Consumer, error) {

	ep, err := e.create(ctx)
	if err != nil {
		return nil, errors.Wrapf(err, "failure creating endpoint")
	}

	cf, ok := ep.(api.ConsumerFactory)
	if !ok {
		return nil, camelerrors.NotImplementedf("scheme %s does not implement consumer", ep.Component().Scheme())
	}

	consumer, err := cf.Consumer()
	if err != nil {
		return nil, errors.Wrapf(err, "error creating consumer")
	}

	for _, o := range e.Outputs() {
		next := o

		consumer.Next(next)
	}
	return consumer, nil
}

func (e *Endpoint) create(ctx api.Context) (api.Endpoint, error) {
	params := make(map[string]interface{})

	u, err := url.Parse(e.URI)
	if err != nil {
		return nil, err
	}

	for k, v := range u.Query() {
		params[k] = v
	}
	for k, v := range e.Parameters {
		params[k] = v
	}

	f, ok := components.Factories[u.Scheme]
	if !ok {
		return nil, camelerrors.NotFoundf("not component for scheme %s", u.Scheme)
	}

	c, err := f(ctx, params)
	if err != nil {
		return nil, err
	}

	return c.Endpoint(params)
}
