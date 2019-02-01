package engage

import (
	"context"

	cfg "vstmp/pkg/config"
	rmq "vstmp/pkg/rabbitmq"

	"vstmp/pkg/actor"
)

const serviceName = "engage"

func engage(ctx context.Context) error {
	// Initialize actors
	actor.Init(ctx)

	// starts rmq service
	rs, err := rmq.NewRmq(ctx, cfg.GetConfig())
	if err != nil {
		return err
	}

	rs.RegisterConsumeHandle(ontapConsumer)
	go rs.Run()

	return nil
}
