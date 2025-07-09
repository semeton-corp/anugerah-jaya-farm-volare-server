package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
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

	CreateWarehouseItem(stockWarehouseItem *entity.WarehouseItem) error
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]entity.WarehouseItem, error)
	GetWarehouseItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (entity.WarehouseItem, error)
	UpdateWarehouseItem(stockWarehouseItem *entity.WarehouseItem) error
	DeleteWarehouseItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) error

	CreateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error
	GetWarehouseOrderItemById(id uint64) (entity.WarehouseOrderItem, error)
	GetWarehouseOrderItems(filter dto.GetWarehouseOrderItemFilter) ([]entity.WarehouseOrderItem, error)
	DeleteWarehouseOrderItem(id uint64) error
	UpdateWarehouseOrderItem(warehouseOrderItem *entity.WarehouseOrderItem) error

	GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.ItemCategory) (entity.Item, error)
	CountStoreRequestItemByWarehouseId(warehouseId uint64) (int64, error)

	GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) ([]entity.WarehouseItemHistory, error)
	GetWarehouseItemHistoryById(id uint64) (entity.WarehouseItemHistory, error)
	CountTotalWarehouseItemHistory(filter dto.GetWarehouseItemHistoryFilter) (int64, error)
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

func (r *WarehouseRepository) CreateWarehouseItem(warehouseItem *entity.WarehouseItem) error {
	if err := r.GetDB().Create(warehouseItem).Error; err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return errx.BadRequest("stock warehouse item already exists")
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errx.BadRequest("invalid warehouse or item")
		}
		return err
	}

	return nil
}

func (r *WarehouseRepository) GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]entity.WarehouseItem, error) {
	var warehouseItems []entity.WarehouseItem
	query := r.GetDB().Model(&entity.WarehouseItem{}).Joins("JOIN items ON items.id = warehouse_items.item_id")

	if filter.WarehouseId != 0 {
		query = query.Where("warehouse_items.warehouse_id = ?", filter.WarehouseId)
	}

	if filter.Category.Value().IsValid() {
		query = query.Where("items.category = ?", filter.Category)
	}

	if filter.ItemNames != nil {
		query = query.Where("items.name IN ?", filter.ItemNames)
	}

	if filter.Units != nil {
		query = query.Where("items.unit IN ?", filter.Units)
	}

	err := query.Preload("Item").Preload("Warehouse.Location").Find(&warehouseItems).Error
	if err != nil {
		return nil, err
	}

	return warehouseItems, nil
}

func (r *WarehouseRepository) GetWarehouseItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, itemId uint64) (entity.WarehouseItem, error) {
	var stockWarehouseItem entity.WarehouseItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("Item").Where("item_id = ? AND warehouse_id = ?", itemId, warehouseId).First(&stockWarehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("stock warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return stockWarehouseItem, nil
}

func (r *WarehouseRepository) UpdateWarehouseItem(warehouseItem *entity.WarehouseItem) error {
	return r.GetDB().Model(entity.WarehouseItem{}).Where("item_id = ? AND warehouse_id = ?", warehouseItem.ItemId, warehouseItem.WarehouseId).Updates(warehouseItem).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, itemId uint64) error {
	return r.GetDB().Where("warehouse_id = ? AND item_id = ?", warehouseId, itemId).Delete(&entity.WarehouseItem{}).Error
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

func (r *WarehouseRepository) GetWarehouseItemByNameAndUnitAndType(name string, unit string, category enum.ItemCategory) (entity.Item, error) {
	var item entity.Item
	err := r.GetDB().Where("name = ? AND unit = ? AND category = ?", name, unit, category).First(&item).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("warehouse item not found")
		}
		return entity.Item{}, err
	}
	return item, nil
}

func (r *WarehouseRepository) CountStoreRequestItemByWarehouseId(warehouseId uint64) (int64, error) {
	var count int64
	err := r.GetDB().Model(&entity.StoreRequestItem{}).Where("warehouse_id = ?", warehouseId).Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *WarehouseRepository) GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) ([]entity.WarehouseItemHistory, error) {
	warehouseItemHistory := make([]entity.WarehouseItemHistory, 0)
	query := r.GetDB().Model(&entity.WarehouseItemHistory{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("Item").Preload("User").Find(&warehouseItemHistory).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemHistory, nil
}

func (r *WarehouseRepository) GetWarehouseItemHistoryById(id uint64) (entity.WarehouseItemHistory, error) {
	var warehouseItemHistory entity.WarehouseItemHistory
	err := r.GetDB().Model(&entity.WarehouseItemHistory{}).Where("id = ?", id).Preload("Item").Preload("User").First(&warehouseItemHistory).Error
	if err != nil {
		return entity.WarehouseItemHistory{}, err
	}

	return warehouseItemHistory, nil
}

func (r *WarehouseRepository) CountTotalWarehouseItemHistory(filter dto.GetWarehouseItemHistoryFilter) (int64, error) {
	var total int64
	query := r.GetDB().Model(&entity.WarehouseItemHistory{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Count(&total).Error
	if err != nil {
		return -1, err
	}

	return total, nil
}
