package enum

type ChickenCategory uint8

const (
	ChickenCategoryUnknown  ChickenCategory = 0
	ChickenCategoryDOC      ChickenCategory = 1
	ChickenCategoryGrower   ChickenCategory = 2
	ChickenCategoryPreLayer ChickenCategory = 3
	ChickenCategoryLayer    ChickenCategory = 4
	ChickenCategoryAfkir    ChickenCategory = 5
)

var (
	ChickenCategoryMap = map[ChickenCategory]string{
		ChickenCategoryDOC:      "DOC",
		ChickenCategoryGrower:   "Grower",
		ChickenCategoryPreLayer: "Pre Layer",
		ChickenCategoryLayer:    "Layer",
		ChickenCategoryAfkir:    "Afkir",
	}
)

func (c ChickenCategory) String() string {
	return ChickenCategoryMap[c]
}

func ValueOfChickenCategory(value string) ChickenCategory {
	for k, v := range ChickenCategoryMap {
		if v == value {
			return k
		}
	}
	return ChickenCategoryUnknown
}

func (c ChickenCategory) IsValid() bool {
	_, ok := ChickenCategoryMap[c]
	return ok
}

type ChickenProductivity uint8

const (
	ChickenProductivityUnknown       ChickenProductivity = 0
	ChickenProductivityProductive    ChickenProductivity = 1
	ChickenProductivityAfkir         ChickenProductivity = 2
	ChickenProductivityNotClassified ChickenProductivity = 3
)

var (
	ChickenProductivityMap = map[ChickenProductivity]string{
		ChickenProductivityProductive:    "Produktif",
		ChickenProductivityAfkir:         "Afkir",
		ChickenProductivityNotClassified: "Tidak Diklasifikasikan",
	}
)

func (c ChickenProductivity) String() string {
	return ChickenProductivityMap[c]
}

func ValueOfChickenProductivity(value string) ChickenProductivity {
	for k, v := range ChickenProductivityMap {
		if v == value {
			return k
		}
	}
	return ChickenProductivityUnknown
}

func (c ChickenProductivity) IsValid() bool {
	_, ok := ChickenProductivityMap[c]
	return ok
}
