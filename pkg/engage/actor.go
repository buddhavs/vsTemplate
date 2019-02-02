package engage

import (
	"context"
	vsactor "vstmp/pkg/actor"
	"vstmp/pkg/ontap"
	"vstmp/pkg/payload"
)

const (
	// KindActor dispatch payload to next actor, base on Kind/Action/Type
	KindActor = "kindActor"
	// OntapActor ontap dispatch actor
	OntapActor = "ontapActor"
)

// kindActor dispatch kind of jobs to next actor
func kindActor(ctx context.Context, actor *vsactor.Actor) {
	actorE := vsactor.GetActor(vsactor.ErrorActor)
	actorD := vsactor.GetActor(vsactor.DebugActor)

	defer vsactor.CleanUp(actorE, actorD)

	for {
		select {
		case <-ctx.Done():
			return
		case d := <-actor.Receive():
			// base on the input then dispatch to next actor
			if v, ok := d.([]byte); ok {
				var m payload.MetaPayload

				if err := vsactor.JSONDecoder(v, &m); err != nil {
					// dispatch to errorActor
					actorE.Send(err)
					continue
				}

				// for debugging
				actorD.Send(m)

				// for small payload, we always copy by value, which is fast
				// and avoid golang's escape analysis kicks in
				switch v := extractKind(m); v {
				case ontapKind:
					if actor := vsactor.GetActor(OntapActor); actor == nil {
						actor = vsactor.NewActor(ctx, OntapActor, 100, ontap.DispatchActor)
						_ = vsactor.RegisterActor(actor)

						defer vsactor.CleanUp(actor)

						actor.Send(m)
					} else {
						defer vsactor.CleanUp(actor)

						actor.Send(m)
					}
				case noSuchActor:
				default:
				}
			}
		}
	}
}
