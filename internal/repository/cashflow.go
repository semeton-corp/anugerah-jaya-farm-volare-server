package repository

import (
	"errors"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
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

	GetChickenProcurementPaymentById(id uint64) (entity.ChickenProcurementPayment, error)
	GetWarehouseItemProcurementPaymentById(id uint64) (entity.WarehouseItemProcurementPayment, error)
	GetWarehouseItemCornProcurementPaymentById(id uint64) (entity.WarehouseItemCornProcurementPayment, error)
	GetUserSalaryPaymentById(id uint64) (entity.UserSalaryPayment, error)

	CreateUserCashAdvance(data *entity.UserCashAdvance) error
	GetUserCashAdvances(filter dto.GetUserCashAdvanceFilter) ([]entity.UserCashAdvance, error)
	GetUserCashAdvance(id uint64) (entity.UserCashAdvance, error)
	UpdateUserCashAdvance(data *entity.UserCashAdvance) error

	CreateUserCashAdvancePayment(data *entity.UserCashAdvancePayment) error
	CreateUserCashAdvancePaymentBatch(payments *[]entity.UserCashAdvancePayment) error
	GetUserCashAdvancePayments(filter dto.GetUserCashAdvancePaymentFilter) ([]entity.UserCashAdvancePayment, error)
	GetUserCashAdvancePayment(id uint64) (entity.UserCashAdvancePayment, error)

	GetStoreSaleCashflows(filter dto.GetStoreSaleFilter) ([]entity.StoreSale, error)
	GetWarehouseSaleCashflows(filter dto.GetWarehouseSaleFilter) ([]entity.WarehouseSale, error)
	GetAfkirChickenSaleCashflows(filter dto.GetAfkirChickenSaleFilter) ([]entity.AfkirChickenSale, error)

	GetStoreSaleCashflow(id uint64) (entity.StoreSale, error)
	GetWarehouseSaleCashflow(id uint64) (entity.WarehouseSale, error)
	GetAfkirChickenSaleCashflow(id uint64) (entity.AfkirChickenSale, error)

	GetUserSalaryPayment(id uint64) (entity.UserSalaryPayment, error)
	CountUserSalaryPayments(filter dto.GetUserSalaryPaymentFilter) (int64, error)
	UpdateUserSalaryPayment(data *entity.UserSalaryPayment) error
	GetUserSalaryPayments(filter dto.GetUserSalaryPaymentFilter) ([]entity.UserSalaryPayment, error)

	GetWarehouseItemProcurementCashflows(filter dto.GetWarehouseItemProcurementFilter) ([]entity.WarehouseItemProcurement, error)
	GetWarehouseItemCornProcurementCashflows(filter dto.GetWarehouseItemCornProcurementFilter) ([]entity.WarehouseItemCornProcurement, error)
	GetChickenProcurementCashflows(filter dto.GetChickenProcurementFilter) ([]entity.ChickenProcurement, error)

	GetChickenProcurementCashflow(id uint64) (entity.ChickenProcurement, error)
	GetWarehouseItemProcurementCashflow(id uint64) (entity.WarehouseItemProcurement, error)
	GetWarehouseItemCornProcurementCashflow(id uint64) (entity.WarehouseItemCornProcurement, error)

	GetCashflowHistories(filter dto.GetCashflowHistoryFilter) ([]entity.CashflowHistory, error)
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
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
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
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
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
		query = query.Where("DATE(created_at) >= ? AND DATE(created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
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

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseSalePayment{}, errx.NotFound("afkir chicken sale payment not fount")
		}

		return entity.WarehouseSalePayment{}, nil
	}

	return payment, nil
}

func (r *CashflowRepository) GetStoreSalePaymentById(id uint64) (entity.StoreSalePayment, error) {
	var payment entity.StoreSalePayment
	err := r.GetDB().
		Preload("StoreSale.Customer").
		Preload("StoreSale.Item").
		Preload("StoreSale.Store.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreSalePayment{}, errx.NotFound("afkir chicken sale payment not fount")
		}

		return entity.StoreSalePayment{}, nil
	}

	return payment, nil
}

