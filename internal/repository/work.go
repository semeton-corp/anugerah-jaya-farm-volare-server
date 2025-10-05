package repository

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
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
	GetAdditionalWorkById(id uint64) (entity.AdditionalWork, error)
	DeleteAdditionalWork(id uint64) error
	GetAdditionalWorks(filter dto.GetAdditonalWorkFilter) ([]entity.AdditionalWork, error)

	CreateAdditionalWorkUser(additionalWorkUser *entity.AdditionalWorkUser) error
	CreateAdditionalWorkUserInBatch(additionalWorkUsers []entity.AdditionalWorkUser) error
	GetAdditionalWorkUserById(id uint64) (entity.AdditionalWorkUser, error)
	UpdateAdditionalWorkUser(additionalWorkUser *entity.AdditionalWorkUser) error
	DeleteAdditionalWorkUser(id uint64) error
	DeleteAdditionalWorkUserByAdditionalWorkIdAndUserIds(additionalWorkId uint64, userIds []uuid.UUID) error

	UpdateAdditionalWorkUserByAdditionalWorkId(id uint64, payload map[string]any) error
	GetAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) ([]entity.AdditionalWorkUser, error)
	CountAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) (int64, error)

	GetDailyWorkUserById(id uint64) (entity.DailyWorkUser, error)
	UpdateDailyWorkUser(dailyWorkUser *entity.DailyWorkUser) error
	GetDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) ([]entity.DailyWorkUser, error)
	CountDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) (int64, error)
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

func (r *WorkRepository) CreateAdditionalWorkUserInBatch(additionalWorks []entity.AdditionalWorkUser) error {
	return r.GetDB().Model(&entity.AdditionalWorkUser{}).CreateInBatches(additionalWorks, len(additionalWorks)).Error
}

