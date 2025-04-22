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
