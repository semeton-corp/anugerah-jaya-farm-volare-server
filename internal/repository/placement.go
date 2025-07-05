package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type PlacementRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IPlacementRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateStorePlacement(data *entity.StorePlacement) error
	CreateWarehousePlacement(data *entity.WarehousePlacement) error
	CreateCagePlacementBatch(data []entity.CagePlacement) error

	GetCagePlacementByUserId(userId uuid.UUID) ([]entity.CagePlacement, error)
	GetStorePlacementByUserId(userId uuid.UUID) (entity.StorePlacement, error)
	GetWarehousePlacementByUserId(userId uuid.UUID) (entity.WarehousePlacement, error)

	DeleteCagePlacementByUserId(userId uuid.UUID) error
	DeleteStorePlacementByUserId(userId uuid.UUID) error
	DeleteWarehousePlacementByUserId(userId uuid.UUID) error
}

func NewPlacementRepository(db *gorm.DB) IPlacementRepository {
	return &PlacementRepository{
		db: db,
	}
}

func (r *PlacementRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *PlacementRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *PlacementRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *PlacementRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *PlacementRepository) CreateStorePlacement(data *entity.StorePlacement) error {
	return r.GetDB().Model(&entity.StorePlacement{}).Create(data).Error
}

func (r *PlacementRepository) CreateWarehousePlacement(data *entity.WarehousePlacement) error {
	return r.GetDB().Model(&entity.WarehousePlacement{}).Create(data).Error
}

func (r *PlacementRepository) CreateCagePlacementBatch(data []entity.CagePlacement) error {
	return r.GetDB().Model(&entity.CagePlacement{}).CreateInBatches(data, len(data)).Error
}

func (r *PlacementRepository) GetCagePlacementByUserId(userId uuid.UUID) ([]entity.CagePlacement, error) {
	data := make([]entity.CagePlacement, 0)
	err := r.GetDB().Model(&entity.CagePlacement{}).Preload("User").Preload("Cage").Where("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) GetStorePlacementByUserId(userId uuid.UUID) (entity.StorePlacement, error) {
	data := new(entity.StorePlacement)
	err := r.GetDB().Model(&entity.StorePlacement{}).Preload("User").Preload("Store").Where("user_id = ?", userId).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StorePlacement{}, errx.NotFound("user not have have placement in store")
		}
		return *data, err
	}

	return *data, nil
}

func (r *PlacementRepository) GetWarehousePlacementByUserId(userId uuid.UUID) (entity.WarehousePlacement, error) {
	data := new(entity.WarehousePlacement)
	err := r.GetDB().Model(&entity.WarehousePlacement{}).Preload("User").Preload("Warehouse").Where("user_id = ?", userId).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehousePlacement{}, errx.NotFound("user not have have placement in warehouse")
		}
		return *data, err
	}

	return *data, nil
}

func (r *PlacementRepository) DeleteCagePlacementByUserId(userId uuid.UUID) error {
	return r.GetDB().Where("user_id = ?", userId).Delete(&entity.CagePlacement{}).Error
}

func (r *PlacementRepository) DeleteStorePlacementByUserId(userId uuid.UUID) error {
	return r.GetDB().Where("user_id = ?", userId).Delete(&entity.StorePlacement{}).Error
}

func (r *PlacementRepository) DeleteWarehousePlacementByUserId(userId uuid.UUID) error {
	return r.GetDB().Where("user_id = ?", userId).Delete(&entity.WarehousePlacement{}).Error
}