func (r *WorkRepository) GetDailyWorkByRoleId(roleId uint64) ([]entity.DailyWork, error) {
	var dailyWorks []entity.DailyWork
	err := r.GetDB().
		Preload("Role").
		Where("role_id = ?", roleId).
		Order("created_at DESC").
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
		Preload("Cage.CagePlacement").
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
	query := r.GetDB().Model(&entity.DailyWorkUser{}).
		Joins("JOIN daily_works ON daily_work_users.daily_work_id = daily_works.id").
		Joins("JOIN users ON daily_work_users.user_id = users.id")

	if filter.Month.Value().IsValid() && filter.Year > 0 {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("daily_work_users.created_at >= ? AND daily_work_users.created_at <= ?", startDate, endDate)
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(daily_work_users.created_at) = ?", filter.Date.Value())
	}

	if filter.WithDeleted != nil && *filter.WithDeleted {
		query = query.Unscoped().Where("daily_works.deleted_at IS NOT NULL OR daily_works.deleted_at IS NULL")
	} else {
		query = query.Unscoped().Where("daily_works.deleted_at IS NULL")
	}

	if filter.Page > 0 {
		query = query.Offset(int(filter.Page-1) * int(constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.
		Where("users.id = ?", userId).
		Preload("DailyWork").
		Preload("User").
		Order("created_at DESC").
		Find(&dailyWorkUsers).Error

	if err != nil {
		return nil, err
	}

	return dailyWorkUsers, nil
}

func (r *WorkRepository) CountDailyWorkUserByUserId(userId uuid.UUID, filter dto.GetDailyWorkUserFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.DailyWorkUser{}).
		Joins("JOIN daily_works ON daily_work_users.daily_work_id = daily_works.id").
		Joins("JOIN users ON daily_work_users.user_id = users.id")

	if filter.Month.Value().IsValid() && filter.Year > 0 {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("DATE(daily_work_users.created_at) >= ? AND DATE(daily_work_users.created_at) <= ?", startDate, endDate)
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(daily_work_users.created_at) = ?", filter.Date.Value())
	}

	if filter.WithDeleted != nil && *filter.WithDeleted {
		query = query.Unscoped().Where("daily_works.deleted_at IS NOT NULL OR daily_works.deleted_at IS NULL")
	} else {
		query = query.Unscoped().Where("daily_works.deleted_at IS NULL")
	}

	err := query.
		Where("users.id = ?", userId).
		Count(&count).Error

	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *WorkRepository) GetAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) ([]entity.AdditionalWorkUser, error) {
	var additionalWorks []entity.AdditionalWorkUser
	query := r.GetDB().Model(&entity.AdditionalWorkUser{}).
		Joins("JOIN additional_works ON additional_work_users.additional_work_id = additional_works.id").
		Joins("JOIN users ON additional_work_users.user_id = users.id")

	if filter.Month.Value().IsValid() && filter.Year > 0 {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("additional_work_users.created_at >= ? AND additional_work_users.created_at <= ?", startDate, endDate)
	}

	if filter.WithDeleted != nil && *filter.WithDeleted {
		query = query.Unscoped().Where("additional_works.deleted_at IS NULL OR additional_works.deleted_at IS NOT NULL")
	} else {
		query = query.Unscoped().Where("additional_works.deleted_at IS NULL")
	}

	if filter.IsAdditionalWorkFull {
		subQuery := r.GetDB().
			Table("additional_work_users").
			Select("additional_work_id").
			Group("additional_work_id").
			Having("COUNT(*) = (SELECT slot FROM additional_works WHERE id = additional_work_users.additional_work_id)")

		query = query.Where("additional_works.id IN (?)", subQuery)
	}

	if filter.Page > 0 {
		query = query.Offset(int(filter.Page-1) * int(constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if filter.WithDoneToday {
		query = query.Where(
			"(additional_work_users.is_done = true AND DATE(additional_work_users.finished_at) = CURRENT_DATE) OR (additional_work_users.is_done = false)",
		)
	}

	err := query.Where("users.id = ?", userId).
		Preload("AdditionalWork").
		Preload("User").
		Preload("AdditionalWork.Store").
		Preload("AdditionalWork.Cage").
		Preload("AdditionalWork.Warehouse").
		Preload("AdditionalWork.Location").
		Order("created_at DESC").
		Find(&additionalWorks).Error

	if err != nil {
		return nil, err
	}

	return additionalWorks, nil
}

func (r *WorkRepository) CountAdditionalWorkUserByUserId(userId uuid.UUID, filter dto.GetAdditionalWorkUserFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.AdditionalWorkUser{}).
		Joins("JOIN additional_works ON additional_work_users.additional_work_id = additional_works.id").
		Joins("JOIN users ON additional_work_users.user_id = users.id")

	if filter.Month.Value().IsValid() && filter.Year > 0 {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("additional_work_users.created_at >= ? AND additional_work_users.created_at <= ?", startDate, endDate)
	}

	if filter.WithDeleted != nil && *filter.WithDeleted {
		query = query.Unscoped().Where("additional_works.deleted_at IS NULL OR additional_works.deleted_at IS NOT NULL")
	} else {
		query = query.Unscoped().Where("additional_works.deleted_at IS NULL")
	}

	if filter.IsAdditionalWorkFull {
		query = query.Where("additional_work_users.is_additional_work_full = ?", filter.IsAdditionalWorkFull)
	}

	err := query.Where("users.id = ?", userId).
		Count(&count).Error

	if err != nil {
		return -1, err
	}
	return count, nil
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
	db := r.GetDB().Model(&entity.AdditionalWork{})

	if filter.LocationId == 0 || !filter.LocationType.Value().IsValid() {
		err := db.
			Preload("AdditionalWorkUsers.User.Role").
			Preload("Location").
			Preload("Cage.CagePlacement").
			Preload("Warehouse").
			Preload("Store").
			Where("deleted_at IS NULL").
			Order("created_at DESC").
			Find(&additionalWorks).Error
		return additionalWorks, err
	}

	locationId := filter.LocationId
	locationType := filter.LocationType.Value()

	typeConditions := []string{}
	typeArgs := []interface{}{}

	switch locationType {
	case enum.LocationTypeCage:
		typeConditions = append(typeConditions, "cage_id IN ?")
		typeArgs = append(typeArgs, filter.PlaceIds)
	case enum.LocationTypeStore:
		typeConditions = append(typeConditions, "store_id IN ?")
		typeArgs = append(typeArgs, filter.PlaceIds)
	case enum.LocationTypeWarehouse:
		typeConditions = append(typeConditions, "warehouse_id IN ?")
		typeArgs = append(typeArgs, filter.PlaceIds)
	}

	var condition string
	var args []interface{}

	if len(typeConditions) > 0 {
		condition = "( (location_id = ? AND location_type IS NULL AND warehouse_id IS NULL AND store_id IS NULL AND cage_id IS NULL) OR (location_id = ? AND location_type = ? AND warehouse_id IS NULL AND store_id IS NULL AND cage_id IS NULL) OR (location_id = ? AND (" + strings.Join(typeConditions, " OR ") + ")) )"
		args = append(args, locationId, locationId, locationType, locationId)
		args = append(args, typeArgs...)
	} else {
		condition = "( (location_id = ? AND location_type IS NULL AND warehouse_id IS NULL AND store_id IS NULL AND cage_id IS NULL) OR (location_id = ? AND location_type = ? AND warehouse_id IS NULL AND store_id IS NULL AND cage_id IS NULL) )"
		args = append(args, locationId, locationId, locationType)
	}

	err := db.
		Preload("AdditionalWorkUsers.User.Role").
		Preload("Location").
		Preload("Cage.CagePlacement").
		Preload("Warehouse").
		Preload("Store").
		Where(condition, args...).
		Where("deleted_at IS NULL").
		Order("created_at DESC").
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
		Preload("AdditionalWork.Store").
		Preload("AdditionalWork.Cage").
		Preload("AdditionalWork.Warehouse").
		Preload("AdditionalWork.Location").
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

func (r *WorkRepository) UpdateAdditionalWorkUser(data *entity.AdditionalWorkUser) error {
	updates := map[string]interface{}{
		"is_done":     data.IsDone,
		"note":        data.Note,
		"finished_at": data.FinishedAt,
		"taken_at":    data.TakenAt,
		"updated_by":  data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.AdditionalWorkUser{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WorkRepository) UpdateDailyWorkUser(dwu *entity.DailyWorkUser) error {
	updates := map[string]interface{}{
		"is_done":     dwu.IsDone,
		"note":        dwu.Note,
		"finished_at": dwu.FinishedAt,
		"updated_by":  dwu.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.DailyWorkUser{}).
		Where("id = ?", dwu.Id).
		Updates(updates).Error
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

func (r *WorkRepository) DeleteAdditionalWorkUserByAdditionalWorkIdAndUserIds(additionalWorkId uint64, userIds []uuid.UUID) error {
	return r.GetDB().Model(&entity.AdditionalWorkUser{}).Where("additional_work_id = ? AND user_id IN ?", additionalWorkId, userIds).Error
}
