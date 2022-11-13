package oni

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

type Runner struct {
	Context   context.Context
	Timeout   time.Duration
	Syscall   []os.Signal
	Consumers []*Consumer
}

func SyscallOpt(syscall ...os.Signal) []os.Signal {
	return syscall
}

func ConsumerOpt(consumers ...*Consumer) []*Consumer {
	return consumers
}

func (r *Runner) Start() <-chan struct{} {
	for _, consumer := range r.Consumers {
		go consumer.Run(r.Context)
	}

	wait := make(chan struct{})
	go func() {
		s := make(chan os.Signal, 1)

		//syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP
		signal.Notify(s, r.Syscall...)
		<-s
		log.Println("shutting down")

		timeoutFunc := time.AfterFunc(r.Timeout, func() {
			log.Printf("timeout %d ms has been elapsed, force exit", r.Timeout.Milliseconds())
			os.Exit(0)
		})

		defer timeoutFunc.Stop()

		var wg sync.WaitGroup

		for i, consumer := range r.Consumers {
			wg.Add(1)
			sequence := i
			c := consumer
			go func() {
				defer wg.Done()

				log.Printf("cleaning up process %d", sequence)

				if err := c.CloseConsumers(); err != nil {
					log.Printf("consumer %d clean up failed: %s", sequence, err.Error())
					return
				}
				if err := c.CloseProducers(); err != nil {
					log.Printf("producer %d clean up failed: %s", sequence, err.Error())
					return
				}

				log.Printf("process %d was shutdown gracefully", sequence)
			}()
		}
		wg.Wait()
		close(wait)
	}()

	return wait
}
