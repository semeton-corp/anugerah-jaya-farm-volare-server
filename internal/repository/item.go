package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type ItemRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IItemRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	CreateItemPrice(eggPrice *entity.ItemPrice) error
	GetItemPrices() ([]entity.ItemPrice, error)
	GetItemPriceById(id uint64) (entity.ItemPrice, error)
	UpdateItemPrice(eggPrice *entity.ItemPrice) error
	DeleteItemPrice(id uint64) error

	CreateItemPriceDiscount(eggPriceDiscount *entity.ItemPriceDiscount) error
	GetItemPriceDiscounts() ([]entity.ItemPriceDiscount, error)
	GetItemPriceDiscountById(id uint64) (entity.ItemPriceDiscount, error)
	UpdateItemPriceDiscount(eggPriceDiscount *entity.ItemPriceDiscount) error
	DeleteItemPriceDiscount(id uint64) error

	CreateWarehouseItem(warehouseItem *entity.Item) error
	GetWarehouseItems(filter dto.GetItemFilter) ([]entity.Item, error)
	GetWarehouseItemById(id uint64) (entity.Item, error)
	UpdateWarehouseItem(warehouseItem *entity.Item) error
	DeleteWarehouseItem(id uint64) error
	GetWarehouseItemByNameAndUnit(name string, unit string) (entity.Item, error)
	GetItemByNameAndUnitAndType(name string, unit string, itemType enum.WarehouseItemCategory) (entity.Item, error)
}

func NewItemRepository(db *gorm.DB) IItemRepository {
	return &ItemRepository{
		db: db,
	}
}

func (r *ItemRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *ItemRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *ItemRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *ItemRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *ItemRepository) CreateItemPrice(eggPrice *entity.ItemPrice) error {
	return r.GetDB().Create(eggPrice).Error
}

func (r *ItemRepository) CreateItemPriceDiscount(eggPriceDiscount *entity.ItemPriceDiscount) error {
	return r.GetDB().Create(eggPriceDiscount).Error
}

func (r *ItemRepository) GetItemPrices() ([]entity.ItemPrice, error) {
	var eggPrice []entity.ItemPrice
	err := r.GetDB().Preload("Item").Find(&eggPrice).Error
	return eggPrice, err
}

func (r *ItemRepository) GetItemPriceDiscounts() ([]entity.ItemPriceDiscount, error) {
	var eggPriceDiscount []entity.ItemPriceDiscount
	err := r.GetDB().Find(&eggPriceDiscount).Error
	return eggPriceDiscount, err
}

func (r *ItemRepository) GetItemPriceById(id uint64) (entity.ItemPrice, error) {
	var eggPrice entity.ItemPrice
	err := r.GetDB().Preload("Item").Where("id = ?", id).First(&eggPrice).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return entity.ItemPrice{}, errx.NotFound("item price not found")
	}
	return eggPrice, nil
}

func (r *ItemRepository) GetItemPriceDiscountById(id uint64) (entity.ItemPriceDiscount, error) {
	var eggPriceDiscount entity.ItemPriceDiscount
	err := r.GetDB().Where("id = ?", id).First(&eggPriceDiscount).Error
	return eggPriceDiscount, err
}

func (r *ItemRepository) UpdateItemPrice(eggPrice *entity.ItemPrice) error {
	return r.GetDB().Where("id = ?", eggPrice.Id).Updates(eggPrice).Error
}

func (r *ItemRepository) UpdateItemPriceDiscount(eggPriceDiscount *entity.ItemPriceDiscount) error {
	return r.GetDB().Where("id = ?", eggPriceDiscount.Id).Updates(eggPriceDiscount).Error
}

func (r *ItemRepository) DeleteItemPrice(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ItemPrice{}).Error
}

func (r *ItemRepository) DeleteItemPriceDiscount(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ItemPriceDiscount{}).Error
}

func (r *ItemRepository) CreateWarehouseItem(warehouseItem *entity.Item) error {
	return r.GetDB().Create(warehouseItem).Error
}

func (r *ItemRepository) GetWarehouseItems(filter dto.GetItemFilter) ([]entity.Item, error) {
	var warehouseItems []entity.Item

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

func (r *ItemRepository) GetWarehouseItemById(id uint64) (entity.Item, error) {
	var warehouseItem entity.Item
	err := r.GetDB().Where("id = ?", id).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("warehouse item not found")
		}
		return entity.Item{}, err
	}
	return warehouseItem, nil
}

func (r *ItemRepository) UpdateWarehouseItem(warehouseItem *entity.Item) error {
	return r.GetDB().Model(entity.Item{}).Where("id = ?", warehouseItem.Id).Updates(warehouseItem).Error
}

func (r *ItemRepository) DeleteWarehouseItem(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Item{}).Error
}

func (r *ItemRepository) GetWarehouseItemByNameAndUnit(name string, unit string) (entity.Item, error) {
	var warehouseItem entity.Item
	err := r.GetDB().Where("name = ? AND unit = ?", name, unit).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("warehouse item not found")
		}
		return entity.Item{}, err
	}
	return warehouseItem, nil
}

func (r *ItemRepository) GetItemByNameAndUnitAndType(name string, unit string, itemType enum.WarehouseItemCategory) (entity.Item, error) {
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
