package repository

import (
	"errors"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"gorm.io/gorm"
)

type CashflowRepository struct {
	db *gorm.DB
	tx *gorm.DB
}

type ICashflowRepository interface {
	UseTx(tx bool)
	Commit() error
	Rollback() error

	GetWarehouseSalePayments(filter dto.GetWarehouseSalePaymentFilter) ([]entity.WarehouseSalePayment, error)
	GetStoreSalePayments(filter dto.GetStoreSalePaymentFilter) ([]entity.StoreSalePayment, error)
	GetAfkirChickenSalePayments(filter dto.GetAfkirChickenSalePaymentFilter) ([]entity.AfkirChickenSalePayment, error)

	GetWarehouseSalePaymentById(id uint64) (entity.WarehouseSalePayment, error)
	GetStoreSalePaymentById(id uint64) (entity.StoreSalePayment, error)
	GetAfkirChickenSalePaymentById(id uint64) (entity.AfkirChickenSalePayment, error)

	CreateExpense(data *entity.Expense) error
	GetExpense(id uint64) (entity.Expense, error)
	GetExpenses(filter dto.GetExpenseFilter) ([]entity.Expense, error)

	GetChickenProcurementPayments(filter dto.GetChickenProcurementPaymentFilter) ([]entity.ChickenProcurementPayment, error)
	GetWarehouseItemCornProcurementPayments(filter dto.GetWarehouseItemCornProcurementPaymentFilter) ([]entity.WarehouseItemCornProcurementPayment, error)
	GetWarehouseItemProcurementPayments(filter dto.GetWarehouseItemProcurementPaymentFilter) ([]entity.WarehouseItemProcurementPayment, error)
	GetUserSalaryPayments(filter dto.GetUserSalaryPaymentFilter) ([]entity.UserSalaryPayment, error)

	GetChickenProcurementPaymentById(id uint64) (entity.ChickenProcurementPayment, error)
	GetWarehouseItemProcurementPaymentById(id uint64) (entity.WarehouseItemProcurementPayment, error)
	GetWarehouseItemCornProcurementPaymentById(id uint64) (entity.WarehouseItemCornProcurementPayment, error)
	GetUserSalaryPaymentById(id uint64) (entity.UserSalaryPayment, error)
}

func NewCashflowRepository(db *gorm.DB) ICashflowRepository {
	return &CashflowRepository{
		db: db,
	}
}

func (r *CashflowRepository) UseTx(tx bool) {
	if tx {
		r.tx = r.db.Begin()
	}
}

func (r *CashflowRepository) Commit() error {
	err := r.GetDB().Commit().Error
	r.tx = nil
	return err
}

func (r *CashflowRepository) Rollback() error {
	if r.tx == nil {
		return nil
	}
	err := r.GetDB().Rollback().Error
	r.tx = nil
	return err
}

func (r *CashflowRepository) GetDB() *gorm.DB {
	if r.tx != nil {
		return r.tx
	}
	return r.db
}

func (r *CashflowRepository) GetWarehouseSalePayments(filter dto.GetWarehouseSalePaymentFilter) ([]entity.WarehouseSalePayment, error) {
	var warehouseSalePayments []entity.WarehouseSalePayment
	query := r.GetDB().Model(&entity.WarehouseSalePayment{})

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= > ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.Preload("WarehouseSale.Customer").
		Preload("WarehouseSale.Item").
		Preload("WarehouseSale.Warehouse.Location").Find(&warehouseSalePayments).Error
	if err != nil {
		return nil, err
	}

	return warehouseSalePayments, nil
}

func (r *CashflowRepository) GetStoreSalePayments(filter dto.GetStoreSalePaymentFilter) ([]entity.StoreSalePayment, error) {
	var storeSalePayments []entity.StoreSalePayment
	query := r.GetDB().Model(&entity.StoreSalePayment{})

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= > ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.Preload("StoreSale.Customer").
		Preload("StoreSale.Item").
		Preload("StoreSale.Store.Location").Find(&storeSalePayments).Error
	if err != nil {
		return nil, err
	}

	return storeSalePayments, nil
}

func (r *CashflowRepository) GetAfkirChickenSalePayments(filter dto.GetAfkirChickenSalePaymentFilter) ([]entity.AfkirChickenSalePayment, error) {
	var AfkirChickenSalePayments []entity.AfkirChickenSalePayment
	query := r.GetDB().Model(&entity.AfkirChickenSalePayment{})

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= > ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.Preload("AfkirChickenSale.AfkirChickenCustomer").
		Preload("AfkirChickenSale.ChickenCage.Cage.Location").Find(&AfkirChickenSalePayments).Error
	if err != nil {
		return nil, err
	}

	return AfkirChickenSalePayments, nil
}

