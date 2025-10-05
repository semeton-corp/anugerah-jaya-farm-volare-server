package enum

type LocationType uint8

const (
	LocationTypeUnknown    LocationType = 0
	LocationTypeCage       LocationType = 1
	LocationTypeStore      LocationType = 2
	LocationTypeWarehouse  LocationType = 3
	LocationTypeSite       LocationType = 4
	LocationTypeUnassigned LocationType = 5
)

var (
	LocationTypeMap = map[LocationType]string{
		LocationTypeCage:       "Kandang",
		LocationTypeStore:      "Toko",
		LocationTypeWarehouse:  "Gudang",
		LocationTypeSite:       "Site",
		LocationTypeUnassigned: "Belum Ditempatkan",
	}
)

func (c LocationType) String() string {
	return LocationTypeMap[c]
}

func ValueOfLocationType(value string) LocationType {
	for k, v := range LocationTypeMap {
		if v == value {
			return k
		}
	}
	return LocationTypeUnknown
}

func (c LocationType) IsValid() bool {
	_, ok := LocationTypeMap[c]
	return ok
}
