package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"gorm.io/gorm"
)

type WorkRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IWorkRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetDailyWorkStaffByStaffId(staffId uuid.UUID, filter dto.GetDailyWorkStaffFilter) ([]entity.DailyWorkUser, error)
	GetAdditionalWorkStaffByStaffId(staffId uuid.UUID, filter dto.GetAdditionalWorkStaffFilter) ([]entity.AdditionalWorkUser, error)

	CreateDailyWork(dailyWork *entity.DailyWork) error
	GetDailyWorkByRoleId(roleId uint64) ([]entity.DailyWork, error)
	GetDailyWorkBasedOnRole(filter dto.GetDailyWorkBasedOnRoleFilter) ([]entity.DailyWorkSummary, error)
	DeleteDailyWork(id uint64) error

	CreateAdditionalWork(additionalWork *entity.AdditionalWork) error
	GetAdditionalWorkById(id uint64) (entity.AdditionalWork, error)
	DeleteAdditionalWork(id uint64) error
	GetAdditionalWorks(filter dto.GetAdditonalWorkFilter) ([]entity.AdditionalWork, error)
	CreateAdditionalWorkStaff(additionalWorkStaff *entity.AdditionalWorkUser) error
	GetAdditionalWorkStaffById(id uint64) (entity.AdditionalWorkUser, error)
	UpdateAdditionalWorkStaff(additionalWorkStaff *entity.AdditionalWorkUser) error
	GetDailyWorkStaffById(id uint64) (entity.DailyWorkUser, error)
	UpdateDailyWorkStaff(dailyWorkStaff *entity.DailyWorkUser) error
}

func NewWorkRepository(db *gorm.DB) IWorkRepository {
	return &WorkRepository{
		db: db,
	}
}

func (r *WorkRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *WorkRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *WorkRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *WorkRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *WorkRepository) CreateDailyWork(dailyWork *entity.DailyWork) error {
	return r.GetDB().Save(dailyWork).Error
}

func (r *WorkRepository) CreateAdditionalWork(additionalWork *entity.AdditionalWork) error {
	return r.GetDB().Save(additionalWork).Error
}

func (r *WorkRepository) GetDailyWorkByRoleId(roleId uint64) ([]entity.DailyWork, error) {
	var dailyWorks []entity.DailyWork
	err := r.GetDB().
		Preload("Role").
		Where("role_id = ?", roleId).
		Find(&dailyWorks).Error

	if err != nil {
		return nil, err
	}

	return dailyWorks, nil
}

func (r *WorkRepository) GetAdditionalWorkById(id uint64) (entity.AdditionalWork, error) {
	var additionalWorks entity.AdditionalWork
	err := r.GetDB().
		Preload("AdditionalWorkStaff.Staff").
		Where("id = ? AND deleted_at IS NULL", id).
		First(&additionalWorks).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AdditionalWork{}, errx.NotFound("additional work not found")
		}
		return entity.AdditionalWork{}, err
	}
	return additionalWorks, nil
}

func (r *WorkRepository) GetDailyWorkStaffByStaffId(staffId uuid.UUID, filter dto.GetDailyWorkStaffFilter) ([]entity.DailyWorkUser, error) {
	var dailyWorkStaffs []entity.DailyWorkUser
	query := r.GetDB()

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value().Format("2006-01-02"))
	}

	if filter.WithDeleted {
		query = query.Preload("DailyWork")
	} else {
		query = query.Preload("DailyWork", "deleted_at IS NULL")
	}

	err := query.Preload("DailyWork").
		Preload("Staff.Account", "id = ?", staffId).
		Order("created_at DESC").
		Find(&dailyWorkStaffs).Error

	if err != nil {
		return nil, err
	}

	return dailyWorkStaffs, nil
}

func (r *WorkRepository) GetAdditionalWorkStaffByStaffId(userId uuid.UUID, filter dto.GetAdditionalWorkStaffFilter) ([]entity.AdditionalWorkUser, error) {
	var additionalWorks []entity.AdditionalWorkUser
	query := r.GetDB()

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if filter.WithDeleted {
		query = query.Preload("AdditionalWork")
	} else {
		query = query.Preload("AdditionalWork", "deleted_at IS NULL")
	}

	err := query.Preload("AdditionalWork").
		Preload("User", "id = ?", userId).
		Order("created_at DESC").
		Find(&additionalWorks).Error

	if err != nil {
		return nil, err
	}
	return additionalWorks, nil
}

