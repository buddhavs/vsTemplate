package rabbitmq

import (
	"context"
	"vstmp/pkg/config"

	"github.com/streadway/amqp"
	// _ "gopkg.in/go-playground/validator.v9"
)

// ConsumeHandle handler for consumer
type ConsumeHandle func(context.Context, *amqp.Channel) error

// RmqStruct is the instance of rabbitmq service
type RmqStruct struct {
	uuid               string
	rmqCfg             config.RmqConnectionType
	rmqQueue           config.RmqQueueType
	rmqConnection      *amqp.Connection
	rmqChannel         *amqp.Channel
	connCloseError     chan *amqp.Error
	channelCancelError chan string
	consumeHandle      ConsumeHandle
}
