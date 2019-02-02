package funactor

import (
	"context"
	"fmt"
	"time"
	vsactor "vstmp/pkg/actor"
)

const (
	funActor = "FunActor"
)

var counter = 0

func calculate(n int) int {
	if n == 1 {
		return 1
	}
	return n * calculate(n-1)
}

func futureActor(ctx context.Context, actor *vsactor.Actor) {
	for {
		select {
		case <-ctx.Done():
			return
		case d := <-actor.Receive():
			v := d.(int)
			if v == 42 {
				fmt.Println("HIT!")
				r := calculate(321)
				fmt.Printf("SHIT r: %v\n\n", r)
				actor.Send(r)
			}
			return
		}
	}
}

func funnyActor(ctx context.Context, actor *vsactor.Actor) {
	var factor *vsactor.Actor

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-actor.Receive():
			fmt.Printf("%v\n", v)

			counter++

			if counter == 42 {
				// give me a future!
				factor = vsactor.NewActor(ctx, "future", 10, futureActor)
				factor.Send(42)
			}
			actor.Send("received Fun! pass on~~")
			if counter > 100 {
				time.Sleep(10 * time.Second)
				// self kill after receive the future :-)
				fmt.Printf("give me my future result:-) %v", <-factor.Receive())
				actor.Cancel()
			}
		}
	}
}
