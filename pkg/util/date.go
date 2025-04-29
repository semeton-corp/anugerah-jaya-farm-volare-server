package util

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func GetStartDayAndEndDayInWeek(date time.Time) (time.Time, time.Time) {
	startDate := date.AddDate(0, 0, -int(date.Weekday()))
	endDate := date.AddDate(0, 0, 6-int(date.Weekday()))
	return startDate, endDate
}

func GetStartDayAndEndDayInMonth(date time.Time) (time.Time, time.Time) {
	startDate := date.AddDate(0, 0, -int(date.Day())+1)
	endDate := date.AddDate(0, 1, -int(date.Day()))
	return startDate, endDate
}

func GetStartedDayAndEndDayOfMonth(month time.Month, year int) (time.Time, time.Time) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.UTC)
	endDate := startDate.AddDate(0, 1, -1)
	return startDate, endDate
}

func GetStartDayAndEndDayByPresenceFilter(presenceStatus enum.PresenceFilter) (time.Time, time.Time) {
	switch presenceStatus {
	case enum.PresenceFilterThisWeek:
		return GetStartDayAndEndDayInWeek(time.Now())
	case enum.PresenceFilterJanuary:
		return GetStartedDayAndEndDayOfMonth(time.January, time.Now().Year())
	default:
		return time.Time{}, time.Time{}
	}
}
