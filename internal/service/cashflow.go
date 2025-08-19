package service

import (
	"fmt"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
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

	if filter.IncomeCategory == constant.IncomeCategoryAll || filter.IncomeCategory == constant.IncomeCategoryStoreEggSale {
		for _, payment := range storeSalePayments {
			incomeResponses = append(incomeResponses, dto.IncomeListResponse{
				ParentId:     payment.StoreSaleId,
				Id:           payment.Id,
				Date:         payment.PaymentDate.Format("2006-01-02"),
				PlaceName:    payment.StoreSale.Store.Name,
				Category:     "Store Egg Sale",
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
				Date:         payment.PaymentDate.Format("2006-01-02"),
				PlaceName:    payment.WarehouseSale.Warehouse.Name,
				Category:     "Warehouse Egg Sale",
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
				Date:         payment.PaymentDate.Format("2006-01-02"),
				PlaceName:    payment.AfkirChickenSale.ChickenCage.Cage.Name, // adjust field
				Category:     "Afkir Chicken Sale",
				ItemName:     "Afkir Chicken",
				ItemUnit:     "Ekor",
				Quantity:     fmt.Sprintf("%v", payment.AfkirChickenSale.TotalSellChicken),
				CustomerName: payment.AfkirChickenSale.AfkirChickenCustomer.Name,
				Nominal:      payment.Nominal.String(),
				PaymentProof: payment.PaymentProof,
			})
		}
	}

	totalIncome := len(warehouseSalePayments) + len(storeSalePayments) + len(afkirChickenSalePayments)
	return dto.IncomeOverviewResponse{
		IncomePies: dto.IncomePieResponse{
			WarehouseEggSalePercentage: float64(len(warehouseSalePayments)) / float64(totalIncome) * 100.0,
			StoreEggSalePercentage:     float64(len(storeSalePayments)) / float64(totalIncome) * 100.0,
			AfkirChickenSalePercentage: float64(len(afkirChickenSalePayments)) / float64(totalIncome) * 100.0,
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
			Date:                payment.PaymentDate.Format("2006-01-02"),
			Time:                payment.PaymentDate.Format("15:04:05"),
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
			Date:                payment.PaymentDate.Format("2006-01-02"),
			Time:                payment.PaymentDate.Format("15:04:05"),
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
			Date:                payment.PaymentDate.Format("2006-01-02"),
			Time:                payment.PaymentDate.Format("15:04:05"),
			Category:            "Afkir Chicken Sale",
			PlaceName:           payment.AfkirChickenSale.ChickenCage.Cage.Name, // adjust if wrong
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
