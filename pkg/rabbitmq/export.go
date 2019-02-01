package rabbitmq

import (
	"context"

	"vstmp/pkg/config"
	"vstmp/pkg/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

// NewRmq creates new rabbitmq instance
// NewXxx always takes context as argument
func NewRmq(ctx context.Context, cfg config.Config) (*RmqStruct, error) {
	// TODO:
	// log using config version.
	// validate it's value by calling cfg.ValidateRmq()
	// return error if validate fails
	rs := RmqStruct{
		ctx:           ctx,
		uuid:          uuid.New().String(),
		rmqCfg:        cfg.GetRmqConnectionConfig(),
		consumeHandle: defaultHandle,
	}

	return &rs, nil
}

// RegisterConsumeHandle register consumer's handle
func (rmq *RmqStruct) RegisterConsumeHandle(handle ConsumeHandle) {
	rmq.consumeHandle = handle
}

// Run starts rabbitmq service
func (rmq *RmqStruct) Run() {
	log.Logger.Info(
		"service starts",
		zap.String("service", serviceName),
		zap.String("uuid", rmq.uuid),
	)

	for {
		select {
		case <-rmq.ctx.Done():
			return
		default:
			for s := range rmq.start() {
				log.Logger.Info(
					"status",
					zap.String("service", serviceName),
					zap.String("uuid", rmq.uuid),
					zap.String("status", s),
				)
			}
		}
	}

}
