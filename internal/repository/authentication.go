package repository

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type AuthenticationRepository struct {
	db *gorm.DB
}

type IAuthenticationRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateAccount(account *entity.Account) error
	GetAccountByEmail(email string) (entity.Account, error)
	GetAccountById(id uuid.UUID) (entity.Account, error)
	UpdateAccount(account *entity.Account) error
}

func NewAuthenticationRepository(db *gorm.DB) IAuthenticationRepository {
	return &AuthenticationRepository{
		db: db,
	}
}

func (r *AuthenticationRepository) UseTx(tx bool) {
	if tx {
		r.db = r.db.Begin()
	}
}

func (r *AuthenticationRepository) Commit() error {
	return r.db.Commit().Error
}

func (r *AuthenticationRepository) Rollback() error {
	return r.db.Rollback().Error
}

func (r *AuthenticationRepository) CreateAccount(account *entity.Account) error {
	if err := r.db.Preload("Role").Create(account).Error; err != nil {
		return err
	}

	return r.db.Preload("Role").First(account, "id = ?", account.Id).Error
}

func (r *AuthenticationRepository) GetAccountByEmail(email string) (entity.Account, error) {
	var account entity.Account
	if err := r.db.Where("email = ?", email).First(&account).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *AuthenticationRepository) GetAccountById(id uuid.UUID) (entity.Account, error) {
	var account entity.Account
	if err := r.db.Where("id = ?", id).First(&account).Error; err != nil {
		return entity.Account{}, err
	}
	return account, nil
}

func (r *AuthenticationRepository) UpdateAccount(account *entity.Account) error {
	return r.db.Save(account).Error
}
