package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

// Note : without ExtraTime
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
	}

	if !presence.StartTime.IsZero() {
		presenceDto.StartTime = presence.StartTime.Format("15:04")
	}

	if !presence.EndTime.IsZero() {
		presenceDto.EndTime = presence.EndTime.Format("15:04")
	}

	return presenceDto
}

// Note : without ExtraTime
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

	if !presence.StartTime.IsZero() {
		presenceDto.StartTime = presence.StartTime.Format("15:04")
	}

	if !presence.EndTime.IsZero() {
		presenceDto.EndTime = presence.EndTime.Format("15:04")
	}

	return presenceDto
}
