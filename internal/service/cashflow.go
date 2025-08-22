package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"github.com/shopspring/decimal"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

type CashflowService struct {
	log              *zap.Logger
	repository       repository.ICashflowRepository
	storeService     IStoreService
	warehouseService IWarehouseService
	chickenService   IChickenService
	userService      IUserService
	workService      IWorkService
}

type ICashflowService interface {
	GetIncomeOverview(filter dto.GetIncomeOverviewFilter) (dto.IncomeOverviewResponse, error)
	GetIncome(incomeCategory string, id uint64) (dto.IncomeResponse, error)

	CreateExpense(request dto.CreateExpenseRequest, userId uuid.UUID) (dto.ExpenseResponse, error)
	GetExpenseOverview(filter dto.GetExpenseOverviewFilter) (dto.ExpenseOverviewResponse, error)
	GetExpense(expenseCategory string, id uint64) (dto.ExpenseResponse, error)

	GetUserCashAdvanceByUserId(userId uuid.UUID) ([]dto.UserCashAdvanceSummaryResponse, error)
	CreateUserCashAdvance(request dto.CreateUserCashAdvanceRequest, userId uuid.UUID) (dto.UserCashAdvanceResponse, error)
	CreateUserCashAdvancePayment(userCashAdvanceId uint64, request dto.CreateUserCashAdvancePaymentRequest, userId uuid.UUID) (dto.UserCashAdvanceResponse, error)

	GetReceiveablesOverview(filter dto.GetReceivablesOverviewFilter) (dto.ReceievablesOverviewResponse, error)
	GetReceiveables(receieveablesCategory string, id uint64) (dto.ReceiveablesResponse, error)

	PayUserSalaryPayment(id uint64, request dto.PayUserSalaryPaymentRequest, userId uuid.UUID) (dto.UserSalaryPaymentResponse, error)

	GetDebtOverview(filter dto.GetDebtOverviewFilter) (dto.DebtOverviewResponse, error)
	GetDebt(debtCategory string, id uint64) (dto.DebtResponse, error)

	GetUserSalarySummary(filter dto.GetUserSalarySummaryFilter) (dto.UserSalarySummaryResponse, error)
	GetUserSalaries(filter dto.GetUserSalaryListFilter) (dto.UserSalaryListPaginationResponse, error)
	GetUserSalaryDetail(id uint64) (dto.UserSalaryDetailResponse, error)

	ExportSalesCashflowToExcel(filter dto.GetSaleCashflowFilter) (*excelize.File, error)
}

func NewCashflowService(log *zap.Logger, repository repository.ICashflowRepository, storeService IStoreService, warehouseService IWarehouseService, chickenService IChickenService, userService IUserService, workService IWorkService) ICashflowService {
	return &CashflowService{
		log:              log,
		repository:       repository,
		storeService:     storeService,
		warehouseService: warehouseService,
		chickenService:   chickenService,
		userService:      userService,
		workService:      workService,
	}
}

func (s *CashflowService) GetIncomeOverview(filter dto.GetIncomeOverviewFilter) (dto.IncomeOverviewResponse, error) {
	incomeResponses := make([]dto.IncomeListResponse, 0)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	warehouseSalePayments, err := s.repository.GetWarehouseSalePayments(dto.GetWarehouseSalePaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse sale payments", zap.Error(err))
		return dto.IncomeOverviewResponse{}, err
	}

	storeSalePayments, err := s.repository.GetStoreSalePayments(dto.GetStoreSalePaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get store sale payments", zap.Error(err))
		return dto.IncomeOverviewResponse{}, err
	}

	afkirChickenSalePayments, err := s.repository.GetAfkirChickenSalePayments(dto.GetAfkirChickenSalePaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get afkir chicken sale payments", zap.Error(err))
		return dto.IncomeOverviewResponse{}, err
	}

	totalPayment := decimal.Zero
	totalWarehouseSalePayment := decimal.Zero
	totalStoreSalePayment := decimal.Zero
	totalAfkirChickenSalePayment := decimal.Zero

	for _, payment := range storeSalePayments {
		totalPayment = totalPayment.Add(payment.Nominal)
		totalStoreSalePayment = totalStoreSalePayment.Add(payment.Nominal)
	}

	for _, payment := range warehouseSalePayments {
		totalPayment = totalPayment.Add(payment.Nominal)
		totalWarehouseSalePayment = totalWarehouseSalePayment.Add(payment.Nominal)
	}

	for _, payment := range afkirChickenSalePayments {
		totalPayment = totalPayment.Add(payment.Nominal)
		totalAfkirChickenSalePayment = totalAfkirChickenSalePayment.Add(payment.Nominal)
	}

	if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryStoreEggSale {
		for _, payment := range storeSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.StoreSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("02 Jan 2006"),
				PlaceName:    payment.StoreSale.Store.Location.Name + " - " + payment.StoreSale.Store.Name,
				Category:     constant.IncomeCategoryStoreEggSale,
				ItemName:     payment.StoreSale.Item.Name,
				ItemUnit:     payment.StoreSale.SaleUnit.String(),
				Quantity:     fmt.Sprintf("%v", payment.StoreSale.Quantity),
				CustomerName: payment.StoreSale.Customer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	} else if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryWarehouseEggSale {
		for _, payment := range warehouseSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.WarehouseSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("02 Jan 2006"),
				PlaceName:    payment.WarehouseSale.Warehouse.Location.Name + " - " + payment.WarehouseSale.Warehouse.Name,
				Category:     constant.IncomeCategoryWarehouseEggSale,
				ItemName:     payment.WarehouseSale.Item.Name,
				ItemUnit:     payment.WarehouseSale.SaleUnit.String(),
				Quantity:     fmt.Sprintf("%v", payment.WarehouseSale.Quantity),
				CustomerName: payment.WarehouseSale.Customer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	} else if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryAfkirChickenSale {
		for _, payment := range afkirChickenSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.AfkirChickenSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("02 Jan 2006"),
				PlaceName:    payment.AfkirChickenSale.ChickenCage.Cage.Location.Name + " - " + payment.AfkirChickenSale.ChickenCage.Cage.Name,
				Category:     constant.IncomeCategoryAfkirChickenSale,
				ItemName:     "Afkir Chicken",
				ItemUnit:     "Ekor",
				Quantity:     fmt.Sprintf("%v", payment.AfkirChickenSale.TotalSellChicken),
				CustomerName: payment.AfkirChickenSale.AfkirChickenCustomer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	}

	return dto.IncomeOverviewResponse{
		IncomePie: dto.IncomePieResponse{
			WarehouseEggSalePercentage: totalWarehouseSalePayment.Div(totalPayment).InexactFloat64() * 100.0,
			StoreEggSalePercentage:     totalStoreSalePayment.Div(totalPayment).InexactFloat64() * 100.0,
			AfkirChickenSalePercentage: totalAfkirChickenSalePayment.Div(totalPayment).InexactFloat64() * 100.0,
		},
		Incomes: incomeResponses,
	}, nil
}

