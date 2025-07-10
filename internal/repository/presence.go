package repository

import (
	"errors"
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