func (r *WorkRepository) DeleteAdditionalWork(id uint64) error {
	err := r.GetDB().Delete(&entity.AdditionalWork{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) GetAdditionalWorks(filter dto.GetAdditonalWorkFilter) ([]entity.AdditionalWork, error) {
	var additionalWorks []entity.AdditionalWork
	err := r.GetDB().
		Preload("AdditionalWorkUser.User").
		Where("deleted_at IS NULL").
		Find(&additionalWorks).Error

	if err != nil {
		return nil, err
	}
	return additionalWorks, nil
}

func (r *WorkRepository) GetDailyWorkBasedOnRole(filter dto.GetDailyWorkBasedOnRoleFilter) ([]entity.DailyWorkSummary, error) {
	var dailyWorkSummaries []entity.DailyWorkSummary

	query := r.GetDB().
		Table("daily_works").
		Select(`daily_works.role_id, roles.name AS role_name,
				COUNT(DISTINCT daily_works.id) AS total_work,
				COUNT(DISTINCT accounts.id) AS total_staff`).
		Joins("LEFT JOIN roles ON daily_works.role_id = roles.id").
		Joins("LEFT JOIN accounts ON accounts.role_id = roles.id").
		Group("daily_works.role_id, roles.name").
		Where("daily_works.deleted_at IS NULL")

	if filter.RoleIds != nil {
		query = query.Where("daily_works.role_id IN ?", filter.RoleIds)
	}

	err := query.Scan(&dailyWorkSummaries).Error

	if err != nil {
		return nil, err
	}
	return dailyWorkSummaries, nil
}

func (r *WorkRepository) CreateAdditionalWorkStaff(additionalWorkStaff *entity.AdditionalWorkUser) error {
	return r.GetDB().Create(additionalWorkStaff).Error
}

func (r *WorkRepository) GetAdditionalWorkStaffById(id uint64) (entity.AdditionalWorkUser, error) {
	var additionalWorkStaff entity.AdditionalWorkUser
	err := r.GetDB().
		Model(&entity.AdditionalWorkUser{}).
		Preload("AdditionalWork", "deleted_at IS NULL").
		Preload("Staff.Account").
		Where("id = ?", id).
		First(&additionalWorkStaff).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AdditionalWorkUser{}, errx.NotFound("additional work staff not found")
		}
		return entity.AdditionalWorkUser{}, err
	}
	return additionalWorkStaff, nil
}

func (r *WorkRepository) UpdateAdditionalWorkStaff(additionalWorkStaff *entity.AdditionalWorkUser) error {
	err := r.GetDB().Model(&entity.AdditionalWorkUser{}).Where("id = ?", additionalWorkStaff.Id).Updates(additionalWorkStaff).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) UpdateDailyWorkStaff(dailyWorkStaff *entity.DailyWorkUser) error {
	err := r.GetDB().Model(&entity.DailyWorkUser{}).Where("id = ?", dailyWorkStaff.Id).Updates(dailyWorkStaff).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) GetDailyWorkStaffById(id uint64) (entity.DailyWorkUser, error) {
	var dailyWorkStaff entity.DailyWorkUser
	err := r.GetDB().
		Preload("DailyWork", "deleted_at IS NULL").
		Preload("Staff.Account").
		Where("id = ?", id).
		First(&dailyWorkStaff).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.DailyWorkUser{}, errx.NotFound("daily work staff not found")
		}
		return entity.DailyWorkUser{}, err
	}
	return dailyWorkStaff, nil
}

func (r *WorkRepository) UpdateDailyWorkStaffStatus(dailyWorkStaff *entity.DailyWorkUser) error {
	err := r.GetDB().Model(&entity.DailyWorkUser{}).Updates(dailyWorkStaff).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) DeleteDailyWork(id uint64) error {
	err := r.GetDB().Delete(&entity.DailyWork{}, id).Error
	if err != nil {
		return err
	}
	return nil
}
