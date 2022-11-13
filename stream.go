// Copyright 2022 coffeehaze. All rights reserved.
// Use of this source code is governed by a MIT style
// license that can be found in the LICENSE file.

package oni

import (
	"context"
	"github.com/segmentio/kafka-go"
	"sync"
)

const (
	implicit consumeMode = iota
	explicit
)

type consumeMode int

type HandlerFunc func(ctx Context) error

type ErrorCallbackFunc func(err error)

type ProducerFunc func() *kafka.Writer

type handler struct {
	HandlerFunc       HandlerFunc
	ErrorCallbackFunc ErrorCallbackFunc
}

type IStream interface {
	addHandler(key string, handlerFunc HandlerFunc, errorCallbackFunc ErrorCallbackFunc)
	addProducer(name string, producerFunc ProducerFunc)
	invoke(wg *sync.WaitGroup)
	closeConsumers() error
	closeProducers() error
	stream()
}

type Stream struct {
	reader    *kafka.Reader
	handlers  map[string][]handler
	producers map[string]ProducerFunc
	fLock     sync.Mutex
	eLock     sync.Mutex
	cm        consumeMode
	ctx       context.Context
}

func NewStream(config kafka.ReaderConfig) *Stream {
	return &Stream{
		reader:    kafka.NewReader(config),
		handlers:  make(map[string][]handler),
		producers: make(map[string]ProducerFunc),
		cm:        implicit,
	}
}

func (s *Stream) closeConsumers() error {
	return s.reader.Close()
}

func (s *Stream) closeProducers() error {
	var err error
	for key, producer := range s.producers {
		err = producer().Close()
		delete(s.producers, key)
		return err
	}
	return nil
}

func (s *Stream) invoke(wg *sync.WaitGroup) {
	s.stream()
	defer wg.Done()
}

func (s *Stream) stream() {
	for {
		var m kafka.Message
		var err error

		switch s.cm {
		case implicit:
			m, err = s.reader.ReadMessage(s.ctx)
		case explicit:
			m, err = s.reader.FetchMessage(s.ctx)
		}
		if err != nil {
			break
		}

		oniCtx := newContext(s.ctx, s.reader, m, s.producers)
		for _, handler := range s.handlers[string(m.Key)] {
			s.fLock.Lock()
			err := handler.HandlerFunc(oniCtx)
			s.fLock.Unlock()
			if err != nil {
				if handler.ErrorCallbackFunc != nil {
					s.eLock.Lock()
					handler.ErrorCallbackFunc(err)
					s.eLock.Unlock()
				}
				break
			}
		}
	}
}

func (s *Stream) addHandler(key string, handlerFunc HandlerFunc, errorCallbackFunc ErrorCallbackFunc) {
	s.handlers[key] = append(s.handlers[key], handler{
		HandlerFunc:       handlerFunc,
		ErrorCallbackFunc: errorCallbackFunc,
	})
}

func (s *Stream) addProducer(name string, producerFunc ProducerFunc) {
	s.producers[name] = producerFunc
}
