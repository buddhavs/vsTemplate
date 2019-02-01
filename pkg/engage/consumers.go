package engage

import (
	"context"
	"errors"

	"vstmp/pkg/actor"

	"github.com/streadway/amqp"
	// TODO:
	// using ratelimit with https://github.com/uber-go/ratelimit
	// if necessary
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

	actorK, cancelK := actor.GetActor(actor.KindActor)

	for {
		select {
		case <-ctx.Done():
			cancelK()
			return errors.New("application ends")
		case d, ok := <-deliveries:
			if ok {
				// passing down the byte slice which is a more
				// compact data structure than string,
				// thus do not decoding the message here
				actorK <- d.Body
				d.Ack(false)
			} else {
				break
			}
		}
	}
}
