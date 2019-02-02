package ontap

type (
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
)
