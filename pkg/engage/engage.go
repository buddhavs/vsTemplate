package engage

import (
	"context"
	vsactor "vstmp/pkg/actor"
	cfg "vstmp/pkg/config"
	"vstmp/pkg/log"
	rmq "vstmp/pkg/rabbitmq"

	"go.uber.org/zap"
)

func engage(ctx context.Context) error {
	// Initialize actors
	vsactor.Init(ctx)

	actorK := vsactor.NewActor(ctx, KindActor, 100, kindActor)

	if err := vsactor.RegisterActor(actorK); err != nil {
		log.Logger.Error(
			"kind actor created error",
			zap.String("service", serviceName),
			zap.String("actor", KindActor),
		)

		return err
	}

	log.Logger.Info(
		"kind actor created",
		zap.String("service", serviceName),
		zap.String("actor", actorK.Name()),
		zap.String("uuid", actorK.UUID()),
	)

	// starts rmq service
	rs, err := rmq.NewRmq(ctx, cfg.GetConfig())
	if err != nil {
		return err
	}

	rs.RegisterConsumeHandle(ontapConsumer)
	go rs.Run()

	return nil
}
