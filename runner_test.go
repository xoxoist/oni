package oni

import (
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/suite"
	"syscall"
	"testing"
)

type TestRunnerSuite struct {
	suite.Suite
}

func TestRunnerTestSuite(t *testing.T) {
	suite.Run(t, new(TestRunnerSuite))
}

func (suite *ContextTestSuite) TestSyscallOpt() {
	suite.Run("TestSyscallOpt", func() {
		sysCalls := SyscallOpt(
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGHUP,
		)

		suite.Assert().Equal(len(sysCalls), 3)
		for i, sysCall := range sysCalls {
			suite.Assert().Equal(sysCall, sysCalls[i])
		}
	})
}

func (suite *ContextTestSuite) TestConsumerOpt() {
	suite.Run("TestConsumerOpt", func() {
		consumers := ConsumerOpt(
			NewConsumer(NewStream(kafka.ReaderConfig{
				Brokers: []string{"localhost:8091"},
				Topic:   "test1",
				GroupID: "consumer-group-test-1",
			})),
			NewConsumer(NewStream(kafka.ReaderConfig{
				Brokers: []string{"localhost:8092"},
				Topic:   "test2",
				GroupID: "consumer-group-test-2",
			})),
			NewConsumer(NewStream(kafka.ReaderConfig{
				Brokers: []string{"localhost:8093"},
				Topic:   "test3",
				GroupID: "consumer-group-test-3",
			})),
		)

		suite.Assert().Equal(len(consumers), 3)
		for i, consumer := range consumers {
			suite.Assert().Equal(consumer.stream.reader.Stats().Topic, consumers[i].stream.reader.Stats().Topic)
			suite.Assert().Equal(consumer.stream.reader.Config().GroupID, consumers[i].stream.reader.Config().GroupID)
			suite.Assert().Equal(consumer.stream.reader.Config().Brokers[0], consumers[i].stream.reader.Config().Brokers[0])
		}
	})
}
