package enum

type ItemCategory uint8

const (
	ItemCategoryUnknown     ItemCategory = 0
	ItemCategoryFeed        ItemCategory = 1
	ItemCategoryEgg         ItemCategory = 2
	ItemCategoryEquipment   ItemCategory = 3
	ItemCategoryRawMaterial ItemCategory = 4
	ItemCategoryChicken     ItemCategory = 5
)

var (
	ItemCategoryMap = map[ItemCategory]string{
		ItemCategoryFeed:        "Pakan",
		ItemCategoryEgg:         "Telur",
		ItemCategoryEquipment:   "Barang",
		ItemCategoryRawMaterial: "Bahan Baku",
		ItemCategoryChicken:     "Ayam",
	}
)

func (c ItemCategory) String() string {
	return ItemCategoryMap[c]
}

func ValueOfWarehouseItemCategory(value string) ItemCategory {
	for k, v := range ItemCategoryMap {
		if v == value {
			return k
		}
	}
	return ItemCategoryUnknown
}

func (c ItemCategory) IsValid() bool {
	_, ok := ItemCategoryMap[c]
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
