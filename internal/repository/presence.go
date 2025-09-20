package repository

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"gorm.io/gorm"
)

type PresenceRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IPresenceRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	UpdateUserPresence(userPresence *entity.UserPresence) error
	GetUserPresenceById(id uint64) (entity.UserPresence, error)
	GetUserPresenceTodayByUserId(userId uuid.UUID) (entity.UserPresence, error)
	GetUserPresencesByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) ([]entity.UserPresence, error)
	GetUserPresenceInRoleIds(roleIds []uint64) ([]entity.UserPresence, error)
	CountTotalUserPresenceByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) (int64, error)
	GetUserPresences(filter dto.GetUserPresenceFilter) ([]entity.UserPresence, error)
	UpdateSubmissionPresenceStatusUserIds(ids []uint64, submissionPresenceStatus enum.SubmissionPresenceStatus) error

	GetLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error)

	GetUserPresenceSummaries(filter dto.GetUserPresenceSummaryFilter) ([]entity.UserPresenceSummary, error)

	GetUserPresenceWorkDetailSummaries(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]entity.UserPresenceWorkDetailSummary, error)

	GetLocationUserPresence(filter dto.GetLocationUserPresenceFilter) ([]entity.UserPresence, error)
}

func NewPresenceRepository(db *gorm.DB) IPresenceRepository {
	return &PresenceRepository{
		db: db,
	}
}

func (r *PresenceRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *PresenceRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *PresenceRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *PresenceRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *PresenceRepository) UpdateUserPresence(data *entity.UserPresence) error {
	updates := map[string]interface{}{
		"user_id":                    data.UserId,
		"start_time":                 data.StartTime,
		"end_time":                   data.EndTime,
		"status":                     data.Status,
		"note":                       data.Note,
		"evidence":                   data.Evidence,
		"submission_presence_status": data.SubmissionPresenceStatus,
		"updated_by":                 data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.UserPresence{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *PresenceRepository) GetUserPresenceById(id uint64) (entity.UserPresence, error) {
	var userPresence entity.UserPresence
	err := r.GetDB().Preload("User.Role").Preload("User.Location").Where("id = ?", id).First(&userPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userPresence, errx.NotFound("user presence not found")
		}
		return userPresence, err
	}
	return userPresence, nil
}

func (r *PresenceRepository) GetUserPresencesByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) ([]entity.UserPresence, error) {
	var userPresences []entity.UserPresence
	query := r.GetDB().Preload("User.Role").Preload("User.Location").Where("user_id = ?", userId)

	if filter.Month.Value().IsValid() && filter.Year > 0 {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if filter.PresenceStatus.Value().IsValid() {
		query = query.Where("status = ?", filter.PresenceStatus.Value())
	}

	err := query.Find(&userPresences).Order("created_at DESC").Error
	if err != nil {
		return userPresences, err
	}
	return userPresences, nil
}

func (r *PresenceRepository) GetUserPresenceInRoleIds(roleIds []uint64) ([]entity.UserPresence, error) {
	var userPresences []entity.UserPresence
	err := r.GetDB().
		Order("created_at DESC").
		Joins("JOIN users ON users.id = user_presences.user_id").
		Where("users.role_id IN ?", roleIds).
		Preload("User.Role").
		Preload("User.Location").
		Find(&userPresences).Error
	if err != nil {
		return nil, err
	}
	return userPresences, nil
}

func (r *PresenceRepository) GetUserPresenceTodayByUserId(userId uuid.UUID) (entity.UserPresence, error) {
	var userPresence entity.UserPresence
	err := r.GetDB().Preload("User.Role").Preload("User.Location").Where("user_id = ? AND DATE(created_at) = ?", userId, time.Now()).First(&userPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return userPresence, errx.NotFound("user presence not found")
		}
		return userPresence, err
	}
	return userPresence, nil
}

