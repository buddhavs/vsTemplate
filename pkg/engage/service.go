package engage

import (
	"context"
)

var serviceName = "engage"

func engage(ctx context.Context) error {
	// starts rmq service
	rs, err := rmq.NewRmq(ctx, cfg.GetConfig())
	if err != nil {
		return err
	}

	rs.RegisterConsumeHandle(ontapConsumer)
	go rs.Run()

	return nil
}
