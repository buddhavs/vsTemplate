package rabbitmq

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"vstmp/pkg/log"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
)

var serviceName = "rabbitmq"

// defaultHandle is a template for handler writers
func defaultHandle(ctx context.Context, channel *amqp.Channel) error {
	log.Logger.Info(
		"default rabbitmq consume handler",
		zap.String("service", serviceName),
	)

	<-ctx.Done()
	return errors.New("service ends")
}

func (rmq *RmqStruct) start() <-chan string {
	status := make(chan string)

	go func() {
		sctx, cancel := context.WithCancel(rmq.ctx)
		var reconnect = true

		defer func() {
			if reconnect {
				log.Logger.Info(
					"re-establish connection",
					zap.String("service", serviceName),
					zap.String("uuid", rmq.uuid),
					zap.String(
						"wait_time",
						(time.Duration(rmq.rmqCfg.Wait)*time.Second).String()),
				)

				time.Sleep(time.Duration(rmq.rmqCfg.Wait) * time.Second)
			}

			// cleanup consumer goroutine
			cancel()
			// cleanup status
			close(status)
		}()

		// create rabbitmq connection
		if err := rmq.createConnect(); err != nil {
			return
		}
		status <- "rabbitmq connection established"

		// create rabbitmq channel
		if err := rmq.createChannel(); err != nil {
			return
		}
		status <- "rabbitmq channel established"

		go rmq.consume(sctx)
		status <- "rabbitmq consumer established"

		err := rmq.catchEvent().(retryError)
		reconnect = err.reconnect
		status <- fmt.Sprintf("amqp event occurred: %s", err.Error())
	}()

	return status
}

func (rmq *RmqStruct) catchEvent() error {
	select {
	case <-rmq.ctx.Done():
		log.Logger.Info(
			"service ends",
			zap.String("service", serviceName),
			zap.String("uuid", rmq.uuid),
		)

		return retryError{
			errors.New("service ends, cleanup connection loop"),
			false,
		}
	case err, _ := <-rmq.connCloseError:
		log.Logger.Warn(
			"lost connection",
			zap.String("service", serviceName),
			zap.String("uuid", rmq.uuid),
			zap.String("error", err.Error()),
		)

		return retryError{
			err,
			true,
		}
	case val, _ := <-rmq.channelCancelError:
		// interestingly, the amqp library won't trigger
		// this event iff we are not using amqp.Channel
		// to declare the queue, which is, auh, easier
		// for us to handle.
		log.Logger.Warn(
			"lost channel",
			zap.String("service", serviceName),
			zap.String("uuid", rmq.uuid),
			zap.String("error", val),
		)

		return retryError{
			errors.New(val),
			true,
		}
	}
}

// createConnect creates amqp connection
func (rmq *RmqStruct) createConnect() error {
	amqpURL := amqp.URI{
		Scheme:   "amqp",
		Host:     rmq.rmqCfg.Host,
		Username: rmq.rmqCfg.Username,
		Password: "XXXXX",
		Port:     rmq.rmqCfg.Port,
		Vhost:    rmq.rmqCfg.Vhost,
	}

	log.Logger.Info(
		"connect URL",
		zap.String("service", serviceName),
		zap.String("uuid", rmq.uuid),
		zap.String("url", amqpURL.String()),
	)

	amqpURL.Password = rmq.rmqCfg.Password

	// tcp connection timeout in 3 seconds
	myconn, err := amqp.DialConfig(
		amqpURL.String(),
		amqp.Config{
			Vhost: rmq.rmqCfg.Vhost,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 3*time.Second)
			},
			Heartbeat: 10 * time.Second,
			Locale:    "en_US"},
	)
	if err != nil {
		log.Logger.Warn(
			"open connection failed",
			zap.String("service", serviceName),
			zap.String("uuid", rmq.uuid),
			zap.String("error", err.Error()),
		)

		return err
	}

	rmq.rmqConnection = myconn
	rmq.connCloseError = make(chan *amqp.Error)
	rmq.rmqConnection.NotifyClose(rmq.connCloseError)
	return nil
}

// createChannel creates amqp channel
func (rmq *RmqStruct) createChannel() error {
	myChannel, err := rmq.rmqConnection.Channel()
	if err != nil {
		log.Logger.Warn(
			"create channel failed",
			zap.String("service", serviceName),
			zap.String("uuid", rmq.uuid),
			zap.String("error", err.Error()),
		)

		return err
	}

	rmq.rmqChannel = myChannel

	// These can be sent from the server when a queue is deleted or
	// when consuming from a mirrored queue where the master has just failed
	// (and was moved to another node)
	rmq.channelCancelError = make(chan string)
	rmq.rmqChannel.NotifyCancel(rmq.channelCancelError)

	return nil
}

func (rmq *RmqStruct) consume(ctx context.Context) {
	if err := rmq.consumeHandle(ctx, rmq.rmqChannel); err != nil {
		log.Logger.Warn(
			"queue consume error",
			zap.String("service", serviceName),
			zap.String("uuid", rmq.uuid),
			zap.String("error", err.Error()),
		)
	}
}
