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
	CreateWarehouseItemInBatch(warehouseItems *[]entity.WarehouseItem) error
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]entity.WarehouseItem, error)
	GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) (entity.WarehouseItem, error)
	UpdateWarehouseItem(stockWarehouseItem *entity.WarehouseItem) error
	DeleteWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) error

	GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.ItemCategory) (entity.Item, error)
	CountStoreRequestItemByWarehouseId(warehouseId uint64) (int64, error)

	GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) ([]entity.WarehouseItemHistory, error)
	GetWarehouseItemHistoryById(id uint64) (entity.WarehouseItemHistory, error)
	CountTotalWarehouseItemHistory(filter dto.GetWarehouseItemHistoryFilter) (int64, error)

	GetWarehouseSalePaymentById(id uint64) (entity.WarehouseSalePayment, error)
	CreateWarehouseSalePayment(warehouseSalePayment *entity.WarehouseSalePayment) error
	UpdateWarehouseSalePayment(warehouseSalePayment *entity.WarehouseSalePayment) error
	DeleteWarehouseSalePayment(id uint64) error

	CountTotalWarehouseSale(filter dto.GetWarehouseSaleFilter) (uint64, error)
	CreateWarehouseSale(warehouseSale *entity.WarehouseSale) error
	GetWarehouseSaleById(id uint64) (entity.WarehouseSale, error)
	GetWarehouseSales(filter dto.GetWarehouseSaleFilter) ([]entity.WarehouseSale, error)
	UpdateWarehouseSale(warehouseSale *entity.WarehouseSale) error
	DeleteWarehouseSale(id uint64) error

	CreateWarehouseSaleQueue(data *entity.WarehouseSaleQueue) error
	GetWarehouseSaleQueueById(id uint64) (entity.WarehouseSaleQueue, error)
	GetWarehouseSaleQueues(filter dto.GetWarehouseSaleQueueFilter) ([]entity.WarehouseSaleQueue, error)
	DeleteWarehouseSaleQueue(id uint64) error

	CreateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error
	GetWarehouseItemProcurementDrafts() ([]entity.WarehouseItemProcurementDraft, error)
	GetWarehouseItemProcurementDraft(id uint64) (entity.WarehouseItemProcurementDraft, error)
	UpdateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error
	DeleteWarehouseItemProcurementDraft(id uint64) error

	CreateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error
	GetWarehouseItemProcurements(filter dto.GetWarehouseItemProcurementFilter) ([]entity.WarehouseItemProcurement, error)
	GetWarehouseItemProcurement(id uint64) (entity.WarehouseItemProcurement, error)
	UpdateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error
	DeleteWarehouseItemProcurement(id uint64) error
	CountWarehouseItemProcurement(filter dto.GetWarehouseItemProcurementFilter) (int64, error)

	CreateWarehouseItemProcurementPayment(data *entity.WarehouseItemProcurementPayment) error
	GetWarehouseItemProcurementPayment(id uint64) (entity.WarehouseItemProcurementPayment, error)
	UpdateWarehouseItemProcurementPayment(data *entity.WarehouseItemProcurementPayment) error
	DeleteWarehouseItemProcurementPayment(id uint64) error

	CreateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error
	GetWarehouseItemCornProcurementDrafts() ([]entity.WarehouseItemCornProcurementDraft, error)
	GetWarehouseItemCornProcurementDraft(id uint64) (entity.WarehouseItemCornProcurementDraft, error)
	UpdateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error
	DeleteWarehouseItemCornProcurementDraf(id uint64) error

	CreateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error
	UpdateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error
	GetWarehouseItemCornProcurement(id uint64) (entity.WarehouseItemCornProcurement, error)
	GetWarehouseItemCornProcurements(filter dto.GetWarehouseItemCornProcurementFilter) ([]entity.WarehouseItemCornProcurement, error)
	DeleteWarehouseItemCornProcurement(id uint64) error
	CountWarehouseItemCornProcurement(filter dto.GetWarehouseItemCornProcurementFilter) (int64, error)

	CreateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error
	UpdateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error
	GetWarehouseItemCornProcurementPayment(id uint64) (entity.WarehouseItemCornProcurementPayment, error)
	DeleteWarehouseItemCornProcurementPayment(id uint64) error

	CreateWarehouseItemCorn(data *entity.WarehouseItemCorn) error
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
			return errx.BadRequest("warehouse item already exists")
		} else if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errx.BadRequest("invalid warehouse or item")
		}
		return err
	}

	return nil
}

