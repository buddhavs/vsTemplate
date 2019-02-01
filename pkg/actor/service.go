package actor

import (
	"encoding/json"
	"errors"
)

const serviceName = "actor"

const (
	// Kind_Action_Type
	noSuchActor = iota - 1
	ontapBackupSchedule
	ontapBackupAdhoc
	ontapRestore
	ontapUpdatePolicy
)

func actorDecoder(d []byte, payload interface{}) error {
	if json.Valid(d) {
		return json.Unmarshal(d, payload)
	}

	err := errors.New("invalid json format")
	return err
}

func ontapExtractType(m metaPayload) int {
	switch m.Type {
	case "schedule":
		return ontapBackupSchedule
	case "adhoc":
		return ontapBackupAdhoc
	case "policy":
		return ontapUpdatePolicy
	default:
		return noSuchActor
	}
}

func ontapExtractAction(m metaPayload) int {
	switch m.Action {
	case "backup":
		return ontapExtractType(m)
	case "update":
		return ontapExtractType(m)
	case "restore":
		return ontapRestore
	default:
		return noSuchActor
	}
}

func extractKind(m metaPayload) int {
	switch m.Kind {
	case "ontap":
		return ontapExtractAction(m)
	default:
		return noSuchActor
	}
}
