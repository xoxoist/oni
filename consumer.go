package oni

import (
	"context"
	"fmt"
	"sync"
)

type IConsumer interface {
	Handler(key string, handlerFunc ...HandlerFunc)
	ErrorHandler(callbackFunc ErrorCallbackFunc)
	Producer(name string, producerFunc ProducerFunc)
	Group(keyGroup string) *Consumer
	Run(ctx context.Context)
	CloseConsumers() error
	CloseProducers() error
	Explicit()
	Implicit()
}

type Consumer struct {
	stream        *Stream
	keyGroup      string
	callbackError ErrorCallbackFunc
}

func NewConsumer(stream *Stream) *Consumer {
	return &Consumer{
		stream: stream,
	}
}

func (c *Consumer) Handler(key string, handlerFunc ...HandlerFunc) {
	if len(c.keyGroup) != 0 {
		key = fmt.Sprintf("%s.%s", c.keyGroup, key)
	}
	for _, f := range handlerFunc {
		c.stream.addHandler(key, f, c.callbackError)
	}
}

func (c *Consumer) Producer(name string, producerFunc ProducerFunc) {
	c.stream.addProducer(name, producerFunc)
}

func (c *Consumer) Group(keyGroup string) *Consumer {
	return &Consumer{
		stream:   c.stream,
		keyGroup: keyGroup,
	}
}

func (c *Consumer) Explicit() {
	c.stream.cm = explicit
}

func (c *Consumer) Implicit() {
	c.stream.cm = implicit
}

func (c *Consumer) ErrorHandler(callbackFunc ErrorCallbackFunc) {
	c.callbackError = callbackFunc
}

func (c *Consumer) Run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(c.stream.handlers))
	c.stream.ctx = ctx
	for range c.stream.handlers {
		go c.stream.invoke(&wg)
	}
	wg.Wait()
}

func (c *Consumer) CloseConsumers() error {
	return c.stream.closeConsumers()
}

func (c *Consumer) CloseProducers() error {
	return c.stream.closeProducers()
}