func (r *WarehouseRepository) CreateWarehouseItemInBatch(warehouseItems *[]entity.WarehouseItem) error {
	return r.GetDB().Model(&entity.WarehouseItem{}).CreateInBatches(warehouseItems, len(*warehouseItems)).Error
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

func (r *WarehouseRepository) GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) (entity.WarehouseItem, error) {
	var stockWarehouseItem entity.WarehouseItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("Item").Where("item_id = ? AND warehouse_id = ?", itemId, warehouseId).First(&stockWarehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return stockWarehouseItem, nil
}

func (r *WarehouseRepository) UpdateWarehouseItem(warehouseItem *entity.WarehouseItem) error {
	return r.GetDB().Model(entity.WarehouseItem{}).Where("item_id = ? AND warehouse_id = ?", warehouseItem.ItemId, warehouseItem.WarehouseId).Updates(warehouseItem).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) error {
	return r.GetDB().Where("warehouse_id = ? AND item_id = ?", warehouseId, itemId).Delete(&entity.WarehouseItem{}).Error
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

	err := query.Count(&total).Error
	if err != nil {
		return -1, err
	}

	return total, nil
}

func (r *WarehouseRepository) CountTotalWarehouseSale(filter dto.GetWarehouseSaleFilter) (uint64, error) {
	var totalData int64
	query := r.GetDB().Model(&entity.WarehouseSale{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.PaymentMethod.Value().IsValid() {
		query = query.Where("payment_method = ?", filter.PaymentMethod.Value())
	}

	err := query.Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), nil
}

func (r *WarehouseRepository) GetWarehouseSalePaymentById(id uint64) (entity.WarehouseSalePayment, error) {
	var warehouseSalePayment entity.WarehouseSalePayment
	err := r.GetDB().First(&warehouseSalePayment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseSalePayment{}, errx.NotFound("warehouse sale payment not found")
		}
		return entity.WarehouseSalePayment{}, err
	}
	return warehouseSalePayment, nil
}

func (r *WarehouseRepository) CreateWarehouseSale(warehouseSale *entity.WarehouseSale) error {
	if err := r.GetDB().Create(warehouseSale).Error; err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errx.NotFound("some resources not found")
		}
		return err
	}

	return nil
}

func (r *WarehouseRepository) GetWarehouseSaleById(id uint64) (entity.WarehouseSale, error) {
	var warehouseSale entity.WarehouseSale
	err := r.GetDB().Preload("Payments").Preload("Warehouse.Location").Preload("Customer").Preload("Item").First(&warehouseSale, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseSale{}, errx.NotFound("warehouse sale not found")
		}
		return entity.WarehouseSale{}, err
	}
	return warehouseSale, nil
}

func (r *WarehouseRepository) GetWarehouseSales(filter dto.GetWarehouseSaleFilter) ([]entity.WarehouseSale, error) {
	var warehouseSales []entity.WarehouseSale
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

	err := query.Preload("Warehouse.Location").Preload("Customer").Preload("Item").Find(&warehouseSales).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return warehouseSales, nil
}

func (r *WarehouseRepository) CreateWarehouseSalePayment(warehouseSalePayment *entity.WarehouseSalePayment) error {
	return r.GetDB().Create(warehouseSalePayment).Error
}

