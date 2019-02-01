package actor

import (
	"context"
	"encoding/json"
	"fmt"
)

func ontapBackupScheduleActor(ctx context.Context, input <-chan interface{}) {
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
		case m := <-input:
			meta := m.(metaPayload)

			var p ontapBackupSchedulePayload
			if err := json.Unmarshal(meta.Payload, &p); err != nil {
				actorE <- err
				continue
			}

			actorD <- fmt.Sprintf(
				"IP: %v, ID: %v, PWD: %v, USER: %v, VolName: %v, VServer: %v",
				p.OntapIP,
				p.OwnerID,
				p.Password,
				p.Username,
				p.VolumeName,
				p.VserverName,
			)
		}
	}
}
