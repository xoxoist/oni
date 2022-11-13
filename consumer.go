// Copyright 2022 coffeehaze. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

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
	run(ctx context.Context)
	closeConsumers() error
	closeProducers() error
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

func (c *Consumer) run(ctx context.Context) {
	var wg sync.WaitGroup
	wg.Add(len(c.stream.handlers))
	c.stream.ctx = ctx
	for range c.stream.handlers {
		go c.stream.invoke(&wg)
	}
	wg.Wait()
}

func (c *Consumer) closeConsumers() error {
	return c.stream.closeConsumers()
}

func (c *Consumer) closeProducers() error {
	return c.stream.closeProducers()
}
