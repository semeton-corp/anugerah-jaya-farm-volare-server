package enum

type WarehouseItemCategory uint8

const (
	WarehouseItemCategoryUnknown     WarehouseItemCategory = 0
	WarehouseItemCategoryFeed        WarehouseItemCategory = 1
	WarehouseItemCategoryEgg         WarehouseItemCategory = 2
	WarehouseItemCategoryEquipment   WarehouseItemCategory = 3
	WarehouseItemCategoryRawMaterial WarehouseItemCategory = 4
)

var (
	WarehouseItemCategoryMap = map[WarehouseItemCategory]string{
		WarehouseItemCategoryFeed:        "Pakan",
		WarehouseItemCategoryEgg:         "Telur",
		WarehouseItemCategoryEquipment:   "Barang",
		WarehouseItemCategoryRawMaterial: "Bahan Baku",
	}
)

func (c WarehouseItemCategory) String() string {
	return WarehouseItemCategoryMap[c]
}

func ValueOfWarehouseItemCategory(value string) WarehouseItemCategory {
	for k, v := range WarehouseItemCategoryMap {
		if v == value {
			return k
		}
	}
	return WarehouseItemCategoryUnknown
}

func (c WarehouseItemCategory) IsValid() bool {
	_, ok := WarehouseItemCategoryMap[c]
	return ok
}

type WarehouseOrderStatus uint8

const (
	WarehouseOrderStatusUnknown WarehouseOrderStatus = 0
	WarehouseOrderStatusInSend  WarehouseOrderStatus = 1
	WarehouseOrderStatusDone    WarehouseOrderStatus = 2
)

var (
	WarehouseOrderStatusMap = map[WarehouseOrderStatus]string{
		WarehouseOrderStatusInSend: "Sedang Dikirim",
		WarehouseOrderStatusDone:   "Selesai",
	}
)

func (c WarehouseOrderStatus) String() string {
	return WarehouseOrderStatusMap[c]
}

func ValueOfWarehouseOrderStatus(value string) WarehouseOrderStatus {
	for k, v := range WarehouseOrderStatusMap {
		if v == value {
			return k
		}
	}
	return WarehouseOrderStatusUnknown
}

func (c WarehouseOrderStatus) IsValid() bool {
	_, ok := WarehouseOrderStatusMap[c]
	return ok
}
