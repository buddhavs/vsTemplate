package rabbitmq

import (
	"context"

	"vstmp/pkg/config"
	"vstmp/pkg/log"

	"github.com/streadway/amqp"
	"go.uber.org/zap"
	// "go.uber.org/zap"
)

var (
	rmqConnection      *amqp.Connection
	rmqChannel         *amqp.Channel
	connCloseError     chan *amqp.Error
	channelCancelError chan string
	consumeHandle      ConsumeHandle
)

// Setup sets rabbitmq configuration
func Setup(cfg config.Config) {
	rmqCfg = cfg.GetRmqConnectionConfig()
	rmqQueue = cfg.GetRmqQueueConfig("cbs_queue_1")
}

// RegisterConsumeHandle register consumer's handle
func RegisterConsumeHandle(handle ConsumeHandle) {
	consumeHandle = handle
}

// Run starts rabbitmq service
func Run(ctx context.Context) {
	log.Logger.Info(
		"rabbitmq service starts",
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for s := range start(ctx) {
				log.Logger.Info(
					"rabbitmq status",
					zap.String("status", s),
				)
			}
		}
	}
}
