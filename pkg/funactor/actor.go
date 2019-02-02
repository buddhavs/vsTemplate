package funactor

import (
	"context"
	"fmt"
	vsactor "vstmp/pkg/actor"
)

const (
	funActor = "FunActor"
)

var counter = 0

// debugActor used for debugging :-)
func funnyActor(ctx context.Context, actor *vsactor.Actor) {
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-actor.Receive():
			fmt.Printf("%v\n", v)
			counter++
			actor.Send("received Fun! pass on~~")
			if counter > 100 {
				// self kill :-)
				actor.Cancel()
			}
		}
	}
}
