package ontap

import (
	"context"
	"encoding/json"
	"fmt"
	vsactor "vstmp/pkg/actor"
	"vstmp/pkg/payload"
)

const (
	// OntapBackupScheduleActor ontap backup scheduled
	OntapBackupScheduleActor = "ontapBackupScheduleActor"
	// OntapBackupAdhocActor ontap backup adhoc
	OntapBackupAdhocActor = "ontapBackupAdhocActor"
	// OntapRestoreActor ontap restore
	OntapRestoreActor = "ontapRestoreActor"
	// OntapUpdatePolicyActor ontap update policy
	OntapUpdatePolicyActor = "ontapUpdatePolicyActor"
)

// DispatchActor ontap's dispatch actor
func DispatchActor(ctx context.Context, actor *vsactor.Actor) {

	actorE := vsactor.GetActor(vsactor.ErrorActor)
	actorD := vsactor.GetActor(vsactor.DebugActor)

	defer vsactor.CleanUp(actorE, actorD)

	for {
		select {
		case <-ctx.Done():
			return
		case d := <-actor.Receive():
			m := d.(payload.MetaPayload)

			switch v := ontapExtractAction(m); v {
			case ontapBackupSchedule:
				if actor := vsactor.GetActor(OntapBackupScheduleActor); actor == nil {
					actor = vsactor.NewActor(
						ctx, OntapBackupScheduleActor, 100, ontapBackupScheduleActor)

					_ = vsactor.RegisterActor(actor)

					defer vsactor.CleanUp(actor)

					actor.Send(m)
				} else {
					defer vsactor.CleanUp(actor)

					actor.Send(m)
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

func ontapBackupScheduleActor(ctx context.Context, actor *vsactor.Actor) {
	actorE := vsactor.GetActor(vsactor.ErrorActor)
	actorD := vsactor.GetActor(vsactor.DebugActor)

	defer vsactor.CleanUp(actorE, actorD)

	for {
		select {
		case <-ctx.Done():
			return
		case m := <-actor.Receive():
			meta := m.(payload.MetaPayload)

			var p ontapBackupSchedulePayload
			if err := json.Unmarshal(meta.Payload, &p); err != nil {
				actorE.Send(err)
				continue
			}

			actorD.Send(
				fmt.Sprintf(
					"IP: %v, ID: %v, PWD: %v, USER: %v, VolName: %v, VServer: %v\n\n",
					p.OntapIP,
					p.OwnerID,
					p.Password,
					p.Username,
					p.VolumeName,
					p.VserverName,
				),
			)
		}
	}
}
