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

	CreateWarehouseItem(warehouseItem *entity.WarehouseItem) error
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]entity.WarehouseItem, error)
	UpdateWarehouseItem(warehouseItem *entity.WarehouseItem) error
	GetWarehouseItemById(id uint64) (entity.WarehouseItem, error)
	GetWarehouseItemByNameAndUnit(name string, unit string) (entity.WarehouseItem, error)
	DeleteWarehouseItem(id uint64) error

	CreateWarehouseStockItem(stockWarehouseItem *entity.WarehouseStockItem) error
	GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]entity.WarehouseStockItem, error)
	GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (entity.WarehouseStockItem, error)

	UpdateWarehouseStockItem(stockWarehouseItem *entity.WarehouseStockItem) error
	DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) error

	CreateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error
	GetWarehouseOrderItemById(id uint64) (entity.WarehouseOrderItem, error)
	GetWarehouseOrderItems(filter dto.GetWarehouseOrderItemFilter) ([]entity.WarehouseOrderItem, error)
	DeleteWarehouseOrderItem(id uint64) error
	UpdateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error

	GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.WarehouseItemCategory) (entity.WarehouseItem, error)
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
	query := r.GetDB().Preload("Location").Find(&warehouses)

	if filter.LocationId > 0 {
		query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Error
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

func (r *WarehouseRepository) CreateWarehouseItem(warehouseItem *entity.WarehouseItem) error {
	return r.GetDB().Create(warehouseItem).Error
}

func (r *WarehouseRepository) GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]entity.WarehouseItem, error) {
	var warehouseItems []entity.WarehouseItem

	query := r.GetDB()

	if filter.Category.Value().IsValid() {
		query = query.Where("category = ?", filter.Category.Value())
	}

	err := query.Find(&warehouseItems).Error
	if err != nil {
		return nil, err
	}

	return warehouseItems, nil
}

func (r *WarehouseRepository) GetWarehouseItemById(id uint64) (entity.WarehouseItem, error) {
	var warehouseItem entity.WarehouseItem
	err := r.GetDB().Where("id = ?", id).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return warehouseItem, nil
}

func (r *WarehouseRepository) UpdateWarehouseItem(warehouseItem *entity.WarehouseItem) error {
	return r.GetDB().Model(entity.WarehouseItem{}).Where("id = ?", warehouseItem.Id).Updates(warehouseItem).Error
}

func (r *WarehouseRepository) DeleteWarehouseItem(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItem{}).Error
}

func (r *WarehouseRepository) GetWarehouseItemByNameAndUnit(name string, unit string) (entity.WarehouseItem, error) {
	var warehouseItem entity.WarehouseItem
	err := r.GetDB().Where("name = ? AND unit = ?", name, unit).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return warehouseItem, nil
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

	if filter.Category.Value().IsValid() {
		query = query.Preload("WarehouseItem", "category = ?", filter.Category)
	} else {
		query = query.Preload("WarehouseItem")
	}

	err := query.Preload("Warehouse.Location").Find(&stockWarehouseItems).Error
	if err != nil {
		return nil, err
	}

	return stockWarehouseItems, nil
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

func (r *WarehouseRepository) GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.WarehouseItemCategory) (entity.WarehouseItem, error) {
	var warehouseItem entity.WarehouseItem
	err := r.GetDB().Where("name = ? AND unit = ? AND type = ?", name, unit, itemType).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return warehouseItem, nil
}
