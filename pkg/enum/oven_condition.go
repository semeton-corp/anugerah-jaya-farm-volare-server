package enum

type OvenCondition uint8

const (
	OvenConditionUnknown  OvenCondition = 0
	OvenConditionOn       OvenCondition = 1
	OvenConditionOff      OvenCondition = 2
	OvenConditionNotInput OvenCondition = 3
)

var (
	OvenConditionMap = map[OvenCondition]string{
		OvenConditionOn:       "Hidup",
		OvenConditionOff:      "Mati",
		OvenConditionNotInput: "-",
	}
)

func (c OvenCondition) String() string {
	return OvenConditionMap[c]
}

func ValueOfOvenCondition(value string) OvenCondition {
	for k, v := range OvenConditionMap {
		if v == value {
			return k
		}
	}
	return OvenConditionUnknown
}

func (c OvenCondition) IsValid() bool {
	_, ok := OvenConditionMap[c]
	return ok
}
