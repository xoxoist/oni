package oni

import (
	"github.com/segmentio/kafka-go"
	"net"
	"time"
)

// BatchTimeoutWriter template
// message was sent will consume after duration defined
// with default round-robin balancer and require ack 0
// which means message do not wait for acknowledgement
// highly recommended when for send message to retry topic
// to avoid fast ping-pong effect between main and retry topic
func BatchTimeoutWriter(addr net.Addr, topic string, duration time.Duration) *kafka.Writer {
	return &kafka.Writer{
		Addr:         addr,
		Topic:        topic,
		BatchTimeout: duration,
	}
}

// BasicWriter template
// message was sent directly to given topic name with
// default round-robin balancer and require ack 0 which
// means message do not wait for acknowledgement
// highly recommended for send basic message to topic
func BasicWriter(addr net.Addr, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:  addr,
		Topic: topic,
	}
}
