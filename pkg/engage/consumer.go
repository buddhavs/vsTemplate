package engage

import (
	"context"
	"errors"

	"vstmp/pkg/actor"
	"vstmp/pkg/log"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	// https://github.com/uber-go/ratelimit
)

// ontapConsumer is the consumer for ontap
func ontapConsumer(ctx context.Context, channel *amqp.Channel) error {
	deliveries, err := channel.Consume(
		"cbs_queue_1",
		"cbs_queue_1",
		false, // autoack
		false, // exclusive is unnecessary in rmqctl's usecase
		false, // nolocal is not supported by rabbitmq
		false, // nowait
		nil,
	)

	if err != nil {
		return errors.New("channel consume creating failed")
	}

	// create Kind actor here
	log.Logger.Info(
		"kind actor created",
		zap.String("service", serviceName),
		zap.String("actor", actor.KindActor),
	)
	actor, cancel := actor.NewActor(ctx, 10, kindActor)

	for {
		select {
		case <-ctx.Done():
			cancel()
			return errors.New("application ends")
		case d, ok := <-deliveries:
			if ok {
				actor <- d.Body
				// fmt.Printf("%s\n", string(d.Body))
				d.Ack(false)
			} else {
				break
			}
		}
	}
}