func (r *WarehouseRepository) UpdateWarehouseSale(warehouseSale *entity.WarehouseSale) error {
	return r.GetDB().Model(entity.WarehouseSale{}).Where("id = ?", warehouseSale.Id).Updates(warehouseSale).Error
}

func (r *WarehouseRepository) UpdateWarehouseSalePayment(warehouseSalePayment *entity.WarehouseSalePayment) error {
	return r.GetDB().Model(entity.StoreSalePayment{}).Where("id = ?", warehouseSalePayment.Id).Updates(warehouseSalePayment).Error
}

func (r *WarehouseRepository) DeleteWarehouseSale(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.StoreSale{}).Error
}

func (r *WarehouseRepository) DeleteWarehouseSalePayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.StoreSalePayment{}).Error
}

func (r *WarehouseRepository) CreateWarehouseSaleQueue(data *entity.WarehouseSaleQueue) error {
	return r.GetDB().Model(&entity.WarehouseSaleQueue{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseSaleQueueById(id uint64) (entity.WarehouseSaleQueue, error) {
	var warehouseSaleQueue entity.WarehouseSaleQueue
	err := r.GetDB().Model(&entity.WarehouseSaleQueue{}).Preload("Store").Preload("Item").Preload("Customer").Where("id = ?", id).First(&warehouseSaleQueue).Error
	if err != nil {
		return entity.WarehouseSaleQueue{}, err
	}

	return warehouseSaleQueue, nil
}

func (r *WarehouseRepository) DeleteWarehouseSaleQueue(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseSaleQueue{}).Error
}

func (r *WarehouseRepository) GetWarehouseSaleQueues(filter dto.GetWarehouseSaleQueueFilter) ([]entity.WarehouseSaleQueue, error) {
	var warehouseSaleQueues []entity.WarehouseSaleQueue
	query := r.GetDB().Model(&entity.WarehouseSaleQueue{})

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	err := query.Preload("Customer").Preload("Item").Find(&warehouseSaleQueues).Error
	if err != nil {
		return nil, err
	}

	return warehouseSaleQueues, nil
}

func (r *WarehouseRepository) CreateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurement{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseItemProcurementDrafts() ([]entity.WarehouseItemProcurementDraft, error) {
	var data []entity.WarehouseItemProcurementDraft
	err := r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *WarehouseRepository) CountWarehouseItemProcurement(filter dto.GetWarehouseItemProcurementFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.WarehouseItemProcurement{})

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}

	err := query.Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *WarehouseRepository) GetWarehouseItemProcurementDraft(id uint64) (entity.WarehouseItemProcurementDraft, error) {
	var data entity.WarehouseItemProcurementDraft
	err := r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.WarehouseItemProcurementDraft{}, err
	}

	return data, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Updates(&data).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemProcurementDraft(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemProcurementDraft{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurement{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseItemProcurements(filter dto.GetWarehouseItemProcurementFilter) ([]entity.WarehouseItemProcurement, error) {
	var data []entity.WarehouseItemProcurement
	query := r.GetDB().Model(&entity.WarehouseItemProcurement{}).Preload("Warehouse").Preload("Item").Preload("Supplier")

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *WarehouseRepository) GetWarehouseItemProcurement(id uint64) (entity.WarehouseItemProcurement, error) {
	var data entity.WarehouseItemProcurement
	err := r.GetDB().Model(&entity.WarehouseItemProcurement{}).Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.WarehouseItemProcurement{}, err
	}

	return data, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurement{}).Updates(data).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemProcurement(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemProcurement{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemProcurementPayment(data *entity.WarehouseItemProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurementPayment{}).Create(&data).Error
}

func (r *WarehouseRepository) GetWarehouseItemProcurementPayment(id uint64) (entity.WarehouseItemProcurementPayment, error) {
	var data entity.WarehouseItemProcurementPayment
	err := r.GetDB().Model(&entity.WarehouseItemProcurementPayment{}).Where("id = ?", id).First(&data).Error
	if err != nil {
		return entity.WarehouseItemProcurementPayment{}, err
	}

	return data, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemProcurementPayment(data *entity.WarehouseItemProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurementPayment{}).Updates(&data).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemProcurementPayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemProcurementPayment{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementDraft{}).Create(&data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurementDrafts() ([]entity.WarehouseItemCornProcurementDraft, error) {
	var warehouseItemCornProcurementDrafts []entity.WarehouseItemCornProcurementDraft
	err := r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Find(&warehouseItemCornProcurementDrafts).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemCornProcurementDrafts, nil
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurementDraft(id uint64) (entity.WarehouseItemCornProcurementDraft, error) {
	var warehouseItemCornProcurementDraft entity.WarehouseItemCornProcurementDraft
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurementDraft{}).Where("id = ?", id).First(&warehouseItemCornProcurementDraft).Error
	if err != nil {
		return entity.WarehouseItemCornProcurementDraft{}, err
	}

	return warehouseItemCornProcurementDraft, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementDraft{}).Where("id = ?", data.Id).Updates(&data).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemCornProcurementDraf(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemCornProcurementDraft{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurement{}).Create(&data).Error
}

func (r *WarehouseRepository) UpdateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurement{}).Where("id = ?", data.Id).Updates(&data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurement(id uint64) (entity.WarehouseItemCornProcurement, error) {
	var warehouseItemCornProcurement entity.WarehouseItemCornProcurement
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurement{}).Where("id = ?", id).First(&warehouseItemCornProcurement).Error
	if err != nil {
		return entity.WarehouseItemCornProcurement{}, err
	}

	return warehouseItemCornProcurement, nil
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurements(filter dto.GetWarehouseItemCornProcurementFilter) ([]entity.WarehouseItemCornProcurement, error) {
	var warehouseItemCornProcurements []entity.WarehouseItemCornProcurement
	query := r.GetDB().Model(&entity.WarehouseItemCornProcurement{})

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Find(&warehouseItemCornProcurements).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemCornProcurements, nil
}

func (r *WarehouseRepository) CountWarehouseItemCornProcurement(filter dto.GetWarehouseItemCornProcurementFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.WarehouseItemCornProcurement{})

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	err := query.Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *WarehouseRepository) DeleteWarehouseItemCornProcurement(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemCornProcurement{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{}).Create(&data).Error
}

func (r *WarehouseRepository) UpdateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{}).Where("id = ?", data.Id).Updates(&data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurementPayment(id uint64) (entity.WarehouseItemCornProcurementPayment, error) {
	var warehouseItemCornProcurementPayment entity.WarehouseItemCornProcurementPayment
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{}).Where("id = ?", id).First(warehouseItemCornProcurementPayment).Error
	if err != nil {
		return entity.WarehouseItemCornProcurementPayment{}, err
	}

	return warehouseItemCornProcurementPayment, nil
}

func (r *WarehouseRepository) DeleteWarehouseItemCornProcurementPayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemCornProcurementPayment{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCorn(data *entity.WarehouseItemCorn) error {
	return r.GetDB().Model(&entity.WarehouseItemCorn{}).Create(&data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCorns() ([]entity.WarehouseItemCorn, error) {
	var warehouseItemCorns []entity.WarehouseItemCorn
	err := r.GetDB().Model(&entity.WarehouseItemCorn{}).Find(&warehouseItemCorns).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemCorns, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemCorn(data *entity.WarehouseItemCorn) error {
	return r.GetDB().Model(&entity.WarehouseItemCorn{}).Where("id = ?", data.Id).Updates(&data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCorn(id uint64) (entity.WarehouseItemCorn, error) {
	var warehouseItemCorn entity.WarehouseItemCorn
	err := r.GetDB().Model(&entity.WarehouseItemCorn{}).Where("id = ?", id).First(&warehouseItemCorn).Error
	if err != nil {
		return entity.WarehouseItemCorn{}, err
	}

	return warehouseItemCorn, nil
}
