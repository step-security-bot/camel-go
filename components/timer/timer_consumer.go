package timer

import (
	"time"

	"github.com/lburgazzoli/camel-go/api"
	"github.com/lburgazzoli/camel-go/camel"
	"github.com/lburgazzoli/camel-go/processor"
)

// ==========================
//
// Producer
//
// ==========================

func newTimerConsumer(endpoint *timerEndpoint) *timerConsumer {
	c := timerConsumer{
		endpoint: endpoint,
		// TODO: this is ugly
		processor: processor.NewProcessingPipeline(func(api.Exchange) {
		}),
	}

	return &c
}

type timerConsumer struct {
	endpoint  *timerEndpoint
	processor api.Processor
	ticker    *time.Ticker
}

func (consumer *timerConsumer) Start() {
	consumer.ticker = time.NewTicker(consumer.endpoint.period)
	go func() {
		var counter uint64

		for t := range consumer.ticker.C {
			exchange := camel.NewExchange(consumer.endpoint.component.context)

			counter++

			exchange.Headers().Bind("timer.fire.time", t.UTC())
			exchange.Headers().Bind("timer.fire.count", counter)
			exchange.SetBody(nil)

			consumer.processor.Publish(exchange)
		}
	}()
}

func (consumer *timerConsumer) Stop() {
	if consumer.ticker != nil {
		consumer.ticker.Stop()
	}
}

func (consumer *timerConsumer) Stage() api.ServiceStage {
	return api.ServiceStageConsumer
}

func (consumer *timerConsumer) Endpoint() api.Endpoint {
	return consumer.endpoint
}

func (consumer *timerConsumer) Processor() api.Processor {
	return consumer.processor
}
