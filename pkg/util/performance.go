package util

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func CalculateKPIScoreUserInMonth(additionalWorkUsers dto.AdditionalWorkUserListPaginationResponse, dailyWorkUsers dto.DailyWorkUserListPaginationResponse, userPresences dto.PresenceListPaginationResponse) (float64, float64) {
	var totalPresent uint64 = 0
	var totalOvertime float64 = 0
	var totalWorkHour float64 = 0
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

	var totalWorkDone uint64 = 0
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

	var presenceScore float64
	if len(userPresences.Presences) == 0 {
		presenceScore = 0
	} else {
		presenceScore = float64(totalWorkHour) / float64(len(userPresences.Presences)*8) * 100
	}

	var workScore float64
	if len(dailyWorkUsers.DailyWorkUsers)+len(additionalWorkUsers.AdditionalWorkUsers) == 0 {
		workScore = 0
	} else {
		workScore = float64(totalWorkDone) / float64(len(dailyWorkUsers.DailyWorkUsers)+len(additionalWorkUsers.AdditionalWorkUsers)) * 100
	}

	return presenceScore, workScore
}
