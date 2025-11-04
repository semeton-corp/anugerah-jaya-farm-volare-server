package service

import (
	"math"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
)

type GeneralService struct {
	log              *zap.Logger
	eggService       IEggService
	storeService     IStoreService
	warehouseService IWarehouseService
	chickenService   IChickenService
	cageService      ICageService
	cashflowService  ICashflowService
}

type IGeneralService interface {
	GetGeneralOverview() (dto.GeneralOverview, error)
}

func NewGeneralService(log *zap.Logger, eggService IEggService, storeService IStoreService, warehouseService IWarehouseService, chickenService IChickenService, cageService ICageService, cashflowService ICashflowService) IGeneralService {
	return &GeneralService{
		log:              log,
		eggService:       eggService,
		storeService:     storeService,
		warehouseService: warehouseService,
		chickenService:   chickenService,
		cageService:      cageService,
		cashflowService:  cashflowService,
	}
}

func (s *GeneralService) GetGeneralOverview() (dto.GeneralOverview, error) {
	eggMonitorings, err := s.eggService.GetEggMonitorings(dto.GetEggMonitoringFilter{
		Date: param.DateParam(time.Now()),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	storeSales, err := s.storeService.GetStoreSales(dto.GetStoreSaleFilter{
		Date: param.DateParam(time.Now().AddDate(0, 0, -1)),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	warehouseSales, err := s.warehouseService.GetWarehouseSales(dto.GetWarehouseSaleFilter{
		Date: param.DateParam(time.Now()),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	chickenMonitorings, err := s.chickenService.GetChickenMonitorings(dto.GetChickenMonitoringFilter{
		Date: param.DateParam(time.Now()),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	chickenCages, err := s.cageService.GetChickenCages(dto.GetChickenCageFilter{})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	warehouseItems, err := s.warehouseService.GetWarehouseItems(dto.GetWarehouseItemFilter{})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	storeItems, err := s.storeService.GetStoreItems(dto.GetStoreItemFilter{})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 23, 59, 59, 0, time.Local)
	startDate := endDate.AddDate(0, 0, -6).Truncate(24 * time.Hour)

	storeSaleInAWeek, err := s.storeService.GetStoreSales(dto.GetStoreSaleFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	warehouseSaleInAWeek, err := s.warehouseService.GetWarehouseSales(dto.GetWarehouseSaleFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	eggMonitoringInAWeek, err := s.eggService.GetEggMonitorings(dto.GetEggMonitoringFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	goodEggInButir := uint64(0)
	goodEggInKg := float64(0)
	for _, eggMonitoring := range eggMonitorings {
		goodEggInButir += eggMonitoring.TotalGoodEgg
		goodEggInKg += eggMonitoring.TotalWeightGoodEgg
	}

	goodEggSaleInKg := float64(0)
	for _, storeSale := range storeSales.StoreSales {
		if storeSale.SaleUnit == enum.SaleUnitIkat.String() {
			goodEggSaleInKg += storeSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if storeSale.SaleUnit == enum.SaleUnitKg.String() {
			goodEggSaleInKg += storeSale.Quantity
		}
	}

	for _, warehouseSale := range warehouseSales.WarehouseSales {
		if warehouseSale.SaleUnit == enum.SaleUnitIkat.String() {
			goodEggSaleInKg += warehouseSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if warehouseSale.SaleUnit == enum.SaleUnitKg.String() {
			goodEggSaleInKg += warehouseSale.Quantity
		}
	}

	totalLiveChicken := uint64(0)
	totalDeatchChicken := uint64(0)
	totalSickChicken := uint64(0)

	chickenMonitoringMap := make(map[uint64]dto.ChickenMonitoringListResponse)
	for _, chickenMonitoring := range chickenMonitorings {
		chickenMonitoringMap[chickenMonitoring.ChickenCage.Id] = chickenMonitoring
	}

	for _, chickenCage := range chickenCages {
		chickenMonitoring, exists := chickenMonitoringMap[chickenCage.Id]
		if !exists {
			totalLiveChicken += chickenCage.TotalChicken
			continue
		}

		totalSickChicken += chickenMonitoring.TotalSickChicken
		totalLiveChicken += chickenMonitoring.TotalLiveChicken
		totalDeatchChicken += chickenMonitoring.TotalDeathChicken
	}

	graphs := make([]dto.ProductionAndSaleEggGraphResponse, 0)
	for day := startDate; !day.After(endDate); day = day.AddDate(0, 0, 1) {
		var production, sale float64
		for _, ss := range storeSaleInAWeek.StoreSales {
			if util.IsSameDate(day, ss.CreatedAt) {
				sale += ss.Quantity
			}
		}

		for _, ws := range warehouseSaleInAWeek.WarehouseSales {
			if util.IsSameDate(day, ws.CreatedAt) {
				sale += ws.Quantity
			}
		}

		for _, em := range eggMonitoringInAWeek {
			if util.IsSameDate(day, em.CreatedAt) {
				production += em.TotalWeightGoodEgg
			}
		}

		graphs = append(graphs, dto.ProductionAndSaleEggGraphResponse{
			Key:        day.Format("02-01-2006"),
			Production: production,
			Sale:       sale,
		})
	}

	income, err := s.cashflowService.GetTotalIncomeProductionInDay(time.Now())
	if err != nil {
		return dto.GeneralOverview{}, err
	}
	expense, err := s.cashflowService.GetTotalExpenseProductionInDay(time.Now())
	if err != nil {
		return dto.GeneralOverview{}, err
	}

	totalSafeStockStoreItem := uint64(0)
	totalNotSafeStockStoreItem := uint64(0)

	for _, e := range storeItems {
		if e.Description == constant.StoreItemDescriptionSafe {
			totalSafeStockStoreItem += 1
		} else {
			totalNotSafeStockStoreItem += 1
		}
	}

	totalSafeStockWarehouseItem := uint64(0)
	totalNotSafeStockWarehouseItem := uint64(0)
	for _, e := range warehouseItems {
		if e.Description == constant.WarehouseItemDescriptionSafe {
			totalSafeStockWarehouseItem += 1
		} else {
			totalNotSafeStockWarehouseItem += 1
		}
	}

	eggSummary := dto.EggSummaryResponse{
		TotalGoodEggProductionInIkat:   math.Floor(goodEggInKg / float64(constant.TotalEggPerIkat)),
		TotalGoodEggProductionInKg:     goodEggInKg,
		TotalGoodEggProductionInKarpet: math.Ceil(float64(goodEggInButir) / float64(constant.TotalEggPerKarpet)),
		TotalGoodEggProductionInButir:  float64(goodEggInButir),
		TotalGoodEggSaleInIkat:         math.Floor(goodEggSaleInKg / float64(constant.TotalEggPerIkat)),
		TotalGoodEggSaleInKg:           goodEggSaleInKg,
	}

	warehouseItemSummary := dto.WarehouseItemSummaryResponse{
		TotalSafeItem:    totalSafeStockWarehouseItem,
		TotalNotSafeItem: totalNotSafeStockWarehouseItem,
	}

	storeItemSummary := dto.StoreItemSummaryResponse{
		TotalSafeItem:    totalSafeStockStoreItem,
		TotalNotSafeItem: totalNotSafeStockStoreItem,
	}

	saleSummary := dto.SaleSummaryResponse{
		Income: income.String(),
		Profit: income.Sub(expense).String(),
	}
	chickenSummary := dto.ChickenSummaryResponse{
		TotalLiveChicken:  totalLiveChicken,
		TotalSickChicken:  totalSickChicken,
		TotalDeathChicken: totalDeatchChicken,
	}

	return dto.GeneralOverview{
		EggSummary:                 eggSummary,
		WarehouseItemSummary:       warehouseItemSummary,
		StoreItemSummary:           storeItemSummary,
		SaleSummary:                saleSummary,
		ChickenSummary:             chickenSummary,
		ProductionAndSaleEggGraphs: graphs,
	}, nil
}
