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

	UpdateUserPresence(staffPresence *entity.UserPresence) error
	GetUserPresenceById(id uint64) (entity.UserPresence, error)
	GetUserPresenceTodayByUserId(staffId uuid.UUID) (entity.UserPresence, error)
	GetUserPresencesByUserId(staffId uuid.UUID, filter dto.GetPresenceFilter) ([]entity.UserPresence, error)
	GetUserPresenceInRoleIds(roleIds []uint64) ([]entity.UserPresence, error)
	CountTotalUserPresenceByUserId(staffId uuid.UUID, filter dto.GetPresenceFilter) (int64, error)
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

func (r *PresenceRepository) UpdateUserPresence(staffPresence *entity.UserPresence) error {
	return r.GetDB().Model(&entity.UserPresence{}).Where("id = ?", staffPresence.Id).Updates(staffPresence).Error
}

func (r *PresenceRepository) GetUserPresenceById(id uint64) (entity.UserPresence, error) {
	var staffPresence entity.UserPresence
	err := r.GetDB().Preload("User.Role").Preload("User.Location").Where("id = ?", id).First(&staffPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staffPresence, errx.NotFound("user presence not found")
		}
		return staffPresence, err
	}
	return staffPresence, nil
}

func (r *PresenceRepository) GetUserPresencesByUserId(staffId uuid.UUID, filter dto.GetPresenceFilter) ([]entity.UserPresence, error) {
	var staffPresences []entity.UserPresence
	query := r.GetDB().Preload("User.Role").Preload("User.Location").Where("user_id = ?", staffId)

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Find(&staffPresences).Order("created_at DESC").Error
	if err != nil {
		return staffPresences, err
	}
	return staffPresences, nil
}

func (r *PresenceRepository) GetUserPresenceInRoleIds(roleIds []uint64) ([]entity.UserPresence, error) {
	var staffPresences []entity.UserPresence
	err := r.GetDB().
		Joins("JOIN users ON users.id = user_presences.user_id").
		Where("users.role_id IN ?", roleIds).
		Preload("User.Role").
		Preload("User.Location").
		Find(&staffPresences).Error
	if err != nil {
		return nil, err
	}
	return staffPresences, nil
}

func (r *PresenceRepository) GetUserPresenceTodayByUserId(staffId uuid.UUID) (entity.UserPresence, error) {
	var staffPresence entity.UserPresence
	err := r.GetDB().Preload("User.Role").Preload("User.Location").Where("user_id = ? AND DATE(created_at) = ?", staffId, time.Now()).First(&staffPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staffPresence, errx.NotFound("user presence not found")
		}
		return staffPresence, err
	}
	return staffPresence, nil
}

func (r *PresenceRepository) CountTotalUserPresenceByUserId(staffId uuid.UUID, filter dto.GetPresenceFilter) (int64, error) {
	var totalData int64
	query := r.GetDB().Model(&entity.UserPresence{}).Where("user_id = ?", staffId)

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
