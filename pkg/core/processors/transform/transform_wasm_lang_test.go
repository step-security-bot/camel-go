// //go:build steps_transform || steps_all

package transform

import (
	"context"
	"testing"
	"time"

	"github.com/lburgazzoli/camel-go/pkg/core/language/wasm"

	"github.com/lburgazzoli/camel-go/pkg/core/language"

	"github.com/lburgazzoli/camel-go/pkg/core/processors"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/stretchr/testify/assert"
)

func TestTransformWASM(t *testing.T) {
	support.Run(t, "wasm_local", func(t *testing.T, ctx context.Context) {
		t.Helper()

		c := camel.GetContext(ctx)

		wg := make(chan camel.Message)

		v := support.NewChannelVerticle(wg)
		err := c.Spawn(v)

		assert.Nil(t, err)

		p := Transform{
			DefaultVerticle: processors.NewDefaultVerticle(),
			Language: language.Language{
				Wasm: &wasm.Wasm{
					Path: "../../../../etc/wasm/fn/simple_process.wasm",
				},
			}}

		p.Next(v.ID())

		id, err := p.Reify(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, id)

		msg, err := message.New()
		assert.Nil(t, err)

		assert.Nil(t, c.Send(id, msg))

		select {
		case msg := <-wg:
			c, ok := msg.Content().([]byte)
			assert.True(t, ok)
			assert.Equal(t, "hello from wasm", string(c))
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})

	support.Run(t, "wasm_registry", func(t *testing.T, ctx context.Context) {
		t.Helper()

		c := camel.GetContext(ctx)

		wg := make(chan camel.Message)

		v := support.NewChannelVerticle(wg)
		err := c.Spawn(v)

		assert.Nil(t, err)

		p := Transform{
			DefaultVerticle: processors.NewDefaultVerticle(),
			Language: language.Language{
				Wasm: &wasm.Wasm{
					Path:  "etc/wasm/fn/simple_process.wasm",
					Image: "docker.io/lburgazzoli/camel-go:latest",
				},
			}}

		p.Next(v.ID())

		id, err := p.Reify(ctx)
		assert.Nil(t, err)
		assert.NotNil(t, id)

		msg, err := message.New()
		assert.Nil(t, err)

		assert.Nil(t, c.Send(id, msg))

		select {
		case msg := <-wg:
			c, ok := msg.Content().([]byte)
			assert.True(t, ok)
			assert.Equal(t, "hello from wasm", string(c))
		case <-time.After(5 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
