package oni

import (
	"context"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type ContextTestSuite struct {
	suite.Suite
}

func TestContextTestSuite(t *testing.T) {
	suite.Run(t, new(ContextTestSuite))
}

func (suite *ContextTestSuite) TestNewContext() {
	suite.Run("TestNewContext", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		suite.Assert().Equal(oniCtx.message.Topic, "test_topic_1")
		suite.Assert().Equal(oniCtx.message.Partition, 1)
		suite.Assert().Equal(oniCtx.message.Offset, int64(1))
		suite.Assert().Equal(oniCtx.message.HighWaterMark, int64(1))
		suite.Assert().Equal(oniCtx.message.Time, t)
		suite.Assert().Equal(string(oniCtx.message.Key), "event.create.bar")
		suite.Assert().Equal(string(oniCtx.message.Value), "{\"content\":\"this is message content\"}")
		for _, header := range oniCtx.message.Headers {
			suite.Assert().Equal(header.Key, "Content-Type")
			suite.Assert().Equal(string(header.Value), "application/json")
		}

		for k, producer := range oniCtx.producers {
			suite.Assert().Equal(k, "producer_func_1")
			suite.Assert().NotNil(producer)
		}
	})
}

func (suite *ContextTestSuite) TestNewContextMessageFunc() {
	suite.Run("TestNewContextMessageFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		suite.Assert().Equal(oniCtx.Message(), oniCtx.message)
		suite.Assert().Equal(oniCtx.Message().Topic, "test_topic_1")
		suite.Assert().Equal(oniCtx.Message().Partition, 1)
		suite.Assert().Equal(oniCtx.Message().Offset, int64(1))
		suite.Assert().Equal(oniCtx.Message().HighWaterMark, int64(1))
		suite.Assert().Equal(oniCtx.Message().Time, t)
		suite.Assert().Equal(string(oniCtx.Message().Key), "event.create.bar")
		suite.Assert().Equal(string(oniCtx.Message().Value), "{\"content\":\"this is message content\"}")
		for _, header := range oniCtx.Message().Headers {
			suite.Assert().Equal(header.Key, "Content-Type")
			suite.Assert().Equal(string(header.Value), "application/json")
		}
	})
}

func (suite *ContextTestSuite) TestNewContextKeyStringFunc() {
	suite.Run("TestNewContextKeyStringFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		suite.Assert().Equal(oniCtx.Message(), oniCtx.message)
		suite.Assert().Equal(oniCtx.KeyString(), "event.create.bar")
	})
}

func (suite *ContextTestSuite) TestNewContextKeyBytesFunc() {
	suite.Run("TestNewContextKeyBytesFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		suite.Assert().Equal(oniCtx.Message(), oniCtx.message)
		suite.Assert().Equal(oniCtx.KeyBytes(), oniCtx.message.Key)
	})
}

func (suite *ContextTestSuite) TestNewContextValueStringFunc() {
	suite.Run("TestNewContextValueStringFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		suite.Assert().Equal(oniCtx.Message(), oniCtx.message)
		suite.Assert().Equal(oniCtx.ValueString(), "{\"content\":\"this is message content\"}")
	})
}

func (suite *ContextTestSuite) TestNewContextValueBytesFunc() {
	suite.Run("TestNewContextKeyBytesFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		suite.Assert().Equal(oniCtx.Message(), oniCtx.message)
		suite.Assert().Equal(oniCtx.ValueBytes(), oniCtx.message.Value)
	})
}

func (suite *ContextTestSuite) TestNewContextShouldBindJSONFunc() {
	suite.Run("TestNewContextShouldBindJSONFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)

		type dummyStruct struct {
			Content string `json:"content"`
		}
		suite.Assert().Equal(oniCtx.Message(), oniCtx.message)
		var ds dummyStruct
		suite.Assert().Nil(oniCtx.ShouldBindJSON(&ds))
		suite.Assert().Equal(ds.Content, "this is message content")
	})
}

func (suite *ContextTestSuite) TestNewContextOuterContextFunc() {
	suite.Run("TestNewContextOuterContextFunc", func() {
		t := time.Now()
		ctx := context.Background()
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return &kafka.Writer{}
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{
			Topic:         "test_topic_1",
			Partition:     1,
			Offset:        1,
			HighWaterMark: 1,
			Key:           []byte("event.create.bar"),
			Value:         []byte("{\"content\":\"this is message content\"}"),
			Headers: []kafka.Header{
				{
					Key:   "Content-Type",
					Value: []byte("application/json"),
				},
			},
			Time: t,
		}, producers)
		suite.Assert().Equal(oniCtx.OuterContext(), ctx)
	})
}

func (suite *ContextTestSuite) TestNewContextCreateKeyValFunc() {
	suite.Run("TestNewContextCreateKeyValFunc", func() {
		ctx := context.Background()
		oniCtx := newContext(ctx, nil, kafka.Message{}, nil)
		suite.Assert().Equal(oniCtx.OuterContext(), ctx)
		oniCtx.CreateKeyVal("test_key", "29198385829")
		suite.Assert().Equal(oniCtx.FindKey("test_key"), "29198385829")
	})
}

func (suite *ContextTestSuite) TestNewContextGetProducerFunc() {
	suite.Run("TestNewContextGetProducerFunc", func() {
		ctx := context.Background()
		w := &kafka.Writer{
			Topic: "21398893434",
		}
		producers := map[string]ProducerFunc{
			"producer_func_1": func() *kafka.Writer {
				return w
			},
		}
		oniCtx := newContext(ctx, nil, kafka.Message{}, producers)
		suite.Assert().Equal(oniCtx.OuterContext(), ctx)
		suite.Assert().Equal(oniCtx.GetProducer("producer_func_1").Topic, w.Topic)
	})
}

func (suite *ContextTestSuite) TestReaderStats() {
	suite.Run("TestReaderStats", func() {
		ctx := context.Background()
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test-topic-11",
			GroupID: "consumer-group-test",
		}))

		oniCtx := newContext(ctx, consumer.stream.reader, kafka.Message{}, nil)
		suite.Assert().Equal(oniCtx.ReaderStats().Topic, "test-topic-11")
	})
}

func (suite *ContextTestSuite) TestReaderConfig() {
	suite.Run("TestReaderConfig", func() {
		ctx := context.Background()
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test-topic-11",
			GroupID: "consumer-group-test",
		}))

		oniCtx := newContext(ctx, consumer.stream.reader, kafka.Message{}, nil)
		suite.Assert().Equal(oniCtx.ReaderConfig().GroupID, "consumer-group-test")
	})
}
