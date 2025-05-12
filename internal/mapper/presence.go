package mapper

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func PresenceToResponse(presence *entity.StaffPresence) dto.PresenceResponse {
	presenceDto := dto.PresenceResponse{
		Id: presence.Id,
		Staff: dto.StaffResponse{
			Id:   presence.Staff.Id.String(),
			Name: presence.Staff.Name,
			Role: dto.RoleResponse{
				Id:   presence.Staff.Account.Role.Id,
				Name: presence.Staff.Account.Role.Name,
			},
		},
		Date:      presence.CreatedAt.Format("02 Januari 2006"),
		IsPresent: presence.IsPresent,
		CreatedAt: presence.CreatedAt,
	}

	if !presence.StartTime.Time.IsZero() {
		presenceDto.StartTime = presence.StartTime.Time.Format("15:04")
	} else {
		presenceDto.StartTime = "-"
	}

	if !presence.EndTime.Time.IsZero() {
		presenceDto.EndTime = presence.EndTime.Time.Format("15:04")
	} else {
		presenceDto.EndTime = "-"
	}

	if !presence.EndTime.Time.IsZero() {
		extraTime := presence.EndTime.Time.Sub(time.Date(presence.CreatedAt.Year(), presence.CreatedAt.Month(), presence.CreatedAt.Day(), 5, 0, 0, 0, time.Local))
		if extraTime > 0 {
			presenceDto.Overtime = extraTime.Hours()
		} else {
			presenceDto.Overtime = 0
		}
	}

	return presenceDto
}

func PresenceToResponseList(presence *entity.StaffPresence) dto.PresenceListResponse {
	presenceDto := dto.PresenceListResponse{
		Id: presence.Id,
		Staff: dto.StaffResponse{
			Id:   presence.Staff.Id.String(),
			Name: presence.Staff.Name,
			Role: dto.RoleResponse{
				Id:   presence.Staff.Account.Role.Id,
				Name: presence.Staff.Account.Role.Name,
			},
		},
		Date:      presence.CreatedAt.Format("02 Januari 2006"),
		IsPresent: presence.IsPresent,
	}

	if !presence.StartTime.Time.IsZero() {
		presenceDto.StartTime = presence.StartTime.Time.Format("15:04")
	}

	if !presence.EndTime.Time.IsZero() {
		presenceDto.EndTime = presence.EndTime.Time.Format("15:04")
	}

	if !presence.EndTime.Time.IsZero() {
		extraTime := presence.EndTime.Time.Sub(time.Date(presence.CreatedAt.Year(), presence.CreatedAt.Month(), presence.CreatedAt.Day(), 5, 0, 0, 0, time.Local))
		if extraTime > 0 {
			presenceDto.Overtime = extraTime.Hours()
		} else {
			presenceDto.Overtime = 0
		}
	}

	return presenceDto
}
