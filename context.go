package oni

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/segmentio/kafka-go"
)

type Context interface {
	ShouldBindJSON(v interface{}) error
	ShouldRetryWith(producerFuncName string) error
	ShouldErrorWith(producerFuncName string) error
	ShouldReturnWith(producerFuncName string) error

	Ack() error
	ValueBytes() []byte
	ValueString() string
	KeyBytes() []byte
	KeyString() string

	Message() kafka.Message
	ReaderStats() kafka.ReaderStats
	ReaderConfig() kafka.ReaderConfig
	GetProducer(producerFuncName string) *kafka.Writer

	OuterContext() context.Context
	FindKey(key string) interface{}
	CreateKeyVal(key string, val interface{})
}

type retry struct {
	OriginKey string `json:"origin_key"`
	Value     string `json:"value"`
}

type octx struct {
	producers    map[string]ProducerFunc
	outerContext context.Context
	message      kafka.Message
	reader       *kafka.Reader
}

func newContext(ctx context.Context, r *kafka.Reader, m kafka.Message, producers map[string]ProducerFunc) *octx {
	return &octx{outerContext: ctx, reader: r, message: m, producers: producers}
}

func (ctx *octx) ShouldBindJSON(v interface{}) error {
	return json.Unmarshal(ctx.message.Value, v)
}

func (ctx *octx) ValueBytes() []byte {
	return ctx.message.Value
}

func (ctx *octx) ValueString() string {
	return string(ctx.message.Value)
}

func (ctx *octx) KeyBytes() []byte {
	return ctx.message.Key
}

func (ctx *octx) KeyString() string {
	return string(ctx.message.Key)
}

func (ctx *octx) Ack() error {
	return ctx.reader.CommitMessages(ctx.outerContext, ctx.message)
}

func (ctx *octx) Message() kafka.Message {
	return ctx.message
}

func (ctx *octx) ReaderStats() kafka.ReaderStats {
	return ctx.reader.Stats()
}

func (ctx *octx) ReaderConfig() kafka.ReaderConfig {
	return ctx.reader.Config()
}

func (ctx *octx) GetProducer(producerFuncName string) *kafka.Writer {
	return ctx.producers[producerFuncName]()
}

func (ctx *octx) FindKey(key string) interface{} {
	return ctx.outerContext.Value(key)
}

func (ctx *octx) CreateKeyVal(key string, val interface{}) {
	ctx.outerContext = context.WithValue(ctx.outerContext, key, val)
}

func (ctx *octx) OuterContext() context.Context {
	return ctx.outerContext
}

func (ctx *octx) ShouldRetryWith(producerFuncName string) error {
	w := ctx.producers[producerFuncName]()

	retryDataByte, err := json.Marshal(retry{
		OriginKey: ctx.KeyString(),
		Value:     ctx.ValueString(),
	})
	if err != nil {
		return err
	}

	err = w.WriteMessages(ctx.outerContext, kafka.Message{
		Key:   []byte(fmt.Sprintf("%s.%s", "retry", ctx.KeyString())),
		Value: retryDataByte,
	})
	if err != nil {
		return err
	}

	return w.Close()
}

func (ctx *octx) ShouldErrorWith(producerFuncName string) error {
	w := ctx.producers[producerFuncName]()

	retryDataByte, err := json.Marshal(retry{
		OriginKey: ctx.KeyString(),
		Value:     ctx.ValueString(),
	})
	if err != nil {
		return err
	}

	err = w.WriteMessages(ctx.outerContext, kafka.Message{
		Key:   []byte(fmt.Sprintf("%s.%s", "failed", ctx.KeyString())),
		Value: retryDataByte,
	})
	if err != nil {
		return err
	}

	return w.Close()
}

func (ctx *octx) ShouldReturnWith(producerFuncName string) error {
	w := ctx.producers[producerFuncName]()

	var retryData retry
	err := json.Unmarshal(ctx.message.Value, &retryData)
	if err != nil {
		return err
	}

	err = w.WriteMessages(ctx.outerContext, kafka.Message{
		Key:   []byte(retryData.OriginKey),
		Value: []byte(retryData.Value),
	})
	if err != nil {
		return err
	}

	return w.Close()
}