func (s *CashflowService) GetIncome(incomeCategory string, id uint64) (dto.IncomeResponse, error) {
	switch incomeCategory {
	case constant.IncomeCategoryWarehouseEggSale:
		payment, err := s.repository.GetWarehouseSalePaymentById(id)
		if err != nil {
			s.log.Error("failed get warehouse sale payment", zap.Error(err))
			return dto.IncomeResponse{}, err
		}

		return dto.IncomeResponse{
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            "Warehouse Egg Sale",
			PlaceName:           payment.WarehouseSale.Warehouse.Name,
			CustomerName:        payment.WarehouseSale.Customer.Name,
			CustomerPhoneNumber: payment.WarehouseSale.Customer.PhoneNumber,
			ItemName:            payment.WarehouseSale.Item.Name,
			ItemUnit:            payment.WarehouseSale.SaleUnit.String(),
			Quantity:            fmt.Sprintf("%v", payment.WarehouseSale.Quantity),
			Nominal:             payment.Nominal.String(),
			PaymentType:         payment.WarehouseSale.PaymentType.String(),
			TotalPrice:          payment.WarehouseSale.TotalPrice.String(),
			PaymentMethod:       payment.PaymentMethod.String(),
			InputBy:             payment.CreatedByUser.Name,
			PaymentProof:        payment.PaymentProof,
		}, nil

	case constant.IncomeCategoryStoreEggSale:
		payment, err := s.repository.GetStoreSalePaymentById(id)
		if err != nil {
			s.log.Error("failed get store sale payment", zap.Error(err))
			return dto.IncomeResponse{}, err
		}

		return dto.IncomeResponse{
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            "Store Egg Sale",
			PlaceName:           payment.StoreSale.Store.Name,
			CustomerName:        payment.StoreSale.Customer.Name,
			CustomerPhoneNumber: payment.StoreSale.Customer.PhoneNumber,
			ItemName:            payment.StoreSale.Item.Name,
			ItemUnit:            payment.StoreSale.SaleUnit.String(),
			Quantity:            fmt.Sprintf("%v", payment.StoreSale.Quantity),
			Nominal:             payment.Nominal.String(),
			PaymentType:         payment.StoreSale.PaymentType.String(),
			TotalPrice:          payment.StoreSale.TotalPrice.String(),
			PaymentMethod:       payment.PaymentMethod.String(),
			InputBy:             payment.CreatedByUser.Name,
			PaymentProof:        payment.PaymentProof,
		}, nil

	case constant.IncomeCategoryAfkirChickenSale:
		payment, err := s.repository.GetAfkirChickenSalePaymentById(id)
		if err != nil {
			s.log.Error("failed get afkir chicken sale payment", zap.Error(err))
			return dto.IncomeResponse{}, err
		}

		return dto.IncomeResponse{
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            "Afkir Chicken Sale",
			PlaceName:           payment.AfkirChickenSale.ChickenCage.Cage.Name,
			CustomerName:        payment.AfkirChickenSale.AfkirChickenCustomer.Name,
			CustomerPhoneNumber: payment.AfkirChickenSale.AfkirChickenCustomer.PhoneNumber,
			ItemName:            "Afkir Chicken",
			ItemUnit:            "Ekor",
			Quantity:            fmt.Sprintf("%v", payment.AfkirChickenSale.TotalSellChicken),
			Nominal:             payment.Nominal.String(),
			PaymentType:         payment.AfkirChickenSale.PaymentType.String(),
			TotalPrice:          payment.AfkirChickenSale.TotalPrice.String(),
			PaymentMethod:       payment.PaymentMethod.String(),
			InputBy:             payment.CreatedByUser.Name,
			PaymentProof:        payment.PaymentProof,
		}, nil

	default:
		return dto.IncomeResponse{}, fmt.Errorf("invalid income category")
	}
}

func (s *CashflowService) CreateExpense(request dto.CreateExpenseRequest, userId uuid.UUID) (dto.ExpenseResponse, error) {
	s.repository.UseTx(false)

	expenseCategory := enum.ValueOfExpenseCategory(request.ExpenseCategory)
	if !expenseCategory.IsValid() {
		s.log.Error("invalid expense category", zap.String("expenseCategory", request.ExpenseCategory))
		return dto.ExpenseResponse{}, errx.BadRequest("invalid expense category")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("invalid nominal", zap.String("nominal", request.Nominal))
		return dto.ExpenseResponse{}, errx.BadRequest("invalid nominal")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.ExpenseResponse{}, errx.BadRequest("invalid payment method")
	}

	locationType := enum.ValueOfLocationWorkType(request.LocationType)
	if !locationType.IsValid() {
		s.log.Error("invalid location", zap.String("location", request.LocationType))
		return dto.ExpenseResponse{}, errx.BadRequest("invalid location")
	}

	if request.ReceiverPhoneNumber != "" && request.ReceiverPhoneNumber[:2] != "08" {
		s.log.Error("invalid phone number", zap.String("phoneNumber", request.LocationType))
		return dto.ExpenseResponse{}, errx.BadRequest("invalid phone number")
	} else if request.ReceiverPhoneNumber == "" {
		request.ReceiverPhoneNumber = "-"
	}

	data := entity.Expense{
		ExpenseCategory:     expenseCategory,
		Name:                request.Name,
		ReceiverName:        request.ReceiverName,
		ReceiverPhoneNumber: request.ReceiverPhoneNumber,
		Nominal:             nominal,
		PaymentMethod:       paymentMethod,
		PaymentProof:        request.PaymentProof,
		Description:         request.Description,
		LocationId:          request.LocationId,
		LocationType:        locationType,
		CreatedBy:           uuid.NullUUID{UUID: userId, Valid: true},
	}

	switch locationType {
	case enum.LocationTypeCage:
		data.CageId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
	case enum.LocationTypeStore:
		data.StoreId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
	case enum.LocationTypeWarehouse:
		data.WarehouseId = sql.NullInt64{Int64: int64(request.PlaceId), Valid: true}
	}

	err = s.repository.CreateExpense(&data)
	if err != nil {
		s.log.Error("failed create expense", zap.Error(err))
		return dto.ExpenseResponse{}, err
	}

	data, err = s.repository.GetExpense(data.Id)
	if err != nil {
		s.log.Error("failed get expense", zap.Error(err))
		return dto.ExpenseResponse{}, err
	}

	response := dto.ExpenseResponse{
		Id:                  data.Id,
		Date:                data.CreatedAt.Format("02 Jan 2006"),
		Time:                data.CreatedAt.Format("15:04"),
		Category:            data.ExpenseCategory.String(),
		Name:                data.Name,
		ReceiverName:        data.ReceiverName,
		ReceiverPhoneNumber: data.ReceiverPhoneNumber,
		Nominal:             data.Nominal.String(),
		PaymentMethod:       data.PaymentMethod.String(),
		PaymentProof:        data.PaymentProof,
		InputBy:             data.CreatedByUser.Name,
	}

	switch data.LocationType {
	case enum.LocationTypeCage:
		response.PlaceName = data.Cage.Name + " - " + data.Location.Name
	case enum.LocationTypeStore:
		response.PlaceName = data.Store.Name + " - " + data.Location.Name
	case enum.LocationTypeWarehouse:
		response.PlaceName = data.Warehouse.Name + " - " + data.Location.Name
	}

	return response, nil
}