func (r *PresenceRepository) CountTotalUserPresenceByUserId(userId uuid.UUID, filter dto.GetPresenceFilter) (int64, error) {
	var totalData int64
	query := r.GetDB().Model(&entity.UserPresence{}).Where("user_id = ?", userId)

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	err := query.Model(&entity.UserPresence{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return totalData, nil
}

func (r *PresenceRepository) GetLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error) {
	var locationPresenceSummaries []entity.LocationPresenceSummary

	if !filter.LocationType.Value().IsValid() {
		return nil, errors.New("invalid location type")
	}

	query := r.GetDB().Model(&entity.LocationPresenceSummary{})

	switch filter.LocationType.Value() {
	case enum.LocationTypeCage:
		query = query.Table("locations").
			Select(`roles.id AS role_id, roles.name AS role_name, 
			        locations.id AS place_id, locations.name AS place_name, 
			        users.id AS user_id, 
			        user_presences.status AS presence_status, 
			        user_presences.submission_presence_status as submission_presence_status`).
			Joins("LEFT JOIN users ON locations.id = users.location_id").
			Joins("LEFT JOIN user_presences ON users.id = user_presences.user_id").
			Joins("LEFT JOIN roles ON users.role_id = roles.id").
			Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM cage_placements)")

	case enum.LocationTypeStore:
		query = query.Table("stores").
			Select(`roles.id AS role_id, roles.name AS role_name, 
			        stores.id AS place_id, stores.name AS place_name, 
			        user_presences.user_id AS user_id, 
			        user_presences.status AS presence_status, 
			        user_presences.submission_presence_status as submission_presence_status`).
			Joins("LEFT JOIN store_placements ON store_placements.store_id = stores.id").
			Joins("LEFT JOIN user_presences ON store_placements.user_id = user_presences.user_id").
			Joins("LEFT JOIN users ON user_presences.user_id = users.id").
			Joins("LEFT JOIN roles ON users.role_id = roles.id").
			Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM store_placements)")

	case enum.LocationTypeWarehouse:
		query = query.Table("warehouses").
			Select(`roles.id AS role_id, roles.name AS role_name, 
			        warehouses.id AS place_id, warehouses.name AS place_name, 
			        user_presences.user_id AS user_id, 
			        user_presences.status AS presence_status, 
			        user_presences.submission_presence_status as submission_presence_status`).
			Joins("LEFT JOIN warehouse_placements ON warehouse_placements.warehouse_id = warehouses.id").
			Joins("LEFT JOIN user_presences ON warehouse_placements.user_id = user_presences.user_id").
			Joins("LEFT JOIN users ON user_presences.user_id = users.id").
			Joins("LEFT JOIN roles ON users.role_id = roles.id").
			Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM warehouse_placements)")

	case enum.LocationTypeSite:
		query = query.Table("locations").
			Select(`roles.id AS role_id, roles.name AS role_name, 
			        locations.id AS place_id, locations.name AS place_name, 
			        users.id AS user_id, 
			        user_presences.status AS presence_status, 
			        user_presences.submission_presence_status as submission_presence_status`).
			Joins("LEFT JOIN users ON locations.id = users.location_id").
			Joins("LEFT JOIN user_presences ON users.id = user_presences.user_id").
			Joins("LEFT JOIN roles ON users.role_id = roles.id").
			Where("roles.name = 'Kepala Kandang'")

	case enum.LocationTypeUnassigned:
		query = query.Table("locations").
			Select(`roles.id AS role_id, roles.name AS role_name, 
			        locations.id AS place_id, locations.name AS place_name, 
			        users.id AS user_id, 
			        user_presences.status AS presence_status, 
			        user_presences.submission_presence_status as submission_presence_status`).
			Joins("LEFT JOIN users ON locations.id = users.location_id").
			Joins("LEFT JOIN user_presences ON users.id = user_presences.user_id").
			Joins("LEFT JOIN roles ON users.role_id = roles.id").
			Where("roles.name NOT IN ('Kepala Kandang', 'Owner') AND user_presences.user_id NOT IN (SELECT DISTINCT user_id FROM warehouse_placements) AND user_presences.user_id NOT IN (SELECT DISTINCT user_id FROM store_placements) AND user_presences.user_id NOT IN (SELECT DISTINCT user_id FROM cage_placements)")
	default:
		return nil, errors.New("unsupported location type")
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(user_presences.created_at) = ?", filter.Date.Value())
	}

	if filter.LocationId > 0 {
		query = query.Where("locations.id = ?", filter.LocationId)
	}

	err := query.Scan(&locationPresenceSummaries).Error
	if err != nil {
		return nil, err
	}

	return locationPresenceSummaries, nil
}

