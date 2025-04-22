package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type WorkService struct {
	log         *zap.Logger
	repository  repository.IWorkRepository
	roleService IRoleService
}

type IWorkService interface {
	CreateAndUpdateDailyWorks(request dto.CreateDailyWorkRequest, accountId uuid.UUID) (dto.DailyWorkResponse, error)
	GetDailyWorksBasedOnRole() ([]dto.DailyWorkListResponse, error)
	GetDailyWorksByRoleId(roleId uint64) (dto.DailyWorkResponse, error)
	UpdateDailyWorkStaff(id uint64, request dto.UpdateDailyWorkStaffRequest, accountId uuid.UUID) (dto.DailyWorkStaffResponse, error)

	CreateAdditionalWork(request dto.CreateAdditionalWorkRequest, accountId uuid.UUID) (dto.AdditionalWorkResponse, error)
	GetAdditionalWorks(filter dto.GetAdditonalWorkFilter) ([]dto.AdditionalWorkListResponse, error)
	GetAdditionalWorkById(id uint64) (dto.AdditionalWorkResponse, error)
	UpdateAdditionalWork(id uint64, request dto.UpdateAdditionalWorkRequest, accountId uuid.UUID) (dto.AdditionalWorkResponse, error)
	DeleteAdditionalWork(id uint64) error
	UpdateAdditionalWorkStaff(id uint64, request dto.UpdateAdditionalWorkStaffRequest, accountId uuid.UUID) (dto.AdditionalWorkStaffResponse, error)
	TakeAdditionalWork(id uint64, staffId uuid.UUID) (dto.AdditionalWorkStaffResponse, error)

	GetStaffWorksByStaffId(staffId uuid.UUID) (dto.WorkStaffResponse, error)
}

func NewWorkService(log *zap.Logger, repository repository.IWorkRepository, roleService IRoleService) IWorkService {
	return &WorkService{
		log:         log,
		repository:  repository,
		roleService: roleService,
	}
}

