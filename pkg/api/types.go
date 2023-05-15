package api

import (
	"context"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"

	"github.com/asynkron/protoactor-go/actor"
)

type ContextKey string

const (
	ContextKeyCamelContext = ContextKey("camel-context")
	ContextKeyActorContext = ContextKey("actor-context")
)

func Wrap(_ context.Context, camelContext Context) context.Context {
	return context.WithValue(context.Background(), ContextKeyCamelContext, camelContext)
}

func ExtractContext(ctx context.Context) Context {
	value := ctx.Value(ContextKeyCamelContext)
	if value == nil {
		panic(fmt.Errorf("unable to get CamelContext from context"))
	}

	answer, ok := value.(Context)
	if !ok {
		panic(fmt.Errorf("type cast error %v", value))
	}

	return answer
}

func ExtractActorContext(ctx context.Context) actor.Context {
	value := ctx.Value(ContextKeyActorContext)
	if value == nil {
		panic(fmt.Errorf("unable to get actor Context from context"))
	}

	answer, ok := value.(actor.Context)
	if !ok {
		panic(fmt.Errorf("type cast error %v", value))
	}

	return answer
}

type Parameters map[string]interface{}

type Closer interface {
	// Close closes the resource.
	Close(context.Context) error
}

type Service interface {
	Start(context.Context) error
	Stop(context.Context) error
}

type Identifiable interface {
	ID() string
}

type Registry interface {
	Get(string) (interface{}, bool)
	Set(string, interface{})
	Del(key string) interface{}
}

type Properties interface {
	AddSource(string) error
	String(string) string
}

//nolint:interfacebloat
type Context interface {
	Identifiable
	Service
	Closer

	Registry() Registry
	Properties() Properties
	TypeConverter() TypeConverter

	LoadRoutes(ctx context.Context, in io.Reader) error

	// Spawn ---
	// TODO: must be hidden
	// TODO: each route must have its own context/supervisor
	Spawn(Verticle) (*actor.PID, error)

	// Send ---
	Send(string, Message) error

	// SendTo ---
	// TODO: must be hidden maybe
	SendTo(*actor.PID, Message) error

	// SendToAs ---
	// TODO: must be hidden maybe
	SendToAs(*actor.PID, *actor.PID, Message) error

	// Receive ---
	Receive(string, time.Duration) (Message, error)

	// Request ---
	// TODO: must be hidden maybe
	RequestTo(*actor.PID, Message, time.Duration) (Message, error)

	Logger() *zap.Logger

	NewMessage() Message
}

type Component interface {
	Identifiable

	Context() Context
	Scheme() string
	Endpoint(Parameters) (Endpoint, error)

	Logger() *zap.Logger
}

type Endpoint interface {
	Identifiable
	Service

	Context() Context
	Component() Component

	Logger() *zap.Logger
}

const (
	MessageAttributeID            string = "camel.apache.org/message.id"
	MessageAttributeTime          string = "camel.apache.org/message.time"
	MessageAttributeType          string = "camel.apache.org/message.type"
	MessageAttributeSource        string = "camel.apache.org/message.source"
	MessageAttributeSubject       string = "camel.apache.org/message.subject"
	MessageAttributeContentType   string = "camel.apache.org/message.content.type"
	MessageAttributeContentSchema string = "camel.apache.org/message.content.schema"
)

//nolint:interfacebloat
type Message interface {
	ID() string
	Time() time.Time

	Context() Context

	Type() string
	Source() string
	Subject() string
	ContentSchema() string
	ContentType() string

	SetType(string)
	SetSource(string)
	SetSubject(string)
	SetContentSchema(string)
	SetContentType(string)

	// Content ---
	// TODO: add options Content(opt.AsType(baz{})).
	Content() interface{}
	SetContent(interface{})

	// Error ---
	Error() error
	SetError(error)

	// Headers ---
	Headers() map[string]any
	SetHeaders(map[string]any)
	Header(string) (any, bool)
	SetHeader(string, any)
	ForEachHeader(func(string, any))

	// Attributes ---
	Attributes() map[string]any
	SetAttributes(map[string]any) error
	Attribute(string) (any, bool)
	SetAttribute(string, any) error
	ForEachAttribute(func(string, any))

	CopyTo(message Message) error
}

type Processor = func(context.Context, Message) error
type Predicate = func(context.Context, Message) (bool, error)

type Producer interface {
	Service
	Verticle

	Endpoint() Endpoint
}

type ProducerFactory interface {
	Producer() (Producer, error)
}

type Consumer interface {
	Service
	Verticle

	Endpoint() Endpoint
}

type ConsumerFactory interface {
	Consumer(pid *actor.PID) (Consumer, error)
}

type Verticle interface {
	Identifiable

	actor.Actor
}

type TypeConverter interface {
	Convert(interface{}, interface{}) (bool, error)
}

type RawJSON map[string]interface{}
