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

	CreateStore(store *entity.Store) error
	UpdateStore(store *entity.Store) error
	DeleteStore(id uint64) error
	GetStoreById(id uint64) (entity.Store, error)
	GetStores(filter dto.GetStoreFilter) ([]entity.Store, error)

	CreateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error
	GetStoreRequestItemById(id uint64) (entity.StoreRequestItem, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]entity.StoreRequestItem, error)
	UpdateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error
	CountTotalStoreRequestItem(filter dto.GetStoreRequestItemFilter) (uint64, error)

	FirstOrCreateStoreItem(storeItem *entity.StoreItem) error
	UpdateStoreItem(storeItem *entity.StoreItem) error
	GetStoreItems(filter dto.GetStoreItemFilter) ([]entity.StoreItem, error)

	CreateStoreSale(storeSale *entity.StoreSale) error
	GetStoreSaleById(id uint64) (entity.StoreSale, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) ([]entity.StoreSale, error)
	UpdateStoreSale(storeSale *entity.StoreSale) error
	CountTotalStoreSale(filter dto.GetStoreSaleFilter) (uint64, error)

	CreateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error
	GetStoreSalePaymentById(id uint64) (entity.StoreSalePayment, error)
	UpdateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error
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

func (r *StoreRepository) CreateStore(store *entity.Store) error {
	return r.GetDB().Model(&entity.Store{}).Create(&store).Error
}

func (r *StoreRepository) UpdateStore(store *entity.Store) error {
	return r.GetDB().Model(&entity.Store{}).Where("id = ?", store.Id).Updates(&store).Error
}

func (r *StoreRepository) DeleteStore(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Store{}).Error
}

func (r *StoreRepository) GetStoreById(id uint64) (entity.Store, error) {
	var data entity.Store
	err := r.GetDB().Model(&entity.Store{}).Preload("Location").Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.Store{}, err
	}

	return data, nil
}

func (r *StoreRepository) GetStores(filter dto.GetStoreFilter) ([]entity.Store, error) {
	stores := make([]entity.Store, 0)
	query := r.GetDB()

	if filter.LocationId > 0 {
		query = query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Preload("StorePlacement").Preload("Location").Find(&stores).Error
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
		ItemId:  storeItem.ItemId,
		StoreId: storeItem.StoreId,
	}).Error
}

func (r *StoreRepository) UpdateStoreItem(storeItem *entity.StoreItem) error {
	return r.GetDB().Model(entity.StoreItem{}).Where("store_id = ? AND warehouse_item_id = ?", storeItem.StoreId, storeItem.ItemId).Updates(storeItem).Error
}

func (r *StoreRepository) GetStoreItems(filter dto.GetStoreItemFilter) ([]entity.StoreItem, error) {
	var storeItems []entity.StoreItem
	query := r.GetDB()

	if filter.StoreId > 0 {
		query = query.Where("store_id = ?", filter.StoreId)
	}

	if filter.Category.Value().IsValid() {
		query = query.Preload("Item", "category = ?", filter.Category)
	} else {
		query = query.Preload("Item")
	}

	err := query.Preload("Store.Location").Find(&storeItems).Error
	if err != nil {
		return nil, err
	}

	return storeItems, nil
}

func (r *StoreRepository) CreateStoreSale(storeSale *entity.StoreSale) error {
	if err := r.GetDB().Create(storeSale).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errx.NotFound("some resources not found")
		}
		return err
	}

	return nil
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

	if filter.PaymentMethod.Value().IsValid() {
		query = query.Where("payment_method = ?", filter.PaymentMethod.Value())
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

func (r *StoreRepository) GetStoreSalePaymentById(id uint64) (entity.StoreSalePayment, error) {
	var storeSalePayment entity.StoreSalePayment
	err := r.GetDB().First(&storeSalePayment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreSalePayment{}, errx.NotFound("store sale payment not found")
		}
		return entity.StoreSalePayment{}, err
	}
	return storeSalePayment, nil
}

func (r *StoreRepository) UpdateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error {
	return r.GetDB().Model(entity.StoreSalePayment{}).Where("id = ?", storeSalePayment.Id).Updates(storeSalePayment).Error
}

func (r *StoreRepository) CountTotalStoreRequestItem(filter dto.GetStoreRequestItemFilter) (uint64, error) {
	var totalData int64
	query := r.GetDB()
	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	err := query.Model(&entity.StoreRequestItem{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), nil
}

func (r *StoreRepository) CountTotalStoreSale(filter dto.GetStoreSaleFilter) (uint64, error) {
	var totalData int64
	query := r.GetDB()

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if filter.PaymentMethod.Value().IsValid() {
		query = query.Where("payment_method = ?", filter.PaymentMethod.Value())
	}

	err := query.Model(&entity.StoreSale{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), nil
}
