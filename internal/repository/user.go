package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IUserRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetUserById(id uuid.UUID) (entity.User, error)
	UpdateUser(user *entity.User) error
	GetUsers(filter *dto.GetUserListFilter) ([]entity.User, error)
	CountTotalUser(filter *dto.GetUserListFilter) (uint64, error)
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *UserRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *UserRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *UserRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *UserRepository) GetUserById(id uuid.UUID) (entity.User, error) {
	var user entity.User
	if err := r.GetDB().Where(&entity.User{Id: id}).Preload("Location").Preload("Role").First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errx.NotFound("user not found")
		}
		return entity.User{}, err
	}
	return user, nil
}

func (r *UserRepository) UpdateUser(user *entity.User) error {
	if err := r.GetDB().Model(&entity.User{}).Where("id = ?", user.Id).Updates(user).Error; err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUsers(filter *dto.GetUserListFilter) ([]entity.User, error) {
	var users []entity.User
	query := r.GetDB().Model(&entity.User{})

	if filter.RoleId > 0 {
		query = query.Where("role_id = ?", filter.RoleId)
	}

	if filter.LocationId > 0 {
		query = query.Where("location_id = ?", filter.LocationId)
	}

	query = query.Preload("Role")

	if err := query.Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CountTotalUser(filter *dto.GetUserListFilter) (uint64, error) {
	var totalData int64

	query := r.GetDB()

	if filter.RoleId > 0 {
		query = query.Where("role_id = ?", filter.RoleId)
	}

	if filter.LocationId > 0 {
		query = query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Model(&entity.User{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), err
}
