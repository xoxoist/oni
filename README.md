# Oni Kafka Framework

<img align="right" width="159px" src="https://raw.githubusercontent.com/coffeehaze/asset/main/oni.png">

[![codecov](https://codecov.io/gh/coffeehaze/oni/branch/master/graph/badge.svg)](https://codecov.io/gh/coffeehaze/oni)
[![Go Report Card](https://goreportcard.com/badge/github.com/coffeehaze/oni)](https://goreportcard.com/report/github.com/coffeehaze/oni)
[![GoDoc](https://pkg.go.dev/badge/github.com/coffeehaze/oni?status.svg)](https://pkg.go.dev/github.com/coffeehaze/oni?tab=doc)
[![Join the chat at https://gitter.im/coffeehaze/oni](https://badges.gitter.im/Join%20Chat.svg)](https://gitter.im/coffeehaze/oni?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)
[![Sourcegraph](https://sourcegraph.com/github.com/coffeehaze/oni/-/badge.svg)](https://sourcegraph.com/github.com/coffeehaze/oni?badge)
[![Release](https://img.shields.io/github/release/coffeehaze/oni.svg?style=flat-square)](https://github.com/coffeehaze/oni/releases)
[![TODOs](https://badgen.net/https/api.tickgit.com/badgen/github.com/coffeehaze/oni)](https://www.tickgit.com/browse?repo=github.com/coffeehaze/oni)

Oni is Kafka Framework written in Go (Golang). that makes you easy to consume and produce kafka messages using robust
API wrapper for [kafka-go](https://github.com/segmentio/kafka-go) thanks to [segmentio](https://github.com/segmentio).
the usage most likely same with Gin / Echo web framework.

Oni Mask art by [@inksyndromeartwork](https://www.freepik.com/author/inksyndromeartwork)

## Contents

- [Oni Kafka Framework](#oni-kafka-framework)
    - [Contents](#contents)
    - [Installation](#installation)
    - [Quick Start](#quick-start)
    - [API Examples](#api-examples)
        - [Stream](#stream)
        - [Consumer](#consumer)
        - [Context](#consumer)

### Installation

1. Required go installed on your machine

```sh
go version
```

2. Get oni and kafka-go

```sh
go get -u github.com/coffeehaze/oni
go get -u github.com/segmentio/kafka-go
```

3. Import oni

```go
import "github.com/coffeehaze/oni"
```

### Quick Start

1. Create package `model` and create `foo.go` file and place this code to it

```go
package model

type Foo struct {
	FooContent string `json:"foo_content"`
}
```

2. Create package `consumer` and create `main.go` file and place this code to it

```go
package main

import (
	"context"
	"fmt"
	"github.com/coffeehaze/oni"
	"github.com/your/projectname/model"
	"github.com/segmentio/kafka-go"
	"syscall"
	"time"
)

func main() {
	ctx := context.Background()
	defer ctx.Done()

	// initialize consumer
	stream := oni.NewStream(kafka.ReaderConfig{
		Brokers: []string{
			"localhost:8097", // kafka brokers 1
			"localhost:8098", // kafka brokers 2
			"localhost:8099", // kafka brokers 3 you can only define one inside array
		},
		Topic: "foos", // topic you want to listen at
	})
	foosConsumer := oni.NewConsumer(stream)
	foosConsumer.Handler(
		"create.foo", // event key you want to map to specific handler function
		func(ctx oni.Context) error {
			var foo model.Foo
			err := ctx.ShouldBindJSON(&foo) // bind message value to struct
			if err != nil {
				return err
			}
			fmt.Println(fmt.Sprintf("key=%s value=%s foo=%s", ctx.KeyString(), ctx.ValueString(), foo))
			return nil
		},
	)

	// initialize oni runner
	// to help start and graceful shutdown all producer and consumer you defined
	oniRunner := oni.Runner{
		Context: ctx,
		Timeout: 15 * time.Second,
		Syscall: oni.SyscallOpt(
			syscall.SIGINT,
			syscall.SIGTERM,
			syscall.SIGHUP,
		),
		Consumers: oni.ConsumerOpt(foosConsumer),
	}
	<-oniRunner.Start()
}

```

3. Create package `producer` and create `main.go` file and place this code to it

```go
package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/your/projectname/model"
	"github.com/segmentio/kafka-go"
	"time"
)

func main() {
	foos := &kafka.Writer{
		Addr:  kafka.TCP("localhost:8097"), // kafka broker
		Topic: "foos",                      // target topic you want to send
	}

	// create json object using json.Marshal
	fooObj := model.Foo{
		FooContent: fmt.Sprintf("This is new foo %d", time.Now().Unix()),
	}
	fooByte, _ := json.Marshal(fooObj)

	err := foos.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte("create.foo"), // target event you want to send at
		Value: fooByte,              // fooObj marshal result as the value
	})
	fmt.Println(err)
}

```

4. Start run `consumer/main.go` and run `producer/main.go` separately

### API Examples

### Stream

- `oni.NewStream(cfg kafka.ReaderConfig)`
    ```go
    // example for oni.NewStream(cfg kafka.ReaderConfig)
    stream := oni.NewStream(kafka.ReaderConfig{
        Brokers: []string{"localhost:8097"},
        Topic: "foos",
    })
    ```

### Consumer

- `oni.NewConsumer(stream *oni.Stream) IConsumer`
    ```go
    // example for oni.NewConsumer(stream *oni.Stream)
    // default consume mode is implicit
    consumer := oni.NewConsumer(stream)
    ```
- `IConsumer.ErrorHandler(callbackFunc ErrorCallbackFunc)`
    ```go
    // set global error handler which will be invoked when oni.HandlerFunc returns error
    // this function should be called before handler creation and only called once 
    consumer.ErrorHandler(func (err error) {
        if err != nil {
            // do error handling logic such as logging
            // or handle special error cases
        }
    })
    ```
- `IConsumer.Implicit()`
    ```go
    // set consume mode to implicit which means every message
    // received by *oni.Stream will be automatically ack or committed
    // this function should be called before handler creation
    consumer.Implicit()
    ```
- `IConsumer.Explicit()`
    ```go
    // set consume mode to explicit which means every message 
    // received by *oni.Stream will be ack or committed manually using Context.Ack()
    // this function should be called before handler creation
    consumer.Explicit()
    ```

- `IConsumer.Group(keyGroup string) *Consumer`
    ```go
    // create new group of consumer key event prefix, for example `event.notification.blast`
    // could be had last suffix like `email.channel` and `sms.channel` so your next handler
    // should be had key event `email.channel` and `sms.channel` and actual handler event key 
    // for each handler should look like this
    // `event.notification.blast.email.channel`
    // `event.notification.blast.sms.channel`
    // and previous *oni.Stream behavior should inherit to new group 
    notificationBlastEvent := consumer.Group("event.notification.blast")
    notificationBlastEvent.Handler("email.channel", func (ctx oni.Context) error {})
    notificationBlastEvent.Handler("sms.channel", func (ctx oni.Context) error {})
    ```
- `IConsumer.Handler(key string, handlerFunc ...HandlerFunc)`
    ```go
    // create handler function for specific key event, for this example is `event.send.email`
    // message that produced to `event.send.email` key will be received by handler function
    // defined in this handler creation and message value and details will be propagated through
    // oni.Context
    consumer.Handler("event.send.email", func (ctx oni.Context) error {
        // put business logic here
        return nil
    })
    ```
- `IConsumer.Producer(name string, producerFunc ProducerFunc)`
    ```go
    // create producer that can be accessed by its name through oni.Context functions that
    // allows business logic to access producer by the name and use it for sending message
    // to targeted topic, can be defined multiple times and only created when producer get
    // called by oni.Context function and closed after sent the message  
    consumer.Producer("notification_producer", func() *kafka.Writer {
        // you can define using *kafka.Writer struct or you can use templates from oni
        // by returning this oni.BasicWriter(addr net.Addr, topic string) function.
        return &kafka.Writer{
            Addr:  kafka.TCP("localhost:8097"),
            Topic: "notification",
        }
    })
    ```

### Context