package service

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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
}

type ICashflowService interface {
	GetIncomeOverview(filter dto.GetIncomeOverviewFilter) (dto.IncomeOverviewResponse, error)
	GetIncome(incomeCategory string, id uint64) (dto.IncomeResponse, error)

	CreateExpense(request dto.CreateExpenseRequest, userId uuid.UUID) (dto.ExpenseResponse, error)
	GetExpenseOverview(filter dto.GetExpenseOverviewFilter) (dto.ExpenseOverviewResponse, error)
	GetExpense(expenseCategory string, id uint64) (dto.ExpenseResponse, error)

	ExportSalesCashflowToExcel(filter dto.GetSaleCashflowFilter) (*excelize.File, error)
}

func NewCashflowService(log *zap.Logger, repository repository.ICashflowRepository, storeService IStoreService, warehouseService IWarehouseService, chickenService IChickenService, userService IUserService) ICashflowService {
	return &CashflowService{
		log:              log,
		repository:       repository,
		storeService:     storeService,
		warehouseService: warehouseService,
		chickenService:   chickenService,
		userService:      userService,
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
	totalWarehouseProcurement := decimal.Zero
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
		totalWarehouseProcurement = totalWarehouseProcurement.Add(p.Nominal)
	}

	for _, p := range warehouseItemCornProcurementPayments {
		totalPayment = totalPayment.Add(p.Nominal)
		totalWarehouseProcurement = totalWarehouseProcurement.Add(p.Nominal)
	}

	for _, p := range userSalaryPayments {
		totalSalary := p.BaseSalary.Add(p.BonusSalary).Add(p.CompentationSalary).Add(p.AdditionalWorkSalary)
		totalPayment = totalPayment.Add(totalSalary)
		totalWarehouseProcurement = totalWarehouseProcurement.Add(totalSalary)
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryChickenProcurement {
		for _, p := range chickenProcurementPayments {
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.PaymentDate.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryChickenProcurement,
				Name:         constant.ExpenseNameChickenProcurement,
				PlaceName:    p.ChickenProcurement.Cage.Location.Name + " - " + p.ChickenProcurement.Cage.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.ChickenProcurement.Supplier.Name,
				PaymentProof: p.PaymentProof,
			})
		}
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryWarehouseProcurement {
		for _, p := range warehouseItemProcurementPayments {
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.PaymentDate.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryWarehouseProcurement,
				Name:         constant.ExpenseNameWarehouseItemProcurement,
				PlaceName:    p.WarehouseItemProcurement.Warehouse.Location.Name + " - " + p.WarehouseItemProcurement.Warehouse.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.WarehouseItemProcurement.Supplier.Name,
				PaymentProof: p.PaymentProof,
			})
		}
		for _, p := range warehouseItemCornProcurementPayments {
			totalPayment = totalPayment.Add(p.Nominal)
			totalWarehouseProcurement = totalWarehouseProcurement.Add(p.Nominal)
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.PaymentDate.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryWarehouseProcurement,
				Name:         constant.ExpenseNameWarehouseItemCornProcurement,
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
				Name:         constant.ExpenseNameSalary,
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
			StaffPercentage:                totalStaffSalary.Div(totalPayment).InexactFloat64() * 100.0,
			ChikckenProcuremtnPercentage:   totalChickenProcurement.Div(totalPayment).InexactFloat64() * 100.0,
			WarehouseProcurementPercentage: totalWarehouseProcurement.Div(totalPayment).InexactFloat64() * 100.0,
			OperationalPercentage:          totalOperational.Div(totalPayment).InexactFloat64() * 100.0,
			OtherPercentage:                totalOther.Div(totalPayment).InexactFloat64() * 100.0,
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
			Category:            "Operational",
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

	case constant.ExpenseCategoryWarehouseProcurement:
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
