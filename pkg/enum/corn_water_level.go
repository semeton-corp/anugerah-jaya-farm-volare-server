package enum

type CornWaterLevel uint8

const (
	CornWaterLevelUnknown         CornWaterLevel = 0
	CornWaterLevelLessThanEqual16 CornWaterLevel = 1
	CornWaterLevelMoreThan16      CornWaterLevel = 2
)

var (
	CornWaterLevelMap = map[CornWaterLevel]string{
		CornWaterLevelLessThanEqual16: "<= 16%",
		CornWaterLevelMoreThan16:      "> 16%",
	}
)

func (c CornWaterLevel) String() string {
	return CornWaterLevelMap[c]
}

func ValueOfCornWaterLevel(value string) CornWaterLevel {
	for k, v := range CornWaterLevelMap {
		if v == value {
			return k
		}
	}
	return CornWaterLevelUnknown
}

func (c CornWaterLevel) IsValid() bool {
	_, ok := CornWaterLevelMap[c]
	return ok
}
