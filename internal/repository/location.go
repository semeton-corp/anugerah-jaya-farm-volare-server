package repository

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"gorm.io/gorm"
)

type LocationRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type ILocationRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetLocations() ([]entity.Location, error)
}

func NewLocationRepository(db *gorm.DB) ILocationRepository {
	return &LocationRepository{
		db: db,
	}
}

func (r *LocationRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *LocationRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *LocationRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *LocationRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *LocationRepository) GetLocations() ([]entity.Location, error) {
	var locations []entity.Location
	err := r.GetDB().Order("created_at DESC").Find(&locations).Error
	return locations, err
}
