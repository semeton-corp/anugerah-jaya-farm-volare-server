package repository

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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
	return r.GetDB().Preload("Role").Create(account).Error
}

func (r *AuthenticationRepository) GetAccountByEmail(email string) (entity.Account, error) {
	var account entity.Account
	if err := r.GetDB().Preload("Role").Where("email = ?", email).First(&account).Error; err != nil {
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
	return r.GetDB().Save(account).Error
}

func (r *AuthenticationRepository) CreateStaff(staff *entity.Staff) error {
	return r.GetDB().Create(staff).Error
}