func (r *PresenceRepository) GetUserPresenceSummaries(filter dto.GetUserPresenceSummaryFilter) ([]entity.UserPresenceSummary, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month))

	db := r.GetDB().Table("user_presences").
		Select(`
			users.id AS user_id,
			users.name AS user_name,
			users.photo_profile AS user_photo_profile,
			users.email AS user_email,
			roles.name AS role_name,
			SUM(CASE WHEN user_presences.status = 1 THEN 1 ELSE 0 END) AS total_present,
			SUM(CASE WHEN user_presences.status = 2 THEN 1 ELSE 0 END) AS total_sick,
			SUM(CASE WHEN user_presences.status = 3 THEN 1 ELSE 0 END) AS total_permission,
			SUM(CASE WHEN user_presences.status = 4 THEN 1 ELSE 0 END) AS total_alpha
		`).
		Joins(`LEFT JOIN users ON user_presences.user_id = users.id`).
		Joins(`LEFT JOIN roles ON users.role_id = roles.id`).
		Where(`DATE(user_presences.created_at) >= ? AND DATE(user_presences.created_at) <= ?`, startDate, endDate).
		Where("roles.id = ?", filter.RoleId).
		Group(`users.id, users.name, users.photo_profile, users.email, roles.name, user_presences.status`)

	switch filter.PlaceType.Value() {
	case enum.LocationTypeCage:
		db = db.Where(`user_id IN (?)`, r.GetDB().Table("cage_placements").
			Select("user_id").
			Joins(`LEFT JOIN cages ON cage_placements.cage_id = cages.id`).
			Where("cages.location_id = ?", filter.PlaceId))
	case enum.LocationTypeStore:
		db = db.Where(`user_id IN (?)`, r.GetDB().Table("store_placements").
			Select("user_id").
			Where("store_id = ?", filter.PlaceId))
	case enum.LocationTypeWarehouse:
		db = db.Where(`user_id IN (?)`, r.GetDB().Table("warehouse_placements").
			Select("user_id").
			Where("warehouse_id = ?", filter.PlaceId))
	default:
		return nil, fmt.Errorf("unsupported place type: %s", filter.PlaceType.Value().String())
	}

	var results []entity.UserPresenceSummary
	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *PresenceRepository) GetUserPresenceWorkDetailSummaries(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]entity.UserPresenceWorkDetailSummary, error) {
	db := r.GetDB().Table("users").
		Select(`
			users.id AS user_id,
			users.name AS user_name,
			users.photo_profile AS user_photo_profile,
			users.email AS user_email,
			roles.name AS role_name,
			user_presences.status AS presence_status,
			user_presences.start_time AS presence_start_time,
			user_presences.end_time AS presence_end_time,
			COUNT(DISTINCT additional_work_users.id) AS total_additional_work_users,
			COUNT(DISTINCT CASE WHEN additional_work_users.is_done THEN additional_work_users.id END) AS total_done_additional_work_users,
			COUNT(DISTINCT daily_work_users.id) AS total_daily_work_users,
			COUNT(DISTINCT CASE WHEN daily_work_users.is_done THEN daily_work_users.id END) AS total_done_daily_work_users
		`).
		Joins(`LEFT JOIN roles ON users.role_id = roles.id`).
		Joins(`LEFT JOIN additional_work_users ON additional_work_users.user_id = users.id`).
		Joins(`LEFT JOIN daily_work_users ON daily_work_users.user_id = users.id`).
		Joins(`LEFT JOIN user_presences ON user_presences.user_id = users.id 
			AND DATE(user_presences.created_at) = ?`, filter.Date.Value()).
		Where("roles.id = ?", filter.RoleId)

	switch filter.PlaceType.Value() {
	case enum.LocationTypeCage:
		db = db.
			Joins(`LEFT JOIN cage_placements ON cage_placements.user_id = users.id`).
			Joins(`LEFT JOIN cages ON cage_placements.cage_id = cages.id`).
			Where("cages.location_id = ?", filter.PlaceId)
	case enum.LocationTypeStore:
		db = db.
			Joins(`LEFT JOIN store_placements ON store_placements.user_id = users.id`).
			Where("store_placements.store_id = ?", filter.PlaceId)
	case enum.LocationTypeWarehouse:
		db = db.
			Joins(`LEFT JOIN warehouse_placements ON warehouse_placements.user_id = users.id`).
			Where("warehouse_placements.warehouse_id = ?", filter.PlaceId)
	default:
		return nil, fmt.Errorf("unsupported place type: %s", filter.PlaceType.Value().String())
	}

	db = db.Group(`users.id, users.name, users.photo_profile, users.email, roles.name, user_presences.status, user_presences.start_time, user_presences.end_time`)

	var results []entity.UserPresenceWorkDetailSummary
	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}

