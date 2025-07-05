package enum

type RequestItemStatus uint8

const (
	RequestItemStatusUnknown      RequestItemStatus = 0
	RequestItemStatusSentOff      RequestItemStatus = 1
	RequestItemStatusPending      RequestItemStatus = 2
	RequestItemStatusRejected     RequestItemStatus = 3
	RequestItemStatusCanceled     RequestItemStatus = 4
	RequestItemStatusArrivedOk    RequestItemStatus = 5
	RequestItemStatusArrivedNotOk RequestItemStatus = 6
)

var (
	RequestItemStatusMap = map[RequestItemStatus]string{
		RequestItemStatusSentOff:      "Sedang Dikirim",
		RequestItemStatusPending:      "Menunggu",
		RequestItemStatusRejected:     "Ditolak",
		RequestItemStatusCanceled:     "Dibatalkan",
		RequestItemStatusArrivedOk:    "Sampai - Sesuai",
		RequestItemStatusArrivedNotOk: "Sampai - Tidak Sesuai",
	}
)

func (c RequestItemStatus) String() string {
	return RequestItemStatusMap[c]
}

func ValueOfRequestItemStatus(value string) RequestItemStatus {
	for k, v := range RequestItemStatusMap {
		if v == value {
			return k
		}
	}
	return RequestItemStatusUnknown
}

func (c RequestItemStatus) IsValid() bool {
	_, ok := RequestItemStatusMap[c]
	return ok
}
