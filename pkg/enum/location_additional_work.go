package enum

type LocationAddionalWork uint8

const (
	LocationAddionalWorkUnknown   LocationAddionalWork = 0
	LocationAddionalWorkCage      LocationAddionalWork = 1
	LocationAddionalWorkStore     LocationAddionalWork = 2
	LocationAddionalWorkWarehouse LocationAddionalWork = 3
)

var (
	LocationAddionalWorkMap = map[LocationAddionalWork]string{
		LocationAddionalWorkCage:      "Kandang",
		LocationAddionalWorkStore:     "Toko",
		LocationAddionalWorkWarehouse: "Gudang",
	}
)

func (c LocationAddionalWork) String() string {
	return LocationAddionalWorkMap[c]
}

func ValueOfLocationAddionalWork(value string) LocationAddionalWork {
	for k, v := range LocationAddionalWorkMap {
		if v == value {
			return k
		}
	}
	return LocationAddionalWorkUnknown
}

func (c LocationAddionalWork) IsValid() bool {
	_, ok := LocationAddionalWorkMap[c]
	return ok
}
