package enum

type LocationWorkType uint8

const (
	LocationWorkTypeUnknown   LocationWorkType = 0
	LocationWorkTypeCage      LocationWorkType = 1
	LocationWorkTypeStore     LocationWorkType = 2
	LocationWorkTypeWarehouse LocationWorkType = 3
)

var (
	LocationWorkTypeMap = map[LocationWorkType]string{
		LocationWorkTypeCage:      "Kandang",
		LocationWorkTypeStore:     "Toko",
		LocationWorkTypeWarehouse: "Gudang",
	}
)

func (c LocationWorkType) String() string {
	return LocationWorkTypeMap[c]
}

func ValueOfLocationWorkType(value string) LocationWorkType {
	for k, v := range LocationWorkTypeMap {
		if v == value {
			return k
		}
	}
	return LocationWorkTypeUnknown
}

func (c LocationWorkType) IsValid() bool {
	_, ok := LocationWorkTypeMap[c]
	return ok
}
