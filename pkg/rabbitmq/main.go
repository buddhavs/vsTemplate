package rabbitmq

import (
	"context"
)

// Run starts rabbitmq service
func Run(ctx context.Context) {
	go reConnect(ctx)
	go reChannel(ctx)
	go catchAmqpEvent(ctx)
	go consume(ctx)

	startConnection <- struct{}{}
}
