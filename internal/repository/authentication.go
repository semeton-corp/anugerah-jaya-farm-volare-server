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

	CreateAccount(account *entity.Account) error
	GetAccountByEmail(email string) (entity.Account, error)
	GetAccountById(id uuid.UUID) (entity.Account, error)
	UpdateAccount(account *entity.Account) error
	CreateStaff(staff *entity.Staff) error
	DeleteAccount(id uuid.UUID) error
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

func (r *AuthenticationRepository) CreateAccount(account *entity.Account) error {
	err := r.GetDB().Create(account).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("email already exists")
		}
	}

	return nil
}

func (r *AuthenticationRepository) GetAccountByEmail(email string) (entity.Account, error) {
	var account entity.Account
	if err := r.GetDB().Preload("Role").Where("email = ?", email).First(&account).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Account{}, errx.BadRequest("password or email is incorrect")
		}
		return entity.Account{}, err
	}
	return account, nil
}

func (r *AuthenticationRepository) GetAccountById(id uuid.UUID) (entity.Account, error) {
	var account entity.Account
	if err := r.GetDB().Preload("Role").Where("id = ?", id).First(&account).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *AuthenticationRepository) UpdateAccount(account *entity.Account) error {
	return r.GetDB().Model(entity.Account{}).Where("id = ?", account.Id).Updates(&account).Error
}

func (r *AuthenticationRepository) CreateStaff(staff *entity.Staff) error {
	return r.GetDB().Create(staff).Error
}

func (r *AuthenticationRepository) DeleteAccount(id uuid.UUID) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Account{}).Error
}
