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

	CreateWarehouseItem(warehouseItem *entity.WarehouseItem) error
	FirstOrCreateWarehouseItem(warehouseItem *entity.WarehouseItem) (entity.WarehouseItem, error)
	CreateWarehouseItemInBatch(warehouseItems *[]entity.WarehouseItem) error
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]entity.WarehouseItem, error)
	GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) (entity.WarehouseItem, error)
	GetWarehouseItemCorn(id uint64) (entity.WarehouseItemCorn, error)
	UpdateWarehouseItem(warehouseItem *entity.WarehouseItem) error
	DeleteWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) error

	GetWarehouseItemByNameAndUnitAndType(name string, unit string, itemType enum.ItemCategory) (entity.Item, error)

	GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) ([]entity.WarehouseItemHistory, error)
	GetWarehouseItemHistoryById(id uint64) (entity.WarehouseItemHistory, error)
	CountTotalWarehouseItemHistory(filter dto.GetWarehouseItemHistoryFilter) (int64, error)

	GetWarehouseSalePaymentById(id uint64) (entity.WarehouseSalePayment, error)
	CreateWarehouseSalePaymentInBatch(data *[]entity.WarehouseSalePayment) error
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

	CreateWarehouseItemProcurementDraftInBatch(data *[]entity.WarehouseItemProcurementDraft) error
	CreateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error
	GetWarehouseItemProcurementDrafts(filter dto.GetWarehouseItemProcurementDraftFilter) ([]entity.WarehouseItemProcurementDraft, error)
	GetWarehouseItemProcurementDraft(id uint64) (entity.WarehouseItemProcurementDraft, error)
	UpdateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error
	DeleteWarehouseItemProcurementDraft(id uint64) error

	CreateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error
	GetWarehouseItemProcurements(filter dto.GetWarehouseItemProcurementFilter) ([]entity.WarehouseItemProcurement, error)
	GetWarehouseItemProcurement(id uint64) (entity.WarehouseItemProcurement, error)
	UpdateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error
	DeleteWarehouseItemProcurement(id uint64) error
	CountWarehouseItemProcurement(filter dto.GetWarehouseItemProcurementFilter) (int64, error)

	CreateWarehouseItemProcurementPaymentInBatch(data *[]entity.WarehouseItemProcurementPayment) error
	CreateWarehouseItemProcurementPayment(data *entity.WarehouseItemProcurementPayment) error
	GetWarehouseItemProcurementPayment(id uint64) (entity.WarehouseItemProcurementPayment, error)
	UpdateWarehouseItemProcurementPayment(data *entity.WarehouseItemProcurementPayment) error
	DeleteWarehouseItemProcurementPayment(id uint64) error

	CreateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error
	GetWarehouseItemCornProcurementDrafts(filter dto.GetWarehouseItemCornProcurementDraftFilter) ([]entity.WarehouseItemCornProcurementDraft, error)
	GetWarehouseItemCornProcurementDraft(id uint64) (entity.WarehouseItemCornProcurementDraft, error)
	UpdateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error
	DeleteWarehouseItemCornProcurementDraft(id uint64) error

	CreateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error
	UpdateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error
	GetWarehouseItemCornProcurement(id uint64) (entity.WarehouseItemCornProcurement, error)
	GetWarehouseItemCornProcurements(filter dto.GetWarehouseItemCornProcurementFilter) ([]entity.WarehouseItemCornProcurement, error)
	DeleteWarehouseItemCornProcurement(id uint64) error
	CountWarehouseItemCornProcurement(filter dto.GetWarehouseItemCornProcurementFilter) (int64, error)
	CountQuantityWarehouseItemCornByWarehouseId(warehouseId uint64) (float64, error)

	CreateWarehouseItemCornProcurementPaymentInBatch(data *[]entity.WarehouseItemCornProcurementPayment) error
	CreateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error
	UpdateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error
	GetWarehouseItemCornProcurementPayment(id uint64) (entity.WarehouseItemCornProcurementPayment, error)
	DeleteWarehouseItemCornProcurementPayment(id uint64) error

	CreateWarehouseItemCorn(data *entity.WarehouseItemCorn) error
	UpdateWarehouseItemCorn(data *entity.WarehouseItemCorn) error
	GetWarehouseItemCorns(filter dto.GetWarehouseItemCornFilter) ([]entity.WarehouseItemCorn, error)

	GetWarehouseItemCornPrices() ([]entity.WarehouseItemCornPrice, error)
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
	query := r.GetDB().Preload("WarehousePlacement").Preload("WarehouseItems").Preload("Location")

	if filter.LocationId > 0 {
		query = query.Where("location_id = ?", filter.LocationId)
	}

	err := query.Order("created_at DESC").Find(&warehouses).Error
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
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Warehouse{}, errx.NotFound("warehouse not found")
		}
		return entity.Warehouse{}, err
	}

	return warehouse, nil
}

