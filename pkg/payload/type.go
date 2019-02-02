package payload

import "encoding/json"

type (
	// Metadata metadata
	Metadata struct {
		Kind   string `json:"kind"`
		Action string `json:"action"`
		Type   string `json:"type"`
	}

	// MetaPayload payload's metadata with payload
	MetaPayload struct {
		Metadata
		Payload json.RawMessage `json:"payload"`
	}
)
