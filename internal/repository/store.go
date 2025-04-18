package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type StoreRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type IStoreRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetStores() ([]entity.Store, error)
	CreateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error
	GetStoreRequestItemById(id uint64) (entity.StoreRequestItem, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]entity.StoreRequestItem, error)
	UpdateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error

	FirstOrCreateStoreItem(storeItem *entity.StoreItem) error
	UpdateStoreItem(storeItem *entity.StoreItem) error
	GetStoreItems() ([]entity.StoreItem, error)

	CreateStoreSale(storeSale *entity.StoreSale) error
	GetStoreSaleById(id uint64) (entity.StoreSale, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) ([]entity.StoreSale, error)
	UpdateStoreSale(storeSale *entity.StoreSale) error

	CreateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error
}

func NewStoreRepository(db *gorm.DB) IStoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (r *StoreRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *StoreRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *StoreRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *StoreRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *StoreRepository) GetStores() ([]entity.Store, error) {
	var stores []entity.Store
	err := r.GetDB().Preload("Location").Find(&stores).Error
	if err != nil {
		return nil, err
	}
	return stores, nil
}

func (r *StoreRepository) CreateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error {
	return r.GetDB().Create(storeRequestItem).Error
}

func (r *StoreRepository) GetStoreRequestItemById(id uint64) (entity.StoreRequestItem, error) {
	var storeRequestItem entity.StoreRequestItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("WarehouseItem").Preload("Store.Location").First(&storeRequestItem, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreRequestItem{}, errx.NotFound("store request item not found")
		}
		return entity.StoreRequestItem{}, err
	}
	return storeRequestItem, nil
}

func (r *StoreRepository) GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]entity.StoreRequestItem, error) {
	var storeRequestItems []entity.StoreRequestItem
	query := r.GetDB()

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("Warehouse.Location").Preload("WarehouseItem").Preload("Store.Location").Find(&storeRequestItems).Order("status ASC").Error
	if err != nil {
		return nil, err
	}
	return storeRequestItems, nil
}

func (r *StoreRepository) UpdateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error {
	return r.GetDB().Model(entity.StoreRequestItem{}).Where("id = ?", storeRequestItem.Id).Updates(storeRequestItem).Error
}

func (r *StoreRepository) FirstOrCreateStoreItem(storeItem *entity.StoreItem) error {
	return r.GetDB().FirstOrCreate(storeItem, entity.StoreItem{
		WarehouseItemId: storeItem.WarehouseItemId,
		StoreId:         storeItem.StoreId,
	}).Error
}

func (r *StoreRepository) UpdateStoreItem(storeItem *entity.StoreItem) error {
	return r.GetDB().Model(entity.StoreItem{}).Where("store_id = ? AND warehouse_item_id = ?", storeItem.StoreId, storeItem.WarehouseItemId).Updates(storeItem).Error
}

func (r *StoreRepository) GetStoreItems() ([]entity.StoreItem, error) {
	var storeItems []entity.StoreItem
	err := r.GetDB().Preload("Store.Location").Preload("WarehouseItem").Find(&storeItems).Error
	if err != nil {
		return nil, err
	}
	return storeItems, nil
}

func (r *StoreRepository) CreateStoreSale(storeSale *entity.StoreSale) error {
	return r.GetDB().Create(storeSale).Error
}

func (r *StoreRepository) GetStoreSaleById(id uint64) (entity.StoreSale, error) {
	var storeSale entity.StoreSale
	err := r.GetDB().Preload("Payments").Preload("Store.Location").Preload("WarehouseItem").First(&storeSale, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreSale{}, errx.NotFound("store sale not found")
		}
		return entity.StoreSale{}, err
	}
	return storeSale, nil
}

func (r *StoreRepository) GetStoreSales(filter dto.GetStoreSaleFilter) ([]entity.StoreSale, error) {
	var storeSales []entity.StoreSale
	query := r.GetDB()

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("Store.Location").Preload("WarehouseItem").Find(&storeSales).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return storeSales, nil
}

func (r *StoreRepository) CreateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error {
	return r.GetDB().Create(storeSalePayment).Error
}

func (r *StoreRepository) UpdateStoreSale(storeSale *entity.StoreSale) error {
	return r.GetDB().Model(entity.StoreSale{}).Where("id = ?", storeSale.Id).Updates(storeSale).Error
}
