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

type PresenceFilter uint8

const (
	PresenceFilterUnknown  PresenceFilter = 0
	PresenceFilterThisWeek PresenceFilter = 1
	PresenceFilterJanuary  PresenceFilter = 2
	// Note : wait confirmation
)

var (
	PresenceFilterMap = map[PresenceFilter]string{
		PresenceFilterThisWeek: "Minggu Ini",
		PresenceFilterJanuary:  "Januari",
	}
)

func (c PresenceFilter) String() string {
	return PresenceFilterMap[c]
}

func ValueOfPresenceFilter(value string) PresenceFilter {
	for k, v := range PresenceFilterMap {
		if v == value {
			return k
		}
	}
	return PresenceFilterUnknown
}

func (c PresenceFilter) IsValid() bool {
	_, ok := PresenceFilterMap[c]
	return ok
}
