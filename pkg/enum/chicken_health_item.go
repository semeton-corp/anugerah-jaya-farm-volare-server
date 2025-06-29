package enum

type ChickenHealthItemType uint8

const (
	ChickenHealthItemTypeUnknown               ChickenHealthItemType = 0
	ChickenHealthItemTypeMedicine              ChickenHealthItemType = 1
	ChickenHealthProductTypeVaccineConditional ChickenHealthItemType = 2
	ChickenHealthProductTypeVaccineRoutine     ChickenHealthItemType = 3
)

var (
	ChickenHealthItemTypeMap = map[ChickenHealthItemType]string{
		ChickenHealthItemTypeMedicine:              "Obat",
		ChickenHealthProductTypeVaccineConditional: "Vaksin Kondisional",
		ChickenHealthProductTypeVaccineRoutine:     "Vaksin Rutin",
	}
)

func (c ChickenHealthItemType) String() string {
	return ChickenHealthItemTypeMap[c]
}

func ValueOfChickenHealthItemType(value string) ChickenHealthItemType {
	for k, v := range ChickenHealthItemTypeMap {
		if v == value {
			return k
		}
	}
	return ChickenHealthItemTypeUnknown
}

func (c ChickenHealthItemType) IsValid() bool {
	_, ok := ChickenHealthItemTypeMap[c]
	return ok
}