func (s *CashflowService) GetExpenseOverview(filter dto.GetExpenseOverviewFilter) (dto.ExpenseOverviewResponse, error) {
	s.repository.UseTx(false)

	expenseResponses := make([]dto.ExpenseListResponse, 0)
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	chickenProcurementPayments, err := s.repository.GetChickenProcurementPayments(dto.GetChickenProcurementPaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get chicken procurement payments", zap.Error(err))
		return dto.ExpenseOverviewResponse{}, err
	}

	isPaid := true
	userSalaryPayments, err := s.repository.GetUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
		IsPaid:    &isPaid,
	})
	if err != nil {
		s.log.Error("failed get user salary payments", zap.Error(err))
		return dto.ExpenseOverviewResponse{}, err
	}

	warehouseItemProcurementPayments, err := s.repository.GetWarehouseItemProcurementPayments(dto.GetWarehouseItemProcurementPaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse item procurement payments", zap.Error(err))
		return dto.ExpenseOverviewResponse{}, err
	}

	warehouseItemCornProcurementPayments, err := s.repository.GetWarehouseItemCornProcurementPayments(dto.GetWarehouseItemCornProcurementPaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement payments", zap.Error(err))
		return dto.ExpenseOverviewResponse{}, err
	}

	expensePayments, err := s.repository.GetExpenses(dto.GetExpenseFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get expenses", zap.Error(err))
		return dto.ExpenseOverviewResponse{}, err
	}

	totalPayment := decimal.Zero
	totalChickenProcurement := decimal.Zero
	totalWarehouseItemProcurement := decimal.Zero
	totalWarehouseItemCornProcurement := decimal.Zero
	totalOperational := decimal.Zero
	totalOther := decimal.Zero
	totalStaffSalary := decimal.Zero

	for _, p := range chickenProcurementPayments {
		totalPayment = totalPayment.Add(p.Nominal)
		totalChickenProcurement = totalChickenProcurement.Add(p.Nominal)
	}

	for _, p := range expensePayments {
		totalPayment = totalPayment.Add(p.Nominal)
		switch p.ExpenseCategory {
		case enum.ExpenseCategoryOperational:
			totalOperational = totalOperational.Add(p.Nominal)
		case enum.ExpenseCategoryOther:
			totalOther = totalOther.Add(p.Nominal)
		}
	}

	for _, p := range warehouseItemProcurementPayments {
		totalPayment = totalPayment.Add(p.Nominal)
		totalWarehouseItemProcurement = totalWarehouseItemProcurement.Add(p.Nominal)
	}

	for _, p := range warehouseItemCornProcurementPayments {
		totalPayment = totalPayment.Add(p.Nominal)
		totalWarehouseItemCornProcurement = totalWarehouseItemCornProcurement.Add(p.Nominal)
	}

	for _, p := range userSalaryPayments {
		totalSalary := p.BaseSalary.Add(p.BonusSalary).Add(p.CompentationSalary).Add(p.AdditionalWorkSalary)
		totalPayment = totalPayment.Add(totalSalary)
		totalWarehouseItemProcurement = totalWarehouseItemProcurement.Add(totalSalary)
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryChickenProcurement {
		for _, p := range chickenProcurementPayments {
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.PaymentDate.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryChickenProcurement,
				Name:         constant.ExpenseTransactionNameChickenProcurement,
				PlaceName:    p.ChickenProcurement.Cage.Location.Name + " - " + p.ChickenProcurement.Cage.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.ChickenProcurement.Supplier.Name,
				PaymentProof: p.PaymentProof,
			})
		}
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryWarehouseItemProcurement {
		for _, p := range warehouseItemProcurementPayments {
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.PaymentDate.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryWarehouseItemProcurement,
				Name:         constant.ExpenseTransactionNameWarehouseItemProcurement,
				PlaceName:    p.WarehouseItemProcurement.Warehouse.Location.Name + " - " + p.WarehouseItemProcurement.Warehouse.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.WarehouseItemProcurement.Supplier.Name,
				PaymentProof: p.PaymentProof,
			})
		}
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryWarehouseItemCornProcurement {
		for _, p := range warehouseItemCornProcurementPayments {
			totalPayment = totalPayment.Add(p.Nominal)
			totalWarehouseItemProcurement = totalWarehouseItemProcurement.Add(p.Nominal)
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.PaymentDate.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryWarehouseItemCornProcurement,
				Name:         constant.ExpenseTransactionNameWarehouseItemCornProcurement,
				PlaceName:    p.WarehouseItemCornProcurement.Warehouse.Location.Name + " - " + p.WarehouseItemCornProcurement.Warehouse.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.WarehouseItemCornProcurement.Supplier.Name,
				PaymentProof: p.PaymentProof,
			})
		}
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryStaff {
		for _, p := range userSalaryPayments {
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.CreatedAt.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryStaff,
				Name:         constant.ExpenseTransactionNameSalary,
				PlaceName:    p.User.Location.Name,
				Nominal:      p.BaseSalary.Add(p.BonusSalary).Add(p.CompentationSalary).Add(p.AdditionalWorkSalary).String(),
				ReceiverName: p.User.Name,
				PaymentProof: p.PaymentProof,
			})
		}
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryOperational {
		for _, p := range expensePayments {
			response := dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.CreatedAt.Format("02 Jan 2006"),
				Category:     p.ExpenseCategory.String(),
				Name:         p.Name,
				PlaceName:    p.Location.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.ReceiverName,
				PaymentProof: p.PaymentProof,
			}

			switch p.LocationType {
			case enum.LocationTypeCage:
				response.PlaceName = p.Cage.Name + " - " + p.Location.Name
			case enum.LocationTypeStore:
				response.PlaceName = p.Store.Name + " - " + p.Location.Name
			case enum.LocationTypeWarehouse:
				response.PlaceName = p.Warehouse.Name + " - " + p.Location.Name
			}

			expenseResponses = append(expenseResponses, response)
		}
	}

	return dto.ExpenseOverviewResponse{
		ExpensePie: dto.ExpensePieResponse{
			StaffPercentage:                        totalStaffSalary.Div(totalPayment).InexactFloat64() * 100.0,
			ChikckenProcuremtnPercentage:           totalChickenProcurement.Div(totalPayment).InexactFloat64() * 100.0,
			WarehouseItemProcurementPercentage:     totalWarehouseItemProcurement.Div(totalPayment).InexactFloat64() * 100.0,
			WarehouseItemCornProcurementPercentage: totalWarehouseItemCornProcurement.Div(totalPayment).InexactFloat64() * 100.0,
			OperationalPercentage:                  totalOperational.Div(totalPayment).InexactFloat64() * 100.0,
			OtherPercentage:                        totalOther.Div(totalPayment).InexactFloat64() * 100.0,
		},
		Expenses: expenseResponses,
	}, nil
}

