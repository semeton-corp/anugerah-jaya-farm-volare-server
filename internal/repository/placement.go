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

	GetCagePlacementByCageId(cageId uint64) ([]entity.CagePlacement, error)
	GetStorePlacementByStoreId(storeId uint64) ([]entity.StorePlacement, error)
	GetWarehousePlacementByWarehouseId(warehouseId uint64) ([]entity.WarehousePlacement, error)

	DeleteCagePlacementByUserIdAndCageId(userId uuid.UUID, cageId uint64) error
	DeleteStorePlacementByUserId(userId uuid.UUID) error
	DeleteWarehousePlacementByUserId(userId uuid.UUID) error

	DeleteCagePlacementByCageId(cageId uint64) error
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
	err := r.GetDB().Model(&entity.StorePlacement{}).Create(data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errx.BadRequest("invalid store or user")
		}
		return err
	}
	return nil
}

func (r *PlacementRepository) CreateWarehousePlacement(data *entity.WarehousePlacement) error {
	return r.GetDB().Model(&entity.WarehousePlacement{}).Create(data).Error
}

func (r *PlacementRepository) CreateCagePlacementBatch(data []entity.CagePlacement) error {
	return r.GetDB().Model(&entity.CagePlacement{}).CreateInBatches(data, len(data)).Error
}

func (r *PlacementRepository) GetCagePlacementByUserId(userId uuid.UUID) ([]entity.CagePlacement, error) {
	data := make([]entity.CagePlacement, 0)
	err := r.GetDB().Model(&entity.CagePlacement{}).Preload("User.Role").Preload("Cage.Location").Where("user_id = ?", userId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) GetStorePlacementByUserId(userId uuid.UUID) (entity.StorePlacement, error) {
	data := new(entity.StorePlacement)
	err := r.GetDB().Model(&entity.StorePlacement{}).Preload("User.Role").Preload("Store.Location").Where("user_id = ?", userId).First(&data).Error
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
	err := r.GetDB().Model(&entity.WarehousePlacement{}).Preload("User.Role").Preload("Warehouse.Location").Where("user_id = ?", userId).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehousePlacement{}, errx.NotFound("user not have have placement in warehouse")
		}
		return *data, err
	}

	return *data, nil
}

func (r *PlacementRepository) DeleteCagePlacementByUserIdAndCageId(userId uuid.UUID, cageId uint64) error {
	return r.GetDB().Where("user_id = ? AND cage_id = ?", userId, cageId).Delete(&entity.CagePlacement{}).Error
}

func (r *PlacementRepository) DeleteStorePlacementByUserId(userId uuid.UUID) error {
	return r.GetDB().Where("user_id = ?", userId).Delete(&entity.StorePlacement{}).Error
}

func (r *PlacementRepository) DeleteWarehousePlacementByUserId(userId uuid.UUID) error {
	return r.GetDB().Where("user_id = ?", userId).Delete(&entity.WarehousePlacement{}).Error
}

func (r *PlacementRepository) GetCagePlacementByCageId(cageId uint64) ([]entity.CagePlacement, error) {
	data := make([]entity.CagePlacement, 0)
	err := r.GetDB().Model(&entity.CagePlacement{}).Preload("User.Location").Preload("User.Role").Preload("Cage.Location").Where("cage_id = ?", cageId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) GetStorePlacementByStoreId(storeId uint64) ([]entity.StorePlacement, error) {
	data := make([]entity.StorePlacement, 0)
	err := r.GetDB().Model(&entity.StorePlacement{}).Preload("User.Location").Preload("User.Role").Preload("Store.Location").Where("store_id = ?", storeId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) GetWarehousePlacementByWarehouseId(warehouseId uint64) ([]entity.WarehousePlacement, error) {
	data := make([]entity.WarehousePlacement, 0)
	err := r.GetDB().Model(&entity.WarehousePlacement{}).Preload("User.Location").Preload("User.Role").Preload("Warehouse.Location").Where("warehouse_id = ?", warehouseId).Find(&data).Error
	if err != nil {
		return data, err
	}

	return data, nil
}

func (r *PlacementRepository) DeleteCagePlacementByCageId(cageId uint64) error {
	return r.GetDB().Where("cage_id = ?", cageId).Delete(&entity.CagePlacement{}).Error
}
