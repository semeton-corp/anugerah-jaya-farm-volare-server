package service

import (
	"database/sql"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
)

type PresenceService struct {
	log             *zap.Logger
	repository      repository.IPresenceRepository
	locationService ILocationService
}

type IPresenceService interface {
	GetCurrentUserPresence(userId uuid.UUID) (dto.PresenceResponse, error)
	GetUserPresencesByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) (dto.PresenceListPaginationResponse, error)
	UpdateUserPresence(id uint64, request dto.UpdateUserPresenceRequest, updatedBy uuid.UUID) (dto.PresenceResponse, error)

	GetRoleLocationPresenceSummaries() ([]dto.RoleLocationPresenceSummaryResponse, error)
	GetUserPresenceSummaries(filter dto.GetUserPresenceSummaryFilter) ([]dto.UserPresenceSummaryResponse, error)
	GetUserPresenceWorkDetailSummaries(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]dto.UserPresenceWorkDetailSummaryResponse, error)
	ApprovalUserPresence(request dto.ApprovalPresenceRequest, userId uuid.UUID) ([]dto.PresenceResponse, error)
}

func NewPresenceService(log *zap.Logger, repository repository.IPresenceRepository, locationService ILocationService) IPresenceService {
	return &PresenceService{
		log:             log,
		repository:      repository,
		locationService: locationService,
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
		Presences: presenceResponses,
	}

	if filter.Page > 0 {
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
		resp.TotalData = uint64(totalData)
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

func (s *PresenceService) GetRoleLocationPresenceSummaries() ([]dto.RoleLocationPresenceSummaryResponse, error) {
	s.repository.UseTx(false)
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	processPresenceSummaries := func(summaries []entity.LocationPresenceSummary, placeType string) map[uint64]dto.RoleLocationPresenceSummaryResponse {
		result := make(map[uint64]dto.RoleLocationPresenceSummaryResponse)
		for _, e := range summaries {
			summary := result[e.PlaceId]
			if summary.PlaceName == "" {
				summary.RoleId = e.RoleId
				summary.RoleName = e.RoleName
				summary.PlaceId = e.PlaceId
				summary.PlaceType = placeType
				summary.PlaceName = placeType + " " + e.PlaceName
			}

			summary.TotalUser += 1
			switch e.PresenceStatus {
			case enum.PresenceStatusPresent:
				summary.TotalPresentUser += 1
			case enum.PresenceStatusSick:
				summary.TotalSickUser += 1
			case enum.PresenceStatusPermission:
				summary.TotalPermissionUser += 1
			case enum.PresenceStatusAlpha:
				summary.TotalAlphaUser += 1
			}
			result[e.PlaceId] = summary
		}
		return result
	}

	cageSummaries, err := s.repository.GetCageLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date: param.DateParam(today),
	})
	if err != nil {
		s.log.Error("failed to get cage location presence summaries", zap.Error(err))
		return nil, err
	}
	cagePresenceMap := processPresenceSummaries(cageSummaries, enum.LocationWorkTypeCage.String())

	storeSummaries, err := s.repository.GetStoreLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date: param.DateParam(today),
	})
	if err != nil {
		s.log.Error("failed to get store location presence summaries", zap.Error(err))
		return nil, err
	}
	storePresenceMap := processPresenceSummaries(storeSummaries, enum.LocationWorkTypeStore.String())

	warehouseSummaries, err := s.repository.GetWarehouseLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date: param.DateParam(today),
	})
	if err != nil {
		s.log.Error("failed to get warehouse location presence summaries", zap.Error(err))
		return nil, err
	}
	warehousePresenceMap := processPresenceSummaries(warehouseSummaries, enum.LocationWorkTypeWarehouse.String())

	response := make([]dto.RoleLocationPresenceSummaryResponse, 0,
		len(cagePresenceMap)+len(storePresenceMap)+len(warehousePresenceMap))

	for _, v := range cagePresenceMap {
		response = append(response, v)
	}
	for _, v := range storePresenceMap {
		response = append(response, v)
	}
	for _, v := range warehousePresenceMap {
		response = append(response, v)
	}

	return response, nil
}

func (s *PresenceService) GetUserPresenceSummaries(filter dto.GetUserPresenceSummaryFilter) ([]dto.UserPresenceSummaryResponse, error) {
	s.repository.UseTx(false)

	response := make([]dto.UserPresenceSummaryResponse, 0)
	userPresenceSummaries, err := s.repository.GetUserPresenceSummaries(filter)
	if err != nil {
		s.log.Error("failed to get user presence summaries", zap.Error(err))
		return nil, err
	}

	for _, userPresenceSummary := range userPresenceSummaries {
		response = append(response, dto.UserPresenceSummaryResponse{
			UserId:              userPresenceSummary.UserId.String(),
			UserName:            userPresenceSummary.UserName,
			UserPhotoProfile:    userPresenceSummary.UserPhotoProfile,
			UserEmail:           userPresenceSummary.UserEmail,
			RoleName:            userPresenceSummary.RoleName,
			TotalPresentUser:    uint64(userPresenceSummary.TotalPresent),
			TotalSickUser:       uint64(userPresenceSummary.TotalSick),
			TotalPermissionUser: uint64(userPresenceSummary.TotalPermission),
			TotalAlphaUser:      uint64(userPresenceSummary.TotalAlpha),
		})
	}

	return response, nil
}