func (s *CashflowService) GetExpense(expenseCategory string, id uint64) (dto.ExpenseResponse, error) {
	s.repository.UseTx(false)

	switch expenseCategory {
	case constant.ExpenseCategoryOperational:
		expense, err := s.repository.GetExpense(id)
		if err != nil {
			s.log.Error("failed get operational expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		return dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.CreatedAt.Format("2006-01-02"),
			Time:                expense.CreatedAt.Format("15:04:05"),
			Category:            constant.ExpenseCategoryOperational,
			PlaceName:           expense.Location.Name,
			Name:                expense.Name,
			ReceiverName:        expense.ReceiverName,
			ReceiverPhoneNumber: expense.ReceiverPhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	case constant.ExpenseCategoryChickenProcurement:
		expense, err := s.repository.GetChickenProcurementPaymentById(id)
		if err != nil {
			s.log.Error("failed get chicken procurement expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		return dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.PaymentDate.Format("2006-01-02"),
			Time:                expense.PaymentDate.Format("15:04:05"),
			Category:            "Chicken Procurement",
			PlaceName:           expense.ChickenProcurement.Cage.Location.Name,
			Name:                "Chicken Procurement",
			ReceiverName:        expense.ChickenProcurement.Supplier.Name,
			ReceiverPhoneNumber: expense.ChickenProcurement.Supplier.PhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	case constant.ExpenseCategoryWarehouseItemProcurement:
		expense, err := s.repository.GetWarehouseItemProcurementPaymentById(id)
		if err != nil {
			s.log.Error("failed get warehouse item procurement expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		return dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.PaymentDate.Format("2006-01-02"),
			Time:                expense.PaymentDate.Format("15:04:05"),
			Category:            "Warehouse Procurement",
			PlaceName:           expense.WarehouseItemProcurement.Warehouse.Location.Name,
			Name:                expense.WarehouseItemProcurement.Item.Name,
			ReceiverName:        expense.WarehouseItemProcurement.Supplier.Name,
			ReceiverPhoneNumber: expense.WarehouseItemProcurement.Supplier.PhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	case constant.ExpenseCategoryWarehouseItemCornProcurement:
		expense, err := s.repository.GetWarehouseItemCornProcurementPaymentById(id)
		if err != nil {
			s.log.Error("failed get warehouse item procurement expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		return dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.PaymentDate.Format("2006-01-02"),
			Time:                expense.PaymentDate.Format("15:04:05"),
			Category:            "Warehouse Procurement",
			PlaceName:           expense.WarehouseItemCornProcurement.Warehouse.Location.Name,
			Name:                constant.Corn,
			ReceiverName:        expense.WarehouseItemCornProcurement.Supplier.Name,
			ReceiverPhoneNumber: expense.WarehouseItemCornProcurement.Supplier.PhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	case constant.ExpenseCategoryStaff:
		expense, err := s.repository.GetUserSalaryPaymentById(id)
		if err != nil {
			s.log.Error("failed get user salary expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		return dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.CreatedAt.Format("2006-01-02"),
			Time:                expense.CreatedAt.Format("15:04:05"),
			Category:            "Staff",
			PlaceName:           "Salary Payment",
			Name:                expense.User.Name,
			ReceiverName:        expense.User.Name,
			ReceiverPhoneNumber: expense.User.PhoneNumber,
			Nominal:             expense.BaseSalary.Add(expense.BonusSalary).Add(expense.CompentationSalary).Add(expense.AdditionalWorkSalary).String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	case constant.ExpenseCategoryOther:
		expense, err := s.repository.GetExpense(id)
		if err != nil {
			s.log.Error("failed get other expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		return dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.CreatedAt.Format("2006-01-02"),
			Time:                expense.CreatedAt.Format("15:04:05"),
			Category:            "Other",
			PlaceName:           expense.Location.Name,
			Name:                expense.Name,
			ReceiverName:        expense.ReceiverName,
			ReceiverPhoneNumber: expense.ReceiverPhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	default:
		return dto.ExpenseResponse{}, fmt.Errorf("invalid expense category")
	}
}

func (s *CashflowService) GetUserCashAdvanceByUserId(userId uuid.UUID) ([]dto.UserCashAdvanceSummaryResponse, error) {
	s.repository.UseTx(false)

	paymentStatuses := []param.PaymentStatusParam{param.PaymentStatusParam(enum.PaymentStatusNotPaid), param.PaymentStatusParam(enum.PaymentStatusUnpaid)}
	data, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		UserId:          userId,
		PaymentStatuses: paymentStatuses,
	})
	if err != nil {
		return nil, err
	}

	response := make([]dto.UserCashAdvanceSummaryResponse, 0)
	for _, e := range data {
		currentPayment := decimal.Zero
		for _, p := range e.Payments {
			currentPayment = currentPayment.Add(p.Nominal)
		}

		isMoreThanDeadlinePaymentDate := false
		if time.Now().After(e.DeadlinePaymentDate) {
			isMoreThanDeadlinePaymentDate = true
		}

		response = append(response, dto.UserCashAdvanceSummaryResponse{
			Id:                            e.Id,
			DeadlinePaymentDate:           e.DeadlinePaymentDate.Format("02 Jan 2006"),
			Nominal:                       e.Nominal.String(),
			RemainingPayment:              e.Nominal.Sub(currentPayment).String(),
			IsMoreThanDeadlinePaymentDate: isMoreThanDeadlinePaymentDate,
		})
	}

	return response, nil
}

func (s *CashflowService) CreateUserCashAdvance(request dto.CreateUserCashAdvanceRequest, userId uuid.UUID) (dto.UserCashAdvanceResponse, error) {
	s.repository.UseTx(false)

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		return dto.UserCashAdvanceResponse{}, errx.BadRequest("invalid nominal format")
	}

	deadlinePaymentDate, err := time.Parse("02-01-2006", request.DeadlinePaymentDate)
	if err != nil {
		return dto.UserCashAdvanceResponse{}, errx.BadRequest("invalid deadline payment date format")
	}

	data := entity.UserCashAdvance{
		UserId:              uuid.MustParse(request.UserId),
		Nominal:             nominal,
		DeadlinePaymentDate: deadlinePaymentDate,
		PaymentStatus:       enum.PaymentStatusNotPaid,
		CreatedBy:           uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateUserCashAdvance(&data)
	if err != nil {
		s.log.Error("failed create user cash advance")
		return dto.UserCashAdvanceResponse{}, err
	}

	data, err = s.repository.GetUserCashAdvance(data.Id)
	if err != nil {
		s.log.Error("failed get user cash advance", zap.Error(err))
		return dto.UserCashAdvanceResponse{}, err
	}

	userCashAdvancePayments := make([]dto.UserCashAdvancePaymentResponse, len(data.Payments))
	remainingPayment := data.Nominal
	for i, storeSalePayment := range data.Payments {
		userCashAdvancePayments[i] = mapper.UserCashAdvancePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		userCashAdvancePayments[i].Remaining = remainingPayment.String()
	}

	response := dto.UserCashAdvanceResponse{
		Id:                      data.Id,
		User:                    mapper.UserToListResponse(&data.User),
		Nominal:                 data.Nominal.String(),
		DeadlinePaymentDate:     data.DeadlinePaymentDate.Format("02 Jan 2006"),
		UserCashAdvancePayments: userCashAdvancePayments,
		RemainingPayment:        remainingPayment.String(),
	}

	return response, nil
}

func (s *CashflowService) CreateUserCashAdvancePayment(userCashAdvanceId uint64, request dto.CreateUserCashAdvancePaymentRequest, userId uuid.UUID) (dto.UserCashAdvanceResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	data, err := s.repository.GetUserCashAdvance(userCashAdvanceId)
	if err != nil {
		s.log.Error("failed get user cash advance", zap.Error(err))
		return dto.UserCashAdvanceResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		return dto.UserCashAdvanceResponse{}, errx.BadRequest("invalid payment date format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.UserCashAdvanceResponse{}, errx.BadRequest("invalid payment method")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		return dto.UserCashAdvanceResponse{}, errx.BadRequest("invalid nominal format")
	}

	currentPayment := nominal
	for _, payment := range data.Payments {
		currentPayment = currentPayment.Add(payment.Nominal)
	}

	if currentPayment.GreaterThan(data.Nominal) {
		return dto.UserCashAdvanceResponse{}, errx.BadRequest("nominal more than needed")
	} else if currentPayment.Equal(data.Nominal) {
		data.PaymentStatus = enum.PaymentStatusPaid
	} else if currentPayment.LessThan(data.Nominal) {
		data.PaymentStatus = enum.PaymentStatusUnpaid
	}

	payment := entity.UserCashAdvancePayment{
		UserCashAdvanceId: userCashAdvanceId,
		Nominal:           nominal,
		PaymentDate:       paymentDate,
		PaymentMethod:     paymentMethod,
		PaymentProof:      request.PaymentProof,
	}

	err = s.repository.CreateUserCashAdvancePayment(&payment)
	if err != nil {
		s.log.Error("failed create user cash advance payment", zap.Error(err))
		return dto.UserCashAdvanceResponse{}, err
	}

	err = s.repository.UpdateUserCashAdvance(&data)
	if err != nil {
		s.log.Error("failed update user cash advance", zap.Error(err))
		return dto.UserCashAdvanceResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.UserCashAdvanceResponse{}, err
	}

	data, err = s.repository.GetUserCashAdvance(data.Id)
	if err != nil {
		s.log.Error("failed get user cash advance", zap.Error(err))
		return dto.UserCashAdvanceResponse{}, err
	}

	userCashAdvancePayments := make([]dto.UserCashAdvancePaymentResponse, len(data.Payments))
	remainingPayment := data.Nominal
	for i, storeSalePayment := range data.Payments {
		userCashAdvancePayments[i] = mapper.UserCashAdvancePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		userCashAdvancePayments[i].Remaining = remainingPayment.String()
	}

	response := dto.UserCashAdvanceResponse{
		Id:                      data.Id,
		User:                    mapper.UserToListResponse(&data.User),
		Nominal:                 data.Nominal.String(),
		DeadlinePaymentDate:     data.DeadlinePaymentDate.Format("02 Jan 2006"),
		UserCashAdvancePayments: userCashAdvancePayments,
		RemainingPayment:        remainingPayment.String(),
	}

	return response, nil
}

func (s *CashflowService) GetReceiveablesOverview(filter dto.GetReceivablesOverviewFilter) (dto.ReceievablesOverviewResponse, error) {
	s.repository.UseTx(false)

	receieveables := make([]dto.ReceiveablesListResponse, 0)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))
	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get store sale cashflows", zap.Error(err))
		return dto.ReceievablesOverviewResponse{}, err
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse sale cashflows", zap.Error(err))
		return dto.ReceievablesOverviewResponse{}, err
	}

	afkirChickenSales, err := s.repository.GetAfkirChickenSaleCashflows(dto.GetAfkirChickenSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get afkir chicken sale cashflows", zap.Error(err))
		return dto.ReceievablesOverviewResponse{}, err
	}

	userCashAdvances, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get user cash advances", zap.Error(err))
		return dto.ReceievablesOverviewResponse{}, err
	}

	totalPayment := decimal.Zero
	totalReceivablesPayment := decimal.Zero
	totalRemainingReceieveablesPayment := decimal.Zero

	for _, e := range storeSales {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentReceieveablesPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentReceieveablesPayment = totalCurrentReceieveablesPayment.Add(p.Nominal)
		}
		totalReceivablesPayment = totalReceivablesPayment.Add(totalCurrentReceieveablesPayment)
		totalRemainingReceieveablesPayment = totalRemainingReceieveablesPayment.Add(e.TotalPrice.Sub(totalCurrentReceieveablesPayment))
	}

	for _, e := range warehouseSales {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentReceieveablesPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentReceieveablesPayment = totalCurrentReceieveablesPayment.Add(p.Nominal)
		}
		totalReceivablesPayment = totalReceivablesPayment.Add(totalCurrentReceieveablesPayment)
		totalRemainingReceieveablesPayment = totalRemainingReceieveablesPayment.Add(e.TotalPrice.Sub(totalCurrentReceieveablesPayment))
	}

	for _, e := range afkirChickenSales {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentReceieveablesPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentReceieveablesPayment = totalCurrentReceieveablesPayment.Add(p.Nominal)
		}
		totalReceivablesPayment = totalReceivablesPayment.Add(totalCurrentReceieveablesPayment)
		totalRemainingReceieveablesPayment = totalRemainingReceieveablesPayment.Add(e.TotalPrice.Sub(totalCurrentReceieveablesPayment))
	}

	for _, e := range userCashAdvances {
		totalPayment = totalPayment.Add(e.Nominal)
		totalCurrentReceieveablesPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentReceieveablesPayment = totalCurrentReceieveablesPayment.Add(p.Nominal)
		}
		totalReceivablesPayment = totalReceivablesPayment.Add(totalCurrentReceieveablesPayment)
		totalRemainingReceieveablesPayment = totalRemainingReceieveablesPayment.Add(e.Nominal.Sub(totalCurrentReceieveablesPayment))
	}

	if filter.ReceieveablesCategory == constant.ReceieveablesCategoryAll || filter.ReceieveablesCategory == constant.ReceieveablesCategoryWarehouseEggSale {
		for _, e := range warehouseSales {
			receieveable := dto.ReceiveablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryWarehouseEggSale,
				PlaceName:           e.Warehouse.Location.Name + " - " + e.Warehouse.Name,
				Name:                e.Customer.Name,
				PhoneNumber:         e.Customer.PhoneNumber,
				TotalNominal:        e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	if filter.ReceieveablesCategory == constant.ReceieveablesCategoryAll || filter.ReceieveablesCategory == constant.ReceieveablesCategoryStoreEggSale {
		for _, e := range storeSales {
			receieveable := dto.ReceiveablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryWarehouseEggSale,
				PlaceName:           e.Store.Location.Name + " - " + e.Store.Name,
				Name:                e.Customer.Name,
				PhoneNumber:         e.Customer.PhoneNumber,
				TotalNominal:        e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	if filter.ReceieveablesCategory == constant.ReceieveablesCategoryAll || filter.ReceieveablesCategory == constant.ReceieveablesCategoryAfkirChickenSale {
		for _, e := range afkirChickenSales {
			receieveable := dto.ReceiveablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryWarehouseEggSale,
				PlaceName:           e.ChickenCage.Cage.Location.Name + " - " + e.ChickenCage.Cage.Name,
				Name:                e.AfkirChickenCustomer.Name,
				PhoneNumber:         e.AfkirChickenCustomer.PhoneNumber,
				TotalNominal:        e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	if filter.ReceieveablesCategory == constant.ReceieveablesCategoryAll || filter.ReceieveablesCategory == constant.ReceieveablesCategoryCashAdvance {
		for _, e := range userCashAdvances {
			receieveable := dto.ReceiveablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryWarehouseEggSale,
				PlaceName:           e.User.Location.Name,
				Name:                e.User.Name,
				PhoneNumber:         e.User.PhoneNumber,
				TotalNominal:        e.Nominal.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.Nominal.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	return dto.ReceievablesOverviewResponse{
		ReceivablesPie: dto.ReceiveablesPieResponse{
			PaidPercentage:   totalReceivablesPayment.Sub(totalPayment).InexactFloat64() * 100.0,
			UnpaidPercentage: totalRemainingReceieveablesPayment.Sub(totalPayment).InexactFloat64() * 100.0,
		},
		Receivables: receieveables,
	}, nil
}

func (s *CashflowService) GetReceiveables(receieveablesCategory string, id uint64) (dto.ReceiveablesResponse, error) {
	s.repository.UseTx(false)

	switch receieveablesCategory {
	case constant.ReceieveablesCategoryWarehouseEggSale:
		data, err := s.repository.GetWarehouseSaleCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceiveablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceieveablesPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.ReceieveablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.ReceiveablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:             data.Warehouse.Location.Name + " - " + data.Warehouse.Name,
			Name:                  data.Customer.Name,
			PhoneNumber:           data.Customer.PhoneNumber,
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           data.PaymentType.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}, nil
	case constant.ReceieveablesCategoryStoreEggSale:
		data, err := s.repository.GetStoreSaleCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceiveablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceieveablesPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.ReceieveablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.ReceiveablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:             data.Store.Location.Name + " - " + data.Store.Name,
			Name:                  data.Customer.Name,
			PhoneNumber:           data.Customer.PhoneNumber,
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           data.PaymentType.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}, nil
	case constant.ReceieveablesCategoryAfkirChickenSale:
		data, err := s.repository.GetAfkirChickenSaleCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceiveablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceieveablesPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.ReceieveablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.ReceiveablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:             data.ChickenCage.Cage.Location.Name + " - " + data.ChickenCage.Cage.Name,
			Name:                  data.AfkirChickenCustomer.Name,
			PhoneNumber:           data.AfkirChickenCustomer.PhoneNumber,
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           data.PaymentType.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}, nil

	case constant.ReceieveablesCategoryCashAdvance:
		data, err := s.repository.GetUserCashAdvance(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceiveablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceieveablesPaymentResponse, 0)
		totalRemainingPayment := data.Nominal
		for _, e := range data.Payments {
			paymentResponse := dto.ReceieveablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.ReceiveablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:             data.User.Location.Name,
			Name:                  data.User.Name,
			PhoneNumber:           data.User.PhoneNumber,
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           enum.PaymentTypeinstallment.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}, nil
	default:
		return dto.ReceiveablesResponse{}, errx.BadRequest("invalid receieveabels category")
	}
}

func (s *CashflowService) PayUserSalaryPayment(id uint64, request dto.PayUserSalaryPaymentRequest, userId uuid.UUID) (dto.UserSalaryPaymentResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	userSalaryPayment, err := s.repository.GetUserSalaryPayment(id)
	if err != nil {
		s.log.Error("failed get user salary payment", zap.Error(err))
		return dto.UserSalaryPaymentResponse{}, err
	}

	baseSalary, err := decimal.NewFromString(request.BaseSalary)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid base salary format")
	}

	bonusSalary, err := decimal.NewFromString(request.BonusSalary)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid bonus salary format")
	}

	compentationSalary, err := decimal.NewFromString(request.CompentationSalary)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid compentation salary format")
	}

	additionalWorkSalary, err := decimal.NewFromString(request.AdditionalWorkSalary)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid additional work salary format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid payment method")
	}

	userSalaryPayment.BaseSalary = baseSalary
	userSalaryPayment.CompentationSalary = compentationSalary
	userSalaryPayment.BonusSalary = bonusSalary
	userSalaryPayment.AdditionalWorkSalary = additionalWorkSalary
	userSalaryPayment.PaymentMethod = paymentMethod
	userSalaryPayment.PaymentProof = request.PaymentProof
	userSalaryPayment.IsPaid = true
	userSalaryPayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateUserSalaryPayment(&userSalaryPayment)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, err
	}

	if request.UserCashAdvancePayments != nil {
		for _, capReq := range request.UserCashAdvancePayments {
			data, err := s.repository.GetUserCashAdvance(capReq.UserCashAdvanceId)
			if err != nil {
				s.log.Error("failed get user cash advance", zap.Error(err))
				return dto.UserSalaryPaymentResponse{}, err
			}

			paymentDate, err := time.Parse("02-01-2006", capReq.PaymentDate)
			if err != nil {
				return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid payment date format")
			}

			paymentMethod := enum.ValueOfPaymentMethod(capReq.PaymentMethod)
			if !paymentMethod.IsValid() {
				return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid payment method")
			}

			nominal, err := decimal.NewFromString(capReq.Nominal)
			if err != nil {
				return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid nominal format")
			}

			currentPayment := nominal
			for _, payment := range data.Payments {
				currentPayment = currentPayment.Add(payment.Nominal)
			}

			if currentPayment.GreaterThan(data.Nominal) {
				return dto.UserSalaryPaymentResponse{}, errx.BadRequest("nominal more than needed")
			} else if currentPayment.Equal(data.Nominal) {
				data.PaymentStatus = enum.PaymentStatusPaid
			} else if currentPayment.LessThan(data.Nominal) {
				data.PaymentStatus = enum.PaymentStatusUnpaid
			}

			payment := entity.UserCashAdvancePayment{
				UserCashAdvanceId: capReq.UserCashAdvanceId,
				Nominal:           nominal,
				PaymentDate:       paymentDate,
				PaymentMethod:     paymentMethod,
				PaymentProof:      capReq.PaymentProof,
			}

			err = s.repository.CreateUserCashAdvancePayment(&payment)
			if err != nil {
				s.log.Error("failed create user cash advance payment", zap.Error(err))
				return dto.UserSalaryPaymentResponse{}, err
			}

			err = s.repository.UpdateUserCashAdvance(&data)
			if err != nil {
				s.log.Error("failed update user cash advance", zap.Error(err))
				return dto.UserSalaryPaymentResponse{}, err
			}
		}
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.UserSalaryPaymentResponse{}, err
	}

	return dto.UserSalaryPaymentResponse{
		Id:                   userSalaryPayment.Id,
		User:                 mapper.UserToListResponse(&userSalaryPayment.User),
		BaseSalary:           userSalaryPayment.BaseSalary.String(),
		BonusSalary:          userSalaryPayment.BonusSalary.String(),
		CompentationSalary:   userSalaryPayment.CompentationSalary.String(),
		AdditionalWorkSalary: userSalaryPayment.AdditionalWorkSalary.String(),
		PaymentProof:         userSalaryPayment.PaymentProof,
		PaymentMethod:        userSalaryPayment.PaymentMethod.String(),
		IsPaid:               userSalaryPayment.IsPaid,
	}, nil
}

func (s *CashflowService) GetDebtOverview(filter dto.GetDebtOverviewFilter) (dto.DebtOverviewResponse, error) {
	s.repository.UseTx(false)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	warehouseItemProcurements, err := s.repository.GetWarehouseItemProcurementCashflows(dto.GetWarehouseItemProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse item procurements", zap.Error(err))
		return dto.DebtOverviewResponse{}, err
	}

	warehouseItemCornProcurements, err := s.repository.GetWarehouseItemCornProcurementCashflows(dto.GetWarehouseItemCornProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse item corn procurements", zap.Error(err))
		return dto.DebtOverviewResponse{}, err
	}

	chickenProcurements, err := s.repository.GetChickenProcurementCashflows(dto.GetChickenProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get chicken procurements", zap.Error(err))
		return dto.DebtOverviewResponse{}, err
	}

	totalPayment := decimal.Zero
	totalDebtPayment := decimal.Zero
	totalRemainingDebtPayment := decimal.Zero

	for _, e := range warehouseItemProcurements {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentDebtPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentDebtPayment = totalCurrentDebtPayment.Add(p.Nominal)
		}
		totalDebtPayment = totalDebtPayment.Add(totalCurrentDebtPayment)
		totalRemainingDebtPayment = totalRemainingDebtPayment.Add(e.TotalPrice.Sub(totalCurrentDebtPayment))
	}

	for _, e := range warehouseItemCornProcurements {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentDebtPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentDebtPayment = totalCurrentDebtPayment.Add(p.Nominal)
		}
		totalDebtPayment = totalDebtPayment.Add(totalCurrentDebtPayment)
		totalRemainingDebtPayment = totalRemainingDebtPayment.Add(e.TotalPrice.Sub(totalCurrentDebtPayment))
	}

	for _, e := range chickenProcurements {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentDebtPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentDebtPayment = totalCurrentDebtPayment.Add(p.Nominal)
		}
		totalDebtPayment = totalDebtPayment.Add(totalCurrentDebtPayment)
		totalRemainingDebtPayment = totalRemainingDebtPayment.Add(e.TotalPrice.Sub(totalCurrentDebtPayment))
	}

	debtResponses := make([]dto.DebtListResponse, 0)
	if filter.DebtCategory == constant.DebtCategoryAll || filter.DebtCategory == constant.DebtCategoryWarehouseItemProcurement {
		for _, e := range warehouseItemProcurements {
			response := dto.DebtListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.DebtCategoryWarehouseItemProcurement,
				PlaceName:           e.Warehouse.Location.Name + " - " + e.Warehouse.Name,
				TransactionName:     constant.DebtTransactionNameWarehouseItemProcurement,
				Name:                e.Supplier.Name,
				PhoneNumber:         e.Supplier.PhoneNumber,
				Nominal:             e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			response.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			debtResponses = append(debtResponses, response)
		}
	}

	if filter.DebtCategory == constant.DebtCategoryAll || filter.DebtCategory == constant.DebtCategoryWarehouseItemCornProcurement {
		for _, e := range warehouseItemCornProcurements {
			response := dto.DebtListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.DebtCategoryWarehouseItemCornProcurement,
				PlaceName:           e.Warehouse.Location.Name + " - " + e.Warehouse.Name,
				TransactionName:     constant.DebtTransactionNameWarehouseItemCornProcurement,
				Name:                e.Supplier.Name,
				PhoneNumber:         e.Supplier.PhoneNumber,
				Nominal:             e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			response.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			debtResponses = append(debtResponses, response)
		}
	}

	if filter.DebtCategory == constant.DebtCategoryAll || filter.DebtCategory == constant.DebtCategoryChickenProcurement {
		for _, e := range chickenProcurements {
			response := dto.DebtListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.DebtCategoryChickenProcurement,
				PlaceName:           e.Cage.Location.Name + " - " + e.Cage.Name,
				TransactionName:     constant.DebtTransactionNameChickenProcurement,
				Name:                e.Supplier.Name,
				PhoneNumber:         e.Supplier.PhoneNumber,
				Nominal:             e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			response.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			debtResponses = append(debtResponses, response)
		}
	}

	return dto.DebtOverviewResponse{
		DebtPie: dto.DebtPieResponse{
			PaidPercentage:   totalDebtPayment.Sub(totalPayment).InexactFloat64() * 100.0,
			UnpaidPercentage: totalRemainingDebtPayment.Sub(totalPayment).InexactFloat64() * 100.0,
		},
		Debts: debtResponses,
	}, nil
}

