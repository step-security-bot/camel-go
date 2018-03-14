package camel

import (
	"github.com/rs/zerolog/log"
)

// NewPipe --
func NewPipe() *Pipe {
	return &Pipe{
		done: make(chan bool, 1),
		In:   nil,
		Next: nil,
	}
}

// NewPipeIn --
func NewPipeIn() *Pipe {
	return &Pipe{
		done: make(chan bool, 1),
		In:   make(chan *Exchange),
		Next: nil,
	}
}

// NewPipeWithNext --
func NewPipeWithNext(pipe *Pipe) *Pipe {
	return &Pipe{
		done: pipe.done,
		In:   nil,
		Next: pipe,
	}
}

// Pipe --
type Pipe struct {
	done chan bool
	In   chan *Exchange
	Next *Pipe
}

// Publish --
func (pipe *Pipe) Publish(exchange *Exchange) *Pipe {
	if pipe.Next != nil && pipe.Next.In != nil {
		pipe.Next.In <- exchange
	}

	return pipe
}

// PublishAsync --
func (pipe *Pipe) PublishAsync(exchange *Exchange) *Pipe {
	go pipe.Publish(exchange)

	return pipe
}

// Process --
func (pipe *Pipe) Process(processor Processor, processors ...Processor) *Pipe {
	next := Pipe{}
	next.In = make(chan *Exchange)
	next.done = pipe.done
	next.Next = pipe.Next

	go func() {
		for {
			select {
			case exchange, ok := <-pipe.In:
				if !ok {
					log.Warn().Msgf("Channel %+v is not ready", pipe.In)
				} else {
					exchange = processor(exchange)

					for _, proc := range processors {
						exchange = proc(exchange)
					}

					next.Publish(exchange)
				}
			case <-pipe.done:
				log.Info().Msg("done")
				return
			}
		}
	}()

	return &next
}
