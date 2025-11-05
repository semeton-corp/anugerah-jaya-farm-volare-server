package util

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func GetStartDayAndEndDayByMonthFilter(month enum.Month, year int) (time.Time, time.Time) {
	switch month {
	case enum.MonthJanuary:
		return GetStartDateAndEndDateInMonth(year, time.January)
	case enum.MonthFebruary:
		return GetStartDateAndEndDateInMonth(year, time.February)
	case enum.MonthMarch:
		return GetStartDateAndEndDateInMonth(year, time.March)
	case enum.MonthApril:
		return GetStartDateAndEndDateInMonth(year, time.April)
	case enum.MonthMay:
		return GetStartDateAndEndDateInMonth(year, time.May)
	case enum.MonthJune:
		return GetStartDateAndEndDateInMonth(year, time.June)
	case enum.MonthJuly:
		return GetStartDateAndEndDateInMonth(year, time.July)
	case enum.MonthAugust:
		return GetStartDateAndEndDateInMonth(year, time.August)
	case enum.MonthSeptember:
		return GetStartDateAndEndDateInMonth(year, time.September)
	case enum.MonthOctober:
		return GetStartDateAndEndDateInMonth(year, time.October)
	case enum.MonthNovember:
		return GetStartDateAndEndDateInMonth(year, time.November)
	case enum.MonthDecember:
		return GetStartDateAndEndDateInMonth(year, time.December)
	default:
		return time.Time{}, time.Time{}
	}
}

type DateRange struct {
	StartDate time.Time
	EndDate   time.Time
	TotalDays int
}

func GetFourWeekRanges(year int, month time.Month) map[int]DateRange {
	weeks := make(map[int]DateRange)

	startOfMonth := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	nextMonth := startOfMonth.AddDate(0, 1, 0)
	daysInMonth := int(nextMonth.Sub(startOfMonth).Hours() / 24)

	baseDaysPerWeek := daysInMonth / 4
	extraDays := daysInMonth % 4

	currentDay := startOfMonth

	for i := 1; i <= 4; i++ {
		daysThisWeek := baseDaysPerWeek
		if i <= extraDays {
			daysThisWeek++
		}

		endDate := currentDay.AddDate(0, 0, daysThisWeek-1)
		endDate = time.Date(
			endDate.Year(), endDate.Month(), endDate.Day(),
			23, 59, 59, 0, endDate.Location(),
		)

		weeks[i] = DateRange{
			StartDate: currentDay,
			EndDate:   endDate,
			TotalDays: daysThisWeek,
		}

		currentDay = endDate.AddDate(0, 0, 1).Truncate(24 * time.Hour)
	}

	return weeks
}

func GetTwelveMonthRanges(year int) map[int]DateRange {
	months := make(map[int]DateRange)

	for month := time.January; month <= time.December; month++ {
		start := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		end := start.AddDate(0, 1, -1)
		end = time.Date(end.Year(), end.Month(), end.Day(), 23, 59, 59, 0, end.Location())

		months[int(month)] = DateRange{
			StartDate: start,
			EndDate:   end,
			TotalDays: end.Day(),
		}
	}

	return months
}

func TotalDaysInMonth(year int, month time.Month) uint64 {
	daysInMonth := time.Date(year, month+1, 0, 0, 0, 0, 0, time.Local).Day()
	return uint64(daysInMonth)
}

func FindWeek(t time.Time, weeks map[int]DateRange) int {
	for i, week := range weeks {
		if !t.Before(week.StartDate) && !t.After(week.EndDate) {
			return i
		}
	}
	return 0
}

func FindMonth(t time.Time, months map[int]DateRange) int {
	for i, month := range months {
		if !t.Before(month.StartDate) && !t.After(month.EndDate) {
			return i
		}
	}
	return 0
}

func GetStartDateAndEndDateInMonth(year int, month time.Month) (time.Time, time.Time) {
	startDate := time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(0, 1, -1)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	return startDate, endDate
}

func GetStartDateAndEndDateInYear(year int) (time.Time, time.Time) {
	startDate := time.Date(year, time.January, 1, 0, 0, 0, 0, time.Local)
	endDate := startDate.AddDate(1, 0, -1)
	endDate = time.Date(endDate.Year(), endDate.Month(), endDate.Day(), 23, 59, 59, 0, endDate.Location())

	return startDate, endDate
}
func IsSameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
}

func IndoMonthName(month int) string {
	switch month {
	case 1:
		return "Januari"
	case 2:
		return "Februari"
	case 3:
		return "Maret"
	case 4:
		return "April"
	case 5:
		return "Mei"
	case 6:
		return "Juni"
	case 7:
		return "Juli"
	case 8:
		return "Agustus"
	case 9:
		return "September"
	case 10:
		return "Oktober"
	case 11:
		return "November"
	case 12:
		return "Desember"
	default:
		return "-"
	}
}

var (
	MapEngMonthToIndoMonth = map[string]string{
		"January":   "Januari",
		"February":  "Februari",
		"March":     "Maret",
		"April":     "April",
		"May":       "Mei",
		"June":      "Juni",
		"July":      "Juli",
		"August":    "Agustus",
		"September": "September",
		"October":   "Oktober",
		"November":  "November",
		"December":  "Desember",
	}
)
