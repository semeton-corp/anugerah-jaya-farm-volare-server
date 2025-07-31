package enum

type ItemCategory uint8

const (
	ItemCategoryUnknown        ItemCategory = 0
	ItemCategoryReadyToEatFeed ItemCategory = 1
	ItemCategoryEgg            ItemCategory = 2
	ItemCategoryEquipment      ItemCategory = 3
	ItemCategoryRawMaterial    ItemCategory = 4
	ItemCategoryChicken        ItemCategory = 5
	ItemCategoryCornMaterial   ItemCategory = 6
)

var (
	ItemCategoryMap = map[ItemCategory]string{
		ItemCategoryEgg:            "Telur",
		ItemCategoryEquipment:      "Barang",
		ItemCategoryRawMaterial:    "Bahan Baku Adukan",
		ItemCategoryChicken:        "Ayam",
		ItemCategoryCornMaterial:   "Bahan Baku Adukan - Jagung",
		ItemCategoryReadyToEatFeed: "Pakan Jadi",
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
