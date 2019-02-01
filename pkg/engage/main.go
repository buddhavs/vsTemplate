package engage

import (
	"context"
	"syscall"
	"time"

	"vstmp/pkg/log"
	"vstmp/pkg/signal"

	"go.uber.org/zap"
)

// Start starts the program
func Start() {
	defer log.Sync()

	ctx, cancel := context.WithCancel(context.Background())

	quitSig := func() {
		cancel()
	}

	signal.RegisterHandler(syscall.SIGQUIT, quitSig)
	signal.RegisterHandler(syscall.SIGTERM, quitSig)
	signal.RegisterHandler(syscall.SIGINT, quitSig)

	engage(ctx)

	<-ctx.Done()

	log.Logger.Info(
		"application ends",
		zap.String("service", serviceName),
		zap.String(
			"wait_time",
			(time.Duration(10)*time.Second).String()),
	)

	time.Sleep(time.Duration(10) * time.Second)
}
