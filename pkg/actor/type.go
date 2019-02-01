package actor

import "encoding/json"

type (
	// --- metadata related ---
	metadata struct {
		Kind   string `json:"kind"`
		Action string `json:"action"`
		Type   string `json:"type"`
	}

	metaPayload struct {
		metadata
		Payload json.RawMessage `json:"payload"`
	}
)
