package service

import (
	"database/sql"
	"fmt"
	"math"
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
	log             *zap.Logger
	repository      repository.ICashflowRepository
	workService     IWorkService
	itemService     IItemService
	presenceService IPresenceService
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

	GetReceiveablesOverview(filter dto.GetReceivablesOverviewFilter) (dto.ReceivablesOverviewResponse, error)
	GetReceiveables(receieveablesCategory string, id uint64) (dto.ReceivablesResponse, error)

	PayUserSalaryPayment(id uint64, request dto.PayUserSalaryPaymentRequest, userId uuid.UUID) (dto.UserSalaryPaymentResponse, error)

	GetDebtOverview(filter dto.GetDebtOverviewFilter) (dto.DebtOverviewResponse, error)
	GetDebt(debtCategory string, id uint64) (dto.DebtResponse, error)

	GetUserSalarySummary(filter dto.GetUserSalarySummaryFilter) (dto.UserSalarySummaryResponse, error)
	GetUserSalaries(filter dto.GetUserSalaryListFilter) (dto.UserSalaryListPaginationResponse, error)
	GetUserSalaryDetail(id uint64) (dto.UserSalaryDetailResponse, error)

	ExportSalesCashflowToExcel(filter dto.GetCashflowSaleReportFilter) (*excelize.File, error)

	GetCashflowSaleOverview(filter dto.GetCashflowSaleOverviewFilter) (dto.CashflowSaleOverviewResponse, error)
	GetCashflowOverview(filter dto.GetCashflowOverviewFilter) (dto.CashflowOverviewResponse, error)

	GetTotalIncomeProductionInMonth(month enum.Month, year uint64) (decimal.Decimal, error)
	GetTotalExpenseProductionInMonth(month enum.Month, year uint64) (decimal.Decimal, error)

	GetCashflowHistories(filter dto.GetCashflowHistoryFilter) ([]dto.CashflowHistoryResponse, error)
}

func NewCashflowService(log *zap.Logger, repository repository.ICashflowRepository, workService IWorkService, itemService IItemService, presenceService IPresenceService) ICashflowService {
	return &CashflowService{
		log:             log,
		repository:      repository,
		workService:     workService,
		presenceService: presenceService,
		itemService:     itemService,
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

	userCashAdvancePayments, err := s.repository.GetUserCashAdvancePayments(dto.GetUserCashAdvancePaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get user cash advance payments", zap.Error(err))
		return dto.IncomeOverviewResponse{}, err
	}

	totalPayment := decimal.Zero
	totalWarehouseSalePayment := decimal.Zero
	totalStoreSalePayment := decimal.Zero
	totalAfkirChickenSalePayment := decimal.Zero
	totalUserCashAdvance := decimal.Zero

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

	for _, payment := range userCashAdvancePayments {
		totalPayment = totalPayment.Add(payment.Nominal)
		totalUserCashAdvance = totalUserCashAdvance.Add(payment.Nominal)
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
				Quantity:     payment.StoreSale.Quantity,
				CustomerName: payment.StoreSale.Customer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	}

	if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryWarehouseEggSale {
		for _, payment := range warehouseSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.WarehouseSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("02 Jan 2006"),
				PlaceName:    payment.WarehouseSale.Warehouse.Location.Name + " - " + payment.WarehouseSale.Warehouse.Name,
				Category:     constant.IncomeCategoryWarehouseEggSale,
				ItemName:     payment.WarehouseSale.Item.Name,
				ItemUnit:     payment.WarehouseSale.SaleUnit.String(),
				Quantity:     payment.WarehouseSale.Quantity,
				CustomerName: payment.WarehouseSale.Customer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	}

	if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryAfkirChickenSale {
		for _, payment := range afkirChickenSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.AfkirChickenSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("02 Jan 2006"),
				PlaceName:    payment.AfkirChickenSale.ChickenCage.Cage.Location.Name + " - " + payment.AfkirChickenSale.ChickenCage.Cage.Name,
				Category:     constant.IncomeCategoryAfkirChickenSale,
				ItemName:     constant.AfkirChicken,
				ItemUnit:     constant.AfkirChickenUnitEkor,
				Quantity:     float64(payment.AfkirChickenSale.TotalSellChicken),
				CustomerName: payment.AfkirChickenSale.AfkirChickenCustomer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	}

	if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryUserCashAdvancePayment {
		for _, payment := range afkirChickenSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.AfkirChickenSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("02 Jan 2006"),
				PlaceName:    payment.AfkirChickenSale.ChickenCage.Cage.Location.Name + " - " + payment.AfkirChickenSale.ChickenCage.Cage.Name,
				Category:     constant.IncomeCategoryAfkirChickenSale,
				ItemName:     constant.AfkirChicken,
				ItemUnit:     constant.AfkirChickenUnitEkor,
				Quantity:     float64(payment.AfkirChickenSale.TotalSellChicken),
				CustomerName: payment.AfkirChickenSale.AfkirChickenCustomer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	}

	warehouseEggSalePercentage := 0.0
	storeEggSalePercentage := 0.0
	afkirChickenSalePercentage := 0.0
	userCashAdvanceSalePercentage := 0.0
	if !totalPayment.IsZero() {
		warehouseEggSalePercentage = totalWarehouseSalePayment.Div(totalPayment).InexactFloat64() * 100.0
		storeEggSalePercentage = totalStoreSalePayment.Div(totalPayment).InexactFloat64() * 100.0
		afkirChickenSalePercentage = totalAfkirChickenSalePayment.Div(totalPayment).InexactFloat64() * 100.0
		userCashAdvanceSalePercentage = totalUserCashAdvance.Div(totalPayment).InexactFloat64() * 100.0
	}

	return dto.IncomeOverviewResponse{
		IncomePie: dto.IncomePieResponse{
			WarehouseEggSalePercentage: warehouseEggSalePercentage,
			StoreEggSalePercentage:     storeEggSalePercentage,
			AfkirChickenSalePercentage: afkirChickenSalePercentage,
			UserCashAdvancePercentage:  userCashAdvanceSalePercentage,
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
			ParentId:            payment.WarehouseSaleId,
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            constant.IncomeCategoryWarehouseEggSale,
			PlaceName:           payment.WarehouseSale.Warehouse.Name,
			CustomerName:        payment.WarehouseSale.Customer.Name,
			CustomerPhoneNumber: payment.WarehouseSale.Customer.PhoneNumber,
			ItemName:            payment.WarehouseSale.Item.Name,
			ItemUnit:            payment.WarehouseSale.SaleUnit.String(),
			Quantity:            payment.WarehouseSale.Quantity,
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
			ParentId:            payment.StoreSaleId,
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            constant.IncomeCategoryStoreEggSale,
			PlaceName:           payment.StoreSale.Store.Name,
			CustomerName:        payment.StoreSale.Customer.Name,
			CustomerPhoneNumber: payment.StoreSale.Customer.PhoneNumber,
			ItemName:            payment.StoreSale.Item.Name,
			ItemUnit:            payment.StoreSale.SaleUnit.String(),
			Quantity:            payment.StoreSale.Quantity,
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
			ParentId:            payment.AfkirChickenSaleId,
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            constant.IncomeCategoryAfkirChickenSale,
			PlaceName:           payment.AfkirChickenSale.ChickenCage.Cage.Name,
			CustomerName:        payment.AfkirChickenSale.AfkirChickenCustomer.Name,
			CustomerPhoneNumber: payment.AfkirChickenSale.AfkirChickenCustomer.PhoneNumber,
			ItemName:            constant.AfkirChicken,
			ItemUnit:            constant.AfkirChickenUnitEkor,
			Quantity:            float64(payment.AfkirChickenSale.TotalSellChicken),
			Nominal:             payment.Nominal.String(),
			PaymentType:         payment.AfkirChickenSale.PaymentType.String(),
			TotalPrice:          payment.AfkirChickenSale.TotalPrice.String(),
			PaymentMethod:       payment.PaymentMethod.String(),
			InputBy:             payment.CreatedByUser.Name,
			PaymentProof:        payment.PaymentProof,
		}, nil

	case constant.IncomeCategoryUserCashAdvancePayment:
		payment, err := s.repository.GetUserCashAdvancePayment(id)
		if err != nil {
			s.log.Error("failed get afkir chicken sale payment", zap.Error(err))
			return dto.IncomeResponse{}, err
		}

		return dto.IncomeResponse{
			ParentId:            payment.UserCashAdvance.Id,
			Id:                  payment.Id,
			Date:                payment.PaymentDate.Format("02 Jan 2006"),
			Time:                payment.PaymentDate.Format("15:04"),
			Category:            constant.IncomeCategoryAfkirChickenSale,
			PlaceName:           payment.UserCashAdvance.User.Location.Name,
			CustomerName:        payment.UserCashAdvance.User.Name,
			CustomerPhoneNumber: payment.UserCashAdvance.User.PhoneNumber,
			ItemName:            constant.AfkirChicken,
			ItemUnit:            constant.AfkirChickenUnitEkor,
			Quantity:            1,
			Nominal:             payment.Nominal.String(),
			PaymentType:         enum.PaymentTypeinstallment.String(),
			TotalPrice:          payment.UserCashAdvance.Nominal.String(),
			PaymentMethod:       payment.PaymentMethod.String(),
			InputBy:             payment.CreatedByUser.Name,
			PaymentProof:        payment.PaymentProof,
		}, nil
	default:
		return dto.IncomeResponse{}, errx.BadRequest("invalid income category")
	}
}

func (s *CashflowService) GetTotalExpenseProductionInMonth(month enum.Month, year uint64) (decimal.Decimal, error) {
	s.repository.UseTx(false)

	totalExpenseProduction := decimal.Zero
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))

	warehouseItemProcurements, err := s.repository.GetWarehouseItemProcurementCashflows(dto.GetWarehouseItemProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse item procurements", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range warehouseItemProcurements {
		totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
	}

	warehouseItemCornProcurements, err := s.repository.GetWarehouseItemCornProcurementCashflows(dto.GetWarehouseItemCornProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse item corn procurements", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range warehouseItemCornProcurements {
		totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
	}

	chickenProcurements, err := s.repository.GetChickenProcurementCashflows(dto.GetChickenProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get chicken procurements", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range chickenProcurements {
		totalExpenseProduction = totalExpenseProduction.Add(e.TotalPrice)
	}

	expense, err := s.repository.GetExpenses(dto.GetExpenseFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get expense", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range expense {
		totalExpenseProduction = totalExpenseProduction.Add(e.Nominal)
	}

	userSalaryPayments, err := s.repository.GetUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get user salary payments", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range userSalaryPayments {
		totalExpenseProduction = totalExpenseProduction.Add(e.BaseSalary.Add(e.AdditionalWorkSalary).Add(e.BonusSalary).Add(e.CompentationSalary).Sub(e.Cashbond))
	}

	return totalExpenseProduction, nil
}

