package enum

type SupplierType uint8

const (
	SupplierTypeUnknown    SupplierType = 0
	SupplierTypeItem       SupplierType = 1
	SupplierTypeDOCChicken SupplierType = 2
)

var (
	SupplierTypeMap = map[SupplierType]string{
		SupplierTypeItem:       "Barang",
		SupplierTypeDOCChicken: "Ayam DOC",
	}
)

func (c SupplierType) String() string {
	return SupplierTypeMap[c]
}

func ValueOfSupplierType(value string) SupplierType {
	for k, v := range SupplierTypeMap {
		if v == value {
			return k
		}
	}
	return SupplierTypeUnknown
}

func (c SupplierType) IsValid() bool {
	_, ok := SupplierTypeMap[c]
	return ok
}
