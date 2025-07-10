package service

import (
	"database/sql"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
)

type PresenceService struct {
	log        *zap.Logger
	repository repository.IPresenceRepository
}

type IPresenceService interface {
	GetCurrentUserPresence(userId uuid.UUID) (dto.PresenceResponse, error)
	GetUserPresencesByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) (dto.PresenceListPaginationResponse, error)
	UpdateUserPresence(id uint64, request dto.UpdateUserPresenceRequest, updatedBy uuid.UUID) (dto.PresenceResponse, error)
}

func NewPresenceService(log *zap.Logger, repository repository.IPresenceRepository) IPresenceService {
	return &PresenceService{
		log:        log,
		repository: repository,
	}
}

func (s *PresenceService) GetCurrentUserPresence(userId uuid.UUID) (dto.PresenceResponse, error) {
	s.repository.UseTx(false)

	userPresence, err := s.repository.GetUserPresenceTodayByUserId(userId)
	if err != nil {
		s.log.Error("failed to get user presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	return mapper.PresenceToResponse(&userPresence), nil
}

func (s *PresenceService) GetUserPresencesByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) (dto.PresenceListPaginationResponse, error) {
	s.repository.UseTx(false)

	userPresence, err := s.repository.GetUserPresencesByUserId(userId, filter)
	if err != nil {
		s.log.Error("failed to get user presences by user id", zap.Error(err))
		return dto.PresenceListPaginationResponse{}, err
	}

	presenceResponses := make([]dto.PresenceListResponse, len(userPresence))
	for i, presence := range userPresence {
		presenceResponses[i] = mapper.PresenceToResponseList(&presence)
	}

	totalData, err := s.repository.CountTotalUserPresenceByUserId(userId, dto.GetPresenceFilter{
		Month: filter.Month,
		Year:  filter.Year,
	})
	if err != nil {
		s.log.Error("failed to count user presence", zap.Error(err))
		return dto.PresenceListPaginationResponse{}, err
	}

	resp := dto.PresenceListPaginationResponse{
		TotalPage: uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData: uint64(totalData),
		Presences: presenceResponses,
	}

	return resp, nil
}

func (s *PresenceService) UpdateUserPresence(id uint64, request dto.UpdateUserPresenceRequest, updatedBy uuid.UUID) (dto.PresenceResponse, error) {
	s.repository.UseTx(false)

	userPresence, err := s.repository.GetUserPresenceById(id)
	if err != nil {
		s.log.Error("failed to get user presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	status := enum.ValueOfPresenceStatus(request.Status)
	if !status.IsValid() {
		s.log.Warn("invalid status presence", zap.String("status", request.Status))
		return dto.PresenceResponse{}, errx.BadRequest("invalid status presence")
	}

	if (userPresence.Status == enum.PresenceStatusPresent && userPresence.EndTime.Time != nil) || userPresence.Status == enum.PresenceStatusSick || userPresence.Status == enum.PresenceStatusPermission {
		return dto.PresenceResponse{}, errx.BadRequest("presence record cannot be updated because it has already been completed or is in a non-present status")
	}

	if status == enum.PresenceStatusPresent && request.StartTime == "" && userPresence.StartTime.Time == nil {
		return dto.PresenceResponse{}, errx.BadRequest("start time is required and end time must be empty when marking presence as present for the first time")
	} else if status == enum.PresenceStatusPresent && request.EndTime == "" && userPresence.EndTime.Time == nil && request.StartTime == "" {
		return dto.PresenceResponse{}, errx.BadRequest("end time is required and start time must be empty when updating presence as present after start time is set")
	} else if (status == enum.PresenceStatusPermission || status == enum.PresenceStatusSick) && request.Evidence == "" && request.Note == "" {
		return dto.PresenceResponse{}, errx.BadRequest("evidence or note is required for permission status")
	}

	switch status {
	case enum.PresenceStatusPresent:
		if !util.IsWithinRadius(userPresence.User.Location.Longitude, userPresence.User.Location.Latitude, request.Longitude, request.Latitude, constant.RadiusPresence) {
			return dto.PresenceResponse{}, errx.BadRequest("location is not within the allowed radius")
		}

		userPresence.Status = status

		var timez datatype.TimeOnly
		if request.StartTime != "" {
			timeParsed, err := time.Parse("15:04", request.StartTime)
			if err != nil {
				return dto.PresenceResponse{}, err
			}

			timez = datatype.TimeOnly{Time: &timeParsed}
			userPresence.StartTime = timez
		} else if request.EndTime != "" {
			timeParsed, err := time.Parse("15:04", request.EndTime)
			if err != nil {
				return dto.PresenceResponse{}, err
			}

			timez = datatype.TimeOnly{Time: &timeParsed}
			userPresence.EndTime = timez
		}
	case enum.PresenceStatusSick, enum.PresenceStatusPermission:
		userPresence.Status = status
		userPresence.Evidence = sql.NullString{String: request.Evidence, Valid: true}
		userPresence.Note = sql.NullString{String: request.Note, Valid: true}
		userPresence.SubmissionPresenceStatus = enum.SubmissionPresenceStatusPending
	}

	userPresence.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateUserPresence(&userPresence)
	if err != nil {
		s.log.Error("failed to update user presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	return mapper.PresenceToResponse(&userPresence), nil
}
