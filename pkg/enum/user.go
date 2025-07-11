package enum

type SalaryInterval uint8

const (
	SalaryIntervalUnknown SalaryInterval = 0
	SalaryIntervalMonthly SalaryInterval = 1
	SalaryIntervalDaily   SalaryInterval = 2
)

var (
	SalaryIntervalMap = map[SalaryInterval]string{
		SalaryIntervalMonthly: "Bulanan",
		SalaryIntervalDaily:   "Harian",
	}
)

func (c SalaryInterval) String() string {
	return SalaryIntervalMap[c]
}

func ValueOfSalaryInterval(value string) SalaryInterval {
	for k, v := range SalaryIntervalMap {
		if v == value {
			return k
		}
	}
	return SalaryIntervalUnknown
}

func (c SalaryInterval) IsValid() bool {
	_, ok := SalaryIntervalMap[c]
	return ok
}
