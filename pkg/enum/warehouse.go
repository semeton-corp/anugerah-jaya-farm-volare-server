package enum

type WarehouseItemCategory uint8

const (
	WarehouseItemCategoryUnknown WarehouseItemCategory = 0
	WarehouseItemCategoryFeed    WarehouseItemCategory = 1
	WarehouseItemCategoryEgg     WarehouseItemCategory = 2
)

var (
	WarehouseItemCategoryMap = map[WarehouseItemCategory]string{
		WarehouseItemCategoryFeed: "Pakan",
		WarehouseItemCategoryEgg:  "Telur",
	}
)

func (c WarehouseItemCategory) String() string {
	return WarehouseItemCategoryMap[c]
}

func ValueOfWarehouseItemCategory(value string) ChickenCategory {
	for k, v := range ChickenCategoryMap {
		if v == value {
			return k
		}
	}
	return ChickenCategoryUnknown
}

func (c WarehouseItemCategory) IsValid() bool {
	_, ok := WarehouseItemCategoryMap[c]
	return ok
}
