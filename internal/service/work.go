package service

import (
	"database/sql"
	"math"
	"slices"
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
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type WorkService struct {
	log         *zap.Logger
	repository  repository.IWorkRepository
	roleService IRoleService
}

type IWorkService interface {
	GetWorkOverview() (dto.WorkOveriew, error)

	SaveDailyWorks(request dto.CreateDailyWorkRequest, userId uuid.UUID) (dto.DailyWorkResponse, error)
	GetDailyWorkSummariesBasedOnRole() ([]dto.DailyWorkListResponse, error)
	GetDailyWorksByRoleId(roleId uint64) (dto.DailyWorkResponse, error)
	DeleteDailyWork(id uint64) error
	GetDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) (dto.DailyWorkUserListPaginationResponse, error)

	CreateAdditionalWork(request dto.CreateAdditionalWorkRequest, userId uuid.UUID) (dto.AdditionalWorkResponse, error)
	GetAdditionalWorks(filter dto.GetAdditonalWorkFilter, currUser uuid.UUID) ([]dto.AdditionalWorkListResponse, error)
	GetAdditionalWorkById(id uint64) (dto.AdditionalWorkResponse, error)
	UpdateAdditionalWork(id uint64, request dto.UpdateAdditionalWorkRequest, userId uuid.UUID) (dto.AdditionalWorkResponse, error)
	DeleteAdditionalWork(id uint64) error
	UpdateAdditionalWorkUser(id uint64, request dto.UpdateAdditionalWorkUserRequest, userId uuid.UUID) (dto.AdditionalWorkUserResponse, error)
	TakeAdditionalWork(id uint64, userId uuid.UUID) (dto.AdditionalWorkUserResponse, error)
	GetAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) (dto.AdditionalWorkUserListPaginationResponse, error)

	GetUserWorksByUserId(userId uuid.UUID) (dto.WorkUserResponse, error)
	UpdateDailyWorkUser(id uint64, request dto.UpdateDailyWorkUserRequest, userId uuid.UUID) (dto.DailyWorkUserResponse, error)
	DeleteAdditionalWorkUser(id uint64) error
}

func NewWorkService(log *zap.Logger, repository repository.IWorkRepository, roleService IRoleService) IWorkService {
	return &WorkService{
		log:         log,
		repository:  repository,
		roleService: roleService,
	}
}

func (w *WorkService) GetWorkOverview() (dto.WorkOveriew, error) {
	dailyWorkSummaries, err := w.GetDailyWorkSummariesBasedOnRole()
	if err != nil {
		return dto.WorkOveriew{}, err
	}

	additionalWorkSummaries, err := w.GetAdditionalWorks(dto.GetAdditonalWorkFilter{}, uuid.Nil)
	if err != nil {
		return dto.WorkOveriew{}, err
	}

	return dto.WorkOveriew{
		AdditionalWorkSummaries: additionalWorkSummaries,
		DailyWorkSummaries:      dailyWorkSummaries,
	}, nil
}

