package rabbitmq

import (
	"context"

	"github.com/streadway/amqp"
	// _ "gopkg.in/go-playground/validator.v9"
)

// ConsumeHandle handler for consumer
type ConsumeHandle func(context.Context, *amqp.Channel) error
