package enum

type LocationType uint8

const (
	LocationTypeUnknown   LocationType = 0
	LocationTypeCage      LocationType = 1
	LocationTypeStore     LocationType = 2
	LocationTypeWarehouse LocationType = 3
	LocationTypeSite      LocationType = 4
)

var (
	LocationWorkTypeMap = map[LocationType]string{
		LocationTypeCage:      "Kandang",
		LocationTypeStore:     "Toko",
		LocationTypeWarehouse: "Gudang",
		LocationTypeSite:      "Site",
	}
)

func (c LocationType) String() string {
	return LocationWorkTypeMap[c]
}

func ValueOfLocationType(value string) LocationType {
	for k, v := range LocationWorkTypeMap {
		if v == value {
			return k
		}
	}
	return LocationTypeUnknown
}

func (c LocationType) IsValid() bool {
	_, ok := LocationWorkTypeMap[c]
	return ok
}
