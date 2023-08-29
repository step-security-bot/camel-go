// //go:build components_http || components_all

package http

import (
	"net/http"

	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	SchemeHTTP       = "http"
	SchemeHTTPS      = "https"
	PropertiesPrefix = "camel.component."

	AttributeStatusMessage = "camel.apache.org/http.status.message"
	AttributeStatusCode    = "camel.apache.org/http.status.code"
	AttributeProto         = "camel.apache.org/http.proto"
	AttributeContentLength = "camel.apache.org/http.content-length"
)

func newComponent(ctx api.Context, scheme string, _ map[string]interface{}) (api.Component, error) {
	component := Component{
		DefaultComponent: components.NewDefaultComponent(ctx, scheme),
	}

	return &component, nil
}

type Component struct {
	components.DefaultComponent
}

func (c *Component) Endpoint(config api.Parameters) (api.Endpoint, error) {
	e := Endpoint{
		DefaultEndpoint: components.NewDefaultEndpoint(c),
	}

	props := c.Context().Properties().View(PropertiesPrefix + c.Scheme()).Parameters()
	for k, v := range config {
		props[k] = v
	}

	if _, err := c.Context().TypeConverter().Convert(&props, &e.config); err != nil {
		return nil, err
	}

	if e.config.Method == "" {
		e.config.Method = http.MethodGet
	}

	return &e, nil
}
