package oni

import (
	"errors"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"testing"
)

type ConsumerTestSuite struct {
	suite.Suite
}

func TestConsumerTestSuite(t *testing.T) {
	suite.Run(t, new(ConsumerTestSuite))
}

func (suite *ContextTestSuite) TestNewConsumer() {
	suite.Run("TestNewConsumer", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		suite.Assert().NotNil(consumer)
	})
}

func (suite *ContextTestSuite) TestHandler() {
	suite.Run("TestHandler", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		consumer.Handler("event.create.test", func(ctx Context) error {
			return nil
		})

		suite.Assert().Equal(consumer.stream.reader.Config().Topic, "test")
		suite.Assert().Equal(consumer.stream.reader.Config().GroupID, "consumer-group-test")
		suite.Assert().NotNil(consumer.stream.handlers["event.create.test"])
	})
}

func (suite *ContextTestSuite) TestProducer() {
	suite.Run("TestProducer", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		consumer.Producer("producer_func_name", func() *kafka.Writer {
			return &kafka.Writer{
				Addr:  kafka.TCP("localhost:8097"),
				Topic: "test-topic",
			}
		})

		suite.Assert().Equal(consumer.stream.producers["producer_func_name"]().Addr.String(), "localhost:8097")
		suite.Assert().Equal(consumer.stream.producers["producer_func_name"]().Topic, "test-topic")
	})
}

func (suite *ContextTestSuite) TestGroup() {
	suite.Run("TestGroup", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		g := consumer.Group("event.create")
		g.Handler("send.test", func(ctx Context) error {
			return nil
		})

		suite.Assert().Equal(consumer.stream.reader.Config().Topic, "test")
		suite.Assert().Equal(consumer.stream.reader.Config().GroupID, "consumer-group-test")
		suite.Assert().Equal(g.keyGroup, "event.create")
	})
}

func (suite *ContextTestSuite) TestExplicit() {
	suite.Run("TestExplicit", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		consumer.Explicit()
		suite.Assert().Equal(consumer.stream.cm, explicit)
	})
}

func (suite *ContextTestSuite) TestImplicit() {
	suite.Run("TestImplicit", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		consumer.Implicit()
		suite.Assert().Equal(consumer.stream.cm, implicit)
	})
}

func (suite *ContextTestSuite) TestErrorHandler() {
	suite.Run("TestErrorHandler", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))

		consumer.ErrorHandler(func(err error) {
			suite.Assert().NotNil(err)
		})
		suite.Assert().NotNil(consumer.callbackError)
		consumer.callbackError(errors.New("error dummy"))
	})
}

func (suite *ContextTestSuite) TestCloseConsumers() {
	suite.Run("TestCloseConsumers", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))
		suite.Assert().Nil(consumer.closeConsumers())
	})
}

func (suite *ContextTestSuite) TestCloseProducers() {
	suite.Run("TestCloseProducers", func() {
		consumer := NewConsumer(NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		}))
		suite.Assert().Nil(consumer.closeProducers())
	})
}

//func (suite *ContextTestSuite) TestNewConsumer4() {
//
//}
