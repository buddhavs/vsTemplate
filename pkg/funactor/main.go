package funactor

import (
	"context"
	vsactor "vstmp/pkg/actor"
)

// Start starts fun!
func Start() {
	ctx, _ := context.WithCancel(context.Background())

	actor := vsactor.NewActor(ctx, "FUN", 10, funnyActor)

	actor.Send("Fun~~")
}
