package repository

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type StaffRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IStaffRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetStaffById(id uuid.UUID) (entity.Staff, error)
	UpdateStaff(staff *entity.Staff) error
}

func NewStaffRepository(db *gorm.DB) IStaffRepository {
	return &StaffRepository{
		db: db,
	}
}

func (r *StaffRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *StaffRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *StaffRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *StaffRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *StaffRepository) GetStaffById(id uuid.UUID) (entity.Staff, error) {
	var staff entity.Staff
	if err := r.GetDB().Where(&entity.Staff{Id: id}).First(&staff).Error; err != nil {
		return entity.Staff{}, err
	}
	return staff, nil
}

func (r *StaffRepository) UpdateStaff(staff *entity.Staff) error {
	if err := r.GetDB().Model(&entity.Staff{}).Where("id = ?", staff.Id).Updates(staff).Error; err != nil {
		return err
	}
	return nil
}
