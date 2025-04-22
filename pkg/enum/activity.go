package enum

type ActivityStatus uint8

const (
	ActivityStatusUnknown ActivityStatus = 0
	ActivityStatusIn      ActivityStatus = 1
	ActivityStatusOut     ActivityStatus = 2
)

var (
	ActivityStatusMap = map[ActivityStatus]string{
		ActivityStatusIn:  "Masuk",
		ActivityStatusOut: "Keluar",
	}
)

func (c ActivityStatus) String() string {
	return ActivityStatusMap[c]
}

func ValueOfActivityStatus(value string) ActivityStatus {
	for k, v := range ActivityStatusMap {
		if v == value {
			return k
		}
	}
	return ActivityStatusUnknown
}

func (c ActivityStatus) IsValid() bool {
	_, ok := ActivityStatusMap[c]
	return ok
}