func (r *CashflowRepository) GetAfkirChickenSalePaymentById(id uint64) (entity.AfkirChickenSalePayment, error) {
	var payment entity.AfkirChickenSalePayment
	err := r.GetDB().
		Preload("AfkirChickenSale.AfkirChickenCustomer").
		Preload("AfkirChickenSale.ChickenCage.Cage.Location").
		Preload("CreatedByUser").
		First(&payment, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AfkirChickenSalePayment{}, errx.NotFound("afkir chicken sale payment not fount")
		}

		return entity.AfkirChickenSalePayment{}, nil
	}
	return payment, nil
}

func (r *CashflowRepository) CreateExpense(data *entity.Expense) error {
	return r.GetDB().Model(&entity.Expense{}).Create(data).Error
}

func (r *CashflowRepository) GetExpense(id uint64) (entity.Expense, error) {
	var data entity.Expense
	err := r.GetDB().Model(&entity.Expense{}).Where("id = ?", id).Preload("Location").Preload("Cage.Location").Preload("Warehouse.Location").Preload("Store.Location").Preload("CreatedByUser").First(&data).Error
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

	err := query.Preload("Location").Preload("Cage.Location").Preload("Warehouse.Location").Preload("Store.Location").Find(&data).Error
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
	query := r.GetDB().Model(&entity.UserSalaryPayment{}).Joins("JOIN users ON user_salary_payments.user_id = users.id")

	if filter.LocationId > 0 {
		query = query.Where("users.location_id = ?", filter.LocationId)
	}

	if filter.RoleId > 0 {
		query = query.Where("users.role_id = ?", filter.LocationId)
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("users.name LIKE ?", keyword)
	}

	if filter.Page > 0 {
		query = query.Offset(int((filter.Page - 1) * constant.PaginationDefaultLimit)).Limit(int(constant.PaginationDefaultLimit))
	}

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(user_salary_payments.created_at) >= ? AND DATE(user_salary_payments.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.IsPaid != nil {
		query = query.Where("is_paid = ?", filter.IsPaid)
	}

	err := query.
		Preload("User.Role").
		Find(&data).Error
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (r *CashflowRepository) CountUserSalaryPayments(filter dto.GetUserSalaryPaymentFilter) (int64, error) {
	var count int64
	query := r.GetDB().Model(&entity.UserSalaryPayment{}).Joins("JOIN users ON user_salary_payments.user_id = users.id")

	if filter.LocationId > 0 {
		query = query.Where("users.location_id = ?", filter.LocationId)
	}

	if filter.RoleId > 0 {
		query = query.Where("users.role_id = ?", filter.LocationId)
	}

	if filter.Keyword != "" {
		keyword := "%" + filter.Keyword + "%"
		query = query.Where("users.name LIKE ?", keyword)
	}

	if !filter.EndDate.Value().IsZero() && !filter.StartDate.Value().IsZero() {
		query = query.Where("DATE(user_salary_payments.created_at) >= ? AND DATE(user_salary_payments.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.IsPaid != nil {
		query = query.Where("is_paid = ?", filter.IsPaid)
	}

	err := query.
		Count(&count).Error
	if err != nil {
		return -1, err
	}
	return count, nil
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

func (r *CashflowRepository) CreateUserCashAdvance(data *entity.UserCashAdvance) error {
	return r.GetDB().Model(&entity.UserCashAdvance{}).Create(data).Error
}

func (r *CashflowRepository) CreateUserCashAdvancePayment(data *entity.UserCashAdvancePayment) error {
	return r.GetDB().Model(&entity.UserCashAdvancePayment{}).Create(data).Error
}

func (r *CashflowRepository) GetUserCashAdvance(id uint64) (entity.UserCashAdvance, error) {
	var data entity.UserCashAdvance
	err := r.GetDB().Model(&entity.UserCashAdvance{}).Where("id = ?", id).Preload("User.Role").Preload("User.Location").Preload("Payments").Preload("CreatedByUser").First(&data).Error
	if err != nil {
		return entity.UserCashAdvance{}, err
	}

	return data, nil
}

func (r *CashflowRepository) GetUserCashAdvances(filter dto.GetUserCashAdvanceFilter) ([]entity.UserCashAdvance, error) {
	var data []entity.UserCashAdvance
	query := r.GetDB().Model(&entity.UserCashAdvance{})

	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	if filter.UserId != uuid.Nil {
		query = query.Where("user_id = ?", filter.UserId)
	}

	if filter.PaymentStatus.Value().IsValid() {
		query = query.Where("payment_status = ?", filter.PaymentStatus.Value())
	}

	if filter.PaymentStatuses != nil {
		paymentStatus := make([]enum.PaymentStatus, 0)
		for _, e := range filter.PaymentStatuses {
			paymentStatus = append(paymentStatus, e.Value())
		}

		query = query.Where("payment_status IN ?", paymentStatus)
	}

	err := query.Preload("User.Location").Preload("Payments").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CashflowRepository) UpdateUserCashAdvance(data *entity.UserCashAdvance) error {
	return r.GetDB().Model(&entity.UserCashAdvance{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"user_id":               data.UserId,
		"nominal":               data.Nominal,
		"deadline_payment_date": data.DeadlinePaymentDate,
		"payment_status":        data.PaymentStatus,
		"updated_by":            data.UpdatedBy,
	}).Error
}

func (r *CashflowRepository) GetStoreSaleCashflows(filter dto.GetStoreSaleFilter) ([]entity.StoreSale, error) {
	var storeSales []entity.StoreSale
	query := r.GetDB().Model(&entity.StoreSale{})

	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	err := query.Preload("Store.Location").Preload("Customer").Preload("Item").Preload("Payments").Find(&storeSales).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return storeSales, nil
}

func (r *CashflowRepository) GetUserCashAdvancePayments(filter dto.GetUserCashAdvancePaymentFilter) ([]entity.UserCashAdvancePayment, error) {
	var data []entity.UserCashAdvancePayment

	query := r.GetDB().Model(&entity.UserCashAdvancePayment{}).
		Joins("LEFT JOIN user_cash_advances ON user_cash_advances.id = user_cash_advance_payments.user_cash_advance_id").
		Joins("LEFT JOIN users ON  users.id = user_cash_advances.user_id")

	if !filter.StartDate.Value().IsZero() && !filter.EndDate.Value().IsZero() {
		query = query.Where("DATE(user_cash_advance_payments.created_at) >= ? AND DATE(user_cash_advance_payments.created_at) <= ?", filter.StartDate.Value(), filter.EndDate.Value())
	}

	if filter.LocationId > 0 {
		query = query.Where("users.location_id = ?", filter.LocationId)
	}

	err := query.Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CashflowRepository) GetUserCashAdvancePayment(id uint64) (entity.UserCashAdvancePayment, error) {
	var data entity.UserCashAdvancePayment
	err := r.GetDB().Model(&entity.UserCashAdvancePayment{}).Where("id = ?").Preload("UserCashAdvance.User.Location").Preload("CreatedByUser").First(&data).Error
	if err != nil {
		return entity.UserCashAdvancePayment{}, err
	}

	return data, nil
}

func (r *CashflowRepository) GetWarehouseSaleCashflows(filter dto.GetWarehouseSaleFilter) ([]entity.WarehouseSale, error) {
	var warehouseSales []entity.WarehouseSale
	query := r.GetDB().Model(&entity.WarehouseSale{})

	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	err := query.Preload("Warehouse.Location").Preload("Customer").Preload("Item").Preload("Payments").Find(&warehouseSales).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return warehouseSales, nil
}

func (r *CashflowRepository) GetAfkirChickenSaleCashflows(filter dto.GetAfkirChickenSaleFilter) ([]entity.AfkirChickenSale, error) {
	var afkirChickenSales []entity.AfkirChickenSale
	query := r.GetDB().Model(&entity.AfkirChickenSale{})

	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	err := query.Preload("ChickenCage.Cage.Location").Preload("AfkirChickenCustomer").Preload("Payments").Find(&afkirChickenSales).Order("created_at DESC").Error
	if err != nil {
		return nil, err
	}
	return afkirChickenSales, nil
}

func (r *CashflowRepository) GetAfkirChickenSaleCashflow(id uint64) (entity.AfkirChickenSale, error) {
	var afkirChickenSale entity.AfkirChickenSale
	err := r.GetDB().Model(&entity.AfkirChickenSale{}).Preload("ChickenCage.Cage.Location").Preload("AfkirChickenCustomer").Preload("Payments").Preload("CreatedByUser").Where("id = ?", id).First(&afkirChickenSale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.AfkirChickenSale{}, errx.NotFound("afkir chicken sale not found")
		}
		return entity.AfkirChickenSale{}, err
	}

	return afkirChickenSale, nil
}

func (r *CashflowRepository) GetWarehouseSaleCashflow(id uint64) (entity.WarehouseSale, error) {
	var warehouseSale entity.WarehouseSale
	err := r.GetDB().Model(&entity.WarehouseSale{}).Where("id = ?", id).Preload("Warehouse.Location").Preload("Customer").Preload("Item").Preload("Payments").Preload("CreatedByUser").First(&warehouseSale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.WarehouseSale{}, errx.NotFound("warehouse sale not found")
		}
		return entity.WarehouseSale{}, err
	}

	return warehouseSale, nil
}

func (r *CashflowRepository) GetStoreSaleCashflow(id uint64) (entity.StoreSale, error) {
	var storeSale entity.StoreSale

	err := r.GetDB().Model(&entity.StoreSale{}).Preload("Store.Location").Preload("Customer").Preload("Item").Preload("Payments").Preload("CreatedByUser").Where("id = ?", id).First(&storeSale).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.StoreSale{}, errx.NotFound("warehouse sale not found")
		}
		return entity.StoreSale{}, err
	}

	return storeSale, nil
}

