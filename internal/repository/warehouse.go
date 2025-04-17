package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type WarehouseRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IWarehouseRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateWarehouseItem(warehouseItem *entity.WarehouseItem) error
	GetWarehouseItem() ([]entity.WarehouseItem, error)

	CreateWarehouseStockItem(stockWarehouseItem *entity.WarehouseStockItem) error
	GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]entity.WarehouseStockItem, error)
	GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (entity.WarehouseStockItem, error)

	UpdateWarehouseStockItem(stockWarehouseItem *entity.WarehouseStockItem) error
	DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) error

	GetWarehouse() ([]entity.Warehouse, error)
}

func NewWarehouseRepository(db *gorm.DB) IWarehouseRepository {
	return &WarehouseRepository{
		db: db,
	}
}

func (r *WarehouseRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *WarehouseRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *WarehouseRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *WarehouseRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *WarehouseRepository) CreateWarehouseItem(warehouseItem *entity.WarehouseItem) error {
	return r.GetDB().Create(warehouseItem).Error
}

func (r *WarehouseRepository) GetWarehouseItem() ([]entity.WarehouseItem, error) {
	var warehouseItems []entity.WarehouseItem
	err := r.GetDB().Find(&warehouseItems).Error
	if err != nil {
		return nil, err
	}
	return warehouseItems, nil
}

func (r *WarehouseRepository) CreateWarehouseStockItem(stockWarehouseItem *entity.WarehouseStockItem) error {
	if err := r.GetDB().Create(stockWarehouseItem).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("stock warehouse item already exists")
		}

		return err
	}

	return nil
}

func (r *WarehouseRepository) GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]entity.WarehouseStockItem, error) {
	var stockWarehouseItems []entity.WarehouseStockItem
	query := r.GetDB()

	if filter.WarehouseId != 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	err := query.Preload("WarehouseItem").Preload("Warehouse.Location").Find(&stockWarehouseItems).Error
	if err != nil {
		return nil, err
	}

	return stockWarehouseItems, nil
}

func (r *WarehouseRepository) GetWarehouse() ([]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	err := r.GetDB().Preload("Location").Find(&warehouses).Error
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r *WarehouseRepository) GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (entity.WarehouseStockItem, error) {
	var stockWarehouseItem entity.WarehouseStockItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("WarehouseItem").Where("warehouse_item_id = ? AND warehouse_id = ?", warehouseItemId, warehouseId).First(&stockWarehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseStockItem{}, errx.NotFound("stock warehouse item not found")
		}
		return entity.WarehouseStockItem{}, err
	}
	return stockWarehouseItem, nil
}

func (r *WarehouseRepository) UpdateWarehouseStockItem(stockWarehouseItem *entity.WarehouseStockItem) error {
	return r.GetDB().Model(entity.WarehouseStockItem{}).Where("warehouse_item_id = ? AND warehouse_id = ?", stockWarehouseItem.WarehouseItemId, stockWarehouseItem.WarehouseId).Updates(stockWarehouseItem).Error
}

func (r *WarehouseRepository) DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) error {
	return r.GetDB().Where("warehouse_id = ? AND warehouse_item_id = ?", warehouseId, warehouseItemId).Delete(&entity.WarehouseStockItem{}).Error
}
