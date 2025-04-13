package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type CageRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type ICageRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetCages() ([]entity.Cage, error)
}

func NewCageRepository(db *gorm.DB) ICageRepository {
	return &CageRepository{
		db: db,
	}
}

func (r *CageRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *CageRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *CageRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *CageRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *CageRepository) GetCages() ([]entity.Cage, error) {
	var (
		cages []entity.Cage
		err   error
	)

	err = r.GetDB().Preload("Location").Find(&cages).Error
	if err != nil {
		return nil, err
	}

	return cages, nil
}
