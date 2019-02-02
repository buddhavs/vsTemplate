package actor

import (
	"context"
	"fmt"
	"sync"
	"vstmp/pkg/log"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type (
	// Actor is the actor
	Actor struct {
		name     string
		uuid     string
		sender   chan<- interface{}
		receiver <-chan interface{}
		cancel   context.CancelFunc
	}
)

// definedActors holds process wide registered actors
var (
	definedActors = make(map[string]*Actor)
	mLock         sync.Mutex
)

// Init pre-defined actors standing by
func Init(ctx context.Context) {
	var actor *Actor

	actor = NewActor(ctx, ErrorActor, 100, errorActor)
	definedActors[actor.name] = actor

	log.Logger.Info(
		"error actor created",
		zap.String("service", serviceName),
		zap.String("actor", actor.name),
		zap.String("uuid", actor.uuid),
	)

	actor = NewActor(ctx, DebugActor, 10, debugActor)
	definedActors[actor.name] = actor

	log.Logger.Info(
		"debug actor created",
		zap.String("service", serviceName),
		zap.String("actor", actor.name),
		zap.String("uuid", actor.uuid),
	)
}

// GetActor returns registered actor
func GetActor(name string) *Actor {
	defer mLock.Unlock()
	mLock.Lock()

	if v, ok := definedActors[name]; ok {
		return v
	}

	return nil
}

// RegisterActor registers an actor
func RegisterActor(actor *Actor) error {
	defer mLock.Unlock()
	mLock.Lock()

	if _, ok := definedActors[actor.name]; ok {
		err := fmt.Errorf(
			"actor: %s uuid: %s already registered",
			actor.name,
			actor.uuid,
		)

		log.Logger.Error(
			"actor registered failed",
			zap.String("service", serviceName),
			zap.String("actor", actor.name),
			zap.String("uuid", actor.uuid),
			zap.String("error", err.Error()),
		)

		return err
	}

	definedActors[actor.name] = actor

	log.Logger.Debug(
		"actor registered",
		zap.String("service", serviceName),
		zap.String("actor", actor.name),
		zap.String("uuid", actor.uuid),
	)

	return nil
}

// NewActor creates an actor
// ctx: caller's context
// buffer: actor channel buffer size
// callback: callback function takes context and actor channel
// return: actor channel for caller to publish message, cancel to notify actor die out
func NewActor(
	ctx context.Context,
	name string,
	buffer int,
	callback func(context.Context, *Actor)) *Actor {

	ctx, cancel := context.WithCancel(ctx)

	pipe := make(chan interface{}, buffer)

	// make actor as heap object, gives us single object per actor.
	actor := &Actor{
		name:     name,
		uuid:     uuid.New().String(),
		sender:   pipe,
		receiver: pipe,
		cancel:   cancel,
	}

	go func() {
		defer actor.close()
		// block call
		// return closes the channel, actor dies
		callback(ctx, actor)
	}()

	return actor
}

// --- Actor interface functions ---

// close closes actor's internal channel
func (actor *Actor) close() {
	close(actor.sender)
}

// Send sends message to actor
func (actor *Actor) Send(message interface{}) {
	defer func() {
		if r := recover(); r != nil {
			log.Logger.Error(
				"actor in closed state",
				zap.String("service", serviceName),
				zap.String("actor", actor.name),
				zap.String("uuid", actor.uuid),
			)
		}
	}()

	actor.sender <- message
}

// Receive receives message from actor
func (actor *Actor) Receive() <-chan interface{} {
	return actor.receiver
}

// Cancel cancels the actor
func (actor *Actor) Cancel() {
	actor.cancel()
}

// Name returns actor's name
func (actor *Actor) Name() string {
	return actor.name
}

// UUID returns actor's UUID
func (actor *Actor) UUID() string {
	return actor.uuid
}

// --- helper functions ---

// CleanUp is a helper function to cancel actors
func CleanUp(actors ...*Actor) {
	for _, actor := range actors {
		actor.cancel()
	}
}
