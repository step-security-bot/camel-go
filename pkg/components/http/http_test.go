package http

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/rs/xid"

	camel "github.com/lburgazzoli/camel-go/pkg/api"
	"github.com/lburgazzoli/camel-go/pkg/core/message"
	"github.com/lburgazzoli/camel-go/pkg/util/tests/support"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	// enable components.
	_ "github.com/lburgazzoli/camel-go/pkg/components/timer"

	// enable processors.
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/process"
	_ "github.com/lburgazzoli/camel-go/pkg/core/processors/to"
)

const simpleHTTPGet = `
- route:
    from:
      uri: "timer:foo"
      steps:
        - process:
            ref: "consumer-1"
        - to:
            uri: "http://localhost:3333/uuid"
        - process:
            ref: "consumer-2"
`

func TestSimpleHTTPGet(t *testing.T) {
	support.Run(t, "http_get", func(t *testing.T, ctx context.Context) {
		t.Helper()

		go func() {
			http.HandleFunc("/uuid", func(w http.ResponseWriter, r *http.Request) {
				answer := map[string]any{
					"uuid": xid.New().String(),
				}

				data, err := json.Marshal(answer)
				require.NoError(t, err)

				w.Header().Set("Content-Type", "application/json")
				_, err = io.WriteString(w, string(data))
				require.NoError(t, err)
			})

			require.NoError(
				t,
				http.ListenAndServe(":3333", nil),
			)
		}()

		wg := make(chan camel.Message)

		c := camel.ExtractContext(ctx)

		c.Registry().Set("consumer-1", func(_ context.Context, message camel.Message) error {
			message.SetHeader("Accept", "application/json")
			return nil
		})
		c.Registry().Set("consumer-2", func(_ context.Context, message camel.Message) error {
			wg <- message
			return nil
		})

		err := c.LoadRoutes(ctx, strings.NewReader(simpleHTTPGet))
		assert.Nil(t, err)

		select {
		case msg := <-wg:
			c, ok := msg.Content().([]byte)
			require.True(t, ok)
			require.NotEmpty(t, c)

			ct, ok := msg.Header("Content-Type")
			require.True(t, ok)
			require.Equal(t, ct, "application/json")

			sc, ok := msg.Attribute(AttributeStatusCode)
			require.True(t, ok)
			require.Equal(t, 200, sc)

			data, err := message.ContentAsBytes(msg)
			require.NoError(t, err)
			require.NotEmpty(t, data)

			m := make(map[string]any)
			err = json.Unmarshal(data, &m)
			require.NoError(t, err)
			require.NotEmpty(t, m["uuid"])

		case <-time.After(60 * time.Second):
			assert.Fail(t, "timeout")
		}
	})
}
