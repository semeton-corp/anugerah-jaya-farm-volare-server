package util

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func CalculateKPIScoreUserInMonthViaDTO(
	additionalWorkUsers dto.AdditionalWorkUserListPaginationResponse,
	dailyWorkUsers dto.DailyWorkUserListPaginationResponse,
	userPresences dto.PresenceListPaginationResponse,
) (float64, float64, uint64) {
	var (
		totalPresent  uint64  = 0
		totalOvertime float64 = 0
		totalWorkHour float64 = 0

		totalWorkDone uint64 = 0

		presenceScore float64 = 0
		workScore     float64 = 0
	)

	for _, userPresence := range userPresences.Presences {
		if userPresence.Status == enum.PresenceStatusPresent.String() {
			totalPresent++

			if userPresence.EndTime != "" {
				startTime, err := time.Parse("15:04", userPresence.StartTime)
				if err != nil {
					continue
				}

				endTime, err := time.Parse("15:04", userPresence.EndTime)
				if err != nil {
					continue
				}

				diffHours := endTime.Sub(startTime).Hours()
				if diffHours > 8 {
					totalWorkHour += 8.0
				} else {
					totalWorkHour += diffHours
				}

				totalOvertime += userPresence.Overtime
			}
		}
	}

	for _, dailyWorkUser := range dailyWorkUsers.DailyWorkUsers {
		if dailyWorkUser.IsDone {
			totalWorkDone++
		}
	}

	for _, additionalWorkUser := range additionalWorkUsers.AdditionalWorkUsers {
		if additionalWorkUser.IsDone {
			totalWorkDone++
		}
	}

	totalPresence := len(userPresences.Presences)
	if totalPresence > 0 {
		presenceScore = (float64(totalWorkHour) / float64(totalPresence*8)) * 100
	}

	totalWork := len(dailyWorkUsers.DailyWorkUsers) + len(additionalWorkUsers.AdditionalWorkUsers)
	if totalWork > 0 {
		workScore = (float64(totalWorkDone) / float64(totalWork)) * 100
	}

	return presenceScore, workScore, uint64(totalPresence) - totalPresent
}

func CalculateKPIScoreUserInMonthViaEntity(
	additionalWorkUsers []entity.AdditionalWorkUser,
	dailyWorkUsers []entity.DailyWorkUser,
	userPresences []entity.UserPresence,
) (float64, float64, uint64) {
	var (
		totalPresent  uint64  = 0
		totalOvertime float64 = 0
		totalWorkHour float64 = 0

		totalWorkDone uint64 = 0

		presenceScore float64 = 0
		workScore     float64 = 0
	)

	for _, userPresence := range userPresences {
		if userPresence.Status == enum.PresenceStatusPresent {
			totalPresent++

			if !userPresence.EndTime.Time.IsZero() {
				startTime := userPresence.StartTime.Time
				endTime := userPresence.EndTime.Time

				diffHours := endTime.Sub(*startTime).Hours()
				if diffHours > 8 {
					totalWorkHour += 8.0
				} else {
					totalWorkHour += diffHours
				}

				endOfWork := time.Date(
					userPresence.CreatedAt.Year(),
					userPresence.CreatedAt.Month(),
					userPresence.CreatedAt.Day(),
					17, 0, 0, 0, time.Local,
				)

				extraTime := userPresence.EndTime.Time.Sub(endOfWork)
				overtime := float64(0)
				if extraTime > 0 {
					overtime = extraTime.Hours()
				}
				totalOvertime += overtime
			}
		}
	}

	for _, dailyWorkUser := range dailyWorkUsers {
		if dailyWorkUser.IsDone {
			totalWorkDone++
		}
	}

	for _, additionalWorkUser := range additionalWorkUsers {
		if additionalWorkUser.IsDone {
			totalWorkDone++
		}
	}

	totalPresence := len(userPresences)
	if totalPresence > 0 {
		presenceScore = (float64(totalWorkHour) / float64(totalPresence*8)) * 100
	}

	totalWork := len(dailyWorkUsers) + len(additionalWorkUsers)
	if totalWork > 0 {
		workScore = (float64(totalWorkDone) / float64(totalWork)) * 100
	}

	return presenceScore, workScore, uint64(totalPresence) - totalPresent
}
