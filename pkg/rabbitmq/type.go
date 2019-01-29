package rabbitmq

import (
	"sync/atomic"

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

	// RmqStruct struct contains config variable of rmqtool
	RmqStruct struct {
		// connection and channel
		Connection *amqp.Connection
		Channel    *amqp.Channel
		// connection configuration
		user              string
		passwd            string
		host              string
		port              uint
		queueName         string
		consumerTag       string
		mode              mode
		exchange          string
		exchangeType      string
		listeningExchange string
		listeningQueue    string
		reconnWait        uint
		rechanWait        uint
		reconn            atomic.Value
		rechan            atomic.Value
		// channel to control rmq behavior
		startConnection chan struct{}
		startChannel    chan struct{}
		canConsume      chan struct{}
		connCloseError  chan *amqp.Error
		chanCloseError  chan *amqp.Error
		chanCancelError chan string
		chanReturnError chan amqp.Return
		// consumer callback function
		consumeCallbackFunc func(rmq *RmqStruct)
		// in memory queue of amqp.Delivery
		ackQueue map[string]*amqp.Delivery
		// bool type at the end of struct
		durable bool
	}
)
