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
	GetItemPriceByItemIdAndSaleUnit(itemId uint64, saleUnit enum.SaleUnit) (entity.ItemPrice, error)

	CreateItemPriceDiscount(eggPriceDiscount *entity.ItemPriceDiscount) error
	GetItemPriceDiscounts() ([]entity.ItemPriceDiscount, error)
	GetItemPriceDiscountById(id uint64) (entity.ItemPriceDiscount, error)
	UpdateItemPriceDiscount(eggPriceDiscount *entity.ItemPriceDiscount) error
	DeleteItemPriceDiscount(id uint64) error

	CreateItem(item *entity.Item) error
	GetItems(filter dto.GetItemFilter) ([]entity.Item, error)
	GetItemById(id uint64) (entity.Item, error)
	UpdateItem(warehouseItem *entity.Item) error
	DeleteItem(id uint64) error
	GetItemByNameAndUnit(name string, unit string) (entity.Item, error)
	GetItemByNameAndUnitAndType(name string, unit string, itemType enum.ItemCategory) (entity.Item, error)
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

func (r *ItemRepository) CreateItemPrice(itemPrice *entity.ItemPrice) error {
	return r.GetDB().Create(itemPrice).Error
}

func (r *ItemRepository) CreateItemPriceDiscount(itemPriceDiscount *entity.ItemPriceDiscount) error {
	return r.GetDB().Create(itemPriceDiscount).Error
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

func (r *ItemRepository) UpdateItemPrice(data *entity.ItemPrice) error {
	updates := map[string]interface{}{
		"category":   data.Category,
		"item_id":    data.ItemId,
		"price":      data.Price,
		"updated_by": data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.ItemPrice{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *ItemRepository) UpdateItemPriceDiscount(data *entity.ItemPriceDiscount) error {
	updates := map[string]interface{}{
		"name":                     data.Name,
		"minimum_transaction_user": data.MinimumTransactionUser,
		"total_discount":           data.TotalDiscount,
		"updated_by":               data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.ItemPriceDiscount{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *ItemRepository) DeleteItemPrice(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ItemPrice{}).Error
}

func (r *ItemRepository) DeleteItemPriceDiscount(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.ItemPriceDiscount{}).Error
}

func (r *ItemRepository) CreateItem(warehouseItem *entity.Item) error {
	return r.GetDB().Create(warehouseItem).Error
}

func (r *ItemRepository) GetItems(filter dto.GetItemFilter) ([]entity.Item, error) {
	var warehouseItems []entity.Item

	query := r.GetDB()

	if filter.Categories != nil {
		categories := make([]enum.ItemCategory, 0)
		for _, e := range filter.Categories {
			categories = append(categories, e.Value())
		}
		query = query.Where("category IN ?", categories)
	}

	if filter.ItemNames != nil {
		query = query.Where("name IN ?", filter.ItemNames)
	}

	err := query.Find(&warehouseItems).Error
	if err != nil {
		return nil, err
	}

	return warehouseItems, nil
}

func (r *ItemRepository) GetItemById(id uint64) (entity.Item, error) {
	var warehouseItem entity.Item
	err := r.GetDB().Where("id = ?", id).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("item not found")
		}
		return entity.Item{}, err
	}
	return warehouseItem, nil
}

func (r *ItemRepository) UpdateItem(data *entity.Item) error {
	updates := map[string]interface{}{
		"name":           data.Name,
		"category":       data.Category,
		"unit":           data.Unit,
		"daily_spending": data.DailySpending,
		"updated_by":     data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.Item{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *ItemRepository) DeleteItem(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Item{}).Error
}

func (r *ItemRepository) GetItemByNameAndUnit(name string, unit string) (entity.Item, error) {
	var warehouseItem entity.Item
	err := r.GetDB().Where("name = ? AND unit = ?", name, unit).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("item not found")
		}
		return entity.Item{}, err
	}
	return warehouseItem, nil
}

func (r *ItemRepository) GetItemByNameAndUnitAndType(name string, unit string, category enum.ItemCategory) (entity.Item, error) {
	var warehouseItem entity.Item
	err := r.GetDB().Where("name = ? AND unit = ? AND category = ?", name, unit, category).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Item{}, errx.NotFound("item not found")
		}
		return entity.Item{}, err
	}
	return warehouseItem, nil
}

func (r *ItemRepository) GetItemPriceByItemIdAndSaleUnit(itemId uint64, saleUnit enum.SaleUnit) (entity.ItemPrice, error) {
	var data entity.ItemPrice
	err := r.GetDB().Model(&entity.ItemPrice{}).Where("item_id = ? AND sale_unit = ?", itemId, saleUnit).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.ItemPrice{}, errx.NotFound("item price not found")
		}
		return entity.ItemPrice{}, err
	}

	return data, nil
}
