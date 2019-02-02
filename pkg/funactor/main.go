package funactor

import (
	"context"
	vsactor "vstmp/pkg/actor"
)

// Start starts fun!
func Start() {
	ctx, _ := context.WithCancel(context.Background())

	actor := vsactor.NewActor(ctx, "fun~", 10, funnyActor)

	actor.Send("start~~")
}