func (s *CashflowService) GetDebt(debtCategory string, id uint64) (dto.DebtResponse, error) {
	switch debtCategory {
	case constant.DebtCategoryChickenProcurement:
		data, err := s.repository.GetChickenProcurementCashflow(id)
		if err != nil {
			s.log.Error("failed get chicken procurement cashflow", zap.Error(err))
			return dto.DebtResponse{}, err
		}

		paymentResponses := make([]dto.DebtPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.DebtPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.DebtResponse{
			Id:                  data.Id,
			Date:                data.CreatedAt.Format("02 Jan 2006"),
			Time:                data.CreatedAt.Format("15:04"),
			Category:            constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:           data.Cage.Location.Name + " - " + data.Cage.Name,
			Name:                data.Supplier.Name,
			PhoneNumber:         data.Supplier.PhoneNumber,
			RemainingPayment:    totalRemainingPayment.String(),
			PaymentType:         data.PaymentType.String(),
			PaymentStatus:       data.PaymentStatus.String(),
			DeadlinePaymentDate: data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:             data.CreatedByUser.Name,
			DebtPayments:        paymentResponses,
		}, nil
	case constant.DebtCategoryWarehouseItemProcurement:
		data, err := s.repository.GetWarehouseItemProcurementCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse item procurement cashflow", zap.Error(err))
			return dto.DebtResponse{}, err
		}

		paymentResponses := make([]dto.DebtPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.DebtPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.DebtResponse{
			Id:                  data.Id,
			Date:                data.CreatedAt.Format("02 Jan 2006"),
			Time:                data.CreatedAt.Format("15:04"),
			Category:            constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:           data.Warehouse.Location.Name + " - " + data.Warehouse.Name,
			Name:                data.Supplier.Name,
			PhoneNumber:         data.Supplier.PhoneNumber,
			RemainingPayment:    totalRemainingPayment.String(),
			PaymentType:         data.PaymentType.String(),
			PaymentStatus:       data.PaymentStatus.String(),
			DeadlinePaymentDate: data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:             data.CreatedByUser.Name,
			DebtPayments:        paymentResponses,
		}, nil

	case constant.DebtCategoryWarehouseItemCornProcurement:
		data, err := s.repository.GetWarehouseItemCornProcurementCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse item corn procurement cashflow", zap.Error(err))
			return dto.DebtResponse{}, err
		}

		paymentResponses := make([]dto.DebtPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.DebtPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			paymentResponse.Remaining = totalRemainingPayment.Sub(e.Nominal).String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.DebtResponse{
			Id:                  data.Id,
			Date:                data.CreatedAt.Format("02 Jan 2006"),
			Time:                data.CreatedAt.Format("15:04"),
			Category:            constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:           data.Warehouse.Location.Name + " - " + data.Warehouse.Name,
			Name:                data.Supplier.Name,
			PhoneNumber:         data.Supplier.PhoneNumber,
			RemainingPayment:    totalRemainingPayment.String(),
			PaymentType:         data.PaymentType.String(),
			PaymentStatus:       data.PaymentStatus.String(),
			DeadlinePaymentDate: data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:             data.CreatedByUser.Name,
			DebtPayments:        paymentResponses,
		}, nil

	default:
		return dto.DebtResponse{}, errx.BadRequest("invalid debt transaction name")
	}
}

func (s *CashflowService) ExportSalesCashflowToExcel(filter dto.GetSaleCashflowFilter) (*excelize.File, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	f := excelize.NewFile()

	storeResp, err := s.storeService.GetStoreSales(dto.GetStoreSaleFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return nil, err
	}

	storeSheet := "Store Sales"
	f.NewSheet(storeSheet)
	headers := []string{
		"ID", "Customer", "Item", "Store", "Quantity", "Sale Unit",
		"Total Price", "Payment Status", "Is Send", "Deadline Payment Date", "Send Date",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(storeSheet, cell, h)
	}

	for row, sale := range storeResp.StoreSales {
		values := []interface{}{
			sale.Id,
			sale.Customer.Name,
			sale.Item.Name,
			sale.Store.Name,
			sale.Quantity,
			sale.SaleUnit,
			sale.TotalPrice.String(),
			sale.PaymentStatus,
			sale.IsSend,
			sale.DeadlinePaymentDate,
			sale.SendDate,
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
			f.SetCellValue(storeSheet, cell, v)
		}
	}

	warehouseResp, err := s.warehouseService.GetWarehouseSales(dto.GetWarehouseSaleFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return nil, err
	}

	warehouseSheet := "Warehouse Sales"
	f.NewSheet(warehouseSheet)
	headers = []string{
		"ID", "Customer", "Item", "Warehouse", "Quantity", "Sale Unit",
		"Total Price", "Payment Status", "Is Send", "Deadline Payment Date", "Send Date",
	}
	for i, h := range headers {
		cell, _ := excelize.CoordinatesToCellName(i+1, 1)
		f.SetCellValue(warehouseSheet, cell, h)
	}

	for row, sale := range warehouseResp.WarehouseSales {
		values := []interface{}{
			sale.Id,
			sale.Customer.Name,
			sale.Item.Name,
			sale.Warehouse.Name,
			sale.Quantity,
			sale.SaleUnit,
			sale.TotalPrice.String(),
			sale.PaymentStatus,
			sale.IsSend,
			sale.DeadlinePaymentDate,
			sale.SendDate,
		}
		for col, v := range values {
			cell, _ := excelize.CoordinatesToCellName(col+1, row+2)
			f.SetCellValue(warehouseSheet, cell, v)
		}
	}

	f.DeleteSheet("Sheet1")

	return f, nil
}

func (s *CashflowService) GetUserSalarySummary(filter dto.GetUserSalarySummaryFilter) (dto.UserSalarySummaryResponse, error) {
	s.repository.UseTx(false)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	userSalaryPayments, err := s.repository.GetUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		LocationId: filter.LocationId,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
	})
	if err != nil {
		return dto.UserSalarySummaryResponse{}, err
	}

	totalStaff := len(userSalaryPayments)
	totalBaseSalary := decimal.Zero
	totalAdditonalWorkSalary := decimal.Zero
	totalBonusSalary := decimal.Zero

	for _, e := range userSalaryPayments {
		totalBaseSalary = totalBaseSalary.Add(totalBaseSalary)

		if e.IsPaid {
			totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(e.AdditionalWorkSalary)
			totalBonusSalary = totalBonusSalary.Add(e.BonusSalary)
		} else {

			withDeleted := true
			additionalWorkSalary := decimal.Zero
			additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(e.UserId,
				dto.GetAdditionalWorkUserFilter{
					Month:       param.MonthParam(filter.Month.Value()),
					Year:        filter.Year,
					WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
				})
			if err != nil {
				return dto.UserSalarySummaryResponse{}, nil
			}

			for _, e := range additionalWorkUsers.AdditionalWorkUsers {
				salary, err := decimal.NewFromString(e.AdditionalWork.Salary)
				if err != nil {
					s.log.Error("failed parse additional work salary", zap.Error(err))
					return dto.UserSalarySummaryResponse{}, err
				}
				additionalWorkSalary = additionalWorkSalary.Add(salary)
			}

			totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(additionalWorkSalary)

			kpiScore, err := s.userService.CalculateKPIScoreUserPerMonth(e.UserId, filter.Year, filter.Month.Value())
			if err != nil {
				return dto.UserSalarySummaryResponse{}, err
			}

			bonusSalary := decimal.Zero
			if kpiScore*0.6 == 60 {
				bonusSalary = bonusSalary.Add(decimal.NewFromFloat(50000))
			}

			diff := kpiScore - 90.0
			if diff > 0 {
				percentage := float64(diff) / 2
				bonusSalary = bonusSalary.Add(decimal.NewFromFloat(percentage).Mul(e.BaseSalary))

			} else if diff < 0 {
				percentage := float64(-diff) / 2
				bonusSalary = bonusSalary.Sub(decimal.NewFromFloat(percentage).Mul(e.BaseSalary))
			}

			totalBonusSalary = totalBonusSalary.Add(bonusSalary)
		}
	}

	return dto.UserSalarySummaryResponse{
		TotalStaff:               uint64(totalStaff),
		TotalBaseSalary:          totalBaseSalary.String(),
		TotalAdditonalWorkSalary: totalAdditonalWorkSalary.String(),
		TotalBonusSalary:         totalBonusSalary.String(),
	}, nil
}

