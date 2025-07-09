package enum

type ItemHistoryStatus uint8

const (
	ItemHistoryStatusUnknown ItemHistoryStatus = 0
	ItemHistoryStatusIn      ItemHistoryStatus = 1
	ItemHistoryStatusOut     ItemHistoryStatus = 2
	ItemHistoryStockUpdated  ItemHistoryStatus = 3
)

var (
	ItemHistoryStatusMap = map[ItemHistoryStatus]string{
		ItemHistoryStockUpdated: "Stok Diperbarui",
		ItemHistoryStatusIn:     "Barang Masuk",
		ItemHistoryStatusOut:    "Barang Keluar",
	}
)

func (c ItemHistoryStatus) String() string {
	return ItemHistoryStatusMap[c]
}

func ValueOfItemHistoryStatus(value string) ItemHistoryStatus {
	for k, v := range ItemHistoryStatusMap {
		if v == value {
			return k
		}
	}
	return ItemHistoryStatusUnknown
}

func (c ItemHistoryStatus) IsValid() bool {
	_, ok := ItemHistoryStatusMap[c]
	return ok
}
