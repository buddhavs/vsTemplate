package engage

import (
	"context"
	"errors"
	"fmt"
	"syscall"
	"time"

	cfg "vstmp/pkg/config"
	"vstmp/pkg/log"
	rmq "vstmp/pkg/rabbitmq"
	"vstmp/pkg/signal"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

func fakeHandle1(ctx context.Context, channel *amqp.Channel) error {
	deliveries, err := channel.Consume(
		"cbs_queue_1",
		"cbs_queue_1",
		false, // autoack
		false, // exclusive is unnecessary in rmqctl's usecase.
		false, // nolocal is not supported by rabbitmq.
		false, // nowait
		nil,
	)
	if err != nil {
		return errors.New("channel consume creating failed")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.New("application ends")
		case d, ok := <-deliveries:
			if ok {
				fmt.Printf("%s\n", string(d.Body))
				d.Ack(false)
			} else {
				break
			}
		}
	}
}

func fakeHandle2(ctx context.Context, channel *amqp.Channel) error {
	deliveries, err := channel.Consume(
		"cbs_queue_2",
		"cbs_queue_2",
		false, // autoack
		false, // exclusive is unnecessary in rmqctl's usecase.
		false, // nolocal is not supported by rabbitmq.
		false, // nowait
		nil,
	)
	if err != nil {
		return errors.New("channel consume creating failed")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.New("application ends")
		case d, ok := <-deliveries:
			if ok {
				fmt.Printf("%s\n", string(d.Body))
				d.Ack(false)
			} else {
				break
			}
		}
	}
}

// Start starts the program.
func Start() {
	defer log.Sync()

	ctx, cancel := context.WithCancel(context.Background())

	quitSig := func() {
		cancel()
	}

	signal.RegisterHandler(syscall.SIGQUIT, quitSig)
	signal.RegisterHandler(syscall.SIGTERM, quitSig)
	signal.RegisterHandler(syscall.SIGINT, quitSig)

	r1 := rmq.NewRmq(cfg.Run(), "cbs_queue_1")
	r1.RegisterConsumeHandle(fakeHandle1)
	go r1.Run(ctx)

	r2 := rmq.NewRmq(cfg.Run(), "cbs_queue_2")
	r2.RegisterConsumeHandle(fakeHandle1)
	go r2.Run(ctx)

	<-ctx.Done()

	log.Logger.Info(
		"application ends",
		zap.String(
			"wait_time",
			(time.Duration(10)*time.Second).String()),
	)

	time.Sleep(time.Duration(10) * time.Second)
}