func (s *CashflowService) GetUserSalaries(filter dto.GetUserSalaryListFilter) (dto.UserSalaryListPaginationResponse, error) {
	s.repository.UseTx(false)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	userSalaryPayments, err := s.repository.GetUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		LocationId: filter.LocationId,
		RoleId:     filter.RoleId,
		Page:       filter.Page,
		Keyword:    filter.Keyword,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
	})
	if err != nil {
		return dto.UserSalaryListPaginationResponse{}, err
	}

	totalData, err := s.repository.CountUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		LocationId: filter.LocationId,
		RoleId:     filter.RoleId,
		Page:       filter.Page,
		Keyword:    filter.Keyword,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed count user salary payments", zap.Error(err))
		return dto.UserSalaryListPaginationResponse{}, err
	}

	responses := make([]dto.UserSalaryListResponse, 0)
	for _, e := range userSalaryPayments {
		responses = append(responses, dto.UserSalaryListResponse{
			Id:             e.Id,
			User:           mapper.UserToListResponse(&e.User),
			SalaryInterval: e.User.SalaryInterval.String(),
			IsPaid:         e.IsPaid,
		})
	}

	response := dto.UserSalaryListPaginationResponse{
		UserSalaries: responses,
	}

	if filter.Page > 0 {
		response.TotalData = uint64(totalData)
		response.TotalPage = uint64(totalData) / constant.PaginationDefaultLimit
	}

	return response, nil
}