func (w *WorkService) CreateAndUpdateDailyWorks(request dto.CreateDailyWorkRequest, accountId uuid.UUID) (dto.DailyWorkResponse, error) {
	w.repository.UseTx(true)
	defer w.repository.Rollback()

	for _, dailyWork := range request.DailyWorkDetail {
		startTime, err := datatype.ParseTimeOnly(dailyWork.StartTime)
		if err != nil {
			w.log.Error("[CreateAndUpdateDailyWork] failed to parse start time", zap.Error(err))
			return dto.DailyWorkResponse{}, err
		}

		endTime, err := datatype.ParseTimeOnly(dailyWork.EndTime)
		if err != nil {
			w.log.Error("[CreateAndUpdateDailyWork] failed to parse end time", zap.Error(err))
			return dto.DailyWorkResponse{}, err
		}

		dailyWorkEntity := entity.DailyWork{
			Id:          dailyWork.Id,
			Description: dailyWork.Description,
			RoleId:      request.RoleId,
			StartTime:   startTime,
			EndTime:     endTime,
		}

		// using save not create because if the daily work already exists, it will be updated
		if err := w.repository.CreateDailyWork(&dailyWorkEntity); err != nil {
			w.log.Error("[CreateAndUpdateDailyWork] failed to create daily work", zap.Error(err))
			return dto.DailyWorkResponse{}, err
		}
	}

	if err := w.repository.Commit(); err != nil {
		w.log.Error("[CreateAndUpdateDailyWork] failed to commit transaction", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	roleResponse, err := w.roleService.GetRoleById(request.RoleId)
	if err != nil {
		w.log.Error("[CreateAndUpdateDailyWork] failed to get role by id", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	dailyWorkResponses := make([]dto.DailyWorkDetailResponse, 0)
	dailyWorkEntity, err := w.repository.GetDailyWorkByRoleId(request.RoleId)
	if err != nil {
		w.log.Error("[CreateAndUpdateDailyWork] failed to get daily work by role id", zap.Error(err))
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

func (w *WorkService) GetDailyWorksBasedOnRole() ([]dto.DailyWorkListResponse, error) {
	dailyWorkSummaries, err := w.repository.GetDailyWorkBasedOnRole()
	if err != nil {
		w.log.Error("[GetDailyWorkBasedOnRole] failed to get daily work based on role", zap.Error(err))
		return nil, err
	}

	dailyWorkResponse := make([]dto.DailyWorkListResponse, len(dailyWorkSummaries))
	for _, dailyWorkSummary := range dailyWorkSummaries {
		dailyWorkResponse = append(dailyWorkResponse, dto.DailyWorkListResponse{
			Role: dto.RoleResponse{
				Id:   dailyWorkSummary.RoleID,
				Name: dailyWorkSummary.RoleName,
			},
			TotalWork:  dailyWorkSummary.TotalWork,
			TotalStaff: dailyWorkSummary.TotalStaff,
		})
	}

	return dailyWorkResponse, nil
}

func (w *WorkService) GetDailyWorksByRoleId(roleId uint64) (dto.DailyWorkResponse, error) {
	dailyWork, err := w.repository.GetDailyWorkByRoleId(roleId)
	if err != nil {
		w.log.Error("[GetDailyWorkByRolId] failed to get daily work by role id", zap.Error(err))
		return dto.DailyWorkResponse{}, err
	}

	dailyWorkResponse := make([]dto.DailyWorkDetailResponse, 0, len(dailyWork))
	for _, dailyWork := range dailyWork {
		dailyWorkResponse = append(dailyWorkResponse, mapper.DailyWorkDetailToResponse(&dailyWork))
	}

	role, err := w.roleService.GetRoleById(roleId)
	if err != nil {
		w.log.Error("[GetDailyWorkByRolId] failed to get role by id", zap.Error(err))
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

// additional work & daily work staffs
func (w *WorkService) GetStaffWorksByStaffId(staffId uuid.UUID) (dto.WorkStaffResponse, error) {
	dailyWorkStaffs, err := w.repository.GetDailyWorkStaffsByStaffId(staffId)
	if err != nil {
		return dto.WorkStaffResponse{}, err
	}

	dailyWorkResponse := make([]dto.DailyWorkStaffResponse, 0, len(dailyWorkStaffs))
	for _, dailyWorkStaff := range dailyWorkStaffs {
		dailyWorkResponse = append(dailyWorkResponse, mapper.DailyWorkStaffToResponse(&dailyWorkStaff))
	}

	additonalWorkStaffs, err := w.repository.GetAdditionalWorkStaffByStaffId(staffId)
	if err != nil {
		return dto.WorkStaffResponse{}, err
	}
	additionalWorkResponse := make([]dto.AdditionalWorkStaffResponse, 0, len(dailyWorkStaffs))
	for _, additionalWork := range additonalWorkStaffs {
		additionalWorkResponse = append(additionalWorkResponse, mapper.AdditionalWorkStaffToResponse(&additionalWork))
	}

	return dto.WorkStaffResponse{
		DailyWorks:      dailyWorkResponse,
		AdditionalWorks: additionalWorkResponse,
	}, nil
}

func (w *WorkService) CreateAdditionalWork(request dto.CreateAdditionalWorkRequest, accountId uuid.UUID) (dto.AdditionalWorkResponse, error) {
	location := enum.ValueOfLocationAddionalWork(request.Location)
	if !location.IsValid() {
		w.log.Error("[CreateAdditonalWork] invalid location", zap.String("location", request.Location))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid location")
	}

	additionalWork := entity.AdditionalWork{
		Description: request.Description,
		Slot:        request.Slot,
		Location:    location,
	}

	if err := w.repository.CreateAdditionalWork(&additionalWork); err != nil {
		w.log.Error("[CreateAdditonalWork] failed to create additional work", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	additionalWork, err := w.repository.GetAdditionalWorkById(additionalWork.Id)
	if err != nil {
		w.log.Error("[CreateAdditonalWork] failed to get additional work by id", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	addtionalWorkStaffResponses := make([]dto.AdditionalWorkStaffInformationResponse, len(additionalWork.AdditionalWorkStaff))
	for i, staff := range additionalWork.AdditionalWorkStaff {
		addtionalWorkStaffResponses[i] = mapper.AdditionalWorkStaffInformationToResponse(&staff)
	}

	additionalWorkResponse := mapper.AdditionalWorkToResponse(&additionalWork)
	additionalWorkResponse.AdditionalWorkStaffInformation = addtionalWorkStaffResponses

	return additionalWorkResponse, nil
}

func (w *WorkService) GetAdditionalWorkById(id uint64) (dto.AdditionalWorkResponse, error) {
	additionalWork, err := w.repository.GetAdditionalWorkById(id)
	if err != nil {
		w.log.Error("[GetAdditionalWorkById] failed to get additional work by role id", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	addtionalWorkStaffResponses := make([]dto.AdditionalWorkStaffInformationResponse, len(additionalWork.AdditionalWorkStaff))
	for i, staff := range additionalWork.AdditionalWorkStaff {
		addtionalWorkStaffResponses[i] = mapper.AdditionalWorkStaffInformationToResponse(&staff)
	}

	additionalWorkResponse := mapper.AdditionalWorkToResponse(&additionalWork)
	additionalWorkResponse.AdditionalWorkStaffInformation = addtionalWorkStaffResponses

	return additionalWorkResponse, nil
}

func (w *WorkService) UpdateAdditionalWork(id uint64, request dto.UpdateAdditionalWorkRequest, accountId uuid.UUID) (dto.AdditionalWorkResponse, error) {
	additionalWork, err := w.repository.GetAdditionalWorkById(id)
	if err != nil {
		w.log.Error("[UpdateAdditonalWork] failed to get additional work by id", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	location := enum.ValueOfLocationAddionalWork(request.Location)
	if !location.IsValid() {
		w.log.Error("[UpdateAdditonalWork] invalid location", zap.String("location", request.Location))
		return dto.AdditionalWorkResponse{}, errx.BadRequest("invalid location")
	}

	additionalWork.Description = request.Description
	additionalWork.Slot = request.Slot
	additionalWork.Location = location
	additionalWork.UpdatedBy = accountId

	if err := w.repository.CreateAdditionalWork(&additionalWork); err != nil {
		w.log.Error("[UpdateAdditonalWork] failed to update additional work", zap.Error(err))
		return dto.AdditionalWorkResponse{}, err
	}

	addtionalWorkStaffResponses := make([]dto.AdditionalWorkStaffInformationResponse, len(additionalWork.AdditionalWorkStaff))
	for i, staff := range additionalWork.AdditionalWorkStaff {
		addtionalWorkStaffResponses[i] = mapper.AdditionalWorkStaffInformationToResponse(&staff)
	}

	additionalWorkResponse := mapper.AdditionalWorkToResponse(&additionalWork)
	additionalWorkResponse.AdditionalWorkStaffInformation = addtionalWorkStaffResponses

	return additionalWorkResponse, nil
}

func (w *WorkService) DeleteAdditionalWork(id uint64) error {
	if err := w.repository.DeleteAdditionalWork(id); err != nil {
		w.log.Error("[DeteleAdditionalWork] failed to delete additional work", zap.Error(err))
		return err
	}

	return nil
}

func (w *WorkService) GetAdditionalWorks(filter dto.GetAdditonalWorkFilter) ([]dto.AdditionalWorkListResponse, error) {
	additionalWorks, err := w.repository.GetAdditionalWorks()
	if err != nil {
		w.log.Error("[GetAdditionalWorks] failed to get additional works", zap.Error(err))
		return nil, err
	}

	additionalWorkResponses := make([]dto.AdditionalWorkListResponse, len(additionalWorks))
	for i, additionalWork := range additionalWorks {
		additionalWorkResponses[i] = mapper.AdditionalWorkToListResponse(&additionalWork)

		if len(additionalWork.AdditionalWorkStaff) == 0 {
			additionalWorkResponses[i].Status = constant.AdditionalWorkStatusNotProgress
		} else {
			additionalWorkResponses[i].Status = constant.AdditionalWorkStatusDone
			for _, staff := range additionalWork.AdditionalWorkStaff {
				if !staff.IsDone {
					additionalWorkResponses[i].Status = constant.AdditionalWorkStatusOnProgress
					break
				}
			}
		}
	}

	if filter.Status == "Available" {
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

func (w *WorkService) CreateAdditionalWorkStaff(id uint64, staffId uuid.UUID) (dto.AdditionalWorkStaffResponse, error) {
	additionalWorkStaff := entity.AdditionalWorkStaff{
		AdditionalWorkId: id,
		StaffId:          staffId,
		IsDone:           false,
	}

	if err := w.repository.CreateAdditionalWorkStaff(&additionalWorkStaff); err != nil {
		w.log.Error("[TakeAdditionalWork] failed to create additional work staff", zap.Error(err))
		return dto.AdditionalWorkStaffResponse{}, err
	}

	additionalWorkStaffResponse := mapper.AdditionalWorkStaffToResponse(&additionalWorkStaff)

	return additionalWorkStaffResponse, nil
}

func (w *WorkService) UpdateAdditionalWorkStaff(id uint64, request dto.UpdateAdditionalWorkStaffRequest, accountId uuid.UUID) (dto.AdditionalWorkStaffResponse, error) {
	additionalWorkStaff, err := w.repository.GetAdditionalWorkStaffById(id)
	if err != nil {
		w.log.Error("[UpdateAdditionalWorkStaff] failed to get additional work staff by id", zap.Error(err))
		return dto.AdditionalWorkStaffResponse{}, err
	}

	if additionalWorkStaff.IsDone {
		w.log.Error("[UpdateAdditionalWorkStaff] additional work staff already done", zap.Error(errx.BadRequest("additional work staff already done")))
		return dto.AdditionalWorkStaffResponse{}, errx.BadRequest("additional work staff already done")
	}

	additionalWorkStaff.IsDone = request.IsDone
	additionalWorkStaff.UpdatedBy = accountId

	if err := w.repository.UpdateAdditionalWorkStaff(&additionalWorkStaff); err != nil {
		w.log.Error("[UpdateAdditionalWorkStaff] failed to update additional work staff", zap.Error(err))
		return dto.AdditionalWorkStaffResponse{}, err
	}

	additionalWorkStaffResponse := mapper.AdditionalWorkStaffToResponse(&additionalWorkStaff)

	return additionalWorkStaffResponse, nil
}

func (w *WorkService) UpdateDailyWorkStaff(id uint64, request dto.UpdateDailyWorkStaffRequest, accountId uuid.UUID) (dto.DailyWorkStaffResponse, error) {
	dailyWorkStaff, err := w.repository.GetDailyWorkStaffById(id)
	if err != nil {
		w.log.Error("[UpdateDailyWorkStaff] failed to get daily work staff by id", zap.Error(err))
		return dto.DailyWorkStaffResponse{}, err
	}

	if dailyWorkStaff.IsDone {
		w.log.Error("[UpdateDailyWorkStaff] daily work staff already done", zap.Error(errx.BadRequest("daily work staff already done")))
		return dto.DailyWorkStaffResponse{}, errx.BadRequest("daily work staff already done")
	}

	dailyWorkStaff.IsDone = request.IsDone
	dailyWorkStaff.UpdatedBy = accountId

	if err := w.repository.UpdateDailyWorkStaff(&dailyWorkStaff); err != nil {
		w.log.Error("[UpdateDailyWorkStaff] failed to update daily work staff", zap.Error(err))
		return dto.DailyWorkStaffResponse{}, err
	}

	dailyWorkStaffResponse := mapper.DailyWorkStaffToResponse(&dailyWorkStaff)

	return dailyWorkStaffResponse, nil
}

func (w *WorkService) TakeAdditionalWork(id uint64, staffId uuid.UUID) (dto.AdditionalWorkStaffResponse, error) {
	additionalWork, err := w.repository.GetAdditionalWorkById(id)
	if err != nil {
		w.log.Error("[TakeAdditionalWork] failed to get additional work by id", zap.Error(err))
		return dto.AdditionalWorkStaffResponse{}, err
	}

	if len(additionalWork.AdditionalWorkStaff) == int(additionalWork.Slot) {
		w.log.Error("[TakeAdditionalWork] additional work already full", zap.Error(errx.BadRequest("additional work already full")))
		return dto.AdditionalWorkStaffResponse{}, errx.BadRequest("additional work already full")
	}

	// Note : Can be improved by checking in the database directly
	for _, staff := range additionalWork.AdditionalWorkStaff {
		if staff.StaffId == staffId {
			w.log.Error("[TakeAdditionalWork] staff already taken additional work", zap.Error(errx.BadRequest("staff already taken additional work")))
			return dto.AdditionalWorkStaffResponse{}, errx.BadRequest("staff already taken additional work")
		}
	}

	additionalWorkStaff := entity.AdditionalWorkStaff{
		AdditionalWorkId: id,
		StaffId:          staffId,
		IsDone:           false,
	}

	if err := w.repository.CreateAdditionalWorkStaff(&additionalWorkStaff); err != nil {
		w.log.Error("[TakeAdditionalWork] failed to create additional work staff", zap.Error(err))
		return dto.AdditionalWorkStaffResponse{}, err
	}

	additionalWorkStaff, err = w.repository.GetAdditionalWorkStaffById(additionalWorkStaff.Id)
	if err != nil {
		w.log.Error("[TakeAdditionalWork] failed to get additional work staff by id", zap.Error(err))
		return dto.AdditionalWorkStaffResponse{}, err
	}

	additionalWorkStaffResponse := mapper.AdditionalWorkStaffToResponse(&additionalWorkStaff)

	return additionalWorkStaffResponse, nil
}
