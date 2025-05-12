package enum

type OverviewGraphTime uint8

const (
	OverviewGraphTimeUnknown   OverviewGraphTime = 0
	OverviewGraphTimeThisWeek  OverviewGraphTime = 1
	OverviewGraphTimeThisMonth OverviewGraphTime = 2
	OverviewGraphTimeThisYear  OverviewGraphTime = 3
)

var (
	OverviewGraphFilterMap = map[OverviewGraphTime]string{
		OverviewGraphTimeThisWeek:  "Minggu Ini",
		OverviewGraphTimeThisMonth: "Bulan Ini",
		OverviewGraphTimeThisYear:  "Tahun Ini",
	}
)

func (c OverviewGraphTime) String() string {
	return OverviewGraphFilterMap[c]
}

func ValueOfOverviewGraphFilter(value string) OverviewGraphTime {
	for k, v := range OverviewGraphFilterMap {
		if v == value {
			return k
		}
	}
	return OverviewGraphTimeUnknown
}

func (c OverviewGraphTime) IsValid() bool {
	_, ok := OverviewGraphFilterMap[c]
	return ok
}
