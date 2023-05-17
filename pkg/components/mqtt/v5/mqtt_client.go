package v5

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/eclipse/paho.golang/paho"
	"github.com/pkg/errors"
)

type Client struct {
	logger *zap.SugaredLogger
	cfg    *Config
	client *paho.Client
}

func (c *Client) Start(ctx context.Context) error {
	cp := &paho.Connect{
		ClientID:   c.cfg.ClientID,
		CleanStart: true,
		Username:   c.cfg.Username,
		Password:   []byte(c.cfg.Password),
		KeepAlive:  30,
	}

	if c.cfg.Keepalive != nil {
		cp.KeepAlive = *c.cfg.Keepalive
	}
	if c.cfg.Username != "" {
		cp.UsernameFlag = true
	}
	if c.cfg.Password != "" {
		cp.PasswordFlag = true
	}

	ca, err := c.client.Connect(ctx, cp)
	if err != nil {
		return errors.Wrapf(err, "failed to connect to %s", c.cfg.Broker)
	}
	if ca.ReasonCode != 0 {
		return fmt.Errorf("failed to connect to %s : %d - %s", c.cfg.Broker, ca.ReasonCode, ca.Properties.ReasonString)
	}

	return nil
}

func (c *Client) Stop(_ context.Context) error {
	if c.client != nil {
		d := &paho.Disconnect{ReasonCode: 0}

		err := c.client.Disconnect(d)
		if err != nil {
			return errors.Wrap(err, "failed to send disconnect")
		}
	}

	return nil
}

func (c *Client) Subscribe(ctx context.Context, topic string) error {
	c.logger.Infof("subscribing to: %s", topic)

	_, err := c.client.Subscribe(ctx, &paho.Subscribe{
		Subscriptions: map[string]paho.SubscribeOptions{
			topic: {
				QoS:     c.cfg.QoS,
				NoLocal: true,
			},
		},
	})

	return err
}

func (c *Client) Publish(ctx context.Context, p *paho.Publish) error {
	_, err := c.client.Publish(ctx, p)

	return err
}
