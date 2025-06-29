package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
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

	CreateWarehouse(warehouse *entity.Warehouse) error
	UpdateWarehouse(warehouse *entity.Warehouse) error
	GetWarehouseById(id uint64) (entity.Warehouse, error)
	DeleteWarehouse(id uint64) error
	GetWarehouses(filter dto.GetWarehouseFilter) ([]entity.Warehouse, error)

	CreateWarehouseStockItem(stockWarehouseItem *entity.WarehouseItem) error
	GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]entity.WarehouseItem, error)
	GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (entity.WarehouseItem, error)

	UpdateWarehouseStockItem(stockWarehouseItem *entity.WarehouseItem) error
	DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) error

	CreateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error
	GetWarehouseOrderItemById(id uint64) (entity.WarehouseOrderItem, error)
	GetWarehouseOrderItems(filter dto.GetWarehouseOrderItemFilter) ([]entity.WarehouseOrderItem, error)
	DeleteWarehouseOrderItem(id uint64) error
	UpdateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error

	GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.WarehouseItemCategory) (entity.Item, error)
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

func (r *WarehouseRepository) GetWarehouses(filter dto.GetWarehouseFilter) ([]entity.Warehouse, error) {
	var warehouses []entity.Warehouse
	query := r.GetDB().Preload("WarehousePlacement").Preload("Location")

	if filter.LocationId > 0 {
		query = query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Find(&warehouses).Error
	if err != nil {
		return nil, err
	}

	return warehouses, nil
}

func (r *WarehouseRepository) CreateWarehouse(warehouse *entity.Warehouse) error {
	return r.GetDB().Create(&warehouse).Error
}

func (r *WarehouseRepository) GetWarehouseById(id uint64) (entity.Warehouse, error) {
	var warehouse entity.Warehouse
	err := r.GetDB().Model(&entity.Warehouse{}).Preload("Location").Where("id = ?", id).First(&warehouse).Error
	if err != nil {
		return entity.Warehouse{}, err
	}

	return warehouse, nil
}

func (r *WarehouseRepository) UpdateWarehouse(warehouse *entity.Warehouse) error {
	return r.GetDB().Model(&entity.Warehouse{}).Where("id = ?", warehouse.Id).Updates(&warehouse).Error
}

func (r *WarehouseRepository) DeleteWarehouse(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Warehouse{}).Error
}

func (r *WarehouseRepository) CreateWarehouseStockItem(stockWarehouseItem *entity.WarehouseItem) error {
	if err := r.GetDB().Create(stockWarehouseItem).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("stock warehouse item already exists")
		}
		return err
	}

	return nil
}

func (r *WarehouseRepository) GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]entity.WarehouseItem, error) {
	var stockWarehouseItems []entity.WarehouseItem
	query := r.GetDB()

	if filter.WarehouseId != 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.Category.Value().IsValid() {
		query = query.Preload("Item", "category = ?", filter.Category)
	} else {
		query = query.Preload("Item")
	}

	err := query.Preload("Warehouse.Location").Find(&stockWarehouseItems).Error
	if err != nil {
		return nil, err
	}

	return stockWarehouseItems, nil
}

func (r *WarehouseRepository) GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (entity.WarehouseItem, error) {
	var stockWarehouseItem entity.WarehouseItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("WarehouseItem").Where("warehouse_item_id = ? AND warehouse_id = ?", warehouseItemId, warehouseId).First(&stockWarehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("stock warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return stockWarehouseItem, nil
}

func (r *WarehouseRepository) UpdateWarehouseStockItem(warehouseItem *entity.WarehouseItem) error {
	return r.GetDB().Model(entity.WarehouseItem{}).Where("item_id = ? AND warehouse_id = ?", warehouseItem.ItemId, warehouseItem.WarehouseId).Updates(warehouseItem).Error
}

func (r *WarehouseRepository) DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) error {
	return r.GetDB().Where("warehouse_id = ? AND warehouse_item_id = ?", warehouseId, warehouseItemId).Delete(&entity.WarehouseItem{}).Error
}

func (r *WarehouseRepository) CreateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error {
	return r.GetDB().Create(warehouseOrderItem).Error
}

func (r *WarehouseRepository) GetWarehouseOrderItemById(id uint64) (entity.WarehouseOrderItem, error) {
	var warehouseOrderItem entity.WarehouseOrderItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("WarehouseItem").Preload("Supplier", func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("WarehouseItemId")
	}).Where("id = ?", id).First(&warehouseOrderItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseOrderItem{}, errx.NotFound("warehouse order item not found")
		}
		return entity.WarehouseOrderItem{}, err
	}
	return warehouseOrderItem, nil
}

func (r *WarehouseRepository) GetWarehouseOrderItems(filter dto.GetWarehouseOrderItemFilter) ([]entity.WarehouseOrderItem, error) {
	var warehouseOrderItems []entity.WarehouseOrderItem
	query := r.GetDB().Model(&entity.WarehouseOrderItem{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(taken_at) = ?", filter.Date.Value())
	}

	if filter.IsTaken {
		query = query.Where("taken_at IS NOT NULL")
	} else {
		query = query.Where("taken_at IS NULL")
	}

	err := query.Preload("Warehouse.Location").Preload("WarehouseItem").Preload("Supplier").Find(&warehouseOrderItems).Error
	if err != nil {
		return nil, err
	}

	return warehouseOrderItems, nil
}

func (r *WarehouseRepository) DeleteWarehouseOrderItem(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseOrderItem{}).Error
}

func (r *WarehouseRepository) UpdateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error {
	return r.GetDB().Model(entity.WarehouseOrderItem{}).Where("id = ?", warehouseOrderItem.Id).Updates(&warehouseOrderItem).Error
}

func (r *WarehouseRepository) GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.WarehouseItemCategory) (entity.Item, error) {
	var warehouseItem entity.Item
	err := r.GetDB().Where("name = ? AND unit = ? AND type = ?", name, unit, itemType).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("warehouse item not found")
		}
		return entity.Item{}, err
	}
	return warehouseItem, nil
}
