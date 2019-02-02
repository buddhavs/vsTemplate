package funactor

import (
	"context"
	"fmt"
	"math/rand"
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
	for {
		select {
		case <-ctx.Done():
			return
		case v := <-actor.Receive():
			fmt.Printf("%v\n", v)

			d, ok := v.(int)

			go func() {
				actor.Send(rand.Intn(100))
			}()

			if ok && d == 42 {
				fmt.Println("bye bye~")
				actor.Cancel()
			}
		}
	}
}
