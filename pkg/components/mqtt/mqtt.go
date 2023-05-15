////go:build components_mqtt || components_all

package mqtt

import (
	"github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/components"
)

const (
	Scheme = "mqtt"

	AttributeMqttMessageID   = "camel.apache.org/mqtt.message.id"
	AttributeMqttMessageQUOS = "camel.apache.org/mqtt.message.qus"
)

func NewComponent(ctx api.Context, _ map[string]interface{}) (api.Component, error) {
	component := Component{
		DefaultComponent: components.NewDefaultComponent(ctx, Scheme),
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

	if _, err := c.Context().TypeConverter().Convert(&config, &e.config); err != nil {
		return nil, err
	}

	return &e, nil
}
