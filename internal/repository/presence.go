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

	UpdateStaffPresence(staffPresence *entity.StaffPresence) error
	GetStaffPresenceById(id uint64) (entity.StaffPresence, error)
	GetStaffPresenceTodayByStaffId(staffId uuid.UUID) (entity.StaffPresence, error)
	GetStaffPresenceByStaffId(staffId uuid.UUID, filter dto.GetPresenceFilter) ([]entity.StaffPresence, error)
	GetStaffPresenceInRoleIds(roleIds []uint64) ([]entity.StaffPresence, error)
	CountTotalStaffPresenceByStaffId(staffId uuid.UUID, filter dto.GetPresenceFilter) (int64, error)
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

func (r *PresenceRepository) UpdateStaffPresence(staffPresence *entity.StaffPresence) error {
	return r.GetDB().Model(&entity.StaffPresence{}).Where("id = ?", staffPresence.Id).Updates(staffPresence).Error
}

func (r *PresenceRepository) GetStaffPresenceById(id uint64) (entity.StaffPresence, error) {
	var staffPresence entity.StaffPresence
	err := r.GetDB().Preload("Staff.Account.Role").Where("id = ?", id).First(&staffPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staffPresence, errx.NotFound("staff presence not found")
		}
		return staffPresence, err
	}
	return staffPresence, nil
}

func (r *PresenceRepository) GetStaffPresenceByStaffId(staffId uuid.UUID, filter dto.GetPresenceFilter) ([]entity.StaffPresence, error) {
	var staffPresences []entity.StaffPresence
	query := r.GetDB().Preload("Staff.Account.Role").Where("staff_id = ?", staffId)

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Find(&staffPresences).Order("created_at DESC").Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staffPresences, errx.NotFound("staff presence not found")
		}
		return staffPresences, err
	}
	return staffPresences, nil
}

func (r *PresenceRepository) GetStaffPresenceInRoleIds(roleIds []uint64) ([]entity.StaffPresence, error) {
	var staffPresences []entity.StaffPresence
	err := r.GetDB().Preload("Staff.Account.Role", "role_id IN ?", roleIds).Find(&staffPresences).Error
	if err != nil {
		return nil, err
	}
	return staffPresences, nil
}

func (r *PresenceRepository) GetStaffPresenceTodayByStaffId(staffId uuid.UUID) (entity.StaffPresence, error) {
	var staffPresence entity.StaffPresence
	err := r.GetDB().Preload("Staff.Account.Role").Where("staff_id = ? AND DATE(created_at) = ?", staffId, time.Now()).First(&staffPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staffPresence, errx.NotFound("staff presence not found")
		}
		return staffPresence, err
	}
	return staffPresence, nil
}

func (r *PresenceRepository) CountTotalStaffPresenceByStaffId(staffId uuid.UUID, filter dto.GetPresenceFilter) (int64, error) {
	var totalData int64
	query := r.GetDB().Model(&entity.StaffPresence{}).Where("staff_id = ?", staffId)

	if filter.Month.Value().IsValid() {
		startDate, endDate := util.GetStartDayAndEndDayByMonthFilter(filter.Month.Value(), int(filter.Year))
		query = query.Where("created_at >= ? AND created_at <= ?", startDate, endDate)
	}

	err := query.Model(&entity.StaffPresence{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return totalData, nil
}
