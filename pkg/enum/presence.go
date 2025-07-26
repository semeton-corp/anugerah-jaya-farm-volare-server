package enum

type PresenceStatus uint8

const (
	PresenceStatusUnknown    PresenceStatus = 0
	PresenceStatusPresent    PresenceStatus = 1
	PresenceStatusSick       PresenceStatus = 2
	PresenceStatusPermission PresenceStatus = 3
	PresenceStatusAlpha      PresenceStatus = 4
)

var (
	PresenceStatusMap = map[PresenceStatus]string{
		PresenceStatusAlpha:      "Alpha",
		PresenceStatusPresent:    "Hadir",
		PresenceStatusSick:       "Sakit",
		PresenceStatusPermission: "Izin",
	}
)

func (c PresenceStatus) String() string {
	return PresenceStatusMap[c]
}

func ValueOfPresenceStatus(value string) PresenceStatus {
	for k, v := range PresenceStatusMap {
		if v == value {
			return k
		}
	}
	return PresenceStatusUnknown
}

func (c PresenceStatus) IsValid() bool {
	_, ok := PresenceStatusMap[c]
	return ok
}

type SubmissionPresenceStatus uint8

const (
	SubmissionPresenceStatusUnknown      SubmissionPresenceStatus = 0
	SubmissionPresenceStatusAccepted     SubmissionPresenceStatus = 1
	SubmissionPresenceStatusPending      SubmissionPresenceStatus = 2
	SubmissionPresenceStatusRejected     SubmissionPresenceStatus = 3
	SubmissionPresenceStatusNoSubmission SubmissionPresenceStatus = 4
)

var (
	SubmissionPresenceStatusMap = map[SubmissionPresenceStatus]string{
		SubmissionPresenceStatusAccepted:     "Disetujui",
		SubmissionPresenceStatusPending:      "Menunggu",
		SubmissionPresenceStatusRejected:     "Ditolak",
		SubmissionPresenceStatusNoSubmission: "-",
	}
)

func (c SubmissionPresenceStatus) String() string {
	return SubmissionPresenceStatusMap[c]
}

func ValueOfWarehouseSubmissionPresenceStatus(value string) SubmissionPresenceStatus {
	for k, v := range SubmissionPresenceStatusMap {
		if v == value {
			return k
		}
	}
	return SubmissionPresenceStatusUnknown
}

func (c SubmissionPresenceStatus) IsValid() bool {
	_, ok := SubmissionPresenceStatusMap[c]
	return ok
}