func (s *PresenceService) GetUserPresenceWorkDetailSummaries(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]dto.UserPresenceWorkDetailSummaryResponse, error) {
	s.repository.UseTx(false)
	filter.Date = param.DateParam(time.Now())
	userPresenceWorkDetailSummaries, err := s.repository.GetUserPresenceWorkDetailSummaries(filter)
	if err != nil {
		s.log.Error("failed to get user presence work detail summaries", zap.Error(err))
		return nil, err
	}

	response := make([]dto.UserPresenceWorkDetailSummaryResponse, 0)

	for _, userPresenceWorkDetailSummary := range userPresenceWorkDetailSummaries {
		workDonePercentage := float64(userPresenceWorkDetailSummary.TotalDoneAdditionalWorkUsers+userPresenceWorkDetailSummary.TotalDoneDailyWorkUsers) / float64(userPresenceWorkDetailSummary.TotalAdditionalWorkUsers+userPresenceWorkDetailSummary.TotalDailyWorkUsers)

		newData := dto.UserPresenceWorkDetailSummaryResponse{
			UserId:             userPresenceWorkDetailSummary.UserId.String(),
			UserName:           userPresenceWorkDetailSummary.UserName,
			UserPhotoProfile:   userPresenceWorkDetailSummary.UserPhotoProfile,
			UserEmail:          userPresenceWorkDetailSummary.UserEmail,
			RoleName:           userPresenceWorkDetailSummary.RoleName,
			PresenceStatus:     userPresenceWorkDetailSummary.PresenceStatus.String(),
			WorkDonePercentage: workDonePercentage,
		}

		if userPresenceWorkDetailSummary.PresenceStartTime.Time != nil {
			newData.ArrivedTime = userPresenceWorkDetailSummary.PresenceStartTime.Time.Format("15:04")
		} else {
			newData.ArrivedTime = "-"
		}

		if userPresenceWorkDetailSummary.PresenceEndTime.Time != nil {
			newData.DepartureTime = userPresenceWorkDetailSummary.PresenceEndTime.Time.Format("15:04")
		} else {
			newData.DepartureTime = "-"
		}

		response = append(response, newData)
	}

	return response, nil
}

func (s *PresenceService) ApprovalUserPresence(request dto.ApprovalPresenceRequest, userId uuid.UUID) ([]dto.PresenceResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	userIds := make([]uuid.UUID, 0)
	for _, userId := range request.AcceptedUserIds {
		userIds = append(userIds, uuid.MustParse(userId))
	}

	for _, userId := range request.RejectedUserIds {
		userIds = append(userIds, uuid.MustParse(userId))
	}

	userPresences, err := s.repository.GetUserPresences(dto.GetUserPresenceFilter{
		UserIds: userIds,
		Date:    param.DateParam(time.Now()),
	})
	if err != nil {
		s.log.Error("failed to get user presences", zap.Error(err))
		return nil, err
	}

	acceptedUserPresenceIds := make([]uint64, 0)
	rejectedUserPresenceIds := make([]uint64, 0)

	for _, userPresence := range userPresences {
		for _, acceptedUserId := range request.AcceptedUserIds {
			if uuid.MustParse(acceptedUserId) == userPresence.UserId {
				userPresence.SubmissionPresenceStatus = enum.SubmissionPresenceStatusAccepted
				acceptedUserPresenceIds = append(acceptedUserPresenceIds, userPresence.Id)
				break
			}
		}

		for _, rejectedUserId := range request.RejectedUserIds {
			if uuid.MustParse(rejectedUserId) == userPresence.UserId {
				userPresence.SubmissionPresenceStatus = enum.SubmissionPresenceStatusRejected
				rejectedUserPresenceIds = append(rejectedUserPresenceIds, userPresence.Id)
				break
			}
		}
	}

	err = s.repository.UpdateSubmissionPresenceStatusUserIds(acceptedUserPresenceIds, enum.SubmissionPresenceStatusAccepted)
	if err != nil {
		s.log.Error("failed update submission presence status by ids", zap.Error(err))
		return nil, err
	}

	err = s.repository.UpdateSubmissionPresenceStatusUserIds(rejectedUserPresenceIds, enum.SubmissionPresenceStatusRejected)
	if err != nil {
		s.log.Error("failed update submission presence status by ids", zap.Error(err))
		return nil, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return nil, err
	}

	userPresenceResponses := make([]dto.PresenceResponse, 0)
	for _, userPresence := range userPresences {
		userPresenceResponses = append(userPresenceResponses, mapper.PresenceToResponse(&userPresence))
	}

	return userPresenceResponses, nil
}
