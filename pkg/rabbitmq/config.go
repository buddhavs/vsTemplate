package rabbitmq

import (
	"sync/atomic"

	config "github.com/buddhavs/vsTemplate/pkg/config"
)

var (
	cleanup = map[string]chan struct{}{
		"reconnect":  make(chan struct{}),
		"rechannel":  make(chan struct{}),
		"catchevent": make(chan struct{}),
	}

	reconn atomic.Value
	rechan atomic.Value
)

// Setup sets rabbitmq configuration
func Setup(cfg config.Config) {

	// initialize atomic reconn
	reconn.Store(false)
	rechan.Store(false)
}
