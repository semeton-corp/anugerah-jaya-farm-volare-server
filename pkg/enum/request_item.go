package enum

type RequestItemStatus uint8

const (
	RequestItemStatusUnknown  RequestItemStatus = 0
	RequestItemStatusSentOff  RequestItemStatus = 1
	RequestItemStatusPending  RequestItemStatus = 2
	RequestItemStatusAccepted RequestItemStatus = 3
	RequestItemStatusRejected RequestItemStatus = 4
)

var (
	RequestItemStatusMap = map[RequestItemStatus]string{
		RequestItemStatusSentOff:  "Dikirim",
		RequestItemStatusAccepted: "Diterima",
		RequestItemStatusPending:  "Menunggu",
		RequestItemStatusRejected: "Ditolak",
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
