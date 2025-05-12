package enum

type Month uint8

const (
	MonthUnknown Month = iota
	MonthJanuary
	MonthFebruary
	MonthMarch
	MonthApril
	MonthMay
	MonthJune
	MonthJuly
	MonthAugust
	MonthSeptember
	MonthOctober
	MonthNovember
	MonthDecember
)

var (
	MonthMap = map[Month]string{
		MonthJanuary:   "Januari",
		MonthFebruary:  "Februari",
		MonthMarch:     "Maret",
		MonthApril:     "April",
		MonthMay:       "Mei",
		MonthJune:      "Juni",
		MonthJuly:      "Juli",
		MonthAugust:    "Agustus",
		MonthSeptember: "September",
		MonthOctober:   "Oktober",
		MonthNovember:  "November",
		MonthDecember:  "Desember",
	}
)

func (c Month) String() string {
	return MonthMap[c]
}

func ValueOfMonth(value string) Month {
	for k, v := range MonthMap {
		if v == value {
			return k
		}
	}
	return MonthUnknown
}

func (c Month) IsValid() bool {
	_, ok := MonthMap[c]
	return ok
}
