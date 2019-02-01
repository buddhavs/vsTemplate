package actor

import "encoding/json"

type (
	// --- metadata related ---
	metadata struct {
		Kind   string `json:"kind"`
		Action string `json:"action"`
		Type   string `json:"type"`
	}

	// --- ontap policy related ---
	ontapPolicySchedule struct {
		SnapmirorLabel  string `json:"snapmirrorLabel"`
		SnapshotsToKeep int    `json:"snapshotsToKeep"`
	}

	ontapSnapMirrorPolicy struct {
		Enabled         bool                `json:"enabled"`
		DailySchedule   ontapPolicySchedule `json:"daily-schedule"`
		WeeklySchedule  ontapPolicySchedule `json:"weekly-schedule"`
		MonthlySchedule ontapPolicySchedule `json:"monthly-schedule"`
	}

	// --- ontap connection related ---
	ontapConnectInfo struct {
		OwnerID  string `json:"ownerId"`
		OntapIP  string `json:"ontapIP"`
		Username string `json:"username"`
		Password string `json:"password"`
	}

	ontapVolumeInfo struct {
		VserverName  string `json:"vserverName"`
		VolumeName   string `json:"volumeName"`
		FileSystemID string `json:"fileSystemId"`
	}

	// --- payload related ---
	ontapBackupSchedulePayload struct {
		ontapConnectInfo
		ontapVolumeInfo
		SnapmirrorPolicy ontapSnapMirrorPolicy `json:"snapmirrorPolicy"`
	}

	metaPayload struct {
		metadata
		Payload json.RawMessage `json:"payload"`
	}
)
