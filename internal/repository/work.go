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

	SaveDailyWork(dailyWork *entity.DailyWork) error
	GetDailyWorkByRoleId(roleId uint64) ([]entity.DailyWork, error)
	GetDailyWorkSummariesByRoleIds(roleIds []uint64) ([]entity.DailyWorkSummary, error)
	DeleteDailyWork(id uint64) error

	SaveAdditionalWork(additionalWork *entity.AdditionalWork) error
	CreateAdditionalWorUserkInBatch(additionalWorks *[]entity.AdditionalWorkUser) error
	GetAdditionalWorkById(id uint64) (entity.AdditionalWork, error)
	DeleteAdditionalWork(id uint64) error
	GetAdditionalWorks(filter dto.GetAdditonalWorkFilter) ([]entity.AdditionalWork, error)

	CreateAdditionalWorkUser(additionalWorkUser *entity.AdditionalWorkUser) error
	GetAdditionalWorkUserById(id uint64) (entity.AdditionalWorkUser, error)
	UpdateAdditionalWorkUser(additionalWorkUser *entity.AdditionalWorkUser) error
	DeleteAdditionalWorkUser(id uint64) error
	UpdateAdditionalWorkUserByAdditionalWorkId(id uint64, payload map[string]any) error
	GetAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) ([]entity.AdditionalWorkUser, error)

	GetDailyWorkUserById(id uint64) (entity.DailyWorkUser, error)
	UpdateDailyWorkUser(dailyWorkUser *entity.DailyWorkUser) error
	GetDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) ([]entity.DailyWorkUser, error)
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

func (r *WorkRepository) SaveDailyWork(dailyWork *entity.DailyWork) error {
	return r.GetDB().Save(dailyWork).Error
}

func (r *WorkRepository) SaveAdditionalWork(additionalWork *entity.AdditionalWork) error {
	return r.GetDB().Save(additionalWork).Error
}

func (r *WorkRepository) CreateAdditionalWorUserkInBatch(additionalWorks *[]entity.AdditionalWorkUser) error {
	return r.GetDB().Model(&entity.AdditionalWorkUser{}).CreateInBatches(additionalWorks, len(*additionalWorks)).Error
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
		Preload("AdditionalWorkUsers.User.Role").
		Preload("Location").
		Preload("Cage").
		Preload("Warehouse").
		Preload("Store").
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

func (r *WorkRepository) GetDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) ([]entity.DailyWorkUser, error) {
	var dailyWorkUsers []entity.DailyWorkUser
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
		Preload("User", "id = ?", userId).
		Order("created_at DESC").
		Find(&dailyWorkUsers).Error

	if err != nil {
		return nil, err
	}

	return dailyWorkUsers, nil
}

func (r *WorkRepository) GetAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) ([]entity.AdditionalWorkUser, error) {
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

	if filter.IsAdditionalWorkFull {
		query = query.Where("is_additional_work_full = ?", filter.IsAdditionalWorkFull)
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
	query := r.GetDB().Model(&entity.AdditionalWork{}).Joins("JOIN additional_work_users ON additional_works.id = additional_work_users.additional_work_id")

	if filter.ExcludeUserIds != nil {
		query = query.Where("additional_work_users.user_id NOT IN ?", filter.ExcludeUserIds)
	}

	err := query.
		Preload("AdditionalWorkUsers.User.Role").
		Preload("Location").
		Preload("Cage").
		Preload("Warehouse").
		Preload("Store").
		Where("deleted_at IS NULL").
		Find(&additionalWorks).Error

	if err != nil {
		return nil, err
	}
	return additionalWorks, nil
}

func (r *WorkRepository) GetDailyWorkSummariesByRoleIds(roleIds []uint64) ([]entity.DailyWorkSummary, error) {
	var dailyWorkSummaries []entity.DailyWorkSummary

	query := r.GetDB().
		Table("daily_works").
		Select(`daily_works.role_id, roles.name AS role_name,
				COUNT(DISTINCT daily_works.id) AS total_work,
				COUNT(DISTINCT users.id) AS total_user`).
		Joins("LEFT JOIN roles ON daily_works.role_id = roles.id").
		Joins("LEFT JOIN users ON users.role_id = roles.id").
		Group("daily_works.role_id, roles.name").
		Where("daily_works.deleted_at IS NULL").
		Where("daily_works.role_id IN ?", roleIds)

	err := query.Scan(&dailyWorkSummaries).Error

	if err != nil {
		return nil, err
	}
	return dailyWorkSummaries, nil
}

func (r *WorkRepository) CreateAdditionalWorkUser(additionalWorkUser *entity.AdditionalWorkUser) error {
	return r.GetDB().Create(additionalWorkUser).Error
}

func (r *WorkRepository) GetAdditionalWorkUserById(id uint64) (entity.AdditionalWorkUser, error) {
	var additionalWorkUser entity.AdditionalWorkUser
	err := r.GetDB().
		Model(&entity.AdditionalWorkUser{}).
		Preload("AdditionalWork", "deleted_at IS NULL").
		Preload("User").
		Where("id = ?", id).
		First(&additionalWorkUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AdditionalWorkUser{}, errx.NotFound("additional work user not found")
		}
		return entity.AdditionalWorkUser{}, err
	}
	return additionalWorkUser, nil
}

func (r *WorkRepository) UpdateAdditionalWorkUser(additionalWorkUser *entity.AdditionalWorkUser) error {
	err := r.GetDB().Model(&entity.AdditionalWorkUser{}).Where("id = ?", additionalWorkUser.Id).Updates(additionalWorkUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) UpdateDailyWorkUser(dailyWorkUser *entity.DailyWorkUser) error {
	err := r.GetDB().Model(&entity.DailyWorkUser{}).Where("id = ?", dailyWorkUser.Id).Updates(dailyWorkUser).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) GetDailyWorkUserById(id uint64) (entity.DailyWorkUser, error) {
	var dailyWorkUser entity.DailyWorkUser
	err := r.GetDB().
		Preload("DailyWork", "deleted_at IS NULL").
		Preload("User").
		Where("id = ?", id).
		First(&dailyWorkUser).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.DailyWorkUser{}, errx.NotFound("daily work user not found")
		}
		return entity.DailyWorkUser{}, err
	}
	return dailyWorkUser, nil
}

func (r *WorkRepository) DeleteDailyWork(id uint64) error {
	err := r.GetDB().Delete(&entity.DailyWork{}, id).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *WorkRepository) DeleteAdditionalWorkUser(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.AdditionalWorkUser{}).Error
}

func (r *WorkRepository) UpdateAdditionalWorkUserByAdditionalWorkId(id uint64, payload map[string]any) error {
	return r.GetDB().Model(&entity.AdditionalWorkUser{}).Where("additional_work_id = ?", id).Updates(payload).Error
}
