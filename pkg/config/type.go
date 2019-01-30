package config

import (
	"github.com/streadway/amqp"
)

type (
	// RmqConnectionType config for amqp connection
	RmqConnectionType struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
		Host     string `validate:"required"`
		Vhost    string `validate:"required"`
		Port     int    `validate:"required"`
		Wait     int    `validate:"isdefault"`
	}

	// RmqQueueType config for consuming queue
	RmqQueueType struct {
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
		// TODO: use https://github.com/hashicorp/go-version
		version       string // v0.0.0
		rmqConnection RmqConnectionType
		rmqQueueMap   map[string]RmqQueueType
	}
)
