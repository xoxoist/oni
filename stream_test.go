package oni

import (
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"testing"
)

type TestStreamSuite struct {
	suite.Suite
}

func TestStreamTestSuite(t *testing.T) {
	suite.Run(t, new(TestStreamSuite))
}

func (suite *ContextTestSuite) TestCloseProducersStream() {
	suite.Run("TestCloseProducersStream", func() {
		s := NewStream(kafka.ReaderConfig{
			Brokers: []string{"localhost:8097"},
			Topic:   "test",
			GroupID: "consumer-group-test",
		})

		suite.Assert().Nil(s.closeProducers())
		suite.Assert().Nil(s.producers["a"])
		suite.Assert().Nil(s.producers["b"])
		suite.Assert().Nil(s.producers["c"])
		suite.Assert().Len(s.producers, 0)
	})
}
