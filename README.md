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

### Installation

1. Required go installed on your machine

```sh
go version
```

2. Get Oni Framework

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

### API Examples