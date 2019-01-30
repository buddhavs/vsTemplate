package config

import (
	"github.com/streadway/amqp"
)

type (
	// ConnectionType config for amqp connection
	ConnectionType struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
		Host     string `validate:"required"`
		Vhost    string `validate:"required"`
		Port     int    `validate:"required"`
		Wait     int    `validate:"isdefault"`
	}

	// QueueType config for consuming queue
	QueueType struct {
		QueueName string `validate:"required"`
		// Consumer is used for active Cancel consumer
		Consumer  string `validate:"required"`
		AutoAck   bool   `validate:"-"`
		Exclusive bool   `validate:"-"`
		// Not supported in rabbitmq
		NoLocal bool       `validate:"-"`
		NoWait  bool       `validate:"-"`
		Args    amqp.Table `validate:"-"`
	}

	// Config program config
	Config struct {
		rmqConnection ConnectionType
		rmqQueueMap   map[string]QueueType
	}
)
