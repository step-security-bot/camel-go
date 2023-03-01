// //go:build components_wasm || components_all
package wasm

import (
	"context"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
	camelerrors "github.com/lburgazzoli/camel-go/pkg/core/errors"
	"github.com/lburgazzoli/camel-go/pkg/util/uuid"
)

type Endpoint struct {
	config Config
	components.DefaultEndpoint
}

func (e *Endpoint) Start(context.Context) error {
	return nil
}

func (e *Endpoint) Stop(context.Context) error {
	return nil
}

func (e *Endpoint) Producer() (api.Producer, error) {
	if e.config.Remaining == "" {
		return nil, camelerrors.MissingParameterf("path", "failure processing %s", Scheme)
	}

	c := Producer{
		id:       uuid.New(),
		endpoint: e,
	}

	return &c, nil
}