func (r *CashflowRepository) GetWarehouseSalePaymentById(id uint64) (entity.WarehouseSalePayment, error) {
	var payment entity.WarehouseSalePayment
	err := r.GetDB().
		Preload("WarehouseSale.Customer").
		Preload("WarehouseSale.Item").
		Preload("WarehouseSale.Warehouse.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}

func (r *CashflowRepository) GetStoreSalePaymentById(id uint64) (entity.StoreSalePayment, error) {
	var payment entity.StoreSalePayment
	err := r.GetDB().
		Preload("StoreSale.Customer").
		Preload("StoreSale.Item").
		Preload("StoreSale.Store.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}

func (r *CashflowRepository) GetAfkirChickenSalePaymentById(id uint64) (entity.AfkirChickenSalePayment, error) {
	var payment entity.AfkirChickenSalePayment
	err := r.GetDB().
		Preload("AfkirChickenSale.AfkirChickenCustomer").
		Preload("AfkirChickenSale.ChickenCage.Cage.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}

func (r *CashflowRepository) CreateExpense(data *entity.Expense) error {
	return r.GetDB().Model(&entity.Expense{}).Create(data).Error
}

func (r *CashflowRepository) GetExpense(id uint64) (entity.Expense, error) {
	var data entity.Expense
	err := r.GetDB().Model(&entity.Expense{}).Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Expense{}, errx.NotFound("expense not found")
		}
		return entity.Expense{}, err
	}

	return data, nil
}

func (r *CashflowRepository) GetExpenses(filter dto.GetExpenseFilter) ([]entity.Expense, error) {
	var data []entity.Expense
	query := r.GetDB().Model(&entity.Expense{})

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.Preload("Cage.Location").Preload("Warehouse.Location").Preload("Store.Location").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CashflowRepository) GetChickenProcurementPayments(filter dto.GetChickenProcurementPaymentFilter) ([]entity.ChickenProcurementPayment, error) {
	var data []entity.ChickenProcurementPayment
	query := r.GetDB().Model(&entity.ChickenProcurementPayment{})

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(payment_date) >= ? AND DATE(payment_date) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.
		Preload("ChickenProcurement.Cage.Location").
		Preload("ChickenProcurement.Supplier").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *CashflowRepository) GetUserSalaryPayments(filter dto.GetUserSalaryPaymentFilter) ([]entity.UserSalaryPayment, error) {
	var data []entity.UserSalaryPayment
	query := r.GetDB().Model(&entity.UserSalaryPayment{})

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.IsPaid != nil {
		query = query.Where("is_paid = ?", filter.IsPaid)
	}

	err := query.
		Preload("User").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *CashflowRepository) GetWarehouseItemProcurementPayments(filter dto.GetWarehouseItemProcurementPaymentFilter) ([]entity.WarehouseItemProcurementPayment, error) {
	var data []entity.WarehouseItemProcurementPayment
	query := r.GetDB().Model(&entity.WarehouseItemProcurementPayment{})

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(payment_date) >= ? AND DATE(payment_date) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.
		Preload("WarehouseItemProcurement.Warehouse.Location").
		Preload("WarehouseItemProcurement.Supplier").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *CashflowRepository) GetWarehouseItemCornProcurementPayments(filter dto.GetWarehouseItemCornProcurementPaymentFilter) ([]entity.WarehouseItemCornProcurementPayment, error) {
	var data []entity.WarehouseItemCornProcurementPayment
	query := r.GetDB().Model(&entity.WarehouseItemCornProcurementPayment{})

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(payment_date) >= ? AND DATE(payment_date) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	err := query.
		Preload("WarehouseItemCornProcurement.Warehouse.Location").
		Preload("WarehouseItemCornProcurement.Supplier").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *CashflowRepository) GetChickenProcurementPaymentById(id uint64) (entity.ChickenProcurementPayment, error) {
	var payment entity.ChickenProcurementPayment
	err := r.db.
		Preload("ChickenProcurement.Cage.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}

func (r *CashflowRepository) GetWarehouseItemProcurementPaymentById(id uint64) (entity.WarehouseItemProcurementPayment, error) {
	var payment entity.WarehouseItemProcurementPayment
	err := r.db.
		Preload("WarehouseItemProcurement.Warehouse.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}

func (r *CashflowRepository) GetWarehouseItemCornProcurementPaymentById(id uint64) (entity.WarehouseItemCornProcurementPayment, error) {
	var payment entity.WarehouseItemCornProcurementPayment
	err := r.db.
		Preload("WarehouseItemCornProcurement.Warehouse.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}

func (r *CashflowRepository) GetUserSalaryPaymentById(id uint64) (entity.UserSalaryPayment, error) {
	var payment entity.UserSalaryPayment
	err := r.db.
		Preload("User").
		Preload("CreatedByUser").
		First(&payment, id).Error
	return payment, err
}
