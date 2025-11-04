package mapper

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func PresenceToResponse(presence *entity.UserPresence) dto.PresenceResponse {
	presenceDto := dto.PresenceResponse{
		Id: presence.Id,
		User: dto.UserResponse{
			Id:   presence.User.Id.String(),
			Name: presence.User.Name,
			Role: dto.RoleResponse{
				Id:   presence.User.Role.Id,
				Name: presence.User.Role.Name,
			},
		},
		Date:             presence.CreatedAt.Format("02-01-2006"),
		CreatedAt:        presence.CreatedAt,
		Status:           presence.Status.String(),
		SubmissionStatus: presence.SubmissionPresenceStatus.String(),
	}

	if presence.StartTime.Time != nil {
		presenceDto.StartTime = presence.StartTime.Time.Format("15:04")
	} else {
		presenceDto.StartTime = "-"
	}

	if presence.EndTime.Time != nil {
		presenceDto.EndTime = presence.EndTime.Time.Format("15:04")
	} else {
		presenceDto.EndTime = "-"
	}

	if presence.EndTime.Time != nil {
		extraTime := presence.EndTime.Time.Sub(time.Date(presence.CreatedAt.Year(), presence.CreatedAt.Month(), presence.CreatedAt.Day(), 5, 0, 0, 0, time.Local))
		if extraTime > 0 {
			presenceDto.Overtime = extraTime.Hours()
		} else {
			presenceDto.Overtime = 0
		}
	}

	if presence.Evidence.Valid {
		presenceDto.Evidence = presence.Evidence.String
	} else {
		presenceDto.Evidence = "-"
	}

	if presence.Note.Valid {
		presenceDto.Note = presence.Note.String
	} else {
		presenceDto.Note = "-"
	}

	return presenceDto
}

func PresenceToResponseList(presence *entity.UserPresence) dto.PresenceListResponse {
	presenceDto := dto.PresenceListResponse{
		Id: presence.Id,
		User: dto.UserResponse{
			Id:   presence.User.Id.String(),
			Name: presence.User.Name,
			Role: dto.RoleResponse{
				Id:   presence.User.Role.Id,
				Name: presence.User.Role.Name,
			},
		},
		Date:             presence.CreatedAt.Format("02-01-2006"),
		Status:           presence.Status.String(),
		SubmissionStatus: presence.SubmissionPresenceStatus.String(),
		CreatedAt:        presence.CreatedAt,
	}

	if presence.StartTime.Time != nil {
		presenceDto.StartTime = presence.StartTime.Time.Format("15:04")
	}

	if presence.EndTime.Time != nil {
		presenceDto.EndTime = presence.EndTime.Time.Format("15:04")
	}

	if presence.EndTime.Time != nil {
		endOfWork := time.Date(
			presence.CreatedAt.Year(),
			presence.CreatedAt.Month(),
			presence.CreatedAt.Day(),
			17, 0, 0, 0, time.Local,
		)

		extraTime := presence.EndTime.Time.Sub(endOfWork)
		if extraTime > 0 {
			presenceDto.Overtime = extraTime.Hours()
		} else {
			presenceDto.Overtime = 0
		}
	}

	if presence.Evidence.Valid {
		presenceDto.Evidence = presence.Evidence.String
	} else {
		presenceDto.Evidence = "-"
	}

	if presence.Note.Valid {
		presenceDto.Note = presence.Note.String
	} else {
		presenceDto.Note = "-"
	}

	return presenceDto
}