func (s *CashflowService) GetTotalIncomeProductionInMonth(month enum.Month, year uint64) (decimal.Decimal, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))

	totalIncomeProduction := decimal.Zero

	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get store sale cashflows", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range storeSales {
		totalIncomeProduction = totalIncomeProduction.Add(e.TotalPrice)
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse sale cashflows", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range warehouseSales {
		totalIncomeProduction = totalIncomeProduction.Add(e.TotalPrice)
	}

	afkirChickenSales, err := s.repository.GetAfkirChickenSaleCashflows(dto.GetAfkirChickenSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get afkir chicken sale cashflows", zap.Error(err))
		return decimal.Zero, err
	}
	for _, e := range afkirChickenSales {
		totalIncomeProduction = totalIncomeProduction.Add(e.TotalPrice)
	}

	return totalIncomeProduction, nil
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

	locationType := enum.ValueOfLocationType(request.LocationType)
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

	userCashAdvances, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		PaymentStatus: param.PaymentStatusParam(enum.PaymentStatusNotPaid),
		StartDate:     param.DateParam(startDate),
		EndDate:       param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get user cash advances", zap.Error(err))
		return dto.ExpenseOverviewResponse{}, err
	}

	totalPayment := decimal.Zero
	totalChickenProcurement := decimal.Zero
	totalWarehouseItemProcurement := decimal.Zero
	totalWarehouseItemCornProcurement := decimal.Zero
	totalOperational := decimal.Zero
	totalOther := decimal.Zero
	totalTax := decimal.Zero
	totalUserSalary := decimal.Zero
	totalUserCashAdvance := decimal.Zero

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
		case enum.ExpenseCategoryTax:
			totalTax = totalTax.Add(p.Nominal)
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

	for _, p := range userCashAdvances {
		totalPayment = totalPayment.Add(p.Nominal)
		totalUserCashAdvance = totalUserCashAdvance.Add(p.Nominal)
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
			if p.ExpenseCategory == enum.ExpenseCategoryOperational {
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
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryOther {
		for _, p := range expensePayments {
			if p.ExpenseCategory == enum.ExpenseCategoryOther {
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
	}

	if filter.ExpenseCategory == constant.ExpenseCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryTax {
		for _, p := range expensePayments {
			if p.ExpenseCategory == enum.ExpenseCategoryTax {
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
	}

	if filter.ExpenseCategory == constant.IncomeCategoryAll || filter.ExpenseCategory == constant.ExpenseCategoryUserCashAdvance {
		for _, p := range userCashAdvances {
			expenseResponses = append(expenseResponses, dto.ExpenseListResponse{
				Id:           p.Id,
				Date:         p.CreatedAt.Format("02 Jan 2006"),
				Category:     constant.ExpenseCategoryUserCashAdvance,
				Name:         constant.ExpenseTransactionNameSalary,
				PlaceName:    p.User.Location.Name,
				Nominal:      p.Nominal.String(),
				ReceiverName: p.User.Name,
				PaymentProof: "-",
			})
		}
	}

	staffPercentage := 0.0
	chickenProcurementPercentage := 0.0
	warehouseItemProcurementPercentage := 0.0
	warehouseItemCornProcurementPercentage := 0.0
	operationalPercentage := 0.0
	otherPercentage := 0.0
	userCashAdvancePercentage := 0.0
	taxPercentage := 0.0

	if !totalPayment.IsZero() {
		staffPercentage = totalUserSalary.Div(totalPayment).InexactFloat64() * 100.0
		chickenProcurementPercentage = totalChickenProcurement.Div(totalPayment).InexactFloat64() * 100.0
		warehouseItemProcurementPercentage = totalWarehouseItemProcurement.Div(totalPayment).InexactFloat64() * 100.0
		warehouseItemCornProcurementPercentage = totalWarehouseItemCornProcurement.Div(totalPayment).InexactFloat64() * 100.0
		operationalPercentage = totalOperational.Div(totalPayment).InexactFloat64() * 100.0
		otherPercentage = totalOther.Div(totalPayment).InexactFloat64() * 100.0
		userCashAdvancePercentage = totalUserCashAdvance.Div(totalPayment).InexactFloat64() * 100.0
		taxPercentage = totalTax.Div(totalPayment).InexactFloat64() * 100.0
	}

	return dto.ExpenseOverviewResponse{
		ExpensePie: dto.ExpensePieResponse{
			StaffPercentage:                        staffPercentage,
			ChikckenProcuremtnPercentage:           chickenProcurementPercentage,
			WarehouseItemProcurementPercentage:     warehouseItemProcurementPercentage,
			WarehouseItemCornProcurementPercentage: warehouseItemCornProcurementPercentage,
			OperationalPercentage:                  operationalPercentage,
			OtherPercentage:                        otherPercentage,
			UserCashAdvancePercentage:              userCashAdvancePercentage,
			TaxPercentage:                          taxPercentage,
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

		response := dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.CreatedAt.Format("2006-01-02"),
			Time:                expense.CreatedAt.Format("15:04:05"),
			Category:            constant.ExpenseCategoryOperational,
			Name:                expense.Name,
			ReceiverName:        expense.ReceiverName,
			ReceiverPhoneNumber: expense.ReceiverPhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}

		switch expense.LocationType {
		case enum.LocationTypeCage:
			response.PlaceName = expense.Cage.Name + " - " + expense.Location.Name
		case enum.LocationTypeStore:
			response.PlaceName = expense.Store.Name + " - " + expense.Location.Name
		case enum.LocationTypeWarehouse:
			response.PlaceName = expense.Warehouse.Name + " - " + expense.Location.Name
		}

		return response, nil

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
			Category:            constant.ExpenseCategoryChickenProcurement,
			PlaceName:           expense.ChickenProcurement.Cage.Location.Name,
			Name:                constant.ExpenseCategoryChickenProcurement,
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
			Category:            constant.ExpenseCategoryWarehouseItemProcurement,
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
			Category:            constant.ExpenseCategoryWarehouseItemCornProcurement,
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
			Category:            constant.ExpenseCategoryStaff,
			PlaceName:           "-",
			Name:                expense.User.Name,
			ReceiverName:        expense.User.Name,
			ReceiverPhoneNumber: expense.User.PhoneNumber,
			Nominal:             expense.BaseSalary.Add(expense.BonusSalary).Add(expense.CompentationSalary).Add(expense.AdditionalWorkSalary).String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}, nil

	case constant.ExpenseCategoryTax:
		expense, err := s.repository.GetExpense(id)
		if err != nil {
			s.log.Error("failed get other expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		response := dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.CreatedAt.Format("2006-01-02"),
			Time:                expense.CreatedAt.Format("15:04:05"),
			Category:            constant.ExpenseCategoryOther,
			PlaceName:           expense.Location.Name,
			Name:                expense.Name,
			ReceiverName:        expense.ReceiverName,
			ReceiverPhoneNumber: expense.ReceiverPhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}

		switch expense.LocationType {
		case enum.LocationTypeCage:
			response.PlaceName = expense.Cage.Name + " - " + expense.Location.Name
		case enum.LocationTypeStore:
			response.PlaceName = expense.Store.Name + " - " + expense.Location.Name
		case enum.LocationTypeWarehouse:
			response.PlaceName = expense.Warehouse.Name + " - " + expense.Location.Name
		}

		return response, nil

	case constant.ExpenseCategoryOther:
		expense, err := s.repository.GetExpense(id)
		if err != nil {
			s.log.Error("failed get other expense", zap.Error(err))
			return dto.ExpenseResponse{}, err
		}

		response := dto.ExpenseResponse{
			Id:                  expense.Id,
			Date:                expense.CreatedAt.Format("2006-01-02"),
			Time:                expense.CreatedAt.Format("15:04:05"),
			Category:            constant.ExpenseCategoryOther,
			PlaceName:           expense.Location.Name,
			Name:                expense.Name,
			ReceiverName:        expense.ReceiverName,
			ReceiverPhoneNumber: expense.ReceiverPhoneNumber,
			Nominal:             expense.Nominal.String(),
			PaymentMethod:       expense.PaymentMethod.String(),
			PaymentProof:        expense.PaymentProof,
			InputBy:             expense.CreatedByUser.Name,
		}

		switch expense.LocationType {
		case enum.LocationTypeCage:
			response.PlaceName = expense.Cage.Name + " - " + expense.Location.Name
		case enum.LocationTypeStore:
			response.PlaceName = expense.Store.Name + " - " + expense.Location.Name
		case enum.LocationTypeWarehouse:
			response.PlaceName = expense.Warehouse.Name + " - " + expense.Location.Name
		}

		return response, nil

	default:
		return dto.ExpenseResponse{}, errx.BadRequest("invalid expense category")
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
		PaymentStatus:           data.PaymentStatus.String(),
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
		data.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
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
		PaymentStatus:           data.PaymentStatus.String(),
		UserCashAdvancePayments: userCashAdvancePayments,
		RemainingPayment:        remainingPayment.String(),
	}

	return response, nil
}

func (s *CashflowService) GetReceiveablesOverview(filter dto.GetReceivablesOverviewFilter) (dto.ReceivablesOverviewResponse, error) {
	s.repository.UseTx(false)

	receieveables := make([]dto.ReceivablesListResponse, 0)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get store sale cashflows", zap.Error(err))
		return dto.ReceivablesOverviewResponse{}, err
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get warehouse sale cashflows", zap.Error(err))
		return dto.ReceivablesOverviewResponse{}, err
	}

	afkirChickenSales, err := s.repository.GetAfkirChickenSaleCashflows(dto.GetAfkirChickenSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get afkir chicken sale cashflows", zap.Error(err))
		return dto.ReceivablesOverviewResponse{}, err
	}

	userCashAdvances, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed get user cash advances", zap.Error(err))
		return dto.ReceivablesOverviewResponse{}, err
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

	if filter.ReceivablesCategory == constant.ReceieveablesCategoryAll || filter.ReceivablesCategory == constant.ReceieveablesCategoryWarehouseEggSale {
		for _, e := range warehouseSales {
			receieveable := dto.ReceivablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryWarehouseEggSale,
				PlaceName:           e.Warehouse.Location.Name + " - " + e.Warehouse.Name,
				Name:                e.Customer.Name,
				PhoneNumber:         e.Customer.PhoneNumber,
				TotalNominal:        e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			if e.PaidDate.Valid {
				receieveable.PaidDate = e.PaidDate.Time.Format("02 Jan 2006")
			} else {
				receieveable.PaidDate = "-"
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	if filter.ReceivablesCategory == constant.ReceieveablesCategoryAll || filter.ReceivablesCategory == constant.ReceieveablesCategoryStoreEggSale {
		for _, e := range storeSales {
			receieveable := dto.ReceivablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryStoreEggSale,
				PlaceName:           e.Store.Location.Name + " - " + e.Store.Name,
				Name:                e.Customer.Name,
				PhoneNumber:         e.Customer.PhoneNumber,
				TotalNominal:        e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			if e.PaidDate.Valid {
				receieveable.PaidDate = e.PaidDate.Time.Format("02 Jan 2006")
			} else {
				receieveable.PaidDate = "-"
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	if filter.ReceivablesCategory == constant.ReceieveablesCategoryAll || filter.ReceivablesCategory == constant.ReceieveablesCategoryAfkirChickenSale {
		for _, e := range afkirChickenSales {
			receieveable := dto.ReceivablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryAfkirChickenSale,
				PlaceName:           e.ChickenCage.Cage.Location.Name + " - " + e.ChickenCage.Cage.Name,
				Name:                e.AfkirChickenCustomer.Name,
				PhoneNumber:         e.AfkirChickenCustomer.PhoneNumber,
				TotalNominal:        e.TotalPrice.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			if e.PaidDate.Valid {
				receieveable.PaidDate = e.PaidDate.Time.Format("02 Jan 2006")
			} else {
				receieveable.PaidDate = "-"
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.TotalPrice.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	if filter.ReceivablesCategory == constant.ReceieveablesCategoryAll || filter.ReceivablesCategory == constant.ReceieveablesCategoryCashAdvance {
		for _, e := range userCashAdvances {
			receieveable := dto.ReceivablesListResponse{
				Id:                  e.Id,
				DeadlinePaymentDate: e.DeadlinePaymentDate.Format("02 Jan 2006"),
				Category:            constant.ReceieveablesCategoryCashAdvance,
				PlaceName:           e.User.Location.Name,
				Name:                e.User.Name,
				PhoneNumber:         e.User.PhoneNumber,
				TotalNominal:        e.Nominal.String(),
				PaymentStatus:       e.PaymentStatus.String(),
			}

			if e.PaidDate.Valid {
				receieveable.PaidDate = e.PaidDate.Time.Format("02 Jan 2006")
			} else {
				receieveable.PaidDate = "-"
			}

			totalCurrentPayment := decimal.Zero
			for _, p := range e.Payments {
				totalCurrentPayment = totalCurrentPayment.Add(p.Nominal)
			}

			receieveable.RemainingPayment = e.Nominal.Sub(totalCurrentPayment).String()

			receieveables = append(receieveables, receieveable)
		}
	}

	paidPercentage := 0.0
	unpaidPercentage := 0.0
	if !totalPayment.IsZero() {
		paidPercentage = totalReceivablesPayment.Div(totalPayment).InexactFloat64() * 100.0
		unpaidPercentage = totalRemainingReceieveablesPayment.Div(totalPayment).InexactFloat64() * 100.0
	}

	return dto.ReceivablesOverviewResponse{
		ReceivablesPie: dto.ReceivablesPieResponse{
			PaidPercentage:   paidPercentage,
			UnpaidPercentage: unpaidPercentage,
		},
		Receivables: receieveables,
	}, nil
}

func (s *CashflowService) GetReceiveables(receieveablesCategory string, id uint64) (dto.ReceivablesResponse, error) {
	s.repository.UseTx(false)

	switch receieveablesCategory {
	case constant.ReceieveablesCategoryWarehouseEggSale:
		data, err := s.repository.GetWarehouseSaleCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceivablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceievablesPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.ReceievablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		response := dto.ReceivablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryWarehouseEggSale,
			PlaceName:             data.Warehouse.Location.Name + " - " + data.Warehouse.Name,
			Name:                  data.Customer.Name,
			PhoneNumber:           data.Customer.PhoneNumber,
			Nominal:               data.TotalPrice.String(),
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           data.PaymentType.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}

		if data.PaidDate.Valid {
			response.PaidDate = data.PaidDate.Time.Format("02 Jan 2006")
		} else {
			response.PaidDate = "-"
		}

		return response, nil
	case constant.ReceieveablesCategoryStoreEggSale:
		data, err := s.repository.GetStoreSaleCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceivablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceievablesPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.ReceievablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		response := dto.ReceivablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryStoreEggSale,
			PlaceName:             data.Store.Location.Name + " - " + data.Store.Name,
			Name:                  data.Customer.Name,
			PhoneNumber:           data.Customer.PhoneNumber,
			Nominal:               data.TotalPrice.String(),
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           data.PaymentType.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}

		if data.PaidDate.Valid {
			response.PaidDate = data.PaidDate.Time.Format("02 Jan 2006")
		} else {
			response.PaidDate = "-"
		}

		return response, nil
	case constant.ReceieveablesCategoryAfkirChickenSale:
		data, err := s.repository.GetAfkirChickenSaleCashflow(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceivablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceievablesPaymentResponse, 0)
		totalRemainingPayment := data.TotalPrice
		for _, e := range data.Payments {
			paymentResponse := dto.ReceievablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		response := dto.ReceivablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryAfkirChickenSale,
			PlaceName:             data.ChickenCage.Cage.Location.Name + " - " + data.ChickenCage.Cage.Name,
			Name:                  data.AfkirChickenCustomer.Name,
			Nominal:               data.TotalPrice.String(),
			PhoneNumber:           data.AfkirChickenCustomer.PhoneNumber,
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           data.PaymentType.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}

		if data.PaidDate.Valid {
			response.PaidDate = data.PaidDate.Time.Format("02 Jan 2006")
		} else {
			response.PaidDate = "-"
		}

		return response, nil

	case constant.ReceieveablesCategoryCashAdvance:
		data, err := s.repository.GetUserCashAdvance(id)
		if err != nil {
			s.log.Error("failed get warehouse cashflow", zap.Error(err))
			return dto.ReceivablesResponse{}, err
		}

		paymentResponses := make([]dto.ReceievablesPaymentResponse, 0)
		totalRemainingPayment := data.Nominal
		for _, e := range data.Payments {
			paymentResponse := dto.ReceievablesPaymentResponse{
				Id:            e.Id,
				Date:          e.PaymentDate.Format("02 Jan 2006"),
				Nominal:       e.Nominal.String(),
				PaymentMethod: e.PaymentMethod.String(),
				PaymentProof:  e.PaymentProof,
			}

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		response := dto.ReceivablesResponse{
			Id:                    data.Id,
			Date:                  data.CreatedAt.Format("02 Jan 2006"),
			Time:                  data.CreatedAt.Format("15:04"),
			Category:              constant.ReceieveablesCategoryCashAdvance,
			PlaceName:             data.User.Location.Name,
			Name:                  data.User.Name,
			PhoneNumber:           data.User.PhoneNumber,
			Nominal:               data.Nominal.String(),
			RemainingPayment:      totalRemainingPayment.String(),
			PaymentType:           enum.PaymentTypeinstallment.String(),
			PaymentStatus:         data.PaymentStatus.String(),
			DeadlinePaymentDate:   data.DeadlinePaymentDate.Format("02 Jan 2006"),
			InputBy:               data.CreatedByUser.Name,
			ReceieveablesPayments: paymentResponses,
		}

		if data.PaidDate.Valid {
			response.PaidDate = data.PaidDate.Time.Format("02 Jan 2006")
		} else {
			response.PaidDate = "-"
		}

		return response, nil
	default:
		return dto.ReceivablesResponse{}, errx.BadRequest("invalid receivables category")
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
	compensationSalary, err := decimal.NewFromString(request.CompentationSalary)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid compensation salary format")
	}
	additionalWorkSalary, err := decimal.NewFromString(request.AdditionalWorkSalary)
	if err != nil {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid additional work salary format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.UserSalaryPaymentResponse{}, errx.BadRequest("invalid payment method")
	}

	totalCashbond := decimal.Zero
	var newPayments []entity.UserCashAdvancePayment

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

			totalCashbond = totalCashbond.Add(nominal)

			currentPayment := nominal
			for _, payment := range data.Payments {
				currentPayment = currentPayment.Add(payment.Nominal)
			}
			switch {
			case currentPayment.GreaterThan(data.Nominal):
				return dto.UserSalaryPaymentResponse{}, errx.BadRequest("nominal more than needed")
			case currentPayment.Equal(data.Nominal):
				data.PaymentStatus = enum.PaymentStatusPaid
				data.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
			default:
				data.PaymentStatus = enum.PaymentStatusUnpaid
			}

			newPayments = append(newPayments, entity.UserCashAdvancePayment{
				UserCashAdvanceId: capReq.UserCashAdvanceId,
				Nominal:           nominal,
				PaymentDate:       paymentDate,
				PaymentMethod:     paymentMethod,
				PaymentProof:      capReq.PaymentProof,
			})

			if err := s.repository.UpdateUserCashAdvance(&data); err != nil {
				s.log.Error("failed batch update user cash advances", zap.Error(err))
				return dto.UserSalaryPaymentResponse{}, err
			}
		}

		if err := s.repository.CreateUserCashAdvancePaymentBatch(&newPayments); err != nil {
			s.log.Error("failed batch create user cash advance payments", zap.Error(err))
			return dto.UserSalaryPaymentResponse{}, err
		}
	}

	userSalaryPayment.BaseSalary = baseSalary
	userSalaryPayment.CompentationSalary = compensationSalary
	userSalaryPayment.BonusSalary = bonusSalary
	userSalaryPayment.AdditionalWorkSalary = additionalWorkSalary
	userSalaryPayment.PaymentMethod = paymentMethod
	userSalaryPayment.PaymentProof = request.PaymentProof
	userSalaryPayment.Cashbond = totalCashbond
	userSalaryPayment.IsPaid = true
	userSalaryPayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if err := s.repository.UpdateUserSalaryPayment(&userSalaryPayment); err != nil {
		return dto.UserSalaryPaymentResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
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
	totalPaidDebtPayment := decimal.Zero
	totalRemainingDebtPayment := decimal.Zero

	for _, e := range warehouseItemProcurements {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentDebtPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentDebtPayment = totalCurrentDebtPayment.Add(p.Nominal)
		}
		totalPaidDebtPayment = totalPaidDebtPayment.Add(totalCurrentDebtPayment)
		totalRemainingDebtPayment = totalRemainingDebtPayment.Add(e.TotalPrice.Sub(totalCurrentDebtPayment))
	}

	for _, e := range warehouseItemCornProcurements {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentDebtPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentDebtPayment = totalCurrentDebtPayment.Add(p.Nominal)
		}
		totalPaidDebtPayment = totalPaidDebtPayment.Add(totalCurrentDebtPayment)
		totalRemainingDebtPayment = totalRemainingDebtPayment.Add(e.TotalPrice.Sub(totalCurrentDebtPayment))
	}

	for _, e := range chickenProcurements {
		totalPayment = totalPayment.Add(e.TotalPrice)
		totalCurrentDebtPayment := decimal.Zero
		for _, p := range e.Payments {
			totalCurrentDebtPayment = totalCurrentDebtPayment.Add(p.Nominal)
		}
		totalPaidDebtPayment = totalPaidDebtPayment.Add(totalCurrentDebtPayment)
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

	paidPercentage := 0.0
	unpaidPercentage := 0.0
	if !totalPayment.IsZero() {
		paidPercentage = totalPaidDebtPayment.Div(totalPayment).InexactFloat64() * 100.0
		unpaidPercentage = totalRemainingDebtPayment.Div(totalPayment).InexactFloat64() * 100.0
	}

	return dto.DebtOverviewResponse{
		DebtPie: dto.DebtPieResponse{
			PaidPercentage:   paidPercentage,
			UnpaidPercentage: unpaidPercentage,
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

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.DebtResponse{
			Id:                  data.Id,
			Date:                data.CreatedAt.Format("02 Jan 2006"),
			Time:                data.CreatedAt.Format("15:04"),
			Category:            constant.DebtCategoryChickenProcurement,
			PlaceName:           data.Cage.Location.Name + " - " + data.Cage.Name,
			TransactionName:     constant.DebtTransactionNameChickenProcurement,
			Name:                data.Supplier.Name,
			PhoneNumber:         data.Supplier.PhoneNumber,
			Nominal:             data.TotalPrice.String(),
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

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.DebtResponse{
			Id:                  data.Id,
			Date:                data.CreatedAt.Format("02 Jan 2006"),
			Time:                data.CreatedAt.Format("15:04"),
			Category:            constant.DebtCategoryWarehouseItemProcurement,
			PlaceName:           data.Warehouse.Location.Name + " - " + data.Warehouse.Name,
			Name:                data.Supplier.Name,
			PhoneNumber:         data.Supplier.PhoneNumber,
			TransactionName:     constant.DebtTransactionNameWarehouseItemProcurement,
			Nominal:             data.TotalPrice.String(),
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

			totalRemainingPayment = totalRemainingPayment.Sub(e.Nominal)
			paymentResponse.Remaining = totalRemainingPayment.String()
			paymentResponses = append(paymentResponses, paymentResponse)
		}

		return dto.DebtResponse{
			Id:                  data.Id,
			Date:                data.CreatedAt.Format("02 Jan 2006"),
			Time:                data.CreatedAt.Format("15:04"),
			Category:            constant.DebtCategoryWarehouseItemCornProcurement,
			PlaceName:           data.Warehouse.Location.Name + " - " + data.Warehouse.Name,
			Name:                data.Supplier.Name,
			PhoneNumber:         data.Supplier.PhoneNumber,
			Nominal:             data.TotalPrice.String(),
			RemainingPayment:    totalRemainingPayment.String(),
			PaymentType:         data.PaymentType.String(),
			PaymentStatus:       data.PaymentStatus.String(),
			TransactionName:     constant.DebtTransactionNameWarehouseItemCornProcurement,
			DeadlinePaymentDate: data.DeadlinePaymentDate.Time.Format("02 Jan 2006"),
			InputBy:             data.CreatedByUser.Name,
			DebtPayments:        paymentResponses,
		}, nil

	default:
		return dto.DebtResponse{}, errx.BadRequest("invalid debt category name")
	}
}

func (s *CashflowService) ExportSalesCashflowToExcel(filter dto.GetCashflowSaleReportFilter) (*excelize.File, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	f := excelize.NewFile()

	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
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

	for row, sale := range storeSales {
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

	for i := 0; i < len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(storeSheet, col, col, 20) // 20 is a good default, you can adjust
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
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

	for row, sale := range warehouseSales {
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

	for i := 0; i < len(headers); i++ {
		col, _ := excelize.ColumnNumberToName(i + 1)
		f.SetColWidth(warehouseSheet, col, col, 20)
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
		s.log.Error("failed get user salary payments", zap.Error(err))
		return dto.UserSalarySummaryResponse{}, err
	}

	var (
		totalUser                = len(userSalaryPayments)
		totalBaseSalary          = decimal.Zero
		totalAdditonalWorkSalary = decimal.Zero
		totalBonusSalary         = decimal.Zero
	)

	for _, userSalaryPayment := range userSalaryPayments {
		totalBaseSalary = totalBaseSalary.Add(userSalaryPayment.User.Salary)

		if userSalaryPayment.IsPaid {
			totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(userSalaryPayment.AdditionalWorkSalary)
			totalBonusSalary = totalBonusSalary.Add(userSalaryPayment.BonusSalary)
		} else {
			var (
				withDeleted          = true
				additionalWorkSalary = decimal.Zero
				bonusSalary          = decimal.Zero
			)

			additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(userSalaryPayment.UserId,
				dto.GetAdditionalWorkUserFilter{
					Month:       param.MonthParam(filter.Month.Value()),
					Year:        filter.Year,
					WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
				})
			if err != nil {
				return dto.UserSalarySummaryResponse{}, err
			}

			dailyWorkUsers, err := s.workService.GetDailyWorkUserByUserId(userSalaryPayment.UserId,
				dto.GetDailyWorkUserFilter{
					Month:       param.MonthParam(filter.Month.Value()),
					Year:        filter.Year,
					WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
				})
			if err != nil {
				return dto.UserSalarySummaryResponse{}, err
			}

			userPresences, err := s.presenceService.GetUserPresencesByUserId(userSalaryPayment.UserId,
				dto.GetPresenceFilter{
					Month: param.MonthParam(filter.Month.Value()),
					Year:  filter.Year,
				})
			if err != nil {
				return dto.UserSalarySummaryResponse{}, err
			}

			presenceScore, workScore, totalNotPresent := util.CalculateKPIScoreUserInMonth(additionalWorkUsers, dailyWorkUsers, userPresences)

			if userSalaryPayment.User.Role.Name == constant.RolePekerjaKandang {
				totalDayInMonth := util.TotalDaysInMonth(int(filter.Year), time.Month(filter.Month.Value()))
				salaryPerDay := userSalaryPayment.User.Salary.Div(decimal.NewFromUint64(totalDayInMonth))
				reduceSalaryCauseNotPresent := salaryPerDay.Mul(decimal.NewFromUint64(totalNotPresent))
				bonusSalary = bonusSalary.Add(reduceSalaryCauseNotPresent)
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

			if presenceScore*0.6 == 60 {
				bonusSalary = bonusSalary.Add(decimal.NewFromFloat(constant.BonusFullPresent))
			}

			kpiPerformance := (presenceScore * 0.6) + (workScore * 0.4)
			if kpiPerformance >= 90 {
				bonusSalary = bonusSalary.Add(decimal.NewFromFloat(constant.BonusGoodPerformancePercentage).Mul(userSalaryPayment.User.Salary))
			} else if kpiPerformance < 75 {
				bonusSalary = bonusSalary.Sub(decimal.NewFromFloat(constant.BonusBadPerformancePercentage).Mul(userSalaryPayment.User.Salary))
			}

			totalBonusSalary = totalBonusSalary.Add(bonusSalary)
		}
	}

	return dto.UserSalarySummaryResponse{
		TotalUser:                uint64(totalUser),
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
		response.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
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

	var (
		additionalWorkUserResponses     = make([]dto.AdditionalWorkUserResponse, 0)
		userCashAdvanceSummaryResponses = make([]dto.UserCashAdvanceSummaryResponse, 0)
		date                            = "-"
		time                            = "-"
		totalAdditonalWorkSalary        = decimal.Zero
		totalBonusSalary                = decimal.Zero
	)

	if userSalaryPayment.IsPaid {
		totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(userSalaryPayment.AdditionalWorkSalary)
		totalBonusSalary = totalBonusSalary.Add(userSalaryPayment.BonusSalary)
		date = userSalaryPayment.CreatedAt.Format("02 Jan 2006")
		time = userSalaryPayment.CreatedAt.Format("15:04")
	} else {
		var (
			withDeleted              = true
			additionalWorkSalaryTemp = decimal.Zero
			bonusSalary              = decimal.Zero
		)

		additionalWorkUsers, err := s.workService.GetAdditionalWorkUserByUserId(userSalaryPayment.UserId,
			dto.GetAdditionalWorkUserFilter{
				Month:       param.MonthParam(enum.ValueOfMonth(userSalaryPayment.CreatedAt.Format("January"))),
				Year:        uint64(userSalaryPayment.CreatedAt.Year()),
				WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
			})
		if err != nil {
			return dto.UserSalaryDetailResponse{}, err
		}

		for _, e := range additionalWorkUsers.AdditionalWorkUsers {
			if e.IsDone {
				additionalWorkUserResponses = append(additionalWorkUserResponses, e)
				salary, err := decimal.NewFromString(e.AdditionalWork.Salary)
				if err != nil {
					s.log.Error("failed parse additional work salary", zap.Error(err))
					return dto.UserSalaryDetailResponse{}, err
				}
				additionalWorkSalaryTemp = additionalWorkSalaryTemp.Add(salary)
			}
		}

		totalAdditonalWorkSalary = totalAdditonalWorkSalary.Add(additionalWorkSalaryTemp)

		userPresences, err := s.presenceService.GetUserPresencesByUserId(userSalaryPayment.UserId,
			dto.GetPresenceFilter{
				Month: param.MonthParam(enum.ValueOfMonth(userSalaryPayment.CreatedAt.Format("January"))),
				Year:  uint64(userSalaryPayment.CreatedAt.Year()),
			})
		if err != nil {
			return dto.UserSalaryDetailResponse{}, err
		}

		dailyWorkUsers, err := s.workService.GetDailyWorkUserByUserId(userSalaryPayment.UserId,
			dto.GetDailyWorkUserFilter{
				Month:       param.MonthParam(enum.ValueOfMonth(userSalaryPayment.CreatedAt.Format("January"))),
				Year:        uint64(userSalaryPayment.CreatedAt.Year()),
				WithDeleted: &withDeleted, // In case the user work is done but the work is deleted
			})
		if err != nil {
			return dto.UserSalaryDetailResponse{}, err
		}

		presenceScore, workScore, totalNotPresent := util.CalculateKPIScoreUserInMonth(additionalWorkUsers, dailyWorkUsers, userPresences)

		if userSalaryPayment.User.Role.Name == constant.RolePekerjaKandang {
			totalDayInMonth := util.TotalDaysInMonth(userSalaryPayment.CreatedAt.Year(), userSalaryPayment.CreatedAt.Month())
			salaryPerDay := userSalaryPayment.User.Salary.Div(decimal.NewFromUint64(totalDayInMonth))
			reduceSalaryCauseNotPresent := salaryPerDay.Mul(decimal.NewFromUint64(totalNotPresent))
			bonusSalary = bonusSalary.Add(reduceSalaryCauseNotPresent)
		}

		if presenceScore*0.6 == 60 {
			bonusSalary = bonusSalary.Add(decimal.NewFromFloat(constant.BonusFullPresent))
		}

		kpiPerformance := (presenceScore * 0.6) + (workScore * 0.4)
		if kpiPerformance >= 90 {
			bonusSalary = bonusSalary.Add(decimal.NewFromFloat(constant.BonusGoodPerformancePercentage).Mul(userSalaryPayment.BaseSalary))
		} else if kpiPerformance < 75 {
			bonusSalary = bonusSalary.Sub(decimal.NewFromFloat(constant.BonusBadPerformancePercentage).Mul(userSalaryPayment.BaseSalary))
		}

		userCashAdvanceSummary, err := s.GetUserCashAdvanceByUserId(userSalaryPayment.UserId)
		if err != nil {
			return dto.UserSalaryDetailResponse{}, err
		}

		userCashAdvanceSummaryResponses = userCashAdvanceSummary
		totalBonusSalary = totalBonusSalary.Add(bonusSalary)
	}

	return dto.UserSalaryDetailResponse{
		User:                     mapper.UserToListResponse(&userSalaryPayment.User),
		Date:                     date,
		Time:                     time,
		SalaryMonth:              userSalaryPayment.CreatedAt.Format("January"),
		AdditionalWorkUsers:      additionalWorkUserResponses,
		UserCashAdvanceSummaries: userCashAdvanceSummaryResponses,
		BaseSalary:               userSalaryPayment.BaseSalary.String(),
		CompentationSalary:       userSalaryPayment.CompentationSalary.String(),
		BonusSalary:              totalBonusSalary.String(),
		AdditionalWorkSalary:     totalAdditonalWorkSalary.String(),
	}, nil
}

func (s *CashflowService) GetCashflowSaleOverview(filter dto.GetCashflowSaleOverviewFilter) (dto.CashflowSaleOverviewResponse, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))
	weeks := util.GetFourWeekRanges(int(filter.Year), time.Month(filter.Month.Value()))

	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		LocationId: filter.LocationId,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed to get store sale", zap.Error(err))
		return dto.CashflowSaleOverviewResponse{}, err
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		LocationId: filter.LocationId,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed to get store sale", zap.Error(err))
		return dto.CashflowSaleOverviewResponse{}, err
	}

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.CashflowSaleOverviewResponse{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.CrackedEgg, constant.UnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.CashflowSaleOverviewResponse{}, err
	}

	brokenEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.BrokenEgg, constant.UnitPlastik, enum.ItemCategoryEgg)
	if err != nil {
		return dto.CashflowSaleOverviewResponse{}, err
	}

	locationSaleSummaries := make([]dto.LocationSaleSummaryResponse, 0)
	eggSaleLocationRevenueSummary, eggSaleLocationReceivablesSummary := s.getEggSaleLocationSummary(&storeSales, &warehouseSales)

	for key := range eggSaleLocationRevenueSummary {
		locationSaleSummaries = append(locationSaleSummaries, dto.LocationSaleSummaryResponse{
			PlaceName:     key,
			Income:        eggSaleLocationRevenueSummary[key].String(),
			Receieveables: eggSaleLocationReceivablesSummary[key].String(),
		})
	}

	eggSaleGraphs, err := s.buildStoreOverviewMonthlyGraph(filter.LocationId, filter.ItemId, filter.Year, filter.Month.Value())
	if err != nil {
		return dto.CashflowSaleOverviewResponse{}, err
	}

	goodEggStoreInKg := float64(0)
	crackedEggStoreInKg := float64(0)
	brokenEggStoreInPlastik := float64(0)

	for _, storeSale := range storeSales {
		if storeSale.SaleUnit.String() == constant.UnitKg && goodEggItem.Id == storeSale.ItemId {
			goodEggStoreInKg += storeSale.Quantity
		} else if storeSale.SaleUnit.String() == constant.UnitIkat && goodEggItem.Id == storeSale.ItemId {
			goodEggStoreInKg += storeSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if storeSale.SaleUnit.String() == constant.UnitKg && crackedEggItem.Id == storeSale.ItemId {
			crackedEggStoreInKg += storeSale.Quantity
		} else if storeSale.SaleUnit.String() == constant.UnitIkat && crackedEggItem.Id == storeSale.ItemId {
			crackedEggStoreInKg += storeSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if storeSale.SaleUnit.String() == constant.UnitPlastik && brokenEggItem.Id == storeSale.ItemId {
			brokenEggStoreInPlastik += storeSale.Quantity
		}
	}

	goodEggWarehouseInKg := float64(0)
	crackedEggWarehouseInKg := float64(0)
	brokenEggWarehouseInPlastik := float64(0)

	for _, warehouseSale := range warehouseSales {
		if warehouseSale.SaleUnit.String() == constant.UnitKg && goodEggItem.Id == warehouseSale.ItemId {
			goodEggWarehouseInKg += warehouseSale.Quantity
		} else if warehouseSale.SaleUnit.String() == constant.UnitIkat && goodEggItem.Id == warehouseSale.ItemId {
			goodEggWarehouseInKg += warehouseSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if warehouseSale.SaleUnit.String() == constant.UnitKg && crackedEggItem.Id == warehouseSale.ItemId {
			crackedEggWarehouseInKg += warehouseSale.Quantity
		} else if warehouseSale.SaleUnit.String() == constant.UnitIkat && crackedEggItem.Id == warehouseSale.ItemId {
			crackedEggWarehouseInKg += warehouseSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if warehouseSale.SaleUnit.String() == constant.UnitPlastik && brokenEggItem.Id == warehouseSale.ItemId {
			brokenEggWarehouseInPlastik += warehouseSale.Quantity
		}
	}

	income, receivables, err := s.getIncomeAndReceivablesPerWeek(filter.LocationId, weeks, startDate, endDate)
	if err != nil {
		s.log.Error("failed get income and receive per week", zap.Error(err))
		return dto.CashflowSaleOverviewResponse{}, err
	}

	expense, debt, err := s.getExpenseAndDebtPerWeek(filter.LocationId, weeks, startDate, endDate)
	if err != nil {
		s.log.Error("failed get expense and debt per week", zap.Error(err))
		return dto.CashflowSaleOverviewResponse{}, err
	}

	prevMonth := filter.Month.Value() - 1
	prevYear := filter.Year

	if prevMonth == 0 {
		prevMonth = 12
		prevYear = filter.Year - 1
	}

	prevMonthWeek := util.GetFourWeekRanges(int(prevYear), time.Month(prevYear))

	incomePreviousMonth, receivablesPreviousMonth, err := s.getIncomeAndReceivablesPerWeek(filter.LocationId, prevMonthWeek, startDate, endDate)
	if err != nil {
		s.log.Error("failed get income and receivable per week", zap.Error(err))
		return dto.CashflowSaleOverviewResponse{}, err
	}

	expensePreviousMonth, debtPreviousMonth, err := s.getExpenseAndDebtPerWeek(filter.LocationId, prevMonthWeek, startDate, endDate)
	if err != nil {
		s.log.Error("failed get expense and debt per week", zap.Error(err))
		return dto.CashflowSaleOverviewResponse{}, err
	}

	totalIncome := decimal.Zero
	totalExpense := decimal.Zero
	totalProfit := decimal.Zero
	cashflowSaleGraphs := make([]dto.CashflowSaleGraphResponse, 0)
	for key := range weeks {
		cashflowSaleGraphs = append(cashflowSaleGraphs, dto.CashflowSaleGraphResponse{
			Key:     fmt.Sprintf("Minggu %d", key),
			Income:  income[key].String(),
			Expense: expense[key].String(),
			Profit:  income[key].Add(receivables[key]).Sub(expense[key].Add(debt[key])).String(),
		})

		totalIncome = totalIncome.Add(income[key])
		totalExpense = totalExpense.Add(expense[key])
		totalProfit = totalProfit.Add(income[key].Add(receivables[key]).Sub(expense[key].Add(debt[key])))
	}

	totalPreviousMonthIncome := decimal.Zero
	totalPreviousMonthExpense := decimal.Zero
	totalPreviousMonthProfit := decimal.Zero
	for key := range weeks {
		totalPreviousMonthIncome = totalPreviousMonthIncome.Add(incomePreviousMonth[key])
		totalPreviousMonthExpense = totalPreviousMonthExpense.Add(expensePreviousMonth[key])
		totalPreviousMonthProfit = totalPreviousMonthProfit.Add(incomePreviousMonth[key].Add(receivablesPreviousMonth[key]).Sub(expensePreviousMonth[key].Add(debtPreviousMonth[key])))
	}

	incomeIncrease, incomeDiff := calculateDiff(totalIncome, totalPreviousMonthIncome)
	profitIncrease, profitDiff := calculateDiff(totalProfit, totalPreviousMonthProfit)
	expenseIncrease, expenseDiff := calculateDiff(totalExpense, totalPreviousMonthExpense)

	return dto.CashflowSaleOverviewResponse{
		CashflowSaleSummary: dto.CashflowSaleSummaryResponse{
			Income:                totalIncome.String(),
			IsIncomeIncrease:      incomeIncrease,
			IncomeDiffPercentage:  incomeDiff,
			Profit:                totalProfit.String(),
			IsProfitIncrease:      profitIncrease,
			ProfitDiffPercentage:  profitDiff,
			Expense:               totalExpense.String(),
			IsExpenseIncrease:     expenseIncrease,
			ExpenseDiffPercentage: expenseDiff,
		},
		EggSaleSummary: dto.EggSaleSummaryResponse{
			TotalGoodEggInKg:        goodEggStoreInKg + goodEggWarehouseInKg,
			TotalGoodEggInIkat:      math.Ceil((goodEggStoreInKg + goodEggWarehouseInKg) / float64(constant.TotalEggPerIkat)),
			TotalCrackedEggInKg:     crackedEggStoreInKg + crackedEggWarehouseInKg,
			TotalCrackedEggInIkat:   math.Ceil((crackedEggStoreInKg + crackedEggWarehouseInKg) / float64(constant.TotalEggPerIkat)),
			TotalBrokenEggInPlastik: brokenEggStoreInPlastik + brokenEggWarehouseInPlastik,
		},
		CashflowSaleGraphs:  cashflowSaleGraphs,
		EggSaleGraphs:       eggSaleGraphs,
		LocationSaleSummary: locationSaleSummaries,
		LocationPieChart: dto.LocationPieChartResponse{
			StorePercentage:     (goodEggStoreInKg + crackedEggStoreInKg) / (goodEggStoreInKg + crackedEggStoreInKg + crackedEggStoreInKg + crackedEggWarehouseInKg),
			WarehousePercentage: (goodEggWarehouseInKg + crackedEggWarehouseInKg) / (goodEggStoreInKg + crackedEggStoreInKg + crackedEggStoreInKg + crackedEggWarehouseInKg),
		},
	}, nil
}

func (s *CashflowService) GetCashflowHistories(filter dto.GetCashflowHistoryFilter) ([]dto.CashflowHistoryResponse, error) {
	s.repository.UseTx(false)

	cashflowHistories, err := s.repository.GetCashflowHistories(dto.GetCashflowHistoryFilter{
		Year:       filter.Year,
		LocationId: filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed get cashflow histories", zap.Error(err))
		return nil, err
	}

	if time.Now().Year() != int(filter.Year) {
		currentCashflowHistory, err := s.getCashflowHistoryInMonth(filter.LocationId, filter.Year, enum.Month(time.Now().Month()))
		if err != nil {
			return nil, err
		}
		cashflowHistories = append(cashflowHistories, currentCashflowHistory)
	}

	cashflowByMonth := make(map[int]entity.CashflowHistory)
	for _, cf := range cashflowHistories {
		month := int(cf.CreatedAt.Month())

		if existing, ok := cashflowByMonth[month]; ok {
			cashflowByMonth[month] = entity.CashflowHistory{
				Income:           existing.Income.Add(cf.Income),
				Profit:           existing.Profit.Add(cf.Profit),
				Expense:          existing.Expense.Add(cf.Expense),
				Cash:             existing.Cash.Add(cf.Cash),
				Receivables:      existing.Receivables.Add(cf.Receivables),
				Debt:             existing.Debt.Add(cf.Debt),
				WarehouseEggSale: existing.WarehouseEggSale.Add(cf.WarehouseEggSale),
				StoreEggSale:     existing.StoreEggSale.Add(cf.StoreEggSale),
				CreatedAt:        existing.CreatedAt,
			}
		} else {
			cashflowByMonth[month] = cf
		}
	}

	cashflowHistoryResponses := make([]dto.CashflowHistoryResponse, 0, 12)

	for month := 1; month <= 12; month++ {
		cf, ok := cashflowByMonth[month]
		if !ok {
			cf = entity.CashflowHistory{
				Income:           decimal.Zero,
				Profit:           decimal.Zero,
				Expense:          decimal.Zero,
				Cash:             decimal.Zero,
				Receivables:      decimal.Zero,
				Debt:             decimal.Zero,
				WarehouseEggSale: decimal.Zero,
				StoreEggSale:     decimal.Zero,
				CreatedAt:        time.Date(int(filter.Year), time.Month(month), 1, 0, 0, 0, 0, time.UTC),
			}
		}

		cashflowHistoryResponses = append(cashflowHistoryResponses, dto.CashflowHistoryResponse{
			LocationId:       cf.LocationId,
			Income:           cf.Income.String(),
			Profit:           cf.Profit.String(),
			Expense:          cf.Expense.String(),
			Cash:             cf.Cash.String(),
			Receivables:      cf.Receivables.String(),
			Debt:             cf.Debt.String(),
			StoreEggSale:     cf.StoreEggSale.String(),
			WarehouseEggSale: cf.WarehouseEggSale.String(),
			CreatedAt:        cf.CreatedAt,
		})
	}

	return cashflowHistoryResponses, nil
}

func (s *CashflowService) GetCashflowOverview(filter dto.GetCashflowOverviewFilter) (dto.CashflowOverviewResponse, error) {
	// Cash -> total price (penjualan ayam, penjualan telur toko, penjualan telur gudang) // Pengeluaran -> total pengeluaran perusahaan (pengadaan barang, pengadaan jagung, pembelian ayam, pembayaran gaji, operasional) -> yang dibayarkan // Pendapatan -> total pendapatan perusahaan (penjualan ayam, penjualan telur toko, penjualan telur gudang) -> yang sudah dibayarkan // Keuntutngan -> total pendapatan + piutang - total pengeluaran + total hiutang

	s.repository.UseTx(false)

	cashflowHistories, err := s.repository.GetCashflowHistories(dto.GetCashflowHistoryFilter{
		Year:       filter.Year,
		LocationId: filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed get cashflow histories", zap.Error(err))
		return dto.CashflowOverviewResponse{}, err
	}

	if time.Now().Year() != int(filter.Year) {
		currentCashflowHistory, err := s.getCashflowHistoryInMonth(filter.LocationId, filter.Year, enum.Month(time.Now().Month()))
		if err != nil {
			return dto.CashflowOverviewResponse{}, err
		}
		cashflowHistories = append(cashflowHistories, currentCashflowHistory)
	}

	cashflowByMonth := make(map[int]entity.CashflowHistory)
	for _, cf := range cashflowHistories {
		month := int(cf.CreatedAt.Month())

		if existing, ok := cashflowByMonth[month]; ok {
			cashflowByMonth[month] = entity.CashflowHistory{
				Income:           existing.Income.Add(cf.Income),
				Profit:           existing.Profit.Add(cf.Profit),
				Expense:          existing.Expense.Add(cf.Expense),
				Cash:             existing.Cash.Add(cf.Cash),
				Receivables:      existing.Receivables.Add(cf.Receivables),
				Debt:             existing.Debt.Add(cf.Debt),
				WarehouseEggSale: existing.WarehouseEggSale.Add(cf.WarehouseEggSale),
				StoreEggSale:     existing.StoreEggSale.Add(cf.StoreEggSale),
				CreatedAt:        existing.CreatedAt,
			}
		} else {
			cashflowByMonth[month] = cf
		}
	}

	cashflowGraphs := make([]dto.CashflowGraphResponse, 0, 12)
	eggSaleCashflowGraphs := make([]dto.EggSaleCashflowGraphResponse, 0, 12)

	totalIncome := decimal.Zero
	totalProfit := decimal.Zero
	totalExpense := decimal.Zero
	totalCash := decimal.Zero
	totalReceivables := decimal.Zero
	totalDebt := decimal.Zero

	for month := 1; month <= 12; month++ {
		cf, ok := cashflowByMonth[month]
		if !ok {
			cf = entity.CashflowHistory{
				Income:           decimal.Zero,
				Profit:           decimal.Zero,
				Expense:          decimal.Zero,
				Cash:             decimal.Zero,
				Receivables:      decimal.Zero,
				Debt:             decimal.Zero,
				WarehouseEggSale: decimal.Zero,
				StoreEggSale:     decimal.Zero,
				CreatedAt:        time.Date(int(filter.Year), time.Month(month), 1, 0, 0, 0, 0, time.UTC),
			}
		}

		cashflowGraphs = append(cashflowGraphs, dto.CashflowGraphResponse{
			Key:     cf.CreatedAt.Format("January"),
			Income:  cf.Income.String(),
			Profit:  cf.Profit.String(),
			Expense: cf.Expense.String(),
			Cash:    cf.Cash.String(),
		})

		eggSaleCashflowGraphs = append(eggSaleCashflowGraphs, dto.EggSaleCashflowGraphResponse{
			Key:              cf.CreatedAt.Format("January"),
			WarehouseEggSale: cf.WarehouseEggSale.String(),
			StoreEggSale:     cf.StoreEggSale.String(),
		})

		totalIncome = totalIncome.Add(cf.Income)
		totalProfit = totalProfit.Add(cf.Profit)
		totalExpense = totalExpense.Add(cf.Expense)
		totalCash = totalCash.Add(cf.Cash)
		totalReceivables = totalReceivables.Add(cf.Receivables)
		totalDebt = totalDebt.Add(cf.Debt)
	}

	previousCashflowHistories, err := s.repository.GetCashflowHistories(dto.GetCashflowHistoryFilter{
		Year:       filter.Year - 1,
		LocationId: filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed get cashflow histories previous year", zap.Error(err))
		return dto.CashflowOverviewResponse{}, err
	}

	totalPreviousIncome := decimal.Zero
	totalPreviousProfit := decimal.Zero
	totalPreviousExpense := decimal.Zero
	totalPreviousCash := decimal.Zero
	totalPreviousReceivables := decimal.Zero
	totalPreviousDebt := decimal.Zero

	for _, prev := range previousCashflowHistories {
		totalPreviousIncome = totalPreviousIncome.Add(prev.Income)
		totalPreviousProfit = totalPreviousProfit.Add(prev.Profit)
		totalPreviousExpense = totalPreviousExpense.Add(prev.Expense)
		totalPreviousCash = totalPreviousCash.Add(prev.Cash)
		totalPreviousReceivables = totalPreviousReceivables.Add(prev.Receivables)
		totalPreviousDebt = totalPreviousDebt.Add(prev.Debt)
	}

	incomeIncrease, incomeDiff := calculateDiff(totalIncome, totalPreviousIncome)
	profitIncrease, profitDiff := calculateDiff(totalProfit, totalPreviousProfit)
	expenseIncrease, expenseDiff := calculateDiff(totalExpense, totalPreviousExpense)
	cashIncrease, cashDiff := calculateDiff(totalCash, totalPreviousCash)
	receivablesIncrease, receivablesDiff := calculateDiff(totalReceivables, totalPreviousReceivables)
	debtIncrease, debtDiff := calculateDiff(totalDebt, totalPreviousDebt)

	return dto.CashflowOverviewResponse{
		CashflowSummary: dto.CashflowSummaryResponse{
			Income:                    totalIncome.String(),
			IsIncomeIncrease:          incomeIncrease,
			IncomeDiffPercentage:      incomeDiff,
			Profit:                    totalProfit.String(),
			IsProfitIncrease:          profitIncrease,
			ProfitDiffPercentage:      profitDiff,
			Expense:                   totalExpense.String(),
			IsExpenseIncrease:         expenseIncrease,
			ExpenseDiffPercentage:     expenseDiff,
			Debt:                      totalDebt.String(),
			IsDebtIncrease:            debtIncrease,
			DebtDiffPercentage:        debtDiff,
			Cash:                      totalCash.String(),
			IsCashIncrease:            cashIncrease,
			CashDiffPercentage:        cashDiff,
			Receivables:               totalReceivables.String(),
			IsReceivablesIncrease:     receivablesIncrease,
			ReceivablesDiffPercentage: receivablesDiff,
		},
		CashflowGraphs:        cashflowGraphs,
		EggSaleCashflowGraphs: eggSaleCashflowGraphs,
	}, nil
}

func (s *CashflowService) getIncomeAndReceivablesPerWeek(locationId uint64, weeks map[int]util.DateRange, startDate, endDate time.Time) (map[int]decimal.Decimal, map[int]decimal.Decimal, error) {
	income := make(map[int]decimal.Decimal)
	receivables := make(map[int]decimal.Decimal)
	for w := range weeks {
		income[w] = decimal.Zero
		receivables[w] = decimal.Zero
	}

	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		LocationId:               locationId,
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return nil, nil, err
	}
	for _, sale := range storeSales {
		var paid decimal.Decimal
		for _, p := range sale.Payments {
			week := util.FindWeek(p.PaymentDate, weeks)
			income[week] = income[week].Add(p.Nominal)
			paid = paid.Add(p.Nominal)
		}
		if sale.DeadlinePaymentDate.Valid {
			week := util.FindWeek(sale.DeadlinePaymentDate.Time, weeks)
			receivables[week] = receivables[week].Add(sale.TotalPrice.Sub(paid))
		}
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		LocationId:               locationId,
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return nil, nil, err
	}
	for _, sale := range warehouseSales {
		var paid decimal.Decimal
		for _, p := range sale.Payments {
			week := util.FindWeek(p.PaymentDate, weeks)
			income[week] = income[week].Add(p.Nominal)
			paid = paid.Add(p.Nominal)
		}
		if sale.DeadlinePaymentDate.Valid {
			week := util.FindWeek(sale.DeadlinePaymentDate.Time, weeks)
			receivables[week] = receivables[week].Add(sale.TotalPrice.Sub(paid))
		}
	}

	afkirSales, err := s.repository.GetAfkirChickenSaleCashflows(dto.GetAfkirChickenSaleFilter{
		LocationId:               locationId,
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return nil, nil, err
	}
	for _, sale := range afkirSales {
		var paid decimal.Decimal
		for _, p := range sale.Payments {
			week := util.FindWeek(p.PaymentDate, weeks)
			income[week] = income[week].Add(p.Nominal)
			paid = paid.Add(p.Nominal)
		}
		if sale.DeadlinePaymentDate.Valid {
			week := util.FindWeek(sale.DeadlinePaymentDate.Time, weeks)
			receivables[week] = receivables[week].Add(sale.TotalPrice.Sub(paid))
		}
	}

	userCashAdvances, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		LocationId:               locationId,
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return nil, nil, err
	}
	for _, adv := range userCashAdvances {
		var paid decimal.Decimal
		for _, p := range adv.Payments {
			week := util.FindWeek(p.PaymentDate, weeks)
			income[week] = income[week].Add(p.Nominal)
			paid = paid.Add(p.Nominal)
		}
		week := util.FindWeek(adv.DeadlinePaymentDate, weeks)
		receivables[week] = receivables[week].Add(adv.Nominal.Sub(paid))
	}

	return income, receivables, nil
}

func (s *CashflowService) getExpenseAndDebtPerWeek(locationId uint64, weeks map[int]util.DateRange, startDate, endDate time.Time) (map[int]decimal.Decimal, map[int]decimal.Decimal, error) {
	expense := make(map[int]decimal.Decimal)
	debt := make(map[int]decimal.Decimal)
	for w := range weeks {
		expense[w] = decimal.Zero
		debt[w] = decimal.Zero
	}

	warehouseItemProcurements, err := s.repository.GetWarehouseItemProcurementCashflows(dto.GetWarehouseItemProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		return nil, nil, err
	}
	for _, p := range warehouseItemProcurements {
		var paid decimal.Decimal
		for _, pay := range p.Payments {
			week := util.FindWeek(pay.PaymentDate, weeks)
			expense[week] = expense[week].Add(pay.Nominal)
			paid = paid.Add(pay.Nominal)
		}
		if p.DeadlinePaymentDate.Valid {
			week := util.FindWeek(p.DeadlinePaymentDate.Time, weeks)
			debt[week] = debt[week].Add(p.TotalPrice.Sub(paid))
		}
	}

	warehouseCornProcurements, err := s.repository.GetWarehouseItemCornProcurementCashflows(dto.GetWarehouseItemCornProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		return nil, nil, err
	}
	for _, p := range warehouseCornProcurements {
		var paid decimal.Decimal
		for _, pay := range p.Payments {
			week := util.FindWeek(pay.PaymentDate, weeks)
			expense[week] = expense[week].Add(pay.Nominal)
			paid = paid.Add(pay.Nominal)
		}
		if p.DeadlinePaymentDate.Valid {
			week := util.FindWeek(p.DeadlinePaymentDate.Time, weeks)
			debt[week] = debt[week].Add(p.TotalPrice.Sub(paid))
		}
	}

	chickenProcurements, err := s.repository.GetChickenProcurementCashflows(dto.GetChickenProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		return nil, nil, err
	}
	for _, p := range chickenProcurements {
		var paid decimal.Decimal
		for _, pay := range p.Payments {
			week := util.FindWeek(pay.PaymentDate, weeks)
			expense[week] = expense[week].Add(pay.Nominal)
			paid = paid.Add(pay.Nominal)
		}
		if p.DeadlinePaymentDate.Valid {
			week := util.FindWeek(p.DeadlinePaymentDate.Time, weeks)
			debt[week] = debt[week].Add(p.TotalPrice.Sub(paid))
		}
	}

	userSalaryPayments, err := s.repository.GetUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		return nil, nil, err
	}
	for _, salary := range userSalaryPayments {
		total := salary.BaseSalary.
			Add(salary.BonusSalary).
			Add(salary.CompentationSalary).
			Add(salary.AdditionalWorkSalary).
			Add(salary.Cashbond)
		week := util.FindWeek(salary.CreatedAt, weeks)
		expense[week] = expense[week].Add(total)
	}

	return expense, debt, nil
}

func (s *CashflowService) buildStoreOverviewMonthlyGraph(locationId uint64, itemId uint64, year uint64, month enum.Month) ([]dto.EggSaleGraphResponse, error) {
	weekMaps := util.GetFourWeekRanges(int(year), time.Month(month))
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))

	monthStoreSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		LocationId: locationId,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		ItemId:     itemId,
	})
	if err != nil {
		s.log.Error("failed to get store sale monthly", zap.Error(err))
		return nil, err
	}

	monthWarehosueSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		LocationId: locationId,
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		ItemId:     itemId,
	})
	if err != nil {
		s.log.Error("failed to get store sale monthly", zap.Error(err))
		return nil, err
	}

	itemSales := make(map[int]float64)
	for _, storeSale := range monthStoreSales {
		week := util.FindWeek(storeSale.CreatedAt, weekMaps)
		if week > 0 {
			if storeSale.SaleUnit.String() == constant.UnitKg {
				itemSales[week] += storeSale.Quantity
			} else if storeSale.SaleUnit.String() == constant.UnitIkat {
				itemSales[week] += storeSale.Quantity * float64(constant.TotalEggPerIkat)
			} else if storeSale.SaleUnit.String() == constant.UnitPlastik {
				itemSales[week] += storeSale.Quantity
			}
		}
	}

	for _, warehosueSale := range monthWarehosueSales {
		week := util.FindWeek(warehosueSale.CreatedAt, weekMaps)
		if week > 0 {
			if warehosueSale.SaleUnit.String() == constant.UnitKg {
				itemSales[week] += warehosueSale.Quantity
			} else if warehosueSale.SaleUnit.String() == constant.UnitIkat {
				itemSales[week] += warehosueSale.Quantity * float64(constant.TotalEggPerIkat)
			} else if warehosueSale.SaleUnit.String() == constant.UnitPlastik {
				itemSales[week] += warehosueSale.Quantity
			}
		}
	}

	keys := util.GetSortedKeys(weekMaps)
	graphs := make([]dto.EggSaleGraphResponse, 0)
	for _, k := range keys {
		graphs = append(graphs, dto.EggSaleGraphResponse{
			Key:   fmt.Sprintf("Minggu %d", k),
			Value: itemSales[k],
		})
	}

	return graphs, nil
}

func (s *CashflowService) getEggSaleLocationSummary(storeSales *[]entity.StoreSale, warehouseSales *[]entity.WarehouseSale) (map[string]decimal.Decimal, map[string]decimal.Decimal) {
	revenueMap := make(map[string]decimal.Decimal)
	receivableMap := make(map[string]decimal.Decimal)

	update := func(loc string, revenue, receivable decimal.Decimal) {
		revenueMap[loc] = revenueMap[loc].Add(revenue)
		receivableMap[loc] = receivableMap[loc].Add(receivable)
	}

	if storeSales != nil {
		for _, sale := range *storeSales {
			loc := "Unknown"
			if sale.Store.Id != 0 && sale.Store.Name != "" {
				loc = sale.Store.Name
			}

			paid := decimal.NewFromInt(0)
			for _, p := range sale.Payments {
				paid = paid.Add(p.Nominal)
			}

			receivable := sale.TotalPrice.Sub(paid)
			update(loc, paid, receivable)
		}
	}

	if warehouseSales != nil {
		for _, sale := range *warehouseSales {
			loc := "Unknown"
			if sale.Warehouse.Id != 0 && sale.Warehouse.Name != "" {
				loc = sale.Warehouse.Name
			}

			paid := decimal.NewFromInt(0)
			for _, p := range sale.Payments {
				paid = paid.Add(p.Nominal)
			}

			receivable := sale.TotalPrice.Sub(paid)
			update(loc, paid, receivable)
		}
	}

	return revenueMap, receivableMap
}

func calculateDiff(current, previous decimal.Decimal) (isIncrease bool, percentage float64) {
	if previous.IsZero() {
		if current.IsZero() {
			return false, 0
		}
		return true, 100
	}

	diff := current.Sub(previous)
	isIncrease = diff.GreaterThan(decimal.Zero)
	percent := diff.Div(previous).Mul(decimal.NewFromInt(100))

	return isIncrease, percent.InexactFloat64()
}

func (s *CashflowService) getCashflowHistoryInMonth(locationId uint64, year uint64, month enum.Month) (entity.CashflowHistory, error) {
	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))

	totalIncome := decimal.Zero
	warehouseSalePayments, err := s.repository.GetWarehouseSalePayments(dto.GetWarehouseSalePaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get warehouse sale payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range warehouseSalePayments {
		totalIncome = totalIncome.Add(e.Nominal)
	}

	storeSalePayments, err := s.repository.GetStoreSalePayments(dto.GetStoreSalePaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get store sale payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range storeSalePayments {
		totalIncome = totalIncome.Add(e.Nominal)
	}

	afkirChickenSalePayments, err := s.repository.GetAfkirChickenSalePayments(dto.GetAfkirChickenSalePaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get afkir chicken sale payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range afkirChickenSalePayments {
		totalIncome = totalIncome.Add(e.Nominal)
	}

	userCashAdvancePayments, err := s.repository.GetUserCashAdvancePayments(dto.GetUserCashAdvancePaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get user cash advance payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range userCashAdvancePayments {
		totalIncome = totalIncome.Add(e.Nominal)
	}

	totalExpense := decimal.Zero
	chickenProcurementPayments, err := s.repository.GetChickenProcurementPayments(dto.GetChickenProcurementPaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get chicken procurement payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range chickenProcurementPayments {
		totalExpense = totalExpense.Add(e.Nominal)
	}

	isPaid := true
	userSalaryPayments, err := s.repository.GetUserSalaryPayments(dto.GetUserSalaryPaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		IsPaid:     &isPaid,
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get user salary payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range userSalaryPayments {
		totalExpense = totalExpense.Add(e.BaseSalary).
			Add(e.AdditionalWorkSalary).
			Add(e.BonusSalary).
			Add(e.CompentationSalary).
			Sub(e.Cashbond)
	}

	warehouseItemProcurementPayments, err := s.repository.GetWarehouseItemProcurementPayments(dto.GetWarehouseItemProcurementPaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get warehouse item procurement payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range warehouseItemProcurementPayments {
		totalExpense = totalExpense.Add(e.Nominal)
	}

	warehouseItemCornProcurementPayments, err := s.repository.GetWarehouseItemCornProcurementPayments(dto.GetWarehouseItemCornProcurementPaymentFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement payments", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range warehouseItemCornProcurementPayments {
		totalExpense = totalExpense.Add(e.Nominal)
	}

	expensePayments, err := s.repository.GetExpenses(dto.GetExpenseFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get expenses", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range expensePayments {
		totalExpense = totalExpense.Add(e.Nominal)
	}

	userCashAdvances, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		s.log.Error("failed get user cash advance", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, e := range userCashAdvances {
		totalExpense = totalExpense.Add(e.Nominal)
	}

	totalReceivables := decimal.Zero
	totalEggStoreSale := decimal.Zero
	totalWarehouseStoreSale := decimal.Zero

	storeSales, err := s.repository.GetStoreSaleCashflows(dto.GetStoreSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get store sale cashflows", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, storeSale := range storeSales {
		totalNominal := storeSale.TotalPrice
		for _, payment := range storeSale.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalReceivables = totalReceivables.Add(totalNominal)
		totalEggStoreSale = totalEggStoreSale.Add(storeSale.TotalPrice)
	}

	warehouseSales, err := s.repository.GetWarehouseSaleCashflows(dto.GetWarehouseSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get warehouse sale cashflows", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, warehouseSale := range warehouseSales {
		totalNominal := warehouseSale.TotalPrice
		for _, payment := range warehouseSale.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalReceivables = totalReceivables.Add(totalNominal)
		totalWarehouseStoreSale = totalWarehouseStoreSale.Add(warehouseSale.TotalPrice)
	}

	afkirChickenSales, err := s.repository.GetAfkirChickenSaleCashflows(dto.GetAfkirChickenSaleFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get afkir chicken sale cashflows", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, sale := range afkirChickenSales {
		totalNominal := sale.TotalPrice
		for _, payment := range sale.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalReceivables = totalReceivables.Add(totalNominal)
	}

	userCashAdvanceReceivables, err := s.repository.GetUserCashAdvances(dto.GetUserCashAdvanceFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get user cash advances", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, adv := range userCashAdvanceReceivables {
		totalNominal := adv.Nominal
		for _, payment := range adv.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalReceivables = totalReceivables.Add(totalNominal)
	}

	totalDebt := decimal.Zero

	warehouseItemProcurements, err := s.repository.GetWarehouseItemProcurementCashflows(dto.GetWarehouseItemProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get warehouse item procurements", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, procurement := range warehouseItemProcurements {
		totalNominal := procurement.TotalPrice
		for _, payment := range procurement.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalDebt = totalDebt.Add(totalNominal)
	}

	warehouseItemCornProcurements, err := s.repository.GetWarehouseItemCornProcurementCashflows(dto.GetWarehouseItemCornProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get warehouse item corn procurements", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, procurement := range warehouseItemCornProcurements {
		totalNominal := procurement.TotalPrice
		for _, payment := range procurement.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalDebt = totalDebt.Add(totalNominal)
	}

	chickenProcurements, err := s.repository.GetChickenProcurementCashflows(dto.GetChickenProcurementFilter{
		DeadlinePaymentStartDate: param.DateParam(startDate),
		DeadlinePaymentEndDate:   param.DateParam(endDate),
		LocationId:               locationId,
	})
	if err != nil {
		s.log.Error("failed get chicken procurements", zap.Error(err))
		return entity.CashflowHistory{}, err
	}
	for _, procurement := range chickenProcurements {
		totalNominal := procurement.TotalPrice
		for _, payment := range procurement.Payments {
			totalNominal = totalNominal.Sub(payment.Nominal)
		}
		totalDebt = totalDebt.Add(totalNominal)
	}

	return entity.CashflowHistory{
		LocationId:       locationId,
		Income:           totalIncome,
		Expense:          totalExpense,
		Receivables:      totalReceivables,
		Debt:             totalDebt,
		Cash:             totalIncome.Add(totalReceivables),
		Profit:           totalIncome.Add(totalReceivables).Sub(totalExpense.Add(totalDebt)),
		WarehouseEggSale: totalWarehouseStoreSale,
		StoreEggSale:     totalEggStoreSale,
		CreatedAt:        time.Now(),
	}, nil
}
