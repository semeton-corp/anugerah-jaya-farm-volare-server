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
	CreateStoreItemsInBatch(storeItems *[]entity.StoreItem) error
	UpdateStoreItem(storeItem *entity.StoreItem) error
	GetStoreItems(filter dto.GetStoreItemFilter) ([]entity.StoreItem, error)
	GetStoreItemByStoreIdAndItemId(storeId uint64, itemId uint64) (entity.StoreItem, error)

	GetStoreItemHistories(filter dto.GetStoreItemHistoryFilter) ([]entity.StoreItemHistory, error)
	GetStoreItemHistoryById(id uint64) (entity.StoreItemHistory, error)
	CountTotalStoreItemHistory(filter dto.GetStoreItemHistoryFilter) (int64, error)

	CreateStoreSale(storeSale *entity.StoreSale) error
	GetStoreSaleById(id uint64) (entity.StoreSale, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) ([]entity.StoreSale, error)
	UpdateStoreSale(storeSale *entity.StoreSale) error
	CountTotalStoreSale(filter dto.GetStoreSaleFilter) (uint64, error)
	DeleteStoreSale(id uint64) error

	CreateStoreSalePaymentInBatch(data *[]entity.StoreSalePayment) error
	CreateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error
	GetStoreSalePaymentById(id uint64) (entity.StoreSalePayment, error)
	UpdateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error
	DeleteStoreSalePayment(id uint64) error

	CreateStoreSaleQueue(data *entity.StoreSaleQueue) error
	GetStoreSaleQueueById(id uint64) (entity.StoreSaleQueue, error)
	GetStoreSaleQueues(filter dto.GetStoreSaleQueueFilter) ([]entity.StoreSaleQueue, error)
	DeleteStoreSaleQueue(id uint64) error

	GetStoreSalePayments(filter dto.GetStoreSalePaymentFilter) ([]entity.StoreSalePayment, error)
	CountTotalStoreSalePayment(filter dto.GetStoreSalePaymentFilter) (uint64, error)
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

func (r *StoreRepository) UpdateStore(data *entity.Store) error {
	updates := map[string]interface{}{
		"name":        data.Name,
		"location_id": data.LocationId,
		"updated_by":  data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.Store{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
}

func (r *StoreRepository) DeleteStore(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.Store{}).Error
}

func (r *StoreRepository) GetStoreById(id uint64) (entity.Store, error) {
	var data entity.Store
	err := r.GetDB().Model(&entity.Store{}).Preload("StoreItems").Preload("StorePlacement").Preload("Location").Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Store{}, errx.BadRequest("store not found")
		}
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

	err := query.Order("created_at DESC").Preload("StorePlacement").Preload("StoreItems").Preload("Location").Find(&stores).Error
	if err != nil {
		return nil, err
	}

	return stores, nil
}

func (r *StoreRepository) CreateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error {
	err := r.GetDB().Create(storeRequestItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrForeignKeyViolated) {
			return errx.BadRequest("invalid warehouse, item, or store")
		}
		return err
	}
	return nil
}

