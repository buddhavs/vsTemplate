package ontap

import "vstmp/pkg/payload"

const (
	noSuchActor = iota - 1
	ontapBackupSchedule
	ontapBackupAdhoc
	ontapRestore
	ontapUpdatePolicy
)

func ontapExtractType(m payload.MetaPayload) int {
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

func ontapExtractAction(m payload.MetaPayload) int {
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
