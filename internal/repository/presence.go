package repository

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
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

	GetCageLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error)
	GetStoreLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error)
	GetWarehouseLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error)
	GetUserPresenceWorkDetailSummary(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]entity.UserPresenceSummary, error)
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

func (r *PresenceRepository) UpdateUserPresence(userPresence *entity.UserPresence) error {
	return r.GetDB().Model(&entity.UserPresence{}).Where("id = ?", userPresence.Id).Updates(userPresence).Error
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

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
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

func (r *PresenceRepository) GetCageLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error) {
	var locationPresenceSummaries []entity.LocationPresenceSummary

	query := r.GetDB().Table("locations").
		Select("locations.id AS place_id, locations.name AS place_name, users.id AS user_id, user_presences.status AS presence_status").
		Joins("LEFT JOIN users ON locations.id = users.location_id").
		Joins("LEFT JOIN user_presences ON users.id = user_presences.user_id").
		Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM cage_placements)")

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(user_presences.created_at) = ?", filter.Date.Value())
	}

	err := query.Scan(&locationPresenceSummaries).Error
	if err != nil {
		return nil, err
	}

	return locationPresenceSummaries, nil
}

func (r *PresenceRepository) GetStoreLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error) {
	var locationPresenceSummaries []entity.LocationPresenceSummary

	query := r.GetDB().Table("stores").
		Select("stores.id AS place_id, stores.name AS place_name, user_presences.user_id AS user_id, user_presences.status AS presence_status").
		Joins("LEFT JOIN store_placements ON store_placements.store_id = stores.id").
		Joins("LEFT JOIN user_presences ON store_placements.user_id = user_presences.user_id").
		Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM store_placements)")

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(user_presences.created_at) = ?", filter.Date.Value())
	}

	err := query.Scan(&locationPresenceSummaries).Error
	if err != nil {
		return nil, err
	}

	return locationPresenceSummaries, nil
}

func (r *PresenceRepository) GetWarehouseLocationPresenceSummaries(filter dto.GetLocationPresenceSummaryFilter) ([]entity.LocationPresenceSummary, error) {
	var locationPresenceSummaries []entity.LocationPresenceSummary

	query := r.GetDB().Table("warehouses").
		Select("warehouses.id AS place_id, warehouses.name AS place_name, user_presences.user_id AS user_id, user_presences.status AS presence_status").
		Joins("LEFT JOIN warehouse_placements ON warehouse_placements.warehouse_id = warehouses.id").
		Joins("LEFT JOIN user_presences ON warehouse_placements.user_id = user_presences.user_id").
		Where("user_presences.user_id IN (SELECT DISTINCT user_id FROM warehouse_placements)")

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(user_presences.created_at) = ?", filter.Date.Value())
	}

	err := query.Scan(&locationPresenceSummaries).Error
	if err != nil {
		return nil, err
	}

	return locationPresenceSummaries, nil
}

func (r *PresenceRepository) GetUserPresenceWorkDetailSummary(filter dto.GetUserPresenceWorkDetailSummaryFilter) ([]entity.UserPresenceSummary, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month))

	db := r.GetDB().Table("user_presences").
		Select(`users.id AS user_id, users.name AS user_name, users.photo_profile AS user_photo_profile, users.email AS user_email, roles.name AS role_name, users.status, COUNT(*) AS total_presence`).
		Joins(`LEFT JOIN users ON user_presences.user_id = users.id`).
		Joins(`LEFT JOIN roles ON users.role_id = roles.id`).
		Where(`DATE(user_presences.created_at) >= ? AND DATE(user_presences.created_at) <= ?`, startDate, endDate).
		Group(`users.id, users.name, users.photo_profile, users.email, roles.name, users.status`)

	switch strings.ToLower(filter.PlaceType) {
	case "cage":
		db = db.Where(`user_id IN (?)`, r.GetDB().Table("cage_placements").
			Select("user_id").
			Joins(`LEFT JOIN cages ON cage_placements.cage_id = cages.id`).
			Where("cages.location_id = ?", filter.PlaceId))
	case "store":
		db = db.Where(`user_id IN (?)`, r.GetDB().Table("store_placements").
			Select("user_id").
			Where("store_id = ?", filter.PlaceId))
	case "warehouse":
		db = db.Where(`user_id IN (?)`, r.GetDB().Table("warehouse_placements").
			Select("user_id").
			Where("warehouse_id = ?", filter.PlaceId))
	default:
		return nil, fmt.Errorf("unsupported place type: %s", filter.PlaceType)
	}

	var results []entity.UserPresenceSummary
	if err := db.Scan(&results).Error; err != nil {
		return nil, err
	}
	return results, nil
}
