package actor

import (
	"context"
	"vstmp/pkg/log"

	"go.uber.org/zap"
)

const (
	// KindActor dispatch to next actor base on Kind/Action/Type
	KindActor = "kind"
	// ErrorActor handle errors
	ErrorActor = "error"
	// DebugActor used for debugging
	DebugActor = "debug"
	// OntapBackupScheduleActor ontap backup scheduled
	OntapBackupScheduleActor = "ontapBackupSchedule"
	// OntapBackupAdhocActor ontap backup adhoc
	OntapBackupAdhocActor = "ontapBackupAdhoc"
	// OntapRestoreActor ontap restore
	OntapRestoreActor = "ontapRestore"
	// OntapUpdatePolicyActor ontap update policy
	OntapUpdatePolicyActor = "ontapUpdatePolicy"
)

// preDefinedActors holds pre-starts actors information
// key: Actor type
// value: slice of interface, [0] channel input, [1] cancel
var definedActors = make(map[string][]interface{})

// Init pre-defined actors standing by
func Init(ctx context.Context) {
	var actor chan<- interface{}
	var cancel context.CancelFunc

	actor, cancel = NewActor(ctx, 100, kindActor)
	definedActors[KindActor] = append(
		[]interface{}{}, actor, cancel)

	log.Logger.Info(
		"kind actor created",
		zap.String("service", serviceName),
		zap.String("actor", KindActor),
	)

	actor, cancel = NewActor(ctx, 100, errorActor)

	definedActors[ErrorActor] = append([]interface{}{}, actor, cancel)

	log.Logger.Info(
		"error actor created",
		zap.String("service", serviceName),
		zap.String("actor", ErrorActor),
	)

	actor, cancel = NewActor(ctx, 10, debugActor)

	definedActors[DebugActor] = append([]interface{}{}, actor, cancel)

	log.Logger.Info(
		"debug actor created",
		zap.String("service", serviceName),
		zap.String("actor", DebugActor),
	)
}

// GetActor gets shared actors
func GetActor(act string) (chan<- interface{}, context.CancelFunc) {
	if v, ok := definedActors[act]; ok {
		return v[0].(chan<- interface{}),
			v[1].(context.CancelFunc)
	}

	return nil, nil
}

// RegisterActor registers actor
func RegisterActor(
	key string,
	actor chan<- interface{},
	cancel context.CancelFunc) {

	definedActors[key] = append(
		[]interface{}{}, actor, cancel)

	log.Logger.Debug(
		"actor registered",
		zap.String("service", serviceName),
		zap.String("actor", key),
	)
}

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
