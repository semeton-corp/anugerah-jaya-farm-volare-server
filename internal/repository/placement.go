package repository

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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

	CreateStorePlacementBatch(data []entity.StorePlacement) error
	CreateWarehousePlacementBatch(data []entity.WarehousePlacement) error
	CreateCagePlacementBatch(data []entity.CagePlacement) error

	GetCagePlacementByUserId(userId uuid.UUID) ([]entity.CagePlacement, error)
	GetStorePlacementByUserId(userId uuid.UUID) ([]entity.StorePlacement, error)
	GetWarehousePlacementByUserId(userId uuid.UUID) ([]entity.WarehousePlacement, error)

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

func (r *PlacementRepository) CreateStorePlacementBatch(data []entity.StorePlacement) error {
	return r.GetDB().Model(&entity.StorePlacement{}).CreateInBatches(data, len(data)).Error
}

func (r *PlacementRepository) CreateWarehousePlacementBatch(data []entity.WarehousePlacement) error {
	return r.GetDB().Model(&entity.WarehousePlacement{}).CreateInBatches(data, len(data)).Error
}

func (r *PlacementRepository) CreateCagePlacementBatch(data []entity.CagePlacement) error {
	return r.GetDB().Model(&entity.CagePlacement{}).CreateInBatches(data, len(data)).Error
}

func (r *PlacementRepository) GetCagePlacementByUserId(userId uuid.UUID) ([]entity.CagePlacement, error) {
	data := make([]entity.CagePlacement, 0)
	err := r.GetDB().Model(&entity.CagePlacement{}).Where("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) GetStorePlacementByUserId(userId uuid.UUID) ([]entity.StorePlacement, error) {
	data := make([]entity.StorePlacement, 0)
	err := r.GetDB().Model(&entity.StorePlacement{}).Where("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) GetWarehousePlacementByUserId(userId uuid.UUID) ([]entity.WarehousePlacement, error) {
	data := make([]entity.WarehousePlacement, 0)
	err := r.GetDB().Model(&entity.WarehousePlacement{}).Where("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
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