func (r *CashflowRepository) GetUserSalaryPayment(id uint64) (entity.UserSalaryPayment, error) {
	var data entity.UserSalaryPayment
	err := r.GetDB().Model(&entity.UserSalaryPayment{}).Preload("User.Role").Preload("CreatedByUser").Where("id = ?", id).First(&data).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.UserSalaryPayment{}, errx.NotFound("user cash advance not found")
		}
		return entity.UserSalaryPayment{}, err
	}

	return data, nil
}

func (r *CashflowRepository) UpdateUserSalaryPayment(data *entity.UserSalaryPayment) error {
	return r.GetDB().Model(&entity.UserSalaryPayment{}).Where("id = ?", data.Id).Updates(map[string]interface{}{
		"user_id":                data.UserId,
		"base_salary":            data.BaseSalary,
		"bonus_salary":           data.BonusSalary,
		"compentation_salary":    data.CompentationSalary,
		"additional_work_salary": data.AdditionalWorkSalary,
		"payment_proof":          data.PaymentProof,
		"payment_method":         data.PaymentMethod,
		"is_paid":                data.IsPaid,
		"updated_by":             data.UpdatedBy,
	}).Error
}

func (r *CashflowRepository) GetWarehouseItemProcurementCashflows(filter dto.GetWarehouseItemProcurementFilter) ([]entity.WarehouseItemProcurement, error) {
	var data []entity.WarehouseItemProcurement

	query := r.GetDB().Model(&entity.WarehouseItemProcurement{})
	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	if filter.PaymentStatuses != nil {
		paymentStatus := make([]enum.PaymentStatus, 0)
		for _, e := range filter.PaymentStatuses {
			paymentStatus = append(paymentStatus, e.Value())
		}

		query = query.Where("payment_status IN ?", paymentStatus)
	}

	err := query.Preload("Warehouse.Location").Preload("Supplier").Preload("Payments").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CashflowRepository) GetWarehouseItemCornProcurementCashflows(filter dto.GetWarehouseItemCornProcurementFilter) ([]entity.WarehouseItemCornProcurement, error) {
	var data []entity.WarehouseItemCornProcurement

	query := r.GetDB().Model(&entity.WarehouseItemCornProcurement{})
	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	if filter.PaymentStatuses != nil {
		paymentStatus := make([]enum.PaymentStatus, 0)
		for _, e := range filter.PaymentStatuses {
			paymentStatus = append(paymentStatus, e.Value())
		}

		query = query.Where("payment_status IN ?", paymentStatus)
	}

	err := query.Preload("Warehouse.Location").Preload("Supplier").Preload("Payments").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CashflowRepository) GetChickenProcurementCashflows(filter dto.GetChickenProcurementFilter) ([]entity.ChickenProcurement, error) {
	var data []entity.ChickenProcurement

	query := r.GetDB().Model(&entity.ChickenProcurement{})
	if !filter.DeadlinePaymentStartDate.Value().IsZero() && !filter.DeadlinePaymentEndDate.Value().IsZero() {
		query = query.Where("DATE(deadline_payment_date) >= ? AND DATE(deadline_payment_date) <= ?", filter.DeadlinePaymentStartDate.Value(), filter.DeadlinePaymentEndDate.Value())
	}

	if filter.PaymentStatuses != nil {
		paymentStatus := make([]enum.PaymentStatus, 0)
		for _, e := range filter.PaymentStatuses {
			paymentStatus = append(paymentStatus, e.Value())
		}

		query = query.Where("payment_status IN ?", paymentStatus)
	}

	err := query.Preload("Cage.Location").Preload("Supplier").Preload("Payments").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (r *CashflowRepository) GetChickenProcurementCashflow(id uint64) (entity.ChickenProcurement, error) {
	var data entity.ChickenProcurement
	err := r.GetDB().Model(&entity.ChickenProcurement{}).Where("id = ?", id).Preload("Cage.Location").Preload("Supplier").Preload("Payments").Preload("CreatedByUser").First(&data).Error
	if err != nil {
		return entity.ChickenProcurement{}, err
	}

	return data, nil
}

func (r *CashflowRepository) GetWarehouseItemCornProcurementCashflow(id uint64) (entity.WarehouseItemCornProcurement, error) {
	var data entity.WarehouseItemCornProcurement
	err := r.GetDB().Model(&entity.WarehouseItemCornProcurement{}).Where("id = ?", id).Preload("Warehouse.Location").Preload("Supplier").Preload("Payments").Preload("CreatedByUser").First(&data).Error
	if err != nil {
		return entity.WarehouseItemCornProcurement{}, err
	}

	return data, nil

}

func (r *CashflowRepository) GetWarehouseItemProcurementCashflow(id uint64) (entity.WarehouseItemProcurement, error) {
	var data entity.WarehouseItemProcurement
	err := r.GetDB().Model(&entity.WarehouseItemProcurement{}).Where("id = ?", id).Preload("Warehouse.Location").Preload("Supplier").Preload("Payments").Preload("CreatedByUser").First(&data).Error
	if err != nil {
		return entity.WarehouseItemProcurement{}, err
	}

	return data, nil

}

func (r *CashflowRepository) CreateUserCashAdvancePaymentBatch(payments *[]entity.UserCashAdvancePayment) error {
	err := r.GetDB().Model(&entity.UserCashAdvancePayment{}).CreateInBatches(payments, len(*payments)).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *CashflowRepository) GetCashflowHistories(filter dto.GetCashflowHistoryFilter) ([]entity.CashflowHistory, error) {
	var data []entity.CashflowHistory

	query := r.GetDB().Model(&entity.CashflowHistory{})

	if filter.Year > 0 {
		query = query.Where("EXTRACT(year FROM created_at) = ?", filter.Year)
	}

	err := query.Order("EXTRACT(month FROM created_at) ASC").Find(&data).Error
	if err != nil {
		return nil, err
	}

	return data, err
}
