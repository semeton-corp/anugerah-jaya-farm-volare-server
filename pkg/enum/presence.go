package enum

type PresenceStatus uint8

const (
	PresenceStatusUnknown    PresenceStatus = 0
	PresenceStatusPresent    PresenceStatus = 1
	PresenceStatusNotPresent PresenceStatus = 2
)

var (
	PresenceStatusMap = map[PresenceStatus]string{
		PresenceStatusPresent:    "Hadir",
		PresenceStatusNotPresent: "Tidak Hadir",
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