func (r *WarehouseRepository) UpdateWarehouse(data *entity.Warehouse) error {
	updates := map[string]interface{}{
		"location_id":   data.LocationId,
		"name":          data.Name,
		"corn_capacity": data.CornCapacity,
		"updated_by":    data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.Warehouse{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
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

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_items.warehouse_id = ?", filter.WarehouseId)
	}

	if filter.LocationId > 0 {
		query = query.Joins("JOIN warehouses ON warehouse_items.warehouse_id = warehouses.id JOIN locations ON locations.id = warehouses.location_id").Where("location_id = ?", filter.LocationId)
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

	if filter.WarehouseIds != nil {
		query = query.Where("warehouse_items.warehouse_id IN ?", filter.WarehouseIds)
	}

	err := query.Order("created_at DESC").Preload("Item").Preload("Warehouse.Location").Find(&warehouseItems).Error
	if err != nil {
		return nil, err
	}

	return warehouseItems, nil
}

func (r *WarehouseRepository) GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) (entity.WarehouseItem, error) {
	var warehouseItem entity.WarehouseItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("Item").Where("item_id = ? AND warehouse_id = ?", itemId, warehouseId).First(&warehouseItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItem{}, errx.NotFound("warehouse item not found")
		}
		return entity.WarehouseItem{}, err
	}
	return warehouseItem, nil
}

