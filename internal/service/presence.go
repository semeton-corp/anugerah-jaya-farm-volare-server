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
	UpdateUserPresence(id uint64, request dto.UpdateUserPresenceRequest, userId uuid.UUID) (dto.PresenceResponse, error)

	GetRoleLocationPresenceSummaries(filter dto.RoleLocationPresenceSummaryFilter) ([]dto.RoleLocationPresenceSummaryResponse, error)
	GetUserPresenceSummaries(filter dto.GetUserPresenceSummaryFilter) ([]dto.UserPresenceSummaryResponse, error)
	GetUserPresenceWorkSummaries(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]dto.UserPresenceWorkDetailSummaryResponse, error)
	ApprovalUserPresence(request dto.ApprovalPresenceRequest, userId uuid.UUID) ([]dto.PresenceResponse, error)

	GetUserPresencePending(filter dto.GetUserPresencePendingFilter) ([]dto.UserPresencePendingResponse, error)
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
		resp.TotalData = uint64(totalData)
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (s *PresenceService) UpdateUserPresence(id uint64, request dto.UpdateUserPresenceRequest, userId uuid.UUID) (dto.PresenceResponse, error) {
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
		returnTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 16, 0, 0, 0, time.Local)
		if time.Now().After(returnTime) && request.StartTime != "" && userPresence.StartTime.Time == nil {
			return dto.PresenceResponse{}, errx.BadRequest("can't presence start more than 16.00 PM")
		}

		if !util.IsWithinRadius(userPresence.User.Location.Longitude, userPresence.User.Location.Latitude, request.Longitude, request.Latitude, constant.RadiusPresence) {
			return dto.PresenceResponse{}, errx.BadRequest("location is not within the allowed radius")
		}

		userPresence.Status = status

		var timez datatype.TimeOnly
		if request.StartTime != "" {
			timeParsed, err := time.Parse("15:04", request.StartTime)
			if err != nil {
				return dto.PresenceResponse{}, errx.BadRequest("invalid time format")
			}

			currTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), timeParsed.Hour(), timeParsed.Minute(), timeParsed.Second(), timeParsed.Nanosecond(), time.Local)
			minPresenceTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 5, 0, 0, 0, time.Local)
			maxPresenceTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 8, 0, 0, 0, time.Local)

			if currTime.Before(minPresenceTime) || currTime.After(maxPresenceTime) {
				return dto.PresenceResponse{}, errx.BadRequest("presence time must be between 05.00 AM and 08.00 AM")
			}

			timez = datatype.TimeOnly{Time: &timeParsed}
			userPresence.StartTime = timez
		} else if request.EndTime != "" {
			timeParsed, err := time.Parse("15:04", request.EndTime)
			if err != nil {
				return dto.PresenceResponse{}, errx.BadRequest("invalid time format")
			}

			timez = datatype.TimeOnly{Time: &timeParsed}
			userPresence.EndTime = timez
		}
	case enum.PresenceStatusSick, enum.PresenceStatusPermission:
		currTime := time.Now()
		maxPresenceSickOrPermissionTime := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 9, 0, 0, 0, time.Local)

		if currTime.After(maxPresenceSickOrPermissionTime) {
			return dto.PresenceResponse{}, errx.BadRequest("presence time for sick or permission must be before 09.00 AM")
		}

		userPresence.Status = status
		userPresence.Evidence = sql.NullString{String: request.Evidence, Valid: true}
		userPresence.Note = sql.NullString{String: request.Note, Valid: true}
		userPresence.SubmissionPresenceStatus = enum.SubmissionPresenceStatusPending
	}

	userPresence.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateUserPresence(&userPresence)
	if err != nil {
		s.log.Error("failed to update user presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	return mapper.PresenceToResponse(&userPresence), nil
}

func (s *PresenceService) GetRoleLocationPresenceSummaries(filter dto.RoleLocationPresenceSummaryFilter) ([]dto.RoleLocationPresenceSummaryResponse, error) {
	s.repository.UseTx(false)
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	processPresenceSummaries := func(summaries []entity.LocationPresenceSummary, placeType string) map[uint64]dto.RoleLocationPresenceSummaryResponse {
		result := make(map[uint64]dto.RoleLocationPresenceSummaryResponse)
		isSickUserPendingExist := false
		isPermissionUserPendingExist := false

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

			if e.PresenceStatus == enum.PresenceStatusSick && e.SubmissionPresenceStatus == enum.SubmissionPresenceStatusPending && !isSickUserPendingExist {
				isSickUserPendingExist = true
				summary.IsSickUserPendingExist = true
			} else if e.PresenceStatus == enum.PresenceStatusPermission && e.SubmissionPresenceStatus == enum.SubmissionPresenceStatusPending && !isPermissionUserPendingExist {
				isPermissionUserPendingExist = true
				summary.IsPermissionUserPendingExist = true
			}

			result[e.PlaceId] = summary
		}

		return result
	}

	cageSummaries, err := s.repository.GetLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date:         param.DateParam(today),
		LocationType: param.LocationTypeParam(enum.LocationTypeCage),
		LocationId:   filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed to get cage location presence summaries", zap.Error(err))
		return nil, err
	}

	cageStaffSummaries := make([]entity.LocationPresenceSummary, 0)
	eggStaffSummaries := make([]entity.LocationPresenceSummary, 0)

	for _, cageSummary := range cageSummaries {
		switch cageSummary.RoleName {
		case constant.RolePekerjaKandang:
			cageStaffSummaries = append(cageStaffSummaries, cageSummary)
		case constant.RolePekerjaTelur:
			eggStaffSummaries = append(eggStaffSummaries, cageSummary)
		}
	}

	cageStaffPresenceMap := processPresenceSummaries(cageStaffSummaries, enum.LocationTypeCage.String())
	eggStaffPresenceMap := processPresenceSummaries(eggStaffSummaries, enum.LocationTypeCage.String())

	headCageSummaries, err := s.repository.GetLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date:         param.DateParam(today),
		LocationType: param.LocationTypeParam(enum.LocationTypeSite),
		LocationId:   filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed to get cage location presence summaries", zap.Error(err))
		return nil, err
	}
	headPresenceMap := processPresenceSummaries(headCageSummaries, enum.LocationTypeSite.String())

	storeSummaries, err := s.repository.GetLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date:         param.DateParam(today),
		LocationType: param.LocationTypeParam(enum.LocationTypeStore),
		LocationId:   filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed to get store location presence summaries", zap.Error(err))
		return nil, err
	}
	storePresenceMap := processPresenceSummaries(storeSummaries, enum.LocationTypeStore.String())

	warehouseSummaries, err := s.repository.GetLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date:         param.DateParam(today),
		LocationType: param.LocationTypeParam(enum.LocationTypeWarehouse),
		LocationId:   filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed to get warehouse location presence summaries", zap.Error(err))
		return nil, err
	}
	warehousePresenceMap := processPresenceSummaries(warehouseSummaries, enum.LocationTypeWarehouse.String())

	unassignedSummaries, err := s.repository.GetLocationPresenceSummaries(dto.GetLocationPresenceSummaryFilter{
		Date:         param.DateParam(today),
		LocationType: param.LocationTypeParam(enum.LocationTypeUnassigned),
		LocationId:   filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed to get cage location presence summaries", zap.Error(err))
		return nil, err
	}

	cageStaffSummariesUnassigned := make([]entity.LocationPresenceSummary, 0)
	eggStaffSummariesUnassigned := make([]entity.LocationPresenceSummary, 0)
	storeStaffSummariesUnassigned := make([]entity.LocationPresenceSummary, 0)
	warehouseStaffSummariesUnassigned := make([]entity.LocationPresenceSummary, 0)

	for _, unassignedSummary := range unassignedSummaries {
		switch unassignedSummary.RoleName {
		case constant.RolePekerjaKandang:
			cageStaffSummariesUnassigned = append(cageStaffSummariesUnassigned, unassignedSummary)
		case constant.RolePekerjaTelur:
			eggStaffSummariesUnassigned = append(eggStaffSummariesUnassigned, unassignedSummary)
		case constant.RolePekerjaToko:
			storeStaffSummariesUnassigned = append(storeStaffSummariesUnassigned, unassignedSummary)
		case constant.RolePekerjaGudang:
			warehouseStaffSummariesUnassigned = append(warehouseStaffSummariesUnassigned, unassignedSummary)
		}
	}

	cageStaffPresenceUnassignedMap := processPresenceSummaries(cageStaffSummariesUnassigned, enum.LocationTypeUnassigned.String())
	eggStaffPresenceUnassignedMap := processPresenceSummaries(eggStaffSummariesUnassigned, enum.LocationTypeUnassigned.String())
	storeStaffPresenceUnassignedMap := processPresenceSummaries(storeStaffSummariesUnassigned, enum.LocationTypeUnassigned.String())
	warehouseStaffPresenceUnassignedMap := processPresenceSummaries(warehouseStaffSummariesUnassigned, enum.LocationTypeUnassigned.String())

	responses := make([]dto.RoleLocationPresenceSummaryResponse, 0)
	for _, v := range cageStaffPresenceUnassignedMap {
		responses = append(responses, v)
	}
	for _, v := range eggStaffPresenceUnassignedMap {
		responses = append(responses, v)
	}
	for _, v := range storeStaffPresenceUnassignedMap {
		responses = append(responses, v)
	}
	for _, v := range warehouseStaffPresenceUnassignedMap {
		responses = append(responses, v)
	}
	for _, v := range eggStaffPresenceMap {
		responses = append(responses, v)
	}
	for _, v := range cageStaffPresenceMap {
		responses = append(responses, v)
	}
	for _, v := range storePresenceMap {
		responses = append(responses, v)
	}
	for _, v := range warehousePresenceMap {
		responses = append(responses, v)
	}
	for _, v := range headPresenceMap {
		responses = append(responses, v)
	}

	return responses, nil
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

func (s *PresenceService) GetUserPresenceWorkSummaries(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]dto.UserPresenceWorkDetailSummaryResponse, error) {
	s.repository.UseTx(false)

	userPresenceWorkDetailSummaries, err := s.repository.GetUserPresenceWorkDetailSummaries(filter)
	if err != nil {
		s.log.Error("failed to get user presence work detail summaries", zap.Error(err))
		return nil, err
	}

	response := make([]dto.UserPresenceWorkDetailSummaryResponse, 0)

	for _, userPresenceWorkDetailSummary := range userPresenceWorkDetailSummaries {
		totalWork := float64(userPresenceWorkDetailSummary.TotalAdditionalWorkUsers + userPresenceWorkDetailSummary.TotalDailyWorkUsers)
		workDonePercentage := float64(0)
		if totalWork != 0 {
			workDonePercentage = (float64(userPresenceWorkDetailSummary.TotalDoneAdditionalWorkUsers+userPresenceWorkDetailSummary.TotalDoneDailyWorkUsers) / float64(userPresenceWorkDetailSummary.TotalAdditionalWorkUsers+userPresenceWorkDetailSummary.TotalDailyWorkUsers)) * 100
		}

		newData := dto.UserPresenceWorkDetailSummaryResponse{
			UserId:                   userPresenceWorkDetailSummary.UserId.String(),
			UserName:                 userPresenceWorkDetailSummary.UserName,
			UserPhotoProfile:         userPresenceWorkDetailSummary.UserPhotoProfile,
			UserEmail:                userPresenceWorkDetailSummary.UserEmail,
			RoleName:                 userPresenceWorkDetailSummary.RoleName,
			PresenceStatus:           userPresenceWorkDetailSummary.PresenceStatus.String(),
			SubmissionPresenceStatus: userPresenceWorkDetailSummary.SubmissionPresenceStatus.String(),
			WorkDonePercentage:       workDonePercentage,
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

	switch request.ApprovalStatus {
	case constant.ApprovalStatusAccepted:
		err := s.repository.UpdateSubmissionPresenceStatusUserByIds(request.UserPresenceIds, enum.SubmissionPresenceStatusAccepted)
		if err != nil {
			s.log.Error("failed update submission presence status by ids", zap.Error(err))
			return nil, err
		}
	case constant.ApprovalStatusRejected:
		err := s.repository.UpdatePresenceStatusAndSubmissionPresenceStatusByUserIds(request.UserPresenceIds, enum.PresenceStatusAlpha, enum.SubmissionPresenceStatusRejected)
		if err != nil {
			s.log.Error("failed update presence status and submission presence status by ids", zap.Error(err))
			return nil, err
		}
	}

	err := s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return nil, err
	}

	userPresences, err := s.repository.GetUserPresences(dto.GetUserPresenceFilter{
		Ids: request.UserPresenceIds,
	})
	if err != nil {
		s.log.Error("failed get user presence", zap.Error(err))
		return nil, err
	}

	userPresenceResponses := make([]dto.PresenceResponse, 0)
	for _, userPresence := range userPresences {
		userPresenceResponses = append(userPresenceResponses, mapper.PresenceToResponse(&userPresence))
	}

	return userPresenceResponses, nil
}

func (s *PresenceService) GetUserPresencePending(filter dto.GetUserPresencePendingFilter) ([]dto.UserPresencePendingResponse, error) {
	s.repository.UseTx(false)
	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	responses := make([]dto.UserPresencePendingResponse, 0)
	if filter.LocationType.Value() == enum.LocationTypeCage {
		userPresences, err := s.repository.GetLocationUserPresence(dto.GetLocationUserPresenceFilter{
			PlaceId:                  filter.PlaceId,
			RoleId:                   filter.RoleId,
			PresenceStatus:           filter.PresenceStatus,
			SubmissionPresenceStatus: filter.SubmissionPresenceStatus,
			Date:                     param.DateParam(today),
			LocationType:             param.LocationTypeParam(enum.LocationTypeCage),
		})
		if err != nil {
			s.log.Error("failed get location user presence", zap.Error(err))
			return nil, err
		}

		for _, e := range userPresences {
			responses = append(responses, dto.UserPresencePendingResponse{
				Id:           e.Id,
				Date:         e.CreatedAt.Format("02-01-2006"),
				Name:         e.User.Name,
				Status:       e.Status.String(),
				Evidence:     e.Evidence.String,
				Note:         e.Note.String,
				LocationType: enum.LocationTypeCage.String(),
			})
		}

	} else if filter.LocationType.Value() == enum.LocationTypeSite {
		userPresences, err := s.repository.GetLocationUserPresence(dto.GetLocationUserPresenceFilter{
			PlaceId:                  filter.PlaceId,
			RoleId:                   filter.RoleId,
			PresenceStatus:           filter.PresenceStatus,
			SubmissionPresenceStatus: filter.SubmissionPresenceStatus,
			Date:                     param.DateParam(today),
			LocationType:             param.LocationTypeParam(enum.LocationTypeSite),
		})
		if err != nil {
			s.log.Error("failed get location user presence", zap.Error(err))
			return nil, err
		}

		for _, e := range userPresences {
			responses = append(responses, dto.UserPresencePendingResponse{
				Id:           e.Id,
				Date:         e.CreatedAt.Format("02-01-2006"),
				Name:         e.User.Name,
				Status:       e.Status.String(),
				Evidence:     e.Evidence.String,
				Note:         e.Note.String,
				LocationType: enum.LocationTypeSite.String(),
			})
		}

	} else if filter.LocationType.Value() == enum.LocationTypeStore {
		userPresences, err := s.repository.GetLocationUserPresence(dto.GetLocationUserPresenceFilter{
			PlaceId:                  filter.PlaceId,
			RoleId:                   filter.RoleId,
			PresenceStatus:           filter.PresenceStatus,
			SubmissionPresenceStatus: filter.SubmissionPresenceStatus,
			Date:                     param.DateParam(today),
			LocationType:             param.LocationTypeParam(enum.LocationTypeStore),
		})
		if err != nil {
			s.log.Error("failed get location user presence", zap.Error(err))
			return nil, err
		}

		for _, e := range userPresences {
			responses = append(responses, dto.UserPresencePendingResponse{
				Id:           e.Id,
				Date:         e.CreatedAt.Format("02-01-2006"),
				Name:         e.User.Name,
				Status:       e.Status.String(),
				Evidence:     e.Evidence.String,
				Note:         e.Note.String,
				LocationType: enum.LocationTypeStore.String(),
			})
		}

	} else if filter.LocationType.Value() == enum.LocationTypeWarehouse {
		userPresences, err := s.repository.GetLocationUserPresence(dto.GetLocationUserPresenceFilter{
			PlaceId:                  filter.PlaceId,
			RoleId:                   filter.RoleId,
			PresenceStatus:           filter.PresenceStatus,
			SubmissionPresenceStatus: filter.SubmissionPresenceStatus,
			Date:                     param.DateParam(today),
			LocationType:             param.LocationTypeParam(enum.LocationTypeWarehouse),
		})
		if err != nil {
			s.log.Error("failed get location user presence", zap.Error(err))
			return nil, err
		}

		for _, e := range userPresences {
			responses = append(responses, dto.UserPresencePendingResponse{
				Id:           e.Id,
				Date:         e.CreatedAt.Format("02-01-2006"),
				Name:         e.User.Name,
				Status:       e.Status.String(),
				Evidence:     e.Evidence.String,
				Note:         e.Note.String,
				LocationType: enum.LocationTypeWarehouse.String(),
			})
		}
	} else if filter.LocationType.Value() == enum.LocationTypeUnassigned {
		userPresences, err := s.repository.GetLocationUserPresence(dto.GetLocationUserPresenceFilter{
			PlaceId:                  filter.PlaceId,
			RoleId:                   filter.RoleId,
			PresenceStatus:           filter.PresenceStatus,
			SubmissionPresenceStatus: filter.SubmissionPresenceStatus,
			Date:                     param.DateParam(today),
			LocationType:             param.LocationTypeParam(enum.LocationTypeUnassigned),
		})
		if err != nil {
			s.log.Error("failed get location user presence", zap.Error(err))
			return nil, err
		}

		for _, e := range userPresences {
			responses = append(responses, dto.UserPresencePendingResponse{
				Id:           e.Id,
				Date:         e.CreatedAt.Format("02-01-2006"),
				Name:         e.User.Name,
				Status:       e.Status.String(),
				Evidence:     e.Evidence.String,
				Note:         e.Note.String,
				LocationType: enum.LocationTypeUnassigned.String(),
			})
		}
	}

	return responses, nil
}
