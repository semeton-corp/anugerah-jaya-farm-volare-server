package enum

type ProcurementStatus uint8

const (
	ProcurementStatusUnknown      ProcurementStatus = 0
	ProcurementStatusSentOff      ProcurementStatus = 1
	ProcurementStatusArrivedOk    ProcurementStatus = 2
	ProcurementStatusArrivedNotOk ProcurementStatus = 3
)

var (
	ProcurementStatusMap = map[ProcurementStatus]string{
		ProcurementStatusSentOff:      "Sedang Dikirim",
		ProcurementStatusArrivedOk:    "Sampai - Sesuai",
		ProcurementStatusArrivedNotOk: "Sampai - Tidak Sesuai",
	}
)

func (c ProcurementStatus) String() string {
	return ProcurementStatusMap[c]
}

func ValueOfProcurementStatus(value string) ProcurementStatus {
	for k, v := range ProcurementStatusMap {
		if v == value {
			return k
		}
	}
	return ProcurementStatusUnknown
}

func (c ProcurementStatus) IsValid() bool {
	_, ok := ProcurementStatusMap[c]
	return ok
}