func (w *WorkService) SaveDailyWorks(request dto.CreateDailyWorkRequest, userId uuid.UUID) (dto.DailyWorkResponse, error) {
	w.repository.UseTx(true)
	defer w.repository.Rollback()

	for _, dailyWork := range request.DailyWorkDetail {
		startTime, err := datatype.ParseTimeOnly(dailyWork.StartTime)
		if err != nil {
			w.log.Error("failed to parse start time", zap.Error(err))
			return dto.DailyWorkResponse{}, err
		}

		endTime, err := datatype.ParseTimeOnly(dailyWork.EndTime)
		if err != nil {
			w.log.Error("failed to parse end time", zap.Error(err))
			return dto.DailyWorkResponse{}, err
		}

		dailyWorkEntity := entity.DailyWork{
			Id:          dailyWork.Id,
			Description: dailyWork.Description,
			RoleId:      request.RoleId,
			StartTime:   startTime,
			EndTime:     endTime,
		}

		if dailyWorkEntity.Id == 0 {
			dailyWorkEntity.CreatedBy = uuid.NullUUID{UUID: userId, Valid: true}
		} else {
			dailyWorkEntity.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
		}

		if err := w.repository.SaveDailyWork(&dailyWorkEntity); err != nil {
			w.log.Error("failed to create daily work", zap.Error(err))
			return dto.DailyWorkResponse{}, err
		}
	}

	if err := w.repository.Commit(); err != nil {
		w.log.Error("failed to commit transaction", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	roleResponse, err := w.roleService.GetRoleById(request.RoleId)
	if err != nil {
		w.log.Error("failed to get role by id", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	dailyWorkResponses := make([]dto.DailyWorkDetailResponse, 0)
	dailyWorkEntity, err := w.repository.GetDailyWorkByRoleId(request.RoleId)
	if err != nil {
		w.log.Error("failed to get daily work by role id", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	for _, dailyWorkEntity := range dailyWorkEntity {
		dailyWorkResponses = append(dailyWorkResponses, mapper.DailyWorkDetailToResponse(&dailyWorkEntity))
	}

	dailyWorkResponse := dto.DailyWorkResponse{
		Role: dto.RoleResponse{
			Id:   roleResponse.Id,
			Name: roleResponse.Name,
		},
		DailyWorks: dailyWorkResponses,
	}

	return dailyWorkResponse, nil
}

func (w *WorkService) GetDailyWorkSummariesBasedOnRole() ([]dto.DailyWorkListResponse, error) {
	roles, err := w.roleService.GetRoles()
	if err != nil {
		return nil, err
	}

	roleIds := make([]uint64, 0)
	for _, role := range roles {
		roleIds = append(roleIds, role.Id)
	}

	dailyWorkSummaries, err := w.repository.GetDailyWorkSummariesByRoleIds(roleIds)
	if err != nil {
		w.log.Error("failed to get daily work based on role", zap.Error(err))
		return nil, err
	}

	dailyWorkResponse := make([]dto.DailyWorkListResponse, len(dailyWorkSummaries))
	for i, dailyWorkSummary := range dailyWorkSummaries {
		dailyWorkResponse[i] = dto.DailyWorkListResponse{
			Role: dto.RoleResponse{
				Id:   dailyWorkSummary.RoleID,
				Name: dailyWorkSummary.RoleName,
			},
			TotalWork: dailyWorkSummary.TotalWork,
			TotalUser: dailyWorkSummary.TotalUser,
		}
	}

	return dailyWorkResponse, nil
}

func (w *WorkService) GetDailyWorksByRoleId(roleId uint64) (dto.DailyWorkResponse, error) {
	dailyWork, err := w.repository.GetDailyWorkByRoleId(roleId)
	if err != nil {
		w.log.Error("failed to get daily work by role id", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	dailyWorkResponse := make([]dto.DailyWorkDetailResponse, 0, len(dailyWork))
	for _, dailyWork := range dailyWork {
		dailyWorkResponse = append(dailyWorkResponse, mapper.DailyWorkDetailToResponse(&dailyWork))
	}

	role, err := w.roleService.GetRoleById(roleId)
	if err != nil {
		w.log.Error("failed to get role by id", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	return dto.DailyWorkResponse{
		Role: dto.RoleResponse{
			Id:   role.Id,
			Name: role.Name,
		},
		DailyWorks: dailyWorkResponse,
	}, nil
}

func (w *WorkService) GetUserWorksByUserId(userId uuid.UUID) (dto.WorkUserResponse, error) {
	withDeleted := false
	dailyWorkUsers, err := w.repository.GetDailyWorkUserByUserId(userId, dto.GetDailyWorkUserFilter{
		Date:        param.DateParam(time.Now()),
		WithDeleted: &withDeleted,
	})
	if err != nil {
		w.log.Error("failed to get additional work user", zap.Error(err))
		return dto.WorkUserResponse{}, err
	}

	dailyWorkResponse := make([]dto.DailyWorkUserResponse, 0)
	for _, dailyWorkUser := range dailyWorkUsers {
		dailyWorkResponse = append(dailyWorkResponse, mapper.DailyWorkUserToResponse(&dailyWorkUser))
	}

	additionalWorkUsers, err := w.repository.GetAdditionalWorkUserByUserId(userId, dto.GetAdditionalWorkUserFilter{
		WithDeleted:          &withDeleted,
		IsAdditionalWorkFull: true,
	})
	if err != nil {
		w.log.Error("failed to get daily work user", zap.Error(err))
		return dto.WorkUserResponse{}, err
	}

	additionalWorkResponse := make([]dto.AdditionalWorkUserResponse, 0)
	for _, additionalWork := range additionalWorkUsers {
		additionalWorkResponse = append(additionalWorkResponse, mapper.AdditionalWorkUserToResponse(&additionalWork))
	}

	return dto.WorkUserResponse{
		DailyWorks:      dailyWorkResponse,
		AdditionalWorks: additionalWorkResponse,
	}, nil
}

func (w *WorkService) CreateAdditionalWork(request dto.CreateAdditionalWorkRequest, userId uuid.UUID) (dto.AdditionalWorkResponse, error) {
	w.repository.UseTx(true)
	defer w.repository.Rollback()

	locationType := enum.ValueOfLocationWorkType(request.LocationType)
	if !locationType.IsValid() {
		w.log.Error("invalid location", zap.String("location", request.LocationType))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid location")
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		w.log.Error("invalid salary", zap.String("salary", request.Salary))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid salary")
	}

	workDate, err := time.Parse("02-01-2006 15:04", request.WorkDate)
	if err != nil {
		w.log.Error("failed to parse work date", zap.Error(err))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid work date")
	}

	additionalWork := entity.AdditionalWork{
		Name:         request.Name,
		LocationId:   request.LocationId,
		LocationType: locationType,
		Description:  request.Description,
		Slot:         request.Slot,
		Salary:       salary,
		WorkDate:     workDate,
		CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	switch locationType {
	case enum.LocationTypeCage:
		additionalWork.CageId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
	case enum.LocationTypeStore:
		additionalWork.StoreId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
	case enum.LocationTypeWarehouse:
		additionalWork.WarehouseId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
	}

	if err := w.repository.SaveAdditionalWork(&additionalWork); err != nil {
		w.log.Error("failed to create additional work", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	IsAdditionalWorkFull := false
	if request.Slot == uint64(len(request.UserIds)) {
		IsAdditionalWorkFull = true
	}

	additionalWorkUsers := make([]entity.AdditionalWorkUser, 0)
	if request.UserIds != nil {
		for _, userIdReq := range request.UserIds {
			additionalWorkUsers = append(additionalWorkUsers, entity.AdditionalWorkUser{
				UserId:               uuid.MustParse(userIdReq),
				AdditionalWorkId:     additionalWork.Id,
				IsAdditionalWorkFull: IsAdditionalWorkFull,
				CreatedBy:            uuid.NullUUID{UUID: userId, Valid: true},
				TakenAt:              sql.NullTime{Time: time.Now(), Valid: true},
			})
		}

		err = w.repository.CreateAdditionalWorkUserInBatch(additionalWorkUsers)
		if err != nil {
			w.log.Error("failed to create additional work user in batch", zap.Error(err))
			return dto.AdditionalWorkResponse{}, err
		}
	}

	err = w.repository.Commit()
	if err != nil {
		w.log.Error("faile to commit transaction", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	additionalWork, err = w.repository.GetAdditionalWorkById(additionalWork.Id)
	if err != nil {
		w.log.Error("failed to get additional work by id", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	addtionalWorkUserResponses := make([]dto.AdditionalWorkUserInformationResponse, len(additionalWork.AdditionalWorkUsers))
	for i, user := range additionalWork.AdditionalWorkUsers {
		addtionalWorkUserResponses[i] = mapper.AdditionalWorkUserInformationToResponse(&user)
	}

	additionalWorkResponse := mapper.AdditionalWorkToResponse(&additionalWork)
	additionalWorkResponse.AdditionalWorkUserInformation = addtionalWorkUserResponses

	return additionalWorkResponse, nil
}

func (w *WorkService) GetAdditionalWorkById(id uint64) (dto.AdditionalWorkResponse, error) {
	additionalWork, err := w.repository.GetAdditionalWorkById(id)
	if err != nil {
		w.log.Error("failed to get additional work by role id", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	addtionalWorkUserResponses := make([]dto.AdditionalWorkUserInformationResponse, len(additionalWork.AdditionalWorkUsers))
	for i, user := range additionalWork.AdditionalWorkUsers {
		addtionalWorkUserResponses[i] = mapper.AdditionalWorkUserInformationToResponse(&user)
	}

	additionalWorkResponse := mapper.AdditionalWorkToResponse(&additionalWork)
	additionalWorkResponse.AdditionalWorkUserInformation = addtionalWorkUserResponses

	return additionalWorkResponse, nil
}

func (w *WorkService) UpdateAdditionalWork(id uint64, request dto.UpdateAdditionalWorkRequest, userId uuid.UUID) (dto.AdditionalWorkResponse, error) {
	w.repository.UseTx(true)
	defer w.repository.Rollback()

	additionalWork, err := w.repository.GetAdditionalWorkById(id)
	if err != nil {
		w.log.Error("failed to get additional work by id", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	locationType := enum.ValueOfLocationWorkType(request.LocationType)
	if !locationType.IsValid() {
		w.log.Error("invalid location", zap.String("location", request.LocationType))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid location")
	}

	salary, err := decimal.NewFromString(request.Salary)
	if err != nil {
		w.log.Error("invalid salary", zap.String("salary", request.Salary))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid salary")
	}

	workDate, err := time.Parse("02-01-2006 15:04", request.WorkDate)
	if err != nil {
		w.log.Error("failed to parse work date", zap.Error(err))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid work date format")
	}

	additionalWork.Name = request.Name
	additionalWork.LocationId = request.LocationId
	additionalWork.Description = request.Description
	additionalWork.Slot = request.Slot
	additionalWork.Salary = salary
	additionalWork.WorkDate = workDate
	additionalWork.LocationType = locationType
	additionalWork.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	switch locationType {
	case enum.LocationTypeCage:
		additionalWork.CageId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
		additionalWork.WarehouseId = sql.NullInt64{}
		additionalWork.StoreId = sql.NullInt64{}
	case enum.LocationTypeStore:
		additionalWork.StoreId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
		additionalWork.WarehouseId = sql.NullInt64{}
		additionalWork.CageId = sql.NullInt64{}
	case enum.LocationTypeWarehouse:
		additionalWork.WarehouseId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
		additionalWork.StoreId = sql.NullInt64{}
		additionalWork.CageId = sql.NullInt64{}
	}

	if err := w.repository.SaveAdditionalWork(&additionalWork); err != nil {
		w.log.Error("failed to update additional work", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	currentUserIds := make([]uuid.UUID, 0)
	for _, e := range additionalWork.AdditionalWorkUsers {
		currentUserIds = append(currentUserIds, e.UserId)
	}

	deleteUserIds := make([]uuid.UUID, 0)
	for _, e := range currentUserIds {
		if !slices.Contains(request.UserIds, e.String()) {
			deleteUserIds = append(deleteUserIds, e)
		}
	}

	newUserIds := make([]uuid.UUID, 0)
	for _, e := range request.UserIds {
		if !slices.Contains(currentUserIds, uuid.MustParse(e)) {
			newUserIds = append(newUserIds, uuid.MustParse(e))
		}
	}

	if deleteUserIds != nil {
		err = w.repository.DeleteAdditionalWorkUserByAdditionalWorkIdAndUserIds(additionalWork.Id, deleteUserIds)
		if err != nil {
			w.log.Error("failed to delete additional work user by additional work id and user ids", zap.Error(err))
			return dto.AdditionalWorkResponse{}, err
		}
	}

	isAdditionalWorkFull := false
	if request.Slot == uint64(len(request.UserIds)) {
		isAdditionalWorkFull = true
	}

	if newUserIds != nil {
		additionalWorkUsers := make([]entity.AdditionalWorkUser, 0)

		for _, userId := range newUserIds {
			additionalWorkUsers = append(additionalWorkUsers, entity.AdditionalWorkUser{
				UserId:               userId,
				AdditionalWorkId:     additionalWork.Id,
				IsAdditionalWorkFull: isAdditionalWorkFull,
				CreatedBy:            uuid.NullUUID{UUID: userId, Valid: true},
			})
		}

		err = w.repository.CreateAdditionalWorkUserInBatch(additionalWorkUsers)
		if err != nil {
			w.log.Error("failed to create additional work user in batch", zap.Error(err))
			return dto.AdditionalWorkResponse{}, err
		}
	}

	err = w.repository.Commit()
	if err != nil {
		w.log.Error("failed to commit trasaction", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	addtionalWorkUserResponses := make([]dto.AdditionalWorkUserInformationResponse, len(additionalWork.AdditionalWorkUsers))
	for i, user := range additionalWork.AdditionalWorkUsers {
		addtionalWorkUserResponses[i] = mapper.AdditionalWorkUserInformationToResponse(&user)
	}

	additionalWorkResponse := mapper.AdditionalWorkToResponse(&additionalWork)
	additionalWorkResponse.AdditionalWorkUserInformation = addtionalWorkUserResponses

	return additionalWorkResponse, nil
}

func (w *WorkService) DeleteAdditionalWork(id uint64) error {
	if err := w.repository.DeleteAdditionalWork(id); err != nil {
		w.log.Error("failed to delete additional work", zap.Error(err))
		return err
	}

	return nil
}

func (w *WorkService) GetAdditionalWorks(filter dto.GetAdditonalWorkFilter, currUser uuid.UUID) ([]dto.AdditionalWorkListResponse, error) {
	additionalWorks, err := w.repository.GetAdditionalWorks(filter)
	if err != nil {
		w.log.Error("failed to get additional works", zap.Error(err))
		return nil, err
	}

	additionalWorkResponses := make([]dto.AdditionalWorkListResponse, len(additionalWorks))
	for i, additionalWork := range additionalWorks {
		isTakenByCurrentUser := false
		takenBy := make([]uuid.UUID, 0)
		for _, additionalWorkUser := range additionalWork.AdditionalWorkUsers {
			takenBy = append(takenBy, additionalWorkUser.UserId)
		}

		if slices.Contains(takenBy, currUser) {
			isTakenByCurrentUser = true
		}

		additionalWorkResponses[i] = mapper.AdditionalWorkToListResponse(&additionalWork)
		additionalWorkResponses[i].IsTakenByCurrentUser = isTakenByCurrentUser
	}

	if filter.Status == constant.AdditionalWorkAvailable {
		availableAdditionalWorkResponse := make([]dto.AdditionalWorkListResponse, 0)
		for _, additionalWorkResponse := range additionalWorkResponses {
			if additionalWorkResponse.RemainingSlot > 0 {
				availableAdditionalWorkResponse = append(availableAdditionalWorkResponse, additionalWorkResponse)
			}
		}

		return availableAdditionalWorkResponse, nil
	}

	return additionalWorkResponses, nil
}

func (w *WorkService) UpdateAdditionalWorkUser(id uint64, request dto.UpdateAdditionalWorkUserRequest, userId uuid.UUID) (dto.AdditionalWorkUserResponse, error) {
	additionalWorkUser, err := w.repository.GetAdditionalWorkUserById(id)
	if err != nil {
		w.log.Error("failed to get additional work user by id", zap.Error(err))
		return dto.AdditionalWorkUserResponse{}, err
	}

	if !additionalWorkUser.IsAdditionalWorkFull {
		w.log.Warn("additional work must be full taken")
		return dto.AdditionalWorkUserResponse{}, errx.BadRequest("additional work must be full taken")
	}

	if additionalWorkUser.IsDone {
		w.log.Warn("additional work user already done")
		return dto.AdditionalWorkUserResponse{}, errx.BadRequest("additional work user already done")
	}

	additionalWorkUser.IsDone = request.IsDone
	additionalWorkUser.Note = request.Note
	additionalWorkUser.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	additionalWorkUser.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := w.repository.UpdateAdditionalWorkUser(&additionalWorkUser); err != nil {
		w.log.Error("failed to update additional work user", zap.Error(err))
		return dto.AdditionalWorkUserResponse{}, err
	}

	additionalWorkUserResponse := mapper.AdditionalWorkUserToResponse(&additionalWorkUser)

	return additionalWorkUserResponse, nil
}

func (w *WorkService) UpdateDailyWorkUser(id uint64, request dto.UpdateDailyWorkUserRequest, userId uuid.UUID) (dto.DailyWorkUserResponse, error) {
	dailyWorkUser, err := w.repository.GetDailyWorkUserById(id)
	if err != nil {
		w.log.Error("failed to get daily work user by id", zap.Error(err))
		return dto.DailyWorkUserResponse{}, err
	}

	if dailyWorkUser.IsDone {
		w.log.Error("daily work user already done", zap.Error(errx.BadRequest("daily work user already done")))
		return dto.DailyWorkUserResponse{}, errx.BadRequest("daily work user already done")
	}

	dailyWorkUser.IsDone = request.IsDone
	dailyWorkUser.Note = request.Note
	dailyWorkUser.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	dailyWorkUser.FinishedAt = sql.NullTime{Time: time.Now(), Valid: true}

	if err := w.repository.UpdateDailyWorkUser(&dailyWorkUser); err != nil {
		w.log.Error("failed to update daily work user", zap.Error(err))
		return dto.DailyWorkUserResponse{}, err
	}

	dailyWorkUserResponse := mapper.DailyWorkUserToResponse(&dailyWorkUser)

	return dailyWorkUserResponse, nil
}

func (w *WorkService) TakeAdditionalWork(id uint64, userId uuid.UUID) (dto.AdditionalWorkUserResponse, error) {
	w.repository.UseTx(true)
	defer w.repository.Rollback()

	additionalWork, err := w.repository.GetAdditionalWorkById(id)
	if err != nil {
		w.log.Error("failed to get additional work by id", zap.Error(err))
		return dto.AdditionalWorkUserResponse{}, err
	}

	if len(additionalWork.AdditionalWorkUsers) == int(additionalWork.Slot) {
		w.log.Error("additional work already full", zap.Error(errx.BadRequest("additional work already full")))
		return dto.AdditionalWorkUserResponse{}, errx.BadRequest("additional work already full")
	}

	// Note : Can be improved by checking in the database directly
	for _, user := range additionalWork.AdditionalWorkUsers {
		if user.UserId == userId {
			w.log.Error("user already taken additional work", zap.Error(errx.BadRequest("user already taken additional work")))
			return dto.AdditionalWorkUserResponse{}, errx.BadRequest("user already taken additional work")
		}
	}

	var (
		additionalWorkUser entity.AdditionalWorkUser
	)
	if len(additionalWork.AdditionalWorkUsers)+1 == int(additionalWork.Slot) {
		additionalWorkUser = entity.AdditionalWorkUser{
			AdditionalWorkId:     id,
			UserId:               userId,
			IsDone:               false,
			IsAdditionalWorkFull: true,
			TakenAt:              sql.NullTime{Time: time.Now(), Valid: true},
		}

		err := w.repository.UpdateAdditionalWorkUserByAdditionalWorkId(id, map[string]any{
			"is_additional_work_full": true,
		})
		if err != nil {
			w.log.Error("failed to update additional work user by additional work id", zap.Error(err))
			return dto.AdditionalWorkUserResponse{}, err
		}

	} else {
		additionalWorkUser = entity.AdditionalWorkUser{
			AdditionalWorkId:     id,
			UserId:               userId,
			IsDone:               false,
			IsAdditionalWorkFull: false,
			TakenAt:              sql.NullTime{Time: time.Now(), Valid: true},
		}
	}

	if err := w.repository.CreateAdditionalWorkUser(&additionalWorkUser); err != nil {
		w.log.Error("failed to create additional work user", zap.Error(err))
		return dto.AdditionalWorkUserResponse{}, err
	}

	err = w.repository.Commit()
	if err != nil {
		w.log.Error("failed to commit transaction", zap.Error(err))
		return dto.AdditionalWorkUserResponse{}, err
	}

	additionalWorkUser, err = w.repository.GetAdditionalWorkUserById(additionalWorkUser.Id)
	if err != nil {
		w.log.Error("failed to get additional work user by id", zap.Error(err))
		return dto.AdditionalWorkUserResponse{}, err
	}

	return mapper.AdditionalWorkUserToResponse(&additionalWorkUser), nil
}

func (w *WorkService) DeleteDailyWork(id uint64) error {
	if err := w.repository.DeleteDailyWork(id); err != nil {
		w.log.Error("failed to delete daily work", zap.Error(err))
		return err
	}

	return nil
}

func (w *WorkService) GetAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) (dto.AdditionalWorkUserListPaginationResponse, error) {
	w.repository.UseTx(false)

	additionalWorkUsers, err := w.repository.GetAdditionalWorkUserByUserId(userId, filter)
	if err != nil {
		w.log.Error("failed to get additional work user", zap.Error(err))
		return dto.AdditionalWorkUserListPaginationResponse{}, err
	}

	additionalWorkResponse := make([]dto.AdditionalWorkUserResponse, 0)
	for _, additionalWork := range additionalWorkUsers {
		additionalWorkResponse = append(additionalWorkResponse, mapper.AdditionalWorkUserToResponse(&additionalWork))
	}

	totalData, err := w.repository.CountAdditionalWorkUserByUserId(userId, filter)
	if err != nil {
		w.log.Error("failed count daily work user by user id", zap.Error(err))
		return dto.AdditionalWorkUserListPaginationResponse{}, err
	}

	resp := dto.AdditionalWorkUserListPaginationResponse{
		AdditionalWorkUsers: additionalWorkResponse,
	}

	if filter.Page > 0 {
		resp.TotalData = uint64(totalData)
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (w *WorkService) GetDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) (dto.DailyWorkUserListPaginationResponse, error) {
	w.repository.UseTx(false)

	dailyWorkUsers, err := w.repository.GetDailyWorkUserByUserId(userId, filter)
	if err != nil {
		w.log.Error("failed to get additional work user", zap.Error(err))
		return dto.DailyWorkUserListPaginationResponse{}, err
	}

	dailyWorkResponse := make([]dto.DailyWorkUserResponse, len(dailyWorkUsers))
	for i, dailyWorkUser := range dailyWorkUsers {
		dailyWorkResponse[i] = mapper.DailyWorkUserToResponse(&dailyWorkUser)
	}

	totalData, err := w.repository.CountDailyWorkUserByUserId(userId, filter)
	if err != nil {
		w.log.Error("failed count daily work user by user id", zap.Error(err))
		return dto.DailyWorkUserListPaginationResponse{}, err
	}

	resp := dto.DailyWorkUserListPaginationResponse{
		DailyWorkUsers: dailyWorkResponse,
	}

	if filter.Page > 0 {
		resp.TotalData = uint64(totalData)
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (w *WorkService) DeleteAdditionalWorkUser(id uint64) error {
	w.repository.UseTx(false)

	err := w.repository.DeleteAdditionalWorkUser(id)
	if err != nil {
		w.log.Error("failed to delete additional work user", zap.Error(err))
		return err
	}

	return nil
}
