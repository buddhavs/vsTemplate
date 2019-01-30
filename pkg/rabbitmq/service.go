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

func start(ctx context.Context) <-chan string {
	status := make(chan string)

	go func() {
		sctx, cancel := context.WithCancel(ctx)

		defer func() {
			log.Logger.Info(
				"re-establish rabbitmq connection",
				zap.String(
					"wait_time",
					(time.Duration(rmqCfg.Wait)*time.Second).String()),
			)

			time.Sleep(time.Duration(rmqCfg.Wait) * time.Second)

			// cleanup consumer goroutine
			cancel()
			// cleanup status
			close(status)
		}()

		// create rabbitmq connection
		if err := createConnect(); err == nil {
			status <- "rabbitmq connection established"
		} else {
			return
		}

		// create rabbitmq channel
		if err := createChannel(); err == nil {
			status <- "rabbitmq channel established"
		} else {
			return
		}

		go consume(sctx)
		status <- "rabbitmq consumer established"

		// block call
		if err := catchEvent(ctx); err != nil {
			status <- fmt.Sprintf("amqp event occured: %s", err.Error())
		}
	}()

	return status
}

func catchEvent(ctx context.Context) error {
	select {
	case <-ctx.Done():
		// TODO: give better error message, e.g pid info
		return errors.New("process ends")
	case err, _ := <-connCloseError:
		log.Logger.Warn(
			"lost rabbitmq connection",
			zap.String("error", err.Error()),
		)

		return err
	case val, _ := <-channelCancelError:
		// interestingly, the amqp library won't trigger
		// this event iff we are not using amqp.Channel
		// to declare the queue, which is, auh, easier
		// for us to handle.
		log.Logger.Warn(
			"lost rabbitmq channel",
			zap.String("error", val),
		)

		return errors.New(val)
	}

	return nil
}

// rmqConnect creates amqp connection
func createConnect() error {
	amqpURL := amqp.URI{
		Scheme:   "amqp",
		Host:     rmqCfg.Host,
		Username: rmqCfg.Username,
		Password: "XXXXX",
		Port:     rmqCfg.Port,
		Vhost:    rmqCfg.Vhost,
	}

	log.Logger.Info(
		"amqp connect URL",
		zap.String("amqp", amqpURL.String()),
	)

	amqpURL.Password = rmqCfg.Password

	// tcp connection timeout in 3 seconds.
	myconn, err := amqp.DialConfig(
		amqpURL.String(),
		amqp.Config{
			Vhost: rmqCfg.Vhost,
			Dial: func(network, addr string) (net.Conn, error) {
				return net.DialTimeout(network, addr, 3*time.Second)
			},
			Heartbeat: 10 * time.Second,
			Locale:    "en_US"},
	)
	if err != nil {
		log.Logger.Warn(
			"Opening amqp connection failed",
			zap.String("error", err.Error()),
		)

		return err
	}

	rmqConnection = myconn
	connCloseError = make(chan *amqp.Error)
	rmqConnection.NotifyClose(connCloseError)
	return nil
}

// rmqChannel creates amqp channel
func createChannel() error {
	myChannel, err := rmqConnection.Channel()
	if err != nil {
		log.Logger.Warn(
			"create amqp channel failed",
			zap.String("error", err.Error()),
		)

		return err
	}

	rmqChannel = myChannel

	// These can be sent from the server when a queue is deleted or
	// when consuming from a mirrored queue where the master has just failed
	// (and was moved to another node).
	channelCancelError = make(chan string)
	rmqChannel.NotifyCancel(channelCancelError)

	return nil
}

func consume(ctx context.Context) {
	if err := consumeHandle(ctx, rmqChannel); err != nil {
		log.Logger.Error(
			"queue handler error",
			zap.String("error", err.Error()),
		)
	}
}
