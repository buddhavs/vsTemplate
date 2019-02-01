package actor

import "context"

const (
	// KindActor as kind actor
	KindActor = "kind"
	// ActionActor as action actor
	ActionActor = "action"
	// TypeActor as type actor
	TypeActor = "type"
	// MonitorActor as monitor actor
	MonitorActor = "monitor"
)

// NewActor creates an actor
// ctx: caller's context
// buffer: actor channel buffer size
// callback: callback function takes context and actor channel
// return: actor channel for caller to publish message, cancel to notify actor die out
func NewActor(
	ctx context.Context,
	buffer int,
	callback func(
		context.Context,
		<-chan interface{}),
) (chan<- interface{}, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)

	actor := make(chan interface{}, buffer)

	go func() {
		defer func() {
			close(actor)
		}()

		// block call
		// return closes the channel, actor dies
		callback(ctx, actor)

	}()

	return actor, cancel
}
