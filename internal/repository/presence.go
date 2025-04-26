package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
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
	return r.GetDB().Model(&entity.StaffPresence{}).Updates(staffPresence).Error
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

func (r *PresenceRepository) GetStaffPresenceByStaffId(staffId uint64) (entity.StaffPresence, error) {
	var staffPresence entity.StaffPresence
	err := r.GetDB().Preload("Staff.Account.Role").Where("staff_id = ?", staffId).First(&staffPresence).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return staffPresence, errx.NotFound("staff presence not found")
		}
		return staffPresence, err
	}
	return staffPresence, nil
}

func (r *PresenceRepository) GetStaffPresenceInRoleIds(roleIds []uint64) ([]entity.StaffPresence, error) {
	var staffPresences []entity.StaffPresence
	err := r.GetDB().Preload("Staff.Account.Role", "role_id IN ?", roleIds).Find(&staffPresences).Error
	if err != nil {
		return nil, err
	}
	return staffPresences, nil
}
