package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type AuthenticationRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IAuthenticationRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateUser(user *entity.User) error
	GetUserByUsername(username string) (entity.User, error)
	GetUserByEmail(email string) (entity.User, error)
	GetUserById(id uuid.UUID) (entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(id uuid.UUID) error
}

func NewAuthenticationRepository(db *gorm.DB) IAuthenticationRepository {
	return &AuthenticationRepository{
		db: db,
	}
}

func (r *AuthenticationRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *AuthenticationRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *AuthenticationRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *AuthenticationRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *AuthenticationRepository) CreateUser(user *entity.User) error {
	err := r.GetDB().Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("user already exists")
		}
	}

	return nil
}

func (r *AuthenticationRepository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	if err := r.GetDB().Preload("Role").Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errx.BadRequest("user not found")
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *AuthenticationRepository) GetUserByEmail(email string) (entity.User, error) {
	var user entity.User
	if err := r.GetDB().Preload("Role").Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, errx.BadRequest("user not found")
		}
		return entity.User{}, err
	}

	return user, nil
}

func (r *AuthenticationRepository) GetUserById(id uuid.UUID) (entity.User, error) {
	var user entity.User
	if err := r.GetDB().Preload("Location").Preload("Role").Where("id = ?", id).First(&user).Error; err != nil {
		return entity.User{}, err
	}
	return user, nil
}

func (r *AuthenticationRepository) UpdateUser(user *entity.User) error {
	updates := map[string]interface{}{
		"username":        user.Username,
		"email":           user.Email,
		"password":        user.Password,
		"location_id":     user.LocationId,
		"role_id":         user.RoleId,
		"photo_profile":   user.PhotoProfile,
		"name":            user.Name,
		"phone_number":    user.PhoneNumber,
		"address":         user.Address,
		"salary_interval": user.SalaryInterval,
		"salary":          user.Salary,
		"updated_by":      user.UpdatedBy,
	}
	return r.GetDB().Model(&entity.User{}).Where("id = ?", user.Id).Updates(updates).Error
}

func (r *AuthenticationRepository) DeleteUser(id uuid.UUID) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.User{}).Error
}
