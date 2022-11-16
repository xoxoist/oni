package oni

import (
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestTemplateSuite struct {
	suite.Suite
}

func TestTemplateTestSuite(t *testing.T) {
	suite.Run(t, new(TestTemplateSuite))
}

func (suite *ContextTestSuite) TestBatchTimeoutWriter() {
	suite.Run("TestBatchTimeoutWriter", func() {
		w := BatchTimeoutWriter(kafka.TCP("localhost:8097"), "2389428", 3*time.Second)
		suite.Assert().Equal(w.Addr.String(), "localhost:8097")
		suite.Assert().Equal(w.Topic, "2389428")
		suite.Assert().Equal(w.BatchTimeout, 3*time.Second)
	})
}

func (suite *ContextTestSuite) TestBasicWriter() {
	suite.Run("TestBasicWriter", func() {
		w := BasicWriter(kafka.TCP("localhost:8097"), "2389428")
		suite.Assert().Equal(w.Addr.String(), "localhost:8097")
		suite.Assert().Equal(w.Topic, "2389428")
	})
}
