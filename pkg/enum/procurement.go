package enum

type ProcurementStatus uint8

const (
	ProcurementStatusUnknown ProcurementStatus = 0
	ProcurementStatusSentOff ProcurementStatus = 1
	ProcurementStatusArrived ProcurementStatus = 2
)

var (
	ProcurmentStatusMap = map[ProcurementStatus]string{
		ProcurementStatusSentOff: "Sedang Dikirim",
		ProcurementStatusArrived: "Sampai",
	}
)

func (c ProcurementStatus) String() string {
	return ProcurmentStatusMap[c]
}

func ValueOfProcurmentStatus(value string) ProcurementStatus {
	for k, v := range ProcurmentStatusMap {
		if v == value {
			return k
		}
	}
	return ProcurementStatusUnknown
}

func (c ProcurementStatus) IsValid() bool {
	_, ok := ProcurmentStatusMap[c]
	return ok
}