func (r *WarehouseRepository) UpdateWarehouseItem(data *entity.WarehouseItem) error {
	updates := map[string]interface{}{
		"quantity":   data.Quantity,
		"expired_at": data.ExpiredAt,
		"updated_by": data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseItem{}).
		Where("item_id = ? AND warehouse_id = ?", data.ItemId, data.WarehouseId).
		Updates(updates).Error
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

func (r *WarehouseRepository) FirstOrCreateWarehouseItem(warehouseItem *entity.WarehouseItem) (entity.WarehouseItem, error) {
	err := r.GetDB().Model(&entity.WarehouseItem{}).FirstOrCreate(warehouseItem, &entity.WarehouseItem{ItemId: warehouseItem.ItemId, WarehouseId: warehouseItem.WarehouseId}).Error
	if err != nil {
		return entity.WarehouseItem{}, err
	}

	return *warehouseItem, nil
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

	err := query.Order("created_at DESC").Preload("User").Find(&warehouseItemHistory).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemHistory, nil
}

func (r *WarehouseRepository) GetWarehouseItemHistoryById(id uint64) (entity.WarehouseItemHistory, error) {
	var warehouseItemHistory entity.WarehouseItemHistory
	err := r.GetDB().Model(&entity.WarehouseItemHistory{}).Where("id = ?", id).Preload("User").First(&warehouseItemHistory).Error
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

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
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
	query := r.GetDB().Model(&entity.WarehouseSale{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	err := query.Preload("Warehouse.Location").Preload("Customer").Preload("Item").Order("created_at DESC").Find(&warehouseSales).Error
	if err != nil {
		return nil, err
	}
	return warehouseSales, nil
}

func (r *WarehouseRepository) CreateWarehouseSalePaymentInBatch(data *[]entity.WarehouseSalePayment) error {
	return r.GetDB().Model(&entity.WarehouseSalePayment{}).CreateInBatches(data, len(*data)).Error
}

func (r *WarehouseRepository) CreateWarehouseSalePayment(warehouseSalePayment *entity.WarehouseSalePayment) error {
	return r.GetDB().Model(&entity.WarehouseSalePayment{}).Create(warehouseSalePayment).Error
}

func (r *WarehouseRepository) UpdateWarehouseSale(data *entity.WarehouseSale) error {
	updates := map[string]interface{}{
		"customer_id":           data.CustomerId,
		"item_id":               data.ItemId,
		"warehouse_id":          data.WarehouseId,
		"quantity":              data.Quantity,
		"sale_unit":             data.SaleUnit,
		"price":                 data.Price,
		"total_price":           data.TotalPrice,
		"discount":              data.Discount,
		"send_date":             data.SendDate,
		"payment_type":          data.PaymentType,
		"payment_status":        data.PaymentStatus,
		"is_send":               data.IsSend,
		"deadline_payment_date": data.DeadlinePaymentDate,
		"updated_by":            data.UpdatedBy,
		"paid_date":             data.PaidDate,
	}

	return r.GetDB().
		Model(&entity.WarehouseSale{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) UpdateWarehouseSalePayment(data *entity.WarehouseSalePayment) error {
	updates := map[string]interface{}{
		"warehouse_sale_id": data.WarehouseSaleId,
		"payment_date":      data.PaymentDate,
		"nominal":           data.Nominal,
		"payment_proof":     data.PaymentProof,
		"payment_method":    data.PaymentMethod,
		"updated_by":        data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseSalePayment{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) DeleteWarehouseSale(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseSale{}).Error
}

func (r *WarehouseRepository) DeleteWarehouseSalePayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseSalePayment{}).Error
}

func (r *WarehouseRepository) CreateWarehouseSaleQueue(data *entity.WarehouseSaleQueue) error {
	return r.GetDB().Model(&entity.WarehouseSaleQueue{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseSaleQueueById(id uint64) (entity.WarehouseSaleQueue, error) {
	var warehouseSaleQueue entity.WarehouseSaleQueue
	err := r.GetDB().Model(&entity.WarehouseSaleQueue{}).Preload("Warehouse.Location").Preload("Item").Preload("Customer").Where("id = ?", id).First(&warehouseSaleQueue).Error
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

	err := query.Order("created_at DESC").Preload("Warehouse.Location").Preload("Customer").Preload("Item").Find(&warehouseSaleQueues).Error
	if err != nil {
		return nil, err
	}

	return warehouseSaleQueues, nil
}

func (r *WarehouseRepository) CreateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseItemProcurementDrafts(filter dto.GetWarehouseItemProcurementDraftFilter) ([]entity.WarehouseItemProcurementDraft, error) {
	var data []entity.WarehouseItemProcurementDraft
	query := r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Joins("LEFT JOIN items ON items.id = warehouse_item_procurement_drafts.id")

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.ItemCategory.Value().IsValid() {
		query = query.Where("items.category = ?", filter.ItemCategory.Value())
	}

	err := query.Order("item_id ASC").Order("price ASC").Order("created_at DESC").Preload("Warehouse.Location").Preload("Item").Preload("Supplier").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *WarehouseRepository) CountWarehouseItemProcurement(filter dto.GetWarehouseItemProcurementFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.WarehouseItemProcurement{})

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.Status.Value().IsValid() {
		query = query.Where("status = ?", filter.Status.Value())
	}

	err := query.Count(&count).Error
	if err != nil {
		return -1, err
	}

	return count, nil
}

func (r *WarehouseRepository) GetWarehouseItemProcurementDraft(id uint64) (entity.WarehouseItemProcurementDraft, error) {
	var data entity.WarehouseItemProcurementDraft
	err := r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).Where("id = ?", id).Preload("Warehouse.Location").Preload("Item").Preload("Supplier").First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItemProcurementDraft{}, errx.NotFound("warehouse item procurement draft not found")
		}
		return entity.WarehouseItemProcurementDraft{}, err
	}

	return data, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemProcurementDraft(data *entity.WarehouseItemProcurementDraft) error {
	updates := map[string]interface{}{
		"warehouse_id":   data.WarehouseId,
		"item_id":        data.ItemId,
		"supplier_id":    data.SupplierId,
		"daily_spending": data.DailySpending,
		"days_need":      data.DaysNeed,
		"price":          data.Price,
		"updated_by":     data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseItemProcurementDraft{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemProcurementDraft(id uint64) error {
	res := r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemProcurementDraft{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errx.NotFound("warehouse item procurement draft not found")
	}
	return nil
}

func (r *WarehouseRepository) CreateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurement{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseItemProcurements(filter dto.GetWarehouseItemProcurementFilter) ([]entity.WarehouseItemProcurement, error) {
	var data []entity.WarehouseItemProcurement
	query := r.GetDB().Model(&entity.WarehouseItemProcurement{}).Preload("Warehouse.Location").Preload("Item").Preload("Supplier")

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus)
	}

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.ProcurementStatus.Value().IsValid() {
		query = query.Where("status = ?", filter.ProcurementStatus.Value())
	}

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Order("status ASC").Order("payment_status DESC").Order("deadline_payment_date ASC").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *WarehouseRepository) GetWarehouseItemProcurement(id uint64) (entity.WarehouseItemProcurement, error) {
	var data entity.WarehouseItemProcurement
	err := r.GetDB().Model(&entity.WarehouseItemProcurement{}).Preload("Warehouse.Location").Preload("Item").Preload("Supplier").Preload("Payments", func(db *gorm.DB) *gorm.DB {
		return db.Order("created_at ASC")
	}).Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItemProcurement{}, errx.NotFound("warehouse item procurement not found")
		}
		return entity.WarehouseItemProcurement{}, err
	}

	return data, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemProcurement(data *entity.WarehouseItemProcurement) error {
	updates := map[string]interface{}{
		"warehouse_id":            data.WarehouseId,
		"item_id":                 data.ItemId,
		"supplier_id":             data.SupplierId,
		"daily_spending":          data.DailySpending,
		"days_need":               data.DaysNeed,
		"quantity":                data.Quantity,
		"receive_quantity":        data.ReceiveQuantity,
		"note":                    data.Note,
		"price":                   data.Price,
		"total_price":             data.TotalPrice,
		"estimation_arrival_date": data.EstimationArrivalDate,
		"is_arrived":              data.IsArrived,
		"taken_at":                data.TakenAt,
		"taken_by":                data.TakenBy,
		"status":                  data.Status,
		"payment_status":          data.PaymentStatus,
		"expired_at":              data.ExpiredAt,
		"deadline_payment_date":   data.DeadlinePaymentDate,
		"payment_type":            data.PaymentType,
		"paid_date":               data.PaidDate,
		"updated_by":              data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseItemProcurement{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemProcurement(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemProcurement{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemProcurementPaymentInBatch(data *[]entity.WarehouseItemProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurementPayment{}).CreateInBatches(data, len(*data)).Error
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
	updates := map[string]interface{}{
		"warehouse_item_procurement_id": data.WarehouseItemProcurementId,
		"payment_date":                  data.PaymentDate,
		"nominal":                       data.Nominal,
		"payment_proof":                 data.PaymentProof,
		"payment_method":                data.PaymentMethod,
		"updated_by":                    data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseItemProcurementPayment{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemProcurementPayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemProcurementPayment{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementDraft{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurementDrafts(filter dto.GetWarehouseItemCornProcurementDraftFilter) ([]entity.WarehouseItemCornProcurementDraft, error) {
	var warehouseItemCornProcurementDrafts []entity.WarehouseItemCornProcurementDraft
	query := r.GetDB().Model(&entity.WarehouseItemCornProcurementDraft{})

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	err := query.Order("price ASC").Order("created_at DESC").Preload("Supplier").Preload("Warehouse.Location").Find(&warehouseItemCornProcurementDrafts).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemCornProcurementDrafts, nil
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurementDraft(id uint64) (entity.WarehouseItemCornProcurementDraft, error) {
	var warehouseItemCornProcurementDraft entity.WarehouseItemCornProcurementDraft
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurementDraft{}).Where("id = ?", id).Preload("Supplier").Preload("Warehouse.Location").First(&warehouseItemCornProcurementDraft).Error
	if err != nil {
		return entity.WarehouseItemCornProcurementDraft{}, err
	}

	return warehouseItemCornProcurementDraft, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemCornProcurementDraft(data *entity.WarehouseItemCornProcurementDraft) error {
	updates := map[string]interface{}{
		"warehouse_id":                    data.WarehouseId,
		"supplier_id":                     data.SupplierId,
		"oven_condition":                  data.OvenCondition,
		"corn_water_level":                data.CornWaterLevel,
		"is_oven_can_operate_in_near_day": data.IsOvenCanOperateInNearDay,
		"quantity":                        data.Quantity,
		"price":                           data.Price,
		"discount":                        data.Discount,
		"updated_by":                      data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseItemCornProcurementDraft{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) DeleteWarehouseItemCornProcurementDraft(id uint64) error {
	res := r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemCornProcurementDraft{})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errx.NotFound("warehouse item corn procurement draft not found")
	}
	return nil
}

func (r *WarehouseRepository) CreateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurement{}).Create(&data).Error
}

func (r *WarehouseRepository) UpdateWarehouseItemCornProcurement(data *entity.WarehouseItemCornProcurement) error {
	updates := map[string]interface{}{
		"warehouse_id":                    data.WarehouseId,
		"supplier_id":                     data.SupplierId,
		"quantity":                        data.Quantity,
		"receive_quantity":                data.ReceiveQuantity,
		"note":                            data.Note,
		"price":                           data.Price,
		"total_price":                     data.TotalPrice,
		"is_arrived":                      data.IsArrived,
		"taken_at":                        data.TakenAt,
		"taken_by":                        data.TakenBy,
		"status":                          data.Status,
		"payment_status":                  data.PaymentStatus,
		"oven_condition":                  data.OvenCondition,
		"corn_water_level":                data.CornWaterLevel,
		"is_oven_can_operate_in_near_day": data.IsOvenCanOperateInNearDay,
		"expired_at":                      data.ExpiredAt,
		"deadline_payment_date":           data.DeadlinePaymentDate,
		"payment_type":                    data.PaymentType,
		"discount":                        data.Discount,
		"updated_by":                      data.UpdatedBy,
		"paid_date":                       data.PaidDate,
	}

	return r.GetDB().
		Model(&entity.WarehouseItemCornProcurement{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurement(id uint64) (entity.WarehouseItemCornProcurement, error) {
	var warehouseItemCornProcurement entity.WarehouseItemCornProcurement
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurement{}).Preload("Warehouse.Location").Preload("Payments").Preload("Supplier").Where("id = ?", id).First(&warehouseItemCornProcurement).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseItemCornProcurement{}, errx.NotFound("failed warehouse item corn procurement")
		}
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

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.ProcurementStatus.Value().IsValid() {
		query = query.Where("status = ?", filter.ProcurementStatus.Value())
	}

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("Warehouse.Location").Order("status ASC").Order("payment_status DESC").Order("deadline_payment_date DESC").Preload("Supplier").Find(&warehouseItemCornProcurements).Error
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

func (r *WarehouseRepository) CreateWarehouseItemCornProcurementPaymentInBatch(data *[]entity.WarehouseItemCornProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{}).CreateInBatches(data, len(*data)).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error {
	return r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{}).Create(data).Error
}

func (r *WarehouseRepository) UpdateWarehouseItemCornProcurementPayment(data *entity.WarehouseItemCornProcurementPayment) error {
	updates := map[string]interface{}{
		"warehouse_item_corn_procurement_id": data.WarehouseItemCornProcurementId,
		"payment_date":                       data.PaymentDate,
		"nominal":                            data.Nominal,
		"payment_proof":                      data.PaymentProof,
		"payment_method":                     data.PaymentMethod,
		"updated_by":                         data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.WarehouseItemCornProcurementPayment{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *WarehouseRepository) GetWarehouseItemCornProcurementPayment(id uint64) (entity.WarehouseItemCornProcurementPayment, error) {
	var warehouseItemCornProcurementPayment entity.WarehouseItemCornProcurementPayment
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{}).Where("id = ?", id).First(&warehouseItemCornProcurementPayment).Error
	if err != nil {
		return entity.WarehouseItemCornProcurementPayment{}, err
	}

	return warehouseItemCornProcurementPayment, nil
}

func (r *WarehouseRepository) DeleteWarehouseItemCornProcurementPayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.WarehouseItemCornProcurementPayment{}).Error
}

func (r *WarehouseRepository) CreateWarehouseItemCorn(data *entity.WarehouseItemCorn) error {
	return r.GetDB().Model(&entity.WarehouseItemCorn{}).Create(data).Error
}

func (r *WarehouseRepository) GetWarehouseItemCorns(filter dto.GetWarehouseItemCornFilter) ([]entity.WarehouseItemCorn, error) {
	var warehouseItemCorns []entity.WarehouseItemCorn
	query := r.GetDB().Model(&entity.WarehouseItemCorn{})

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.FromNewest {
		query = query.Order("created_at DESC")
	} else {
		query = query.Order("created_at ASC")
	}

	if filter.WithZeroQuantity != nil && *filter.WithZeroQuantity {
		query = query.Or("quantity = 0")
	} else {
		query = query.Where("quantity <> 0")
	}

	err := query.Preload("Warehouse.Location").Preload("Supplier").Find(&warehouseItemCorns).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemCorns, nil
}

func (r *WarehouseRepository) UpdateWarehouseItemCorn(data *entity.WarehouseItemCorn) error {
	return r.GetDB().Model(&entity.WarehouseItemCorn{}).Where("id = ?", data.Id).Updates(map[string]any{
		"warehouse_id": data.WarehouseId,
		"supplier_id":  data.SupplierId,
		"quantity":     data.Quantity,
		"expired_at":   data.ExpiredAt,
		"order_date":   data.OrderDate,
		"updated_by":   data.UpdatedBy,
	}).Error
}

func (r *WarehouseRepository) GetWarehouseItemCorn(id uint64) (entity.WarehouseItemCorn, error) {
	var warehouseItemCorn entity.WarehouseItemCorn
	err := r.GetDB().Model(&entity.WarehouseItemCorn{}).Preload("Warehouse.Location").Preload("Supplier").Where("id = ?", id).First(&warehouseItemCorn).Error
	if err != nil {
		return entity.WarehouseItemCorn{}, err
	}

	return warehouseItemCorn, nil
}

func (r *WarehouseRepository) GetWarehouseItemCornPrices() ([]entity.WarehouseItemCornPrice, error) {
	var warehouseItemCornPrices []entity.WarehouseItemCornPrice
	err := r.GetDB().Model(&entity.WarehouseItemCornPrice{}).Find(&warehouseItemCornPrices).Error
	if err != nil {
		return nil, err
	}

	return warehouseItemCornPrices, nil
}

func (r *WarehouseRepository) CreateWarehouseItemProcurementDraftInBatch(data *[]entity.WarehouseItemProcurementDraft) error {
	return r.GetDB().Model(&entity.WarehouseItemProcurementDraft{}).CreateInBatches(data, len(*data)).Error
}

func (r *WarehouseRepository) CountQuantityWarehouseItemCornByWarehouseId(warehouseId uint64) (float64, error) {
	var total float64
	err := r.db.Model(&entity.WarehouseItemCorn{}).Select("SUM(quantity)").Where("warehouse_id = ?", warehouseId).Row().Scan(&total)
	if err != nil {
		return 0, err
	}
	return total, nil
}