func (r *StoreRepository) GetStoreRequestItemById(id uint64) (entity.StoreRequestItem, error) {
	var storeRequestItem entity.StoreRequestItem
	err := r.GetDB().Preload("Warehouse.Location").Preload("Store.Location").Preload("Item").Preload("CreatedByUser").First(&storeRequestItem, id).Error
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
	query := r.GetDB().Model(&entity.StoreRequestItem{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.StoreId > 0 {
		query = query.Where("store_id = ?", filter.StoreId)
	}

	if filter.WarehouseId > 0 {
		query = query.Where("warehouse_id = ?", filter.WarehouseId)
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Preload("Warehouse.Location").Preload("Store.Location").Preload("Item").Find(&storeRequestItems).Order("created_at ASC").Error
	if err != nil {
		return nil, err
	}
	return storeRequestItems, nil
}

func (r *StoreRepository) UpdateStoreRequestItem(storeRequestItem *entity.StoreRequestItem) error {
	return r.GetDB().Model(entity.StoreRequestItem{}).Where("id = ?", storeRequestItem.Id).Updates(map[string]any{
		"quantity":       storeRequestItem.Quantity,
		"updated_by":     storeRequestItem.UpdatedBy,
		"status":         storeRequestItem.Status,
		"warehouse_note": storeRequestItem.WarehouseNote,
		"store_note":     storeRequestItem.StoreNote,
	}).Error
}

func (r *StoreRepository) FirstOrCreateStoreItem(storeItem *entity.StoreItem) error {
	return r.GetDB().FirstOrCreate(storeItem, entity.StoreItem{
		ItemId:  storeItem.ItemId,
		StoreId: storeItem.StoreId,
	}).Error
}

func (r *StoreRepository) CreateStoreItemsInBatch(storeItems *[]entity.StoreItem) error {
	return r.GetDB().Model(&entity.StoreItem{}).CreateInBatches(storeItems, len(*storeItems)).Error
}

func (r *StoreRepository) UpdateStoreItem(storeItem *entity.StoreItem) error {
	return r.GetDB().Model(entity.StoreItem{}).
		Where("store_id = ? AND item_id = ?", storeItem.StoreId, storeItem.ItemId).
		Updates(map[string]interface{}{
			"quantity": storeItem.Quantity,
		}).Error
}

func (r *StoreRepository) GetStoreItemByStoreIdAndItemId(storeId uint64, itemId uint64) (entity.StoreItem, error) {
	var storeItem entity.StoreItem
	err := r.GetDB().Model(&entity.StoreItem{}).Preload("Store.Location").Preload("Item").Where("store_id = ? AND item_id = ?", storeId, itemId).First(&storeItem).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreItem{}, errx.BadRequest("store item not found")
		}
		return entity.StoreItem{}, err
	}

	return storeItem, nil
}

func (r *StoreRepository) GetStoreItems(filter dto.GetStoreItemFilter) ([]entity.StoreItem, error) {
	var storeItems []entity.StoreItem
	query := r.GetDB().Model(&entity.StoreItem{}).Joins("JOIN items ON items.id = store_items.item_id")

	if filter.StoreId != 0 {
		query = query.Where("store_items.store_id = ?", filter.StoreId)
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

	if filter.StoreIds != nil {
		query = query.Where("store_items.store_id IN ?", filter.StoreIds)
	}

	err := query.Order("created_at DESC").Preload("Item").Preload("Store.Location").Find(&storeItems).Error
	if err != nil {
		return nil, err
	}

	return storeItems, nil
}

func (r *StoreRepository) GetStoreItemHistories(filter dto.GetStoreItemHistoryFilter) ([]entity.StoreItemHistory, error) {
	storeItemHistory := make([]entity.StoreItemHistory, 0)
	query := r.GetDB().Model(&entity.StoreItemHistory{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	err := query.Order("created_at DESC").Preload("User").Find(&storeItemHistory).Error
	if err != nil {
		return nil, err
	}

	return storeItemHistory, nil
}

func (r *StoreRepository) GetStoreItemHistoryById(id uint64) (entity.StoreItemHistory, error) {
	var storeItemHistory entity.StoreItemHistory
	err := r.GetDB().Model(&entity.StoreItemHistory{}).Where("id = ?", id).Preload("User").First(&storeItemHistory).Error
	if err != nil {
		return entity.StoreItemHistory{}, err
	}

	return storeItemHistory, nil
}

func (r *StoreRepository) CountTotalStoreItemHistory(filter dto.GetStoreItemHistoryFilter) (int64, error) {
	var total int64
	query := r.GetDB().Model(&entity.StoreItemHistory{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	err := query.Count(&total).Error
	if err != nil {
		return -1, err
	}

	return total, nil
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
	err := r.GetDB().Preload("Payments").Preload("Customer").Preload("Store.Location").Preload("Item").First(&storeSale, id).Error
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
	query := r.GetDB().Model(&entity.StoreSale{})

	if !filter.Date.Value().IsZero() {
		query = query.Where("DATE(created_at) = ?", filter.Date.Value())
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.ItemId > 0 {
		query = query.Where("item_id = ?", filter.ItemId)
	}

	if filter.StoreId > 0 {
		query = query.Where("store_id = ?", filter.StoreId)
	}

	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.PaymentStatuses != nil {
		paymentStatus := make([]enum.PaymentStatus, 0)
		for _, e := range filter.PaymentStatuses {
			paymentStatus = append(paymentStatus, e.Value())
		}

		query = query.Where("payment_status IN ?", paymentStatus)
	}

	err := query.Preload("Store.Location").Preload("Customer").Preload("Item").Preload("Payments").Find(&storeSales).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return storeSales, nil
}

func (r *StoreRepository) CreateStoreSalePaymentInBatch(data *[]entity.StoreSalePayment) error {
	return r.GetDB().Model(&entity.StoreSalePayment{}).CreateInBatches(data, len(*data)).Error
}

func (r *StoreRepository) CreateStoreSalePayment(storeSalePayment *entity.StoreSalePayment) error {
	return r.GetDB().Model(&entity.StoreSalePayment{}).Create(&storeSalePayment).Error
}

func (r *StoreRepository) UpdateStoreSale(data *entity.StoreSale) error {
	updates := map[string]interface{}{
		"customer_id":           data.CustomerId,
		"item_id":               data.ItemId,
		"store_id":              data.StoreId,
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
		Model(&entity.StoreSale{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
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

func (r *StoreRepository) UpdateStoreSalePayment(data *entity.StoreSalePayment) error {
	updates := map[string]interface{}{
		"store_sale_id":  data.StoreSaleId,
		"payment_date":   data.PaymentDate,
		"nominal":        data.Nominal,
		"payment_proof":  data.PaymentProof,
		"payment_method": data.PaymentMethod,
		"updated_by":     data.UpdatedBy,
	}

	return r.GetDB().
		Model(&entity.StoreSalePayment{}).
		Where("id = ?", data.Id).
		Updates(updates).Error
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

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.StoreId > 0 {
		query = query.Where("store_id = ?", filter.StoreId)
	}

	if filter.ItemId > 0 {
		query = query.Where("item_id = ?", filter.ItemId)
	}

	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.PaymentStatuses != nil {
		paymentStatus := make([]enum.PaymentStatus, 0)
		for _, e := range filter.PaymentStatuses {
			paymentStatus = append(paymentStatus, e.Value())
		}

		query = query.Where("payment_status IN ?", paymentStatus)
	}

	err := query.Model(&entity.StoreSale{}).Count(&totalData).Error
	if err != nil {
		return 0, err
	}

	return uint64(totalData), nil
}

func (r *StoreRepository) DeleteStoreSale(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.StoreSale{}).Error
}

func (r *StoreRepository) DeleteStoreSalePayment(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.StoreSalePayment{}).Error
}

func (r *StoreRepository) CreateStoreSaleQueue(data *entity.StoreSaleQueue) error {
	return r.GetDB().Model(&entity.StoreSaleQueue{}).Create(data).Error
}

func (r *StoreRepository) GetStoreSaleQueueById(id uint64) (entity.StoreSaleQueue, error) {
	var storeSaleQueue entity.StoreSaleQueue
	err := r.GetDB().Model(&entity.StoreSaleQueue{}).Preload("Store.Location").Preload("Item").Preload("Customer").Where("id = ?", id).First(&storeSaleQueue).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreSaleQueue{}, errx.NotFound("store sale queue not found")
		}
		return entity.StoreSaleQueue{}, err
	}

	return storeSaleQueue, nil
}

func (r *StoreRepository) DeleteStoreSaleQueue(id uint64) error {
	return r.GetDB().Where("id = ?", id).Delete(&entity.StoreSaleQueue{}).Error
}

func (r *StoreRepository) GetStoreSaleQueues(filter dto.GetStoreSaleQueueFilter) ([]entity.StoreSaleQueue, error) {
	var storeSaleQueues []entity.StoreSaleQueue
	query := r.GetDB().Model(&entity.StoreSaleQueue{})

	if filter.StoreId > 0 {
		query = query.Where("store_id = ?", filter.StoreId)
	}

	err := query.Order("created_at DESC").Preload("Store.Location").Preload("Item").Preload("Customer").Find(&storeSaleQueues).Error
	if err != nil {
		return nil, err
	}

	return storeSaleQueues, nil
}

func (r *StoreRepository) GetStoreSalePayments(filter dto.GetStoreSalePaymentFilter) ([]entity.StoreSalePayment, error) {
	var storeSalePayments []entity.StoreSalePayment
	query := r.GetDB().Model(&entity.StoreSalePayment{}).Joins("LEFT JOIN store_sales ON store_sales.id = store_sale_payments.store_sale_id")

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(store_sale_payments.created_at) >= ? AND DATE(store_sale_payments.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.StoreId > 0 {
		query = query.Where("store_sales.store_id = ?", filter.StoreId)
	}

	if filter.Page > 0 {
		query = query.Limit(int(constant.PaginationDefaultLimit)).Offset((int(filter.Page) - 1) * int(constant.PaginationDefaultLimit))
	}

	err := query.Order("store_sale_payments.created_at DESC").Preload("StoreSale.Customer").
		Preload("StoreSale.Item").
		Preload("StoreSale.Store.Location").Find(&storeSalePayments).Error
	if err != nil {
		return nil, err
	}

	return storeSalePayments, nil
}

func (r *StoreRepository) CountTotalStoreSalePayment(filter dto.GetStoreSalePaymentFilter) (uint64, error) {
	var count int64
	query := r.GetDB().Model(&entity.StoreSalePayment{}).Joins("LEFT JOIN store_sales ON store_sales.id = store_sale_payments.store_sale_id")

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(store_sale_payments.created_at) >= ? AND DATE(store_sale_payments.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.StoreId > 0 {
		query = query.Where("store_sales.store_id = ?", filter.StoreId)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return uint64(count), nil
}
