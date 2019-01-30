package engage

import (
	"context"
	"errors"
	"fmt"

	cfg "vstmp/pkg/config"
	"vstmp/pkg/log"
	rmq "vstmp/pkg/rabbitmq"

	"github.com/streadway/amqp"
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

	for d := range deliveries {
		fmt.Printf("%s\n", string(d.Body))
		d.Ack(false)
	}

	return errors.New("delivery channel closed")
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

	for d := range deliveries {
		fmt.Printf("%s\n", string(d.Body))
		d.Ack(false)
	}

	return errors.New("delivery channel closed")
}

// Start starts the program.
func Start() {
	defer log.Sync()

	ctx := context.Background()

	rmq.Setup(cfg.Run())
	rmq.RegisterConsumeHandle(fakeHandle1)

	// block call
	rmq.Run(ctx)
}
