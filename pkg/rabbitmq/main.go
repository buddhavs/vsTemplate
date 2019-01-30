package rabbitmq

import (
	"context"

	"vstmp/pkg/config"
	"vstmp/pkg/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NewRmq creates new rabbitmq instance
func NewRmq(cfg config.Config, queue string) *RmqStruct {
	// validate it's value.
	rs := RmqStruct{
		uuid:     uuid.New().String(),
		rmqCfg:   cfg.GetRmqConnectionConfig(),
		rmqQueue: cfg.GetRmqQueueConfig(queue),
	}

	return &rs
}

// RegisterConsumeHandle register consumer's handle
func (rmq *RmqStruct) RegisterConsumeHandle(handle ConsumeHandle) {
	rmq.consumeHandle = handle
}

// Run starts rabbitmq service
func (rmq *RmqStruct) Run(ctx context.Context) {
	log.Logger.Info(
		"rabbitmq service starts",
		zap.String("uuid", rmq.uuid),
	)

	for {
		select {
		case <-ctx.Done():
			return
		default:
			for s := range rmq.start(ctx) {
				log.Logger.Info(
					"rabbitmq status",
					zap.String("status", s),
				)
			}
		}
	}

}
