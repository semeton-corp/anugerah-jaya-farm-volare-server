package repository

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
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

	GetRoleByName(name string) (entity.Role, error)

	GetUserOverviewLists(filter *dto.GetUserOverviewListFilter) ([]entity.User, error)
	CountTotalUserOverviewList(filter *dto.GetUserOverviewListFilter) (uint64, error)

	GetUserSalaryPaymentSpesificMonth(userId uuid.UUID, month time.Month, year uint64) (entity.UserSalaryPayment, error)
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

	if filter.ExcluseRoleIds != nil {
		query = query.Where("role_id NOT IN ?", filter.ExcluseRoleIds)
	}

	if err := query.Preload("Role").Find(&users).Error; err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) CountTotalUserOverviewList(filter *dto.GetUserOverviewListFilter) (uint64, error) {
	var totalData int64

	query := r.db.Model(&entity.User{})

	if filter.RoleId > 0 {
		query = query.Where("role_id = ?", filter.RoleId)
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", keyword, keyword)
	}

	err := query.Model(&entity.User{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), err
}

func (r *UserRepository) GetUserOverviewLists(filter *dto.GetUserOverviewListFilter) ([]entity.User, error) {
	users := make([]entity.User, 0)
	query := r.db.Model(&entity.User{})

	if filter.RoleId > 0 {
		query = query.Where("role_id = ?", filter.RoleId)
	}

	if filter.Page > 0 {
		query = query.Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("name ILIKE ? OR email ILIKE ?", keyword, keyword)
	}

	err := query.Preload("Role").Find(&users).Error
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) GetRoleByName(name string) (entity.Role, error) {
	var data entity.Role
	err := r.GetDB().Model(&entity.Role{}).Where("name = ?", name).First(&data).Error
	if err != nil {
		return entity.Role{}, err
	}

	return data, nil
}

func (r *UserRepository) GetUserSalaryPaymentSpesificMonth(userId uuid.UUID, month time.Month, year uint64) (entity.UserSalaryPayment, error) {
	var userSalaryPayment entity.UserSalaryPayment
	err := r.GetDB().Model(&entity.UserSalaryPayment{}).
		Where("user_id = ? AND EXTRACT(month FROM created_at) = ? AND EXTRACT(year FROM created_at) = ?", userId, int(month), int(year)).
		First(&userSalaryPayment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.UserSalaryPayment{}, errx.NotFound("user salary payment not found")
		}
		return entity.UserSalaryPayment{}, err
	}

	return userSalaryPayment, nil
}