func (r *PresenceRepository) GetLocationUserPresence(filter dto.GetLocationUserPresenceFilter) ([]entity.UserPresence, error) {
	var data []entity.UserPresence

	query := r.GetDB().Model(&entity.UserPresence{})

	if filter.LocationType.Value().IsValid() {
		switch filter.LocationType.Value() {
		case enum.LocationTypeCage:
			query = query.
				Joins("LEFT JOIN users ON users.id = user_presences.user_id").
				Where("user_presences.user_id IN (SELECT DISTINCT cage_placements.user_id FROM cage_placements) AND users.location_id = ? AND users.role_id = ? AND user_presences.status = ? AND user_presences.submission_presence_status = ?", filter.PlaceId, filter.RoleId, filter.PresenceStatus.Value(), filter.SubmissionPresenceStatus.Value())

		case enum.LocationTypeStore:
			query = query.
				Joins("LEFT JOIN store_placements ON store_placements.user_id = user_presences.user_id").
				Joins("LEFT JOIN stores ON store_placements.store_id = stores.id").
				Joins("LEFT JOIN users ON user_presences.user_id = users.id").
				Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM store_placements) AND stores.id = ? AND users.role_id = ? AND user_presences.status = ? AND user_presences.submission_presence_status = ?",
					filter.PlaceId, filter.RoleId, filter.PresenceStatus.Value(), filter.SubmissionPresenceStatus.Value())

		case enum.LocationTypeWarehouse:
			query = query.
				Joins("LEFT JOIN warehouse_placements ON warehouse_placements.user_id = user_presences.user_id").
				Joins("LEFT JOIN warehouses ON warehouse_placements.store_id = warehouses.id").
				Joins("LEFT JOIN users ON user_presences.user_id = users.id").
				Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM warehouse_placements) AND warehouses.id = ? AND users.role_id = ? AND user_presences.status = ? AND user_presences.submission_presence_status = ?", filter.PlaceId, filter.RoleId, filter.PresenceStatus.Value(), filter.SubmissionPresenceStatus.Value())

		case enum.LocationTypeSite:
			query = query.
				Joins("LEFT JOIN users ON users.id = user_presences.user_id").
				Where("users.location_id = ? AND users.role_id = ? AND user_presences.status = ? AND user_presences.submission_presence_status = ?",
					filter.PlaceId, filter.RoleId, filter.PresenceStatus.Value(), filter.SubmissionPresenceStatus.Value())

		case enum.LocationTypeUnassigned:
			query = query.
				Joins("LEFT JOIN users ON users.id = user_presences.user_id").
				Where("roles.name NOT IN ('Kepala Kandang', 'Owner') AND user_presences.user_id NOT IN (SELECT DISTINCT user_id FROM warehouse_placements) AND user_presences.user_id NOT IN (SELECT DISTINCT user_id FROM store_placements) AND user_presences.user_id NOT IN (SELECT DISTINCT user_id FROM cage_placements) users.location_id = ? AND users.role_id = ? AND user_presences.status = ? AND user_presences.submission_presence_status = ?",
					filter.PlaceId, filter.RoleId, filter.PresenceStatus.Value(), filter.SubmissionPresenceStatus.Value())
		default:
			return nil, errors.New("unsupported location type")
		}
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(user_presences.created_at) = ?", filter.Date.Value())
	}

	err := query.Preload("user_presences.created_at DESC").Preload("User").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *PresenceRepository) GetUserPresences(filter dto.GetUserPresenceFilter) ([]entity.UserPresence, error) {
	var userPresences []entity.UserPresence

	query := r.GetDB().Model(&entity.UserPresence{})

	if filter.UserIds != nil {
		query = query.Where("user_id IN ?", filter.UserIds)
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Ids != nil {
		query = query.Where("id In ?", filter.Ids)
	}

	err := query.Order("created_at DESC").Preload("User").Find(&userPresences).Error
	if err != nil {
		return nil, err
	}

	return userPresences, nil
}

func (r *PresenceRepository) UpdateSubmissionPresenceStatusUserIds(ids []uint64, submissionPresenceStatus enum.SubmissionPresenceStatus) error {
	return r.GetDB().Model(entity.UserPresence{}).Where("id IN ?", ids).Updates(map[string]any{
		"submission_presence_status": submissionPresenceStatus,
	}).Error
}
