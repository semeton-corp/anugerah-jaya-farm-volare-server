package repository

import (
	"strings"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
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
	GetStaffs(filter *dto.GetStaffFilter) ([]entity.Staff, error)
	CountTotalStaff(filter *dto.GetStaffFilter) (uint64, error)
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
	if err := r.GetDB().Where(&entity.Staff{Id: id}).Preload("Account.Role").First(&staff).Error; err != nil {
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

func (r *StaffRepository) GetStaffs(filter *dto.GetStaffFilter) ([]entity.Staff, error) {
	var staffs []entity.Staff
	query := r.GetDB().Model(&entity.Staff{})

	if filter.Keyword != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Keyword)+"%")
	}

	if filter.RoleId != "" {
		query = query.
			Joins("JOIN accounts ON accounts.id = staffs.account_id").
			Where("accounts.role_id = ?", filter.RoleId)
	}

	if filter.Page != 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	query = query.Preload("Account.Role")

	if err := query.Find(&staffs).Error; err != nil {
		return nil, err
	}

	return staffs, nil
}

func (r *StaffRepository) CountTotalStaff(filter *dto.GetStaffFilter) (uint64, error) {
	var totalData int64

	query := r.GetDB()

	if filter.Keyword != "" {
		query = query.Where("LOWER(name) LIKE ?", "%"+strings.ToLower(filter.Keyword)+"%")
	}

	if filter.RoleId != "" {
		query = query.
			Joins("JOIN accounts ON accounts.id = staffs.account_id").
			Where("accounts.role_id = ?", filter.RoleId)
	}

	err := query.Model(&entity.Staff{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), err
}
