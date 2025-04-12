package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IRoleRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetRoles() ([]entity.Role, error)
}

func NewRoleRepository(db *gorm.DB) IRoleRepository {
	return &RoleRepository{
		db: db,
	}
}

func (r *RoleRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *RoleRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *RoleRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *RoleRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *RoleRepository) GetRoles() ([]entity.Role, error) {
	var (
		roles []entity.Role
		err   error
	)

	err = r.GetDB().Find(&roles).Error
	if err != nil {
		return nil, err
	}

	return roles, nil
}