func (s *CashflowService) GetUserSalaryDetail(id uint64) (dto.UserSalaryDetailResponse, error) {
	s.repository.UseTx(false)

	userSalaryPayment, err := s.repository.GetUserSalaryPayment(id)
	if err != nil {
		s.log.Error("failed get user salary payment")
		return dto.UserSalaryDetailResponse{}, err
	}

	userCashAdvanceSummary, err := s.GetUserCashAdvanceByUserId(userSalaryPayment.UserId)
	if err != nil {
		return dto.UserSalaryDetailResponse{}, err
	}

	additionalWorkUserResponses := make([]dto.AdditionalWorkUserResponse, 0)

	totalAdditonalWorkSalary := decimal.Zero
	totalBonusSalary := decimal.Zero
	if userSalaryPayment.IsPaid {
		totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(userSalaryPayment.AdditionalWorkSalary)
		totalBonusSalary = totalBonusSalary.Add(userSalaryPayment.BonusSalary)
	} else {

		withDeleted := true
		additionalWorkSalary := decimal.Zero
		additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(userSalaryPayment.UserId,
			dto.GetAdditionalWorkUserFilter{
				Month:       param.MonthParam(enum.ValueOfMonth(userSalaryPayment.CreatedAt.Format("Januari"))),
				Year:        uint64(userSalaryPayment.CreatedAt.Year()),
				WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
			})
		if err != nil {
			return dto.UserSalaryDetailResponse{}, nil
		}

		for _, e := range additionalWorkUsers.AdditionalWorkUsers {
			if e.IsDone {
				additionalWorkUserResponses = append(additionalWorkUserResponses, e)
				salary, err := decimal.NewFromString(e.AdditionalWork.Salary)
				if err != nil {
					s.log.Error("failed parse additional work salary", zap.Error(err))
					return dto.UserSalaryDetailResponse{}, err
				}
				additionalWorkSalary = additionalWorkSalary.Add(salary)
			}
		}

		totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(additionalWorkSalary)

		kpiScore, err := s.userService.CalculateKPIScoreUserPerMonth(userSalaryPayment.UserId, uint64(userSalaryPayment.CreatedAt.Year()), enum.ValueOfMonth(userSalaryPayment.CreatedAt.Format("Januari")))
		if err != nil {
			return dto.UserSalaryDetailResponse{}, err
		}

		bonusSalary := decimal.Zero
		if kpiScore*0.6 == 60 {
			bonusSalary = bonusSalary.Add(decimal.NewFromFloat(50000))
		}

		diff := kpiScore - 90.0
		if diff > 0 {
			percentage := float64(diff) / 2
			bonusSalary = bonusSalary.Add(decimal.NewFromFloat(percentage).Mul(userSalaryPayment.BaseSalary))

		} else if diff < 0 {
			percentage := float64(-diff) / 2
			bonusSalary = bonusSalary.Sub(decimal.NewFromFloat(percentage).Mul(userSalaryPayment.BaseSalary))
		}

		totalBonusSalary = totalBonusSalary.Add(bonusSalary)
	}

	return dto.UserSalaryDetailResponse{
		AdditionalWorkUsers:      additionalWorkUserResponses,
		UserCashAdvanceSummaries: userCashAdvanceSummary,
		BaseSalary:               userSalaryPayment.BaseSalary.String(),
		CompentationSalary:       userSalaryPayment.CompentationSalary.String(),
		BonusSalary:              totalBonusSalary.String(),
		AdditionalWorkSalary:     totalAdditonalWorkSalary.String(),
	}, nil
}
