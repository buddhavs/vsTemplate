package actor

import (
	"context"
	"fmt"
	"vstmp/pkg/log"

	"go.uber.org/zap"
)

// debugActor used for debugging :-)
func debugActor(ctx context.Context, input <-chan interface{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-input:
			fmt.Printf("print whatever we got: %v\n", v)
		}
	}
}

// errorActor processing errors
func errorActor(ctx context.Context, input <-chan interface{}) {
	for {
		select {
		case <-ctx.Done():
			return
		case err := <-input:
			e := err.(error)
			log.Logger.Error(
				"error logged",
				zap.String("service", serviceName),
				zap.String("actor", ErrorActor),
				zap.String("error", e.Error()),
			)
		}
	}
}

// kindActor dispatch kind of jobs to next actor
func kindActor(ctx context.Context, input <-chan interface{}) {
	actorE, cancelE := GetActor(ErrorActor)
	actorD, cancelD := GetActor(DebugActor)
	cancelSlice := append([]context.CancelFunc{}, cancelE, cancelD)

	defer func() {
		for _, c := range cancelSlice {
			c()
		}
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case d := <-input:
			// base on the input then dispatch to next actor
			if v, ok := d.([]byte); ok {
				var m metaPayload

				if err := actorDecoder(v, &m); err != nil {
					// dispatch to errorActor
					actorE <- err
					continue
				}

				// for debugging
				actorD <- m

				// for small payload, we always copy by value, which is fast
				// and avoid golang's escape analysis kicks in
				switch v := extractKind(m); v {
				case ontapBackupSchedule:
					if actor, cancel := GetActor(OntapBackupScheduleActor); actor == nil {

						actor, cancel = NewActor(ctx, 100, ontapBackupScheduleActor)
						RegisterActor(OntapBackupScheduleActor, actor, cancel)
						cancelSlice = append(cancelSlice, cancel)

						actor <- m

					} else {
						actor <- m
					}
				case ontapBackupAdhoc:
				case ontapRestore:
				case ontapUpdatePolicy:
				case noSuchActor:
				default:
				}
			}
		}
	}
}
