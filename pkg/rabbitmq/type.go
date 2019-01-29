package rabbitmq

import (
	"github.com/streadway/amqp"
	// _ "gopkg.in/go-playground/validator.v9"
)

type (
	connectionType struct {
		Username string `validate:"required"`
		Password string `validate:"required"`
		Host     string `validate:"required"`
		Vhost    string `validate:"required"`
		Port     int    `validate:"required"`
	}

	publishType struct {
		connectionType
		ExchangeName string `validate:"required"`
		Key          string `validate:"required"`
		Mandatory    bool   `validate:"-"`
		Immediate    bool   `validate:"-"`
		// Message      amqp.Publishing `validate:"required"`
	}

	consumeType struct {
		connectionType
		QueueName string `validate:"required"`
		// Consumer is not used in our case. Set as default string value.
		Consumer  string     `validate:"isdefault"`
		AutoAck   bool       `validate:"-"`
		Exclusive bool       `validate:"-"`
		NoLocal   bool       `validate:"-"`
		NoWait    bool       `validate:"-"`
		Args      amqp.Table `validate:"-"`
	}
)
