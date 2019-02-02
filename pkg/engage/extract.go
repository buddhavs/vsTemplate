package engage

import (
	"vstmp/pkg/payload"
)

const (
	// Kind_Action_Type
	noSuchActor = iota - 1
	ontapKind
)

// --- extract kind info ---

// extractKind extract payload's metadata Kind
func extractKind(m payload.MetaPayload) int {
	switch m.Kind {
	case "ontap":
		return ontapKind
	default:
		return noSuchActor
	}
}
