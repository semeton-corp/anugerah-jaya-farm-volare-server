package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/cache"
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
	"go.uber.org/zap"
)

type ChickenService struct {
	log              *zap.Logger
	repository       repository.IChickenRepository
	eggService       IEggService
	cageService      ICageService
	cacheService     cache.ICache
	warehouseService IWarehouseService
	itemService      IItemService
	cashflowService  ICashflowService
}

type IChickenService interface {
	CreateChickenMonitoring(request dto.CreateChickenMonitoringRequest, accoundId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	GetChickenMonitorings(filter dto.GetChickenMonitoringFilter) ([]dto.ChickenMonitoringListResponse, error)
	GetChickenMonitoringById(id uint64) (dto.ChickenMonitoringResponse, error)
	UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, accoundId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	DeleteChickenMonitoring(id uint64) error

	CreateChickenHealthItem(request dto.CreateChickenHealthItemRequest, createdBy uuid.UUID) (dto.ChickenHealthItemResponse, error)
	GetChickenHealthItems(filter dto.GetChickenHealthItemFilter) ([]dto.ChickenHealthItemResponse, error)
	GetChickenHealthItemById(id uint64) (dto.ChickenHealthItemResponse, error)
	UpdateChickenHealthItem(id uint64, request dto.UpdateChickenHealthItemRequest, updatedBy uuid.UUID) (dto.ChickenHealthItemResponse, error)
	DeleteChickenHealthItem(id uint64) error

	CreateChickenHealthMonitoring(request dto.CreateChickenHealthMonitoringRequest, createdBy uuid.UUID) (dto.ChickenHealthMonitoringResponse, error)
	UpdateChickenHealthMonitoring(id uint64, request dto.UpdateChickenHealthMonitoringRequest, updatedBy uuid.UUID) (dto.ChickenHealthMonitoringResponse, error)
	DeleteChickenHealthMonitoring(id uint64) error
	GetChickenHealthMonitoringDetails(chickenCageId uint64) (dto.ChickenHealthMonitoringDetailResponse, error)
	GetChickenHealthMonitoringById(id uint64) (dto.ChickenHealthMonitoringResponse, error)

	GetChickenOverview(filter dto.GetChickenOverviewFilter) (dto.ChickenOverviewResponse, error)

	CreateChickenProcurementDraft(request dto.CreateChickenProcurementDraftRequest, userId uuid.UUID) (dto.ChickenProcurementDraftResponse, error)
	UpdateChickenProcurementDraft(id uint64, request dto.UpdateChickenProcurementDraftRequest, userId uuid.UUID) (dto.ChickenProcurementDraftResponse, error)
	GetChickenProcurementDrafts() ([]dto.ChickenProcurementDraftResponse, error)
	GetChickenProcurementDraft(id uint64) (dto.ChickenProcurementDraftResponse, error)
	DeleteChickenProcurementDraft(id uint64) error
	ConfirmationChickenProcurementDraft(id uint64, request dto.ConfirmationChickenProcurementRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error)

	GetChickenProcurements(filter dto.GetChickenProcurementFilter) (dto.ChickenProcurementListPaginationResponse, error)
	GetChickenProcurement(id uint64) (dto.ChickenProcurementResponse, error)

	ArrivalConfirmationChickenProcurement(id uint64, request dto.ArrivalConfirmationChickenProcurementRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error)

	CreateChickenProcurementPayment(chickenProcurementId uint64, request dto.CreateChickenProcurementPaymentRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error)
	UpdateChickenProcurementPayment(chickenProcurementId uint64, id uint64, request dto.UpdateChickenProcurementPaymentRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error)
	DeleteChickenProcurementPayment(chickenProcurementId uint64, id uint64, userId uuid.UUID) error

	CreateAfkirChickenCustomer(request dto.CreateAfkirChickenCustomerRequest, userId uuid.UUID) (dto.AfkirChickenCustomerResponse, error)
	GetAfkirChickenCustomers() ([]dto.AfkirChickenCustomerListResponse, error)
	GetAfkirChickenCustomer(id uint64) (dto.AfkirChickenCustomerResponse, error)
	UpdateAfkirChickenCustomer(id uint64, request dto.UpdateAfkirChickenCustomerRequest, userId uuid.UUID) (dto.AfkirChickenCustomerResponse, error)
	DeleteAfkirChickenCustomer(id uint64) error

	CreateAkfirChickenSaleDraft(request dto.CreateAfkirChickenSaleDraftRequest, userId uuid.UUID) (dto.AfkirChickenSaleDraftResponse, error)
	GetAfkirChickenSaleDrafts() ([]dto.AfkirChickenSaleDraftResponse, error)
	GetAfkirChickenSaleDraft(id uint64) (dto.AfkirChickenSaleDraftResponse, error)
	UpdateAfkirChickenSaleDraft(id uint64, request dto.UpdateAfkirChickenSaleDraftRequest, userId uuid.UUID) (dto.AfkirChickenSaleDraftResponse, error)
	DeleteAfkirChickenSaleDraft(id uint64) error
	ConfirmationAfkirChickenSaleDraft(id uint64, request dto.CreateAfkirChickenSaleRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error)

	CreateAfkirChickenSale(request dto.CreateAfkirChickenSaleRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error)
	GetAfkirChickenSales(filter dto.GetAfkirChickenSaleFilter) (dto.AfkirChickenSaleListPaginationResponse, error)
	GetAkfirChickenSale(id uint64) (dto.AfkirChickenSaleResponse, error)

	CreateAfkirChickenSalePayment(afkirChickenSaleId uint64, request dto.CreateAfkirChickenSalePaymentRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error)
	UpdateAfkirChickenSalePayment(afkirChickenSaleId uint64, id uint64, request dto.UpdateAfkirChickenSalePaymentRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error)
	DeleteAfkirChickenSalePayment(afkirChickenSaleId uint64, id uint64) error

	GetChickenPerformances(filter dto.GetChickenPerformanceFilter) ([]dto.ChickenPerformanceResponse, error)

	GetKPIScoreChickenInMonth(locationId uint64, month enum.Month, year uint64) (float64, error)
	GetKPIScoreChickenPerWeek(locationId uint64, month enum.Month, year uint64) (map[int]float64, error)

	GetChickenAndWarehouseOverview(filter dto.GetChickenAndWarehouseOverviewFilter) (dto.ChickenAndWarehouseOverviewResponse, error)
	GetChickenAndCompanyOverview(filter dto.GetChickenAndCompanyOverviewFilter) (dto.ChickenAndCompanyOverviewResponse, error)
}

func NewChickenService(log *zap.Logger, repository repository.IChickenRepository, eggService IEggService, cageService ICageService, itemService IItemService, cashflowService ICashflowService, warehouseService IWarehouseService, cacheService cache.ICache) IChickenService {
	return &ChickenService{
		log:              log,
		repository:       repository,
		eggService:       eggService,
		cageService:      cageService,
		itemService:      itemService,
		warehouseService: warehouseService,
		cashflowService:  cashflowService,
		cacheService:     cacheService,
	}
}

func (s *ChickenService) CreateChickenMonitoring(request dto.CreateChickenMonitoringRequest, userId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	s.repository.UseTx(false)

	count, err := s.repository.CountChickenMonitoringByChickenCageIdToday(request.ChickenCageId)
	if err != nil {
		s.log.Error("failed to count chicken monitoring by cage id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	if count > 0 {
		return dto.ChickenMonitoringResponse{}, errx.BadRequest("chicken monitoring already exists for today")
	}

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring := entity.ChickenMonitoring{
		ChickenCageId:     request.ChickenCageId,
		TotalChicken:      chickenCage.TotalChicken,
		TotalDeathChicken: request.TotalDeathChicken,
		TotalSickChicken:  request.TotalSickChicken,
		TotalFeed:         request.TotalFeed,
		Note:              request.Note,
		CreatedBy:         uuid.NullUUID{UUID: userId, Valid: true},
	}

	if chickenCage.TotalChicken < request.TotalSickChicken {
		return dto.ChickenMonitoringResponse{}, errx.BadRequest("total chicken is less than total sick chicken")
	} else if chickenCage.TotalChicken < request.TotalDeathChicken {
		return dto.ChickenMonitoringResponse{}, errx.BadRequest("total chicken is less than total death chicken")
	}

	currentChicken := chickenCage.TotalChicken - request.TotalDeathChicken
	_, err = s.cageService.UpdateChickenCage(request.ChickenCageId, dto.UpdateChickenCageRequest{
		TotalChicken: currentChicken,
	}, userId)
	if err != nil {
		return dto.ChickenMonitoringResponse{}, err
	}

	err = s.repository.CreateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		s.log.Error("failed to create chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	mortalityRate := 0.0
	if (chickenMonitoring.TotalChicken) == 0 {
		mortalityRate = 0
	} else {
		mortalityRate = float64((chickenMonitoring.TotalDeathChicken / (chickenMonitoring.TotalChicken)) * 100.0)
	}

	chickenMonitoring, err = s.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	if mortalityRate < 0.5 {
		notificationJsonParsed, err := json.Marshal(entity.Notification{
			CageId:              sql.NullInt64{Int64: int64(chickenMonitoring.ChickenCage.CageId), Valid: true},
			NotificationContext: pq.StringArray{constant.ChickenMonitoringNotificationContext},
			Description:         fmt.Sprintf(constant.ChickenStatusNotification, chickenMonitoring.ChickenCage.Cage.Name, mortalityRate),
		})
		if err != nil {
			s.log.Error("failed to parse struct into json", zap.Error(err))
			return dto.ChickenMonitoringResponse{}, errx.BadRequest("failed parsed struct into json")
		}

		s.cacheService.Publish(context.Background(), constant.NotificationTopic, string(notificationJsonParsed))
	}

	return mapper.ChickenMonitoringToResponse(&chickenMonitoring), nil
}

func (s *ChickenService) GetChickenMonitoringById(id uint64) (dto.ChickenMonitoringResponse, error) {
	chickenMonitoring, err := s.repository.GetChickenMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	return mapper.ChickenMonitoringToResponse(&chickenMonitoring), nil
}

func (s *ChickenService) GetChickenMonitorings(filter dto.GetChickenMonitoringFilter) ([]dto.ChickenMonitoringListResponse, error) {
	chickenMonitorings, err := s.repository.GetChickenMonitorings(&filter)
	if err != nil {
		s.log.Error("failed to get chicken monitorings", zap.Error(err))
		return []dto.ChickenMonitoringListResponse{}, err
	}

	chickenMonitoringsResponse := make([]dto.ChickenMonitoringListResponse, len(chickenMonitorings))
	for i, chickenMonitoring := range chickenMonitorings {
		chickenMonitoringsResponse[i] = mapper.ChickenMonitoringToListResponse(&chickenMonitoring)
	}

	return chickenMonitoringsResponse, nil
}

func (s *ChickenService) UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, userId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	s.repository.UseTx(false)
	chickenMonitoring, err := s.repository.GetChickenMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring.TotalSickChicken = request.TotalSickChicken
	chickenMonitoring.TotalDeathChicken = request.TotalDeathChicken
	chickenMonitoring.TotalFeed = request.TotalFeed
	chickenMonitoring.Note = request.Note
	chickenMonitoring.UpdateBy = uuid.NullUUID{UUID: userId, Valid: true}
	chickenMonitoring.ChickenCageId = request.ChickenCageId

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.ChickenMonitoringResponse{}, err
	}

	currentChicken := chickenCage.TotalChicken - request.TotalDeathChicken
	_, err = s.cageService.UpdateChickenCage(request.ChickenCageId, dto.UpdateChickenCageRequest{
		TotalChicken: currentChicken,
	}, userId)
	if err != nil {
		return dto.ChickenMonitoringResponse{}, err
	}

	err = s.repository.UpdateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		s.log.Error("failed to update chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = s.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	return mapper.ChickenMonitoringToResponse(&chickenMonitoring), nil
}

func (c *ChickenService) DeleteChickenMonitoring(id uint64) error {
	err := c.repository.DeleteChickenMonitoring(id)
	if err != nil {
		c.log.Error("failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (c *ChickenService) GetChickenOverview(filter dto.GetChickenOverviewFilter) (dto.ChickenOverviewResponse, error) {
	c.repository.UseTx(false)

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	currentChickenMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		Date:       param.DateParam(today),
		LocationId: filter.LocationId,
		CageId:     filter.CageId,
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return dto.ChickenOverviewResponse{}, err
	}

	currentEggMonitoring, err := c.eggService.GetEggMonitorings(dto.GetEggMonitoringFilter{
		Date:       param.DateParam(today),
		LocationId: filter.LocationId,
		CageId:     filter.CageId,
	})
	if err != nil {
		c.log.Error("failed to get egg monitorings", zap.Error(err))
		return dto.ChickenOverviewResponse{}, err
	}

	var totalEgg uint64
	for _, eggMonitoring := range currentEggMonitoring {
		totalEgg += eggMonitoring.TotalAllEgg
	}

	var totalDOCChicken, totalGrowerChicken, totalPreLayerChicken, totalLayerChicken, totalAfkirChicken uint64
	var totalLiveChicken, totalSickChicken, totalDeathChicken uint64

	for _, cm := range currentChickenMonitorings {
		totalLiveChicken += cm.ChickenCage.TotalChicken - cm.TotalSickChicken
		totalSickChicken += cm.TotalSickChicken
		totalDeathChicken += cm.TotalDeathChicken

		count := cm.ChickenCage.TotalChicken
		switch cm.ChickenCage.Cage.ChickenCategory {
		case enum.ChickenCategoryDOC:
			totalDOCChicken += count
		case enum.ChickenCategoryGrower:
			totalGrowerChicken += count
		case enum.ChickenCategoryPreLayer:
			totalPreLayerChicken += count
		case enum.ChickenCategoryLayer:
			totalLayerChicken += count
		case enum.ChickenCategoryAfkir:
			totalAfkirChicken += count
		}
	}

	chickenGraphs := make([]dto.ChickenGraphResponse, 0)

	switch filter.OverviewGraphTime.Value() {
	case enum.OverviewGraphTimeThisWeek:
		chickenGraphs, err = c.buildWeeklyGraph(filter.LocationId, filter.CageId)
	case enum.OverviewGraphTimeThisMonth:
		chickenGraphs, err = c.buildMonthlyGraph(filter.LocationId, filter.CageId)
	case enum.OverviewGraphTimeThisYear:
		chickenGraphs, err = c.buildYearlyGraph(filter.LocationId, filter.CageId)
	}
	if err != nil {
		return dto.ChickenOverviewResponse{}, err
	}

	totalChicken := float64(totalLiveChicken + totalDeathChicken + totalSickChicken)
	mortalityRate := float64(0)
	hdpRate := float64(0)
	kpiChicken := float64(0)

	if totalChicken != 0 {
		mortalityRate = float64(totalDeathChicken) / totalChicken
		hdpRate = float64(totalEgg) / (totalChicken - float64(totalDeathChicken))
		kpiChicken = (mortalityRate + hdpRate) / 2
	}

	return dto.ChickenOverviewResponse{
		ChickenDetail: dto.ChickenDetailOverview{
			TotalLiveChicken:    totalLiveChicken,
			TotalSickChicken:    totalSickChicken,
			TotalDeathChicken:   totalDeathChicken,
			TotalKPIPerformance: kpiChicken,
		},
		ChickenPie: dto.ChickenBarChartResponse{
			ChickenDOC:       float64(totalDOCChicken),
			ChickenGrower:    float64(totalGrowerChicken),
			ChickentPreLayer: float64(totalPreLayerChicken),
			ChickenLayer:     float64(totalLayerChicken),
			ChickenAfkir:     float64(totalAfkirChicken),
		},
		ChickenGraphs: chickenGraphs,
	}, nil
}

func (c *ChickenService) GetKPIScoreChickenInMonth(locationId uint64, month enum.Month, year uint64) (float64, error) {
	c.repository.UseTx(false)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))
	chickenMonitoringInMonth, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return -1, err
	}

	eggMonitoringInMonth, err := c.eggService.GetEggMonitorings(dto.GetEggMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		c.log.Error("failed to get egg monitorings", zap.Error(err))
		return -1, err
	}

	totalChicken := uint64(0)
	totalDeathChicken := uint64(0)
	totalEgg := uint64(0)

	for _, chickenMonitoring := range chickenMonitoringInMonth {
		totalChicken += chickenMonitoring.TotalChicken
		totalDeathChicken += chickenMonitoring.TotalDeathChicken
	}

	for _, eggMonitoring := range eggMonitoringInMonth {
		totalEgg += eggMonitoring.TotalAllEgg
	}

	mortalityRate := float64(0)
	hdpRate := float64(0)
	kpiChicken := float64(0)

	if totalChicken != 0 || totalChicken-totalDeathChicken > 0 {
		mortalityRate = float64(totalDeathChicken) / float64(totalChicken)
		hdpRate = float64(totalEgg) / float64((totalChicken - totalDeathChicken))
		kpiChicken = ((mortalityRate + hdpRate) / 2) * 100
	}

	return kpiChicken, nil
}

func (c *ChickenService) GetKPIScoreChickenPerWeek(locationId uint64, month enum.Month, year uint64) (map[int]float64, error) {
	c.repository.UseTx(false)

	weekRanges := util.GetFourWeekRanges(int(year), time.Month(month))

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(year), time.Month(month))
	chickenMonitoringInMonth, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return nil, err
	}

	eggMonitoringInMonth, err := c.eggService.GetEggMonitorings(dto.GetEggMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
	})
	if err != nil {
		c.log.Error("failed to get egg monitorings", zap.Error(err))
		return nil, err
	}

	totalChickenInWeek := make(map[int]uint64)
	totalDeathChickenInWeek := make(map[int]uint64)
	totalEggInWeek := make(map[int]uint64)
	for _, chickenMonitoring := range chickenMonitoringInMonth {
		week := util.FindWeek(chickenMonitoring.CreatedAt, weekRanges)
		if week == 0 {
			continue
		}

		totalChickenInWeek[week] += chickenMonitoring.TotalChicken
		totalDeathChickenInWeek[week] += chickenMonitoring.TotalDeathChicken
	}

	for _, eggMonitoring := range eggMonitoringInMonth {
		week := util.FindWeek(eggMonitoring.CreatedAt, weekRanges)
		if week == 0 {
			continue
		}
		totalEggInWeek[week] += eggMonitoring.TotalAllEgg
	}

	kpiChickenInWeek := make(map[int]float64)
	keys := util.GetSortedKeys(weekRanges)

	for key := range keys {
		if totalChickenInWeek[key] == 0 {
			continue
		}

		mortalityRate := float64(totalDeathChickenInWeek[key]) / float64(totalChickenInWeek[key])
		hdpRate := float64(totalEggInWeek[key]) / float64((totalChickenInWeek[key] - totalDeathChickenInWeek[key]))
		kpiChickenInWeek[key] = ((mortalityRate + hdpRate) / 2) * 100
	}

	return kpiChickenInWeek, nil
}

func (c *ChickenService) buildWeeklyGraph(locationId uint64, cageId uint64) ([]dto.ChickenGraphResponse, error) {
	endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	startDate := endDate.AddDate(0, 0, -6)

	weekMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
		CageId:     cageId,
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return nil, err
	}

	graphs := make([]dto.ChickenGraphResponse, 0)
	for day := startDate; !day.After(endDate); day = day.AddDate(0, 0, 1) {
		var sickSum, deathSum uint64
		for _, cm := range weekMonitorings {
			if util.IsSameDate(day, cm.CreatedAt) {
				sickSum += cm.TotalSickChicken
				deathSum += cm.TotalDeathChicken
			}
		}
		graphs = append(graphs, dto.ChickenGraphResponse{
			Key:          day.Format("02 Jan 2006"),
			SickChicken:  sickSum,
			DeathChicken: deathSum,
		})
	}
	return graphs, nil
}

func (c *ChickenService) buildMonthlyGraph(locationId uint64, cageId uint64) ([]dto.ChickenGraphResponse, error) {
	weekMaps := util.GetFourWeekRanges(time.Now().Year(), time.Now().Month())
	startDate, endDate := util.GetStartDateAndEndDateInMonth(time.Now().Year(), time.Now().Month())

	monthMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
		CageId:     cageId,
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return nil, err
	}

	totalSick, totalDeath := make(map[int]uint64), make(map[int]uint64)
	for _, cm := range monthMonitorings {
		week := util.FindWeek(cm.CreatedAt, weekMaps)
		if week > 0 {
			totalSick[week] += cm.TotalSickChicken
			totalDeath[week] += cm.TotalDeathChicken
		}
	}

	keys := util.GetSortedKeys(weekMaps)
	graphs := make([]dto.ChickenGraphResponse, 0)
	for _, k := range keys {
		graphs = append(graphs, dto.ChickenGraphResponse{
			Key:          fmt.Sprintf("Minggu %d", k),
			SickChicken:  totalSick[k],
			DeathChicken: totalDeath[k],
		})
	}

	return graphs, nil
}

func (c *ChickenService) buildYearlyGraph(locationId uint64, cageId uint64) ([]dto.ChickenGraphResponse, error) {
	monthMaps := util.GetTwelveMonthRanges(time.Now().Year())
	startDate, endDate := util.GetStartDateAndEndDateInYear(time.Now().Year())

	yearMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate:  param.DateParam(startDate),
		EndDate:    param.DateParam(endDate),
		LocationId: locationId,
		CageId:     cageId,
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return nil, err
	}

	totalSick, totalDeath := make(map[int]uint64), make(map[int]uint64)
	for _, cm := range yearMonitorings {
		month := util.FindMonth(cm.CreatedAt, monthMaps)
		if month > 0 {
			totalSick[month] += cm.TotalSickChicken
			totalDeath[month] += cm.TotalDeathChicken
		}
	}

	keys := util.GetSortedKeys(monthMaps)
	graphs := make([]dto.ChickenGraphResponse, 0)
	for _, k := range keys {
		graphs = append(graphs, dto.ChickenGraphResponse{
			Key:          util.IndoMonthName(k),
			SickChicken:  totalSick[k],
			DeathChicken: totalDeath[k],
		})
	}
	return graphs, nil
}

func (s *ChickenService) CreateChickenHealthItem(request dto.CreateChickenHealthItemRequest, createdBy uuid.UUID) (dto.ChickenHealthItemResponse, error) {
	s.repository.UseTx(false)

	chickenHealthitemType := enum.ValueOfChickenHealthItemType(request.Type)
	if !chickenHealthitemType.IsValid() {
		s.log.Error("invalid chicken health type item", zap.String("type", request.Type))
		return dto.ChickenHealthItemResponse{}, errx.BadRequest("invalid chicken health Item")
	}

	data := entity.ChickenHealthItem{
		Name:      request.Name,
		Type:      chickenHealthitemType,
		CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
		Note:      request.Note,
	}

	if request.ChickenAge != nil {
		data.ChickenAge = sql.NullInt64{Int64: int64(*request.ChickenAge), Valid: true}
	}

	err := s.repository.CreateChickenHealthItem(&data)
	if err != nil {
		s.log.Error("failed to create chicken health item", zap.Error(err))
		return dto.ChickenHealthItemResponse{}, err
	}

	data, err = s.repository.GetChickenHealthItemById(data.Id)
	if err != nil {
		s.log.Error("failed to get chicken health item", zap.Error(err))
		return dto.ChickenHealthItemResponse{}, err
	}

	return mapper.ChickenHealthItemToResponse(&data), nil
}

func (s *ChickenService) GetChickenHealthItems(filter dto.GetChickenHealthItemFilter) ([]dto.ChickenHealthItemResponse, error) {
	s.repository.UseTx(false)
	chickenHealthItemResponses := make([]dto.ChickenHealthItemResponse, 0)
	chickenHealthItems, err := s.repository.GetChickenHealthItems(filter)
	if err != nil {
		s.log.Error("failed to get chicken health items", zap.Error(err))
		return chickenHealthItemResponses, err
	}

	for _, chickenHealthItem := range chickenHealthItems {
		chickenHealthItemResponses = append(chickenHealthItemResponses, mapper.ChickenHealthItemToResponse(&chickenHealthItem))
	}

	return chickenHealthItemResponses, nil
}

func (s *ChickenService) GetChickenHealthItemById(id uint64) (dto.ChickenHealthItemResponse, error) {
	s.repository.UseTx(false)

	chickenHealthItem, err := s.repository.GetChickenHealthItemById(id)
	if err != nil {
		s.log.Error("failed to get chicken health item by id", zap.Error(err))
		return dto.ChickenHealthItemResponse{}, err
	}

	return mapper.ChickenHealthItemToResponse(&chickenHealthItem), nil
}

func (s *ChickenService) UpdateChickenHealthItem(id uint64, request dto.UpdateChickenHealthItemRequest, updatedBy uuid.UUID) (dto.ChickenHealthItemResponse, error) {
	s.repository.UseTx(false)

	chickenHealthItem, err := s.repository.GetChickenHealthItemById(id)
	if err != nil {
		s.log.Error("failed to get chicken health item by id", zap.Error(err))
		return dto.ChickenHealthItemResponse{}, err
	}

	chickenHealthItemType := enum.ValueOfChickenHealthItemType(request.Type)
	if !chickenHealthItemType.IsValid() {
		s.log.Error("invalid chicken health item type", zap.String("type", request.Type))
		return dto.ChickenHealthItemResponse{}, errx.Unauthorized("invalid chicken item health type")
	}

	chickenHealthItem.Name = request.Name
	if request.ChickenAge != nil {
		chickenHealthItem.ChickenAge = sql.NullInt64{Int64: int64(*request.ChickenAge), Valid: true}
	}
	chickenHealthItem.Type = chickenHealthItemType
	chickenHealthItem.Note = request.Note

	err = s.repository.UpdateChickenHealthItem(&chickenHealthItem)
	if err != nil {
		s.log.Error("failed to update chicken health item", zap.Error(err))
		return dto.ChickenHealthItemResponse{}, err
	}

	chickenHealthItem, err = s.repository.GetChickenHealthItemById(id)
	if err != nil {
		s.log.Error("failed to get chicken health item by id", zap.Error(err))
		return dto.ChickenHealthItemResponse{}, err
	}

	return mapper.ChickenHealthItemToResponse(&chickenHealthItem), nil
}

func (s *ChickenService) DeleteChickenHealthItem(id uint64) error {
	s.repository.UseTx(false)
	err := s.repository.DeleteChickenHealthItem(id)
	if err != nil {
		s.log.Error("failed to delete chicken health item", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) CreateChickenHealthMonitoring(request dto.CreateChickenHealthMonitoringRequest, createdBy uuid.UUID) (dto.ChickenHealthMonitoringResponse, error) {
	s.repository.UseTx(false)

	chickenHealthMonitoringType := enum.ValueOfChickenHealthItemType(request.Type)
	if !chickenHealthMonitoringType.IsValid() {
		s.log.Warn("invalid chicken health monitoring type")
		return dto.ChickenHealthMonitoringResponse{}, errx.BadRequest("invalid chicken health monitoring type")
	}

	if chickenHealthMonitoringType == enum.ChickenHealthItemTypeMedicine && request.Disease == nil {
		return dto.ChickenHealthMonitoringResponse{}, errx.BadRequest("disease is required, since you choose medicine type")
	}

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.ChickenHealthMonitoringResponse{}, nil
	}

	data := entity.ChickenHealthMonitoring{
		ChickenCageId:  request.ChickenCageId,
		HealthItemName: request.HealthItemName,
		Type:           chickenHealthMonitoringType,
		Dose:           request.Dose,
		Unit:           request.Unit,
		ChickenAge:     chickenCage.ChickenAge,
		CreatedBy:      uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	if request.Disease != nil {
		data.Disease = sql.NullString{String: *request.Disease, Valid: true}
	}

	_, err = s.cageService.UpdateChickenCage(chickenCage.Id, dto.UpdateChickenCageRequest{
		TotalChicken:                   chickenCage.TotalChicken,
		LatestChickenAgeVaccineRoutine: &chickenCage.ChickenAge,
	}, createdBy)
	if err != nil {
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	err = s.repository.CreateChickenHealthMonitoring(&data)
	if err != nil {
		s.log.Error("failed to create chicken health monitoring", zap.Error(err))
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	data, err = s.repository.GetChickenHealthMonitoringById(data.Id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	return mapper.ChickenHealthMonitoringToResponse(&data), nil
}

func (s *ChickenService) UpdateChickenHealthMonitoring(id uint64, request dto.UpdateChickenHealthMonitoringRequest, updatedBy uuid.UUID) (dto.ChickenHealthMonitoringResponse, error) {
	s.repository.UseTx(false)

	chickenHealthMonitoring, err := s.repository.GetChickenHealthMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	chickenHealthMonitoringType := enum.ValueOfChickenHealthItemType(request.Type)
	if !chickenHealthMonitoringType.IsValid() {
		s.log.Warn("invalid chicken health monitoring type")
		return dto.ChickenHealthMonitoringResponse{}, errx.BadRequest("invalid chicken health monitoring type")
	}

	if chickenHealthMonitoringType == enum.ChickenHealthItemTypeMedicine && request.Disease == nil {
		return dto.ChickenHealthMonitoringResponse{}, errx.BadRequest("disease is required, since you choose medicine type")
	}

	chickenHealthMonitoring.ChickenCageId = request.ChickenCageId
	chickenHealthMonitoring.HealthItemName = request.HealthItemName
	chickenHealthMonitoring.Dose = request.Dose
	chickenHealthMonitoring.Unit = request.Unit
	chickenHealthMonitoring.Type = chickenHealthMonitoringType

	if request.Disease != nil {
		chickenHealthMonitoring.Disease = sql.NullString{String: *request.Disease, Valid: true}
	} else {
		chickenHealthMonitoring.Disease = sql.NullString{}
	}

	err = s.repository.UpdateChickenHealthMonitoring(&chickenHealthMonitoring)
	if err != nil {
		s.log.Error("failed to update chicken monitoring", zap.Error(err))
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	chickenHealthMonitoring, err = s.repository.GetChickenHealthMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	return mapper.ChickenHealthMonitoringToResponse(&chickenHealthMonitoring), nil
}

func (s *ChickenService) DeleteChickenHealthMonitoring(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteChickenHealthMonitoring(id)
	if err != nil {
		s.log.Error("failed to delete chicken health monitoring")
		return err
	}

	return nil
}

func (s *ChickenService) GetChickenHealthMonitoringDetails(chickenCageId uint64) (dto.ChickenHealthMonitoringDetailResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.cageService.GetChickenCageById(chickenCageId)
	if err != nil {
		s.log.Error("failed to get chicken cage by id", zap.Error(err))
		return dto.ChickenHealthMonitoringDetailResponse{}, err
	}

	chickenHealthMonitorings, err := s.repository.GetChickenHealthMonitoringByChickenCageId(chickenCageId)
	if err != nil {
		s.log.Error("failed to get chicken health monitoring by chicken cage id")
		return dto.ChickenHealthMonitoringDetailResponse{}, err
	}

	chickenHealthMonitoringResponses := make([]dto.ChickenHealthMonitoringResponse, 0)
	for _, e := range chickenHealthMonitorings {
		chickenHealthMonitoringResponses = append(chickenHealthMonitoringResponses, mapper.ChickenHealthMonitoringToResponse(&e))
	}

	return dto.ChickenHealthMonitoringDetailResponse{
		ChickenCage:              chickenCage,
		ChickenHealthMonitorings: chickenHealthMonitoringResponses,
	}, nil
}

func (s *ChickenService) GetChickenHealthMonitoringById(id uint64) (dto.ChickenHealthMonitoringResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetChickenHealthMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get chicken health monitoring by id", zap.Error(err))
		return dto.ChickenHealthMonitoringResponse{}, err
	}

	return mapper.ChickenHealthMonitoringToResponse(&data), nil
}

func (s *ChickenService) CreateChickenProcurementDraft(request dto.CreateChickenProcurementDraftRequest, userId uuid.UUID) (dto.ChickenProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	cage, err := s.cageService.GetCageById(request.CageId)
	if err != nil {
		return dto.ChickenProcurementDraftResponse{}, err
	}

	if cage.IsUsed {
		return dto.ChickenProcurementDraftResponse{}, errx.BadRequest("cage is in used by another chicken")
	}

	totalPrice, err := decimal.NewFromString(request.TotalPrice)
	if err != nil {
		s.log.Error("failed to parse price from string", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	data := entity.ChickenProcurementDraft{
		CageId:     request.CageId,
		SupplierId: request.SupplierId,
		Quantity:   request.Quantity,
		TotalPrice: totalPrice,
		CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateChickenProcurementDraft(&data)
	if err != nil {
		return dto.ChickenProcurementDraftResponse{}, err
	}

	chickenProcurementDraft, err := s.repository.GetChickenProcurementDraft(data.Id)
	if err != nil {
		s.log.Error("failed to created chicken procurement draft", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	return mapper.ChickenProcurementDraftToResponse(&chickenProcurementDraft), nil
}

func (s *ChickenService) GetChickenProcurementDraft(id uint64) (dto.ChickenProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetChickenProcurementDraft(id)
	if err != nil {
		s.log.Error("failed get chicken procurement draft", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	return mapper.ChickenProcurementDraftToResponse(&data), nil
}

func (s *ChickenService) UpdateChickenProcurementDraft(id uint64, request dto.UpdateChickenProcurementDraftRequest, userId uuid.UUID) (dto.ChickenProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	cage, err := s.cageService.GetCageById(request.CageId)
	if err != nil {
		return dto.ChickenProcurementDraftResponse{}, err
	}

	if cage.IsUsed {
		return dto.ChickenProcurementDraftResponse{}, errx.BadRequest("cage is in used by another chicken")
	}

	totalPrice, err := decimal.NewFromString(request.TotalPrice)
	if err != nil {
		s.log.Error("failed to parse price from string", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	chickenProcurementDraft, err := s.repository.GetChickenProcurementDraft(id)
	if err != nil {
		s.log.Error("failed get chicken procurement by id", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	chickenProcurementDraft.CageId = request.CageId
	chickenProcurementDraft.SupplierId = request.SupplierId
	chickenProcurementDraft.Quantity = request.Quantity
	chickenProcurementDraft.TotalPrice = totalPrice
	chickenProcurementDraft.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateChickenProcurementDraft(&chickenProcurementDraft)
	if err != nil {
		s.log.Error("failed update chicken procurement draft", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	chickenProcurementDraft, err = s.repository.GetChickenProcurementDraft(id)
	if err != nil {
		s.log.Error("failed to created chicken procurement draft", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	return mapper.ChickenProcurementDraftToResponse(&chickenProcurementDraft), nil
}

func (s *ChickenService) GetChickenProcurementDrafts() ([]dto.ChickenProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	chickenProcurementDrafts, err := s.repository.GetChickenProcurementDrafts()
	if err != nil {
		s.log.Error("failed to get chicken procurement drafts")
		return nil, err
	}

	response := make([]dto.ChickenProcurementDraftResponse, 0)
	for _, chickenProcurementDraft := range chickenProcurementDrafts {
		response = append(response, mapper.ChickenProcurementDraftToResponse(&chickenProcurementDraft))
	}

	return response, nil
}

func (s *ChickenService) DeleteChickenProcurementDraft(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteChickenProcurementDraft(id)
	if err != nil {
		s.log.Error("failed delete chicken procurement draft", zap.Error(err))
		return err
	}
	return nil
}

func (s *ChickenService) GetChickenProcurements(filter dto.GetChickenProcurementFilter) (dto.ChickenProcurementListPaginationResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetChickenProcurements(filter)
	if err != nil {
		s.log.Error("failed get chicken procurements", zap.Error(err))
		return dto.ChickenProcurementListPaginationResponse{}, err
	}

	chickenProcurementResponse := make([]dto.ChickenProcurementListResponse, 0)
	for _, e := range data {
		chickenProcurementResponse = append(chickenProcurementResponse, mapper.ChickenProcurementToListResponse(&e))
	}

	count, err := s.repository.CountChickenProcurement(filter)
	if err != nil {
		s.log.Error("failed count chicken procurement", zap.Error(err))
		return dto.ChickenProcurementListPaginationResponse{}, err
	}

	response := dto.ChickenProcurementListPaginationResponse{
		ChickenProcurements: chickenProcurementResponse,
	}

	if filter.Page > 0 {
		response.TotalData = uint64(count)
		response.TotalPage = uint64(math.Ceil(float64(count) / float64(constant.PaginationDefaultLimit)))
	}

	return response, nil
}

func (s *ChickenService) GetChickenProcurement(id uint64) (dto.ChickenProcurementResponse, error) {
	s.repository.UseTx(false)

	chickenProcurement, err := s.repository.GetChickenProcurement(id)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	paymentResponses := make([]dto.ChickenProcurementPaymentResponse, 0)
	remainingPayment := chickenProcurement.TotalPrice
	for _, payment := range chickenProcurement.Payments {
		newPaymentResponse := mapper.ChickenProcurementPaymentToResponse(&payment)
		remainingPayment = remainingPayment.Sub(payment.Nominal)
		newPaymentResponse.Remaining = remainingPayment.String()

		paymentResponses = append(paymentResponses, newPaymentResponse)
	}

	chickenProcurementResponse := mapper.ChickenProcurementToResponse(&chickenProcurement)
	chickenProcurementResponse.Payments = paymentResponses
	chickenProcurementResponse.RemainingPayment = remainingPayment.String()

	return chickenProcurementResponse, nil
}

func (s *ChickenService) ConfirmationChickenProcurementDraft(id uint64, request dto.ConfirmationChickenProcurementRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurementDraft, err := s.repository.GetChickenProcurementDraft(id)
	if err != nil {
		s.log.Error("failed get chicken procurement draft", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	cage, err := s.cageService.GetCageById(chickenProcurementDraft.CageId)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	if cage.IsUsed {
		return dto.ChickenProcurementResponse{}, errx.BadRequest("cage is in used by another chicken")
	}

	totalPrice, err := decimal.NewFromString(request.TotalPrice)
	if err != nil {
		s.log.Error("failed to parse price from string", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	estimateArrivalDate, err := time.Parse("02-01-2006", request.EstimateArrivalDate)
	if err != nil {
		s.log.Error("failed to parse estimate arrival date", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid estimation arrival date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment type")
	}

	chickenProcurement := entity.ChickenProcurement{
		CageId:                chickenProcurementDraft.CageId,
		SupplierId:            chickenProcurementDraft.SupplierId,
		Quantity:              request.Quantity,
		TotalPrice:            totalPrice,
		Status:                enum.ProcurementStatusSentOff,
		PaymentType:           paymentType,
		PaymentStatus:         enum.PaymentStatusNotPaid,
		EstimationArrivalDate: estimateArrivalDate,
		CreatedBy:             uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.DeadlinePaymentDate != nil {
		deadlinePaymentDate, err := time.Parse("02-01-2006", *request.DeadlinePaymentDate)
		if err != nil {
			s.log.Error("failed parse deadline payment date", zap.Error(err))
			return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid deadline payment date format")
		}

		chickenProcurement.DeadlinePaymentDate = sql.NullTime{Time: deadlinePaymentDate, Valid: true}
	}

	chickenProcurementPayments := make([]entity.ChickenProcurementPayment, 0)
	totalPayment := decimal.Zero
	for _, payment := range request.Payments {
		paymentDate, err := time.Parse("02-01-2006", payment.PaymentDate)
		if err != nil {
			s.log.Error("failed to parse payment date", zap.Error(err))
			return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment date format")
		}

		nominal, err := decimal.NewFromString(payment.Nominal)
		if err != nil {
			s.log.Error("failed to parse nominal from string", zap.Error(err))
			return dto.ChickenProcurementResponse{}, err
		}

		paymentMethod := enum.ValueOfPaymentMethod(payment.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("invalid payment method", zap.String("paymentMethod", payment.PaymentMethod))
			return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment method")
		}

		chickenProcurementPayments = append(chickenProcurementPayments, entity.ChickenProcurementPayment{
			PaymentDate:   paymentDate,
			PaymentMethod: paymentMethod,
			PaymentProof:  payment.PaymentProof,
			Nominal:       nominal,
			CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
		})

		totalPayment = totalPayment.Add(nominal)
	}

	if totalPayment.Equal(chickenProcurement.TotalPrice) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.LessThan(chickenProcurement.TotalPrice) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.ChickenProcurementResponse{}, errx.BadRequest("nominal greater than total price")
	}

	err = s.repository.CreateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed to create chicken procurement", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	if len(chickenProcurementPayments) > 0 {
		for i := range chickenProcurementPayments {
			chickenProcurementPayments[i].ChickenProcurementId = chickenProcurement.Id
		}

		err = s.repository.CreateChickenProcurementPaymentInBatch(&chickenProcurementPayments)
		if err != nil {
			s.log.Error("failed to create chicken procurement payments in batch", zap.Error(err))
			return dto.ChickenProcurementResponse{}, err
		}
	}

	err = s.repository.DeleteChickenProcurementDraft(id)
	if err != nil {
		s.log.Error("failed delete chicken procurement draft", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	isUsed := true
	_, err = s.cageService.UpdateCage(chickenProcurement.CageId, dto.UpdateCageRequest{
		Name:            chickenProcurement.Cage.Name,
		Capacity:        chickenProcurement.Cage.Capacity,
		LocationId:      chickenProcurement.Cage.LocationId,
		ChickenCategory: chickenProcurement.Cage.ChickenCategory.String(),
		IsUsed:          &isUsed,
	}, userId)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurement, err = s.repository.GetChickenProcurement(chickenProcurement.Id)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	paymentResponses := make([]dto.ChickenProcurementPaymentResponse, 0)
	remainingPayment := chickenProcurement.TotalPrice
	for _, payment := range chickenProcurement.Payments {
		newPaymentResponse := mapper.ChickenProcurementPaymentToResponse(&payment)
		remainingPayment = remainingPayment.Sub(payment.Nominal)
		newPaymentResponse.Remaining = remainingPayment.String()

		paymentResponses = append(paymentResponses, newPaymentResponse)
	}

	chickenProcurementResponse := mapper.ChickenProcurementToResponse(&chickenProcurement)
	chickenProcurementResponse.Payments = paymentResponses
	chickenProcurementResponse.RemainingPayment = remainingPayment.String()

	return chickenProcurementResponse, nil
}

func (s *ChickenService) ArrivalConfirmationChickenProcurement(id uint64, request dto.ArrivalConfirmationChickenProcurementRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error) {
	s.repository.UseTx(false)

	chickenProcurement, err := s.repository.GetChickenProcurement(id)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurement.ReceiveQuantity = sql.NullInt64{Int64: int64(request.Quantity), Valid: true}
	chickenProcurement.Note = request.Note
	chickenProcurement.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	chickenProcurement.TakenBy = uuid.NullUUID{UUID: userId, Valid: true}
	chickenProcurement.IsArrived = true

	if chickenProcurement.Quantity != request.Quantity {
		chickenProcurement.Status = enum.ProcurementStatusArrivedNotOk
	} else {
		chickenProcurement.Status = enum.ProcurementStatusArrivedOk
	}

	_, err = s.cageService.CreateChickenCage(dto.CreateChickenCageRequest{
		CageId:               chickenProcurement.CageId,
		ChickenProcurementId: &chickenProcurement.Id,
		TotalChicken:         uint64(chickenProcurement.ReceiveQuantity.Int64),
	}, userId)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	err = s.repository.UpdateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed update chicken procurement", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurement, err = s.repository.GetChickenProcurement(id)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	paymentResponses := make([]dto.ChickenProcurementPaymentResponse, 0)
	remainingPayment := chickenProcurement.TotalPrice
	for _, payment := range chickenProcurement.Payments {
		newPaymentResponse := mapper.ChickenProcurementPaymentToResponse(&payment)
		remainingPayment = remainingPayment.Sub(payment.Nominal)
		newPaymentResponse.Remaining = remainingPayment.String()

		paymentResponses = append(paymentResponses, newPaymentResponse)
	}

	chickenProcurementResponse := mapper.ChickenProcurementToResponse(&chickenProcurement)
	chickenProcurementResponse.Payments = paymentResponses
	chickenProcurementResponse.RemainingPayment = remainingPayment.String()

	return chickenProcurementResponse, nil
}

func (s *ChickenService) CreateChickenProcurementPayment(chickenProcurementId uint64, request dto.CreateChickenProcurementPaymentRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurement, err := s.repository.GetChickenProcurement(chickenProcurementId)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	if chickenProcurement.PaymentStatus == enum.PaymentStatusPaid {
		return dto.ChickenProcurementResponse{}, errx.BadRequest("chicken procurement is already paid")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid nominal format")
	}

	chickenProcurementPayment := entity.ChickenProcurementPayment{
		ChickenProcurementId: chickenProcurementId,
		PaymentProof:         request.PaymentProof,
		Nominal:              nominal,
		PaymentDate:          paymentDate,
		PaymentMethod:        paymentMethod,
		CreatedBy:            uuid.NullUUID{UUID: userId, Valid: true},
	}

	totalPayment := nominal
	for _, payment := range chickenProcurement.Payments {
		totalPayment = totalPayment.Add(payment.Nominal)
	}

	if totalPayment.Equal(chickenProcurement.TotalPrice) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(chickenProcurement.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.UpdateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed update chicken procurement", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	err = s.repository.CreateChickenProcurementPayment(&chickenProcurementPayment)
	if err != nil {
		s.log.Error("failed create chicken procurement payment", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurement, err = s.repository.GetChickenProcurement(chickenProcurement.Id)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	paymentResponses := make([]dto.ChickenProcurementPaymentResponse, 0)
	remainingPayment := chickenProcurement.TotalPrice
	for _, payment := range chickenProcurement.Payments {
		newPaymentResponse := mapper.ChickenProcurementPaymentToResponse(&payment)
		remainingPayment = remainingPayment.Sub(payment.Nominal)
		newPaymentResponse.Remaining = remainingPayment.String()

		paymentResponses = append(paymentResponses, newPaymentResponse)
	}

	chickenProcurementResponse := mapper.ChickenProcurementToResponse(&chickenProcurement)
	chickenProcurementResponse.Payments = paymentResponses
	chickenProcurementResponse.RemainingPayment = remainingPayment.String()

	return chickenProcurementResponse, nil
}

func (s *ChickenService) UpdateChickenProcurementPayment(chickenProcurementId uint64, id uint64, request dto.UpdateChickenProcurementPaymentRequest, userId uuid.UUID) (dto.ChickenProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurement, err := s.repository.GetChickenProcurement(chickenProcurementId)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurementPayment, err := s.repository.GetChickenProcurementPayment(id)
	if err != nil {
		s.log.Error("failed to get chicken procurement payment by id", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	if chickenProcurement.PaymentStatus == enum.PaymentStatusPaid {
		return dto.ChickenProcurementResponse{}, errx.BadRequest("chicken procurement is already paid")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid nominal format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("invalid payment method")
	}

	totalPayment := nominal
	for _, payment := range chickenProcurement.Payments {
		if payment.Id != chickenProcurementPayment.Id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.Equal(chickenProcurement.TotalPrice) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(chickenProcurement.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.ChickenProcurementResponse{}, errx.BadRequest("total payment is greater than total price")
	} else if totalPayment.LessThan(chickenProcurement.TotalPrice) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusUnpaid
	}

	chickenProcurementPayment.PaymentMethod = paymentMethod
	chickenProcurementPayment.Nominal = nominal
	chickenProcurementPayment.PaymentProof = request.PaymentProof
	chickenProcurementPayment.PaymentDate = paymentDate
	chickenProcurementPayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed update chicken procurement", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	err = s.repository.UpdateChickenProcurementPayment(&chickenProcurementPayment)
	if err != nil {
		s.log.Error("failed create chicken procurement payment", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurement, err = s.repository.GetChickenProcurement(chickenProcurement.Id)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	chickenProcurement, err = s.repository.GetChickenProcurement(chickenProcurement.Id)
	if err != nil {
		return dto.ChickenProcurementResponse{}, err
	}

	paymentResponses := make([]dto.ChickenProcurementPaymentResponse, 0)
	remainingPayment := chickenProcurement.TotalPrice
	for _, payment := range chickenProcurement.Payments {
		newPaymentResponse := mapper.ChickenProcurementPaymentToResponse(&payment)
		remainingPayment = remainingPayment.Sub(payment.Nominal)
		newPaymentResponse.Remaining = remainingPayment.String()

		paymentResponses = append(paymentResponses, newPaymentResponse)
	}

	chickenProcurementResponse := mapper.ChickenProcurementToResponse(&chickenProcurement)
	chickenProcurementResponse.Payments = paymentResponses
	chickenProcurementResponse.RemainingPayment = remainingPayment.String()

	return chickenProcurementResponse, nil
}

func (s *ChickenService) DeleteChickenProcurementPayment(chickenProcurementId uint64, id uint64, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurement, err := s.repository.GetChickenProcurement(chickenProcurementId)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return err
	}

	if chickenProcurement.PaymentStatus == enum.PaymentStatusPaid {
		return errx.BadRequest("chicken procurement is already paid")
	}

	totalPayment := decimal.Zero
	for _, payment := range chickenProcurement.Payments {
		if payment.Id != id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.LessThan(chickenProcurement.TotalPrice) && totalPayment.GreaterThan(decimal.Zero) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusUnpaid
		chickenProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	} else if totalPayment.LessThan(decimal.Zero) {
		s.log.Error("delete this payment make minus", zap.Error(err))
		return errx.BadRequest("delete this payment make minus")
	}

	err = s.repository.UpdateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed update chicken procurement", zap.Error(err))
		return err
	}

	err = s.repository.DeleteChickenProcurementPayment(id)
	if err != nil {
		s.log.Error("failed to delete chicken procurement payment")
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) CreateAfkirChickenCustomer(request dto.CreateAfkirChickenCustomerRequest, userId uuid.UUID) (dto.AfkirChickenCustomerResponse, error) {
	s.repository.UseTx(false)

	data := entity.AfkirChickenCustomer{
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		Address:     request.Address,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateAfkirChickenCustomer(&data)
	if err != nil {
		s.log.Error("failed create afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenCustomerResponse{}, err
	}

	data, err = s.repository.GetAfkirChickenCustomer(data.Id)
	if err != nil {
		s.log.Error("failed get afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenCustomerResponse{}, err
	}

	return mapper.AfkirChickenCustomerToResponse(&data), nil
}

func (s *ChickenService) GetAfkirChickenCustomers() ([]dto.AfkirChickenCustomerListResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetAfkirChickenCustomers()
	if err != nil {
		s.log.Error("failed get afkir chicken customers", zap.Error(err))
		return nil, err
	}

	response := make([]dto.AfkirChickenCustomerListResponse, 0)
	for _, e := range data {
		response = append(response, mapper.AfkirChickenCustomerToListResponse(&e))
	}

	return response, nil
}

func (s *ChickenService) GetAfkirChickenCustomer(id uint64) (dto.AfkirChickenCustomerResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetAfkirChickenCustomer(id)
	if err != nil {
		return dto.AfkirChickenCustomerResponse{}, err
	}

	return mapper.AfkirChickenCustomerToResponse(&data), nil
}

func (s *ChickenService) UpdateAfkirChickenCustomer(id uint64, request dto.UpdateAfkirChickenCustomerRequest, userId uuid.UUID) (dto.AfkirChickenCustomerResponse, error) {
	s.repository.UseTx(false)

	afkirChickenCustomer, err := s.repository.GetAfkirChickenCustomer(id)
	if err != nil {
		s.log.Error("failed get afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenCustomerResponse{}, err
	}

	afkirChickenCustomer.Name = request.Name
	afkirChickenCustomer.PhoneNumber = request.PhoneNumber
	afkirChickenCustomer.Address = request.Address
	afkirChickenCustomer.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateAfkirChickenCustomer(&afkirChickenCustomer)
	if err != nil {
		s.log.Error("failed update afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenCustomerResponse{}, err
	}

	return mapper.AfkirChickenCustomerToResponse(&afkirChickenCustomer), nil
}

func (s *ChickenService) DeleteAfkirChickenCustomer(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteAfkirChickenCustomer(id)
	if err != nil {
		s.log.Error("failed delete afkir chicken customer", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) CreateAkfirChickenSaleDraft(request dto.CreateAfkirChickenSaleDraftRequest, userId uuid.UUID) (dto.AfkirChickenSaleDraftResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	if chickenCage.TotalChicken < request.TotalSellChicken {
		return dto.AfkirChickenSaleDraftResponse{}, errx.BadRequest("total sell chicken must be less than or equal total chicken")
	}

	pricePerChicken, err := decimal.NewFromString(request.PricePerChicken)
	if err != nil {
		s.log.Error("failed parse price from string", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	totalPrice := pricePerChicken.Mul(decimal.NewFromUint64(request.TotalSellChicken))

	data := entity.AfkirChickenSaleDraft{
		ChickenCageId:          request.ChickenCageId,
		AfkirChickenCustomerId: request.AfkirChickenCustomerId,
		TotalSellChicken:       request.TotalSellChicken,
		PricePerChicken:        pricePerChicken,
		TotalPrice:             totalPrice,
	}

	err = s.repository.CreateAfkirChickenSaleDraft(&data)
	if err != nil {
		s.log.Error("failed create afkir chicken sale draft", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	data, err = s.repository.GetAfkirChickenSaleDraft(data.Id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale draft", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	return mapper.AfkirChickenSaleDraftToResponse(&data), nil
}

func (s *ChickenService) GetAfkirChickenSaleDrafts() ([]dto.AfkirChickenSaleDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetAfkirChickenSaleDrafts()
	if err != nil {
		s.log.Error("failed get afkir chicken sale drafts", zap.Error(err))
		return nil, err
	}

	response := make([]dto.AfkirChickenSaleDraftResponse, 0)
	for _, e := range data {
		response = append(response, mapper.AfkirChickenSaleDraftToResponse(&e))
	}

	return response, nil
}

func (s *ChickenService) GetAfkirChickenSaleDraft(id uint64) (dto.AfkirChickenSaleDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetAfkirChickenSaleDraft(id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale draft", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	return mapper.AfkirChickenSaleDraftToResponse(&data), nil
}

func (s *ChickenService) UpdateAfkirChickenSaleDraft(id uint64, request dto.UpdateAfkirChickenSaleDraftRequest, userId uuid.UUID) (dto.AfkirChickenSaleDraftResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	if chickenCage.TotalChicken < request.TotalSellChicken {
		return dto.AfkirChickenSaleDraftResponse{}, errx.BadRequest("total sell chicken must be less than or equal total chicken")
	}

	data, err := s.repository.GetAfkirChickenSaleDraft(id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale draft", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	pricePerChicken, err := decimal.NewFromString(request.PricePerChicken)
	if err != nil {
		s.log.Error("failed parse price from string", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	data.ChickenCageId = request.ChickenCageId
	data.AfkirChickenCustomerId = request.AfkirChickenCustomerId
	data.TotalSellChicken = request.TotalSellChicken
	data.PricePerChicken = pricePerChicken
	data.TotalPrice = pricePerChicken.Mul(decimal.NewFromUint64(request.TotalSellChicken))

	err = s.repository.UpdateAfkirChickenSaleDraft(&data)
	if err != nil {
		s.log.Error("failed update afkir chicken sale draft", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	data, err = s.repository.GetAfkirChickenSaleDraft(id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale draft", zap.Error(err))
		return dto.AfkirChickenSaleDraftResponse{}, err
	}

	return mapper.AfkirChickenSaleDraftToResponse(&data), nil
}

func (s *ChickenService) DeleteAfkirChickenSaleDraft(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteAfkirChickenSaleDraft(id)
	if err != nil {
		s.log.Error("failed delete afkir chickan sale draft", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) CreateAfkirChickenSale(request dto.CreateAfkirChickenSaleRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	pricePerChicken, err := decimal.NewFromString(request.PricePerChicken)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid price format")
	}

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment type: %s", request.PaymentType))
	}

	totalPrice := pricePerChicken.Mul(decimal.NewFromUint64(request.TotalSellChicken))

	dateNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	afkirSale := entity.AfkirChickenSale{
		AfkirChickenCustomerId: request.AfkirChickenCustomerId,
		ChickenCageId:          request.ChickenCageId,
		TotalSellChicken:       request.TotalSellChicken,
		PricePerChicken:        pricePerChicken,
		TotalPrice:             totalPrice,
		ChickenAge:             chickenCage.ChickenAge,
		PaymentStatus:          enum.PaymentStatusNotPaid,
		PaymentType:            paymentType,
		CreatedBy:              uuid.NullUUID{UUID: userId, Valid: true},
	}

	payments := make([]entity.AfkirChickenSalePayment, 0, len(request.Payments))
	totalPayment := decimal.Zero

	for _, p := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(p.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment method: %s", p.PaymentMethod))
		}

		paymentDate, err := time.Parse("02-01-2006", p.PaymentDate)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid payment date format")
		}

		nominal, err := decimal.NewFromString(p.Nominal)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid nominal format")
		}

		totalPayment = totalPayment.Add(nominal)
		payments = append(payments, entity.AfkirChickenSalePayment{
			AfkirChickenSaleId: 0,
			Nominal:            nominal,
			PaymentDate:        paymentDate,
			PaymentMethod:      paymentMethod,
			PaymentProof:       p.PaymentProof,
			CreatedBy:          uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !afkirSale.TotalPrice.Equal(totalPayment) {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}
		afkirSale.PaymentStatus = enum.PaymentStatusPaid
	} else {
		if totalPayment.Equal(totalPrice) {
			afkirSale.PaymentStatus = enum.PaymentStatusPaid
		} else if totalPayment.GreaterThan(decimal.Zero) {
			afkirSale.PaymentStatus = enum.PaymentStatusUnpaid
		} else {
			afkirSale.PaymentStatus = enum.PaymentStatusNotPaid
		}
	}

	if afkirSale.PaymentStatus != enum.PaymentStatusPaid {
		afkirSale.DeadlinePaymentDate = sql.NullTime{Time: dateNow.AddDate(0, 0, 7), Valid: true}
	}

	err = s.repository.CreateAfkirChickenSale(&afkirSale)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	if len(payments) > 0 {
		for i := range payments {
			payments[i].AfkirChickenSaleId = afkirSale.Id
		}
		err = s.repository.CreateAfkirChickenSalePaymentInBatch(&payments)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, err
		}
	}

	isUsed := false
	currentChicken := chickenCage.TotalChicken - request.TotalSellChicken
	if currentChicken > 0 {
		isUsed = true
	}

	_, err = s.cageService.UpdateCage(chickenCage.Cage.Id, dto.UpdateCageRequest{
		Name:            chickenCage.Cage.Name,
		Capacity:        chickenCage.Cage.Capacity,
		LocationId:      chickenCage.Cage.Location.Id,
		ChickenCategory: chickenCage.Cage.ChickenCategory,
		IsUsed:          &isUsed,
	}, userId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	_, err = s.cageService.UpdateChickenCage(chickenCage.Id, dto.UpdateChickenCageRequest{
		TotalChicken: currentChicken,
	}, userId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	if !isUsed {
		_, err = s.cageService.CreateChickenCage(dto.CreateChickenCageRequest{
			CageId:               chickenCage.Cage.Id,
			ChickenProcurementId: nil,
			TotalChicken:         0,
		}, userId)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, err
		}
	}

	afkirChickenCustomer, err := s.repository.GetAfkirChickenCustomer(request.AfkirChickenCustomerId)
	if err != nil {
		s.log.Error("failed get afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirChickenCustomer.LatestPrice = pricePerChicken

	err = s.repository.UpdateAfkirChickenCustomer(&afkirChickenCustomer)
	if err != nil {
		s.log.Error("failed update afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	if err = s.repository.Commit(); err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirSale, err = s.repository.GetAfkirChickenSale(afkirSale.Id)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	resPayments := make([]dto.AfkirChickenSalePaymentResponse, len(afkirSale.Payments))
	remainingPayment := afkirSale.TotalPrice
	for i, pay := range afkirSale.Payments {
		resPayments[i] = mapper.AfkirChickenSalePaymentToResponse(&pay)
		remainingPayment = remainingPayment.Sub(pay.Nominal)
		resPayments[i].Remaining = remainingPayment.String()
	}

	resp := mapper.AfkirChickenSaleToResponse(&afkirSale)
	resp.RemainingPayment = remainingPayment.String()
	resp.Payments = resPayments

	return resp, nil
}

func (s *ChickenService) GetAfkirChickenSales(filter dto.GetAfkirChickenSaleFilter) (dto.AfkirChickenSaleListPaginationResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetAfkirChickenSales(filter)
	if err != nil {
		s.log.Error("failed get afkir chicken sales", zap.Error(err))
		return dto.AfkirChickenSaleListPaginationResponse{}, err
	}

	totalData, err := s.repository.CountChickenAfkirChickenSale(filter)
	if err != nil {
		return dto.AfkirChickenSaleListPaginationResponse{}, err
	}

	afkirChickenSaleResponse := make([]dto.AfkirChickenSaleListResponse, 0)
	for _, e := range data {
		afkirChickenSaleResponse = append(afkirChickenSaleResponse, mapper.AfkirChickenSaleToListResponse(&e))
	}

	response := dto.AfkirChickenSaleListPaginationResponse{
		AfkirChickenSales: afkirChickenSaleResponse,
	}

	if filter.Page > 0 {
		response.TotalData = uint64(totalData)
		response.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return response, nil
}

func (s *ChickenService) GetAkfirChickenSale(id uint64) (dto.AfkirChickenSaleResponse, error) {
	s.repository.UseTx(false)

	afkirSale, err := s.repository.GetAfkirChickenSale(id)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	// Map payments with remaining calculation
	resPayments := make([]dto.AfkirChickenSalePaymentResponse, len(afkirSale.Payments))
	remainingPayment := afkirSale.TotalPrice
	for i, pay := range afkirSale.Payments {
		resPayments[i] = mapper.AfkirChickenSalePaymentToResponse(&pay)
		remainingPayment = remainingPayment.Sub(pay.Nominal)
		resPayments[i].Remaining = remainingPayment.String()
	}

	resp := mapper.AfkirChickenSaleToResponse(&afkirSale)
	resp.RemainingPayment = remainingPayment.String()
	resp.Payments = resPayments

	return resp, nil
}

func (s *ChickenService) CreateAfkirChickenSalePayment(afkirChickenSaleId uint64, request dto.CreateAfkirChickenSalePaymentRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	afkirChickenSale, err := s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed parse nominal from string", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed parse payment date", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment method : %s", request.PaymentMethod))
	}

	payment := entity.AfkirChickenSalePayment{
		AfkirChickenSaleId: afkirChickenSaleId,
		Nominal:            nominal,
		PaymentDate:        paymentDate,
		PaymentMethod:      paymentMethod,
		PaymentProof:       request.PaymentProof,
		CreatedBy:          uuid.NullUUID{UUID: userId, Valid: true},
	}

	totalCurrentPayment := decimal.Zero
	for _, e := range afkirChickenSale.Payments {
		totalCurrentPayment = totalCurrentPayment.Add(e.Nominal)
	}

	if totalCurrentPayment.Add(nominal).Equal(afkirChickenSale.TotalPrice) {
		afkirChickenSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalCurrentPayment.Add(nominal).GreaterThan(afkirChickenSale.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.UpdateAfkirChickenSale(&afkirChickenSale)
	if err != nil {
		s.log.Error("failed update afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	err = s.repository.CreateAfkirChickenSalePayment(&payment)
	if err != nil {
		s.log.Error("failed create afkir chicken sale payment", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirSale, err := s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	resPayments := make([]dto.AfkirChickenSalePaymentResponse, len(afkirSale.Payments))
	remainingPayment := afkirSale.TotalPrice
	for i, pay := range afkirSale.Payments {
		resPayments[i] = mapper.AfkirChickenSalePaymentToResponse(&pay)
		remainingPayment = remainingPayment.Sub(pay.Nominal)
		resPayments[i].Remaining = remainingPayment.String()
	}

	resp := mapper.AfkirChickenSaleToResponse(&afkirSale)
	resp.RemainingPayment = remainingPayment.String()
	resp.Payments = resPayments

	return resp, nil
}

func (s *ChickenService) UpdateAfkirChickenSalePayment(afkirChickenSaleId uint64, id uint64, request dto.UpdateAfkirChickenSalePaymentRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	afkirChickenSalePayment, err := s.repository.GetAfkirChickenSalePaymentById(id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale payment", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirChickenSale, err := s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed parse nominal from string", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed parse payment date", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment method : %s", request.PaymentMethod))
	}

	afkirChickenSalePayment.Nominal = nominal
	afkirChickenSalePayment.PaymentDate = paymentDate
	afkirChickenSalePayment.PaymentMethod = paymentMethod
	afkirChickenSalePayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	totalCurrentPrice := decimal.Zero
	for _, e := range afkirChickenSale.Payments {
		if e.Id != afkirChickenSalePayment.Id {
			totalCurrentPrice = totalCurrentPrice.Add(e.Nominal)
		}
	}

	if totalCurrentPrice.Add(nominal).Equal(afkirChickenSale.TotalPrice) {
		afkirChickenSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalCurrentPrice.Add(nominal).LessThan(afkirChickenSale.TotalPrice) {
		afkirChickenSale.PaymentStatus = enum.PaymentStatusUnpaid
	} else if totalCurrentPrice.Add(nominal).GreaterThan(afkirChickenSale.TotalPrice) {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("nominal is greater than total price")
	}

	err = s.repository.UpdateAfkirChickenSale(&afkirChickenSale)
	if err != nil {
		s.log.Error("failed update afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	err = s.repository.UpdateAfkirChickenSalePayment(&afkirChickenSalePayment)
	if err != nil {
		s.log.Error("failed update afkir chicken sale payment", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirSale, err := s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	// Map payments with remaining calculation
	resPayments := make([]dto.AfkirChickenSalePaymentResponse, len(afkirSale.Payments))
	remainingPayment := afkirSale.TotalPrice
	for i, pay := range afkirSale.Payments {
		resPayments[i] = mapper.AfkirChickenSalePaymentToResponse(&pay)
		remainingPayment = remainingPayment.Sub(pay.Nominal)
		resPayments[i].Remaining = remainingPayment.String()
	}

	resp := mapper.AfkirChickenSaleToResponse(&afkirSale)
	resp.RemainingPayment = remainingPayment.String()
	resp.Payments = resPayments

	return resp, nil
}

func (s *ChickenService) DeleteAfkirChickenSalePayment(afkirChickenSaleId uint64, id uint64) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	afkirChickenSalePayment, err := s.repository.GetAfkirChickenSalePaymentById(id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale payment", zap.Error(err))
		return err
	}

	afkirChickenSale, err := s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return err
	}

	totalCurrentPrice := decimal.Zero
	for _, e := range afkirChickenSale.Payments {
		if e.Id != afkirChickenSalePayment.Id {
			totalCurrentPrice = totalCurrentPrice.Add(e.Nominal)
		}
	}

	if totalCurrentPrice.LessThan(afkirChickenSale.TotalPrice) {
		afkirChickenSale.PaymentStatus = enum.PaymentStatusUnpaid
	} else if totalCurrentPrice.LessThan(decimal.Zero) {
		return errx.BadRequest("nominal is less than 0")
	}

	err = s.repository.UpdateAfkirChickenSale(&afkirChickenSale)
	if err != nil {
		s.log.Error("failed update afkir chicken sale", zap.Error(err))
		return err
	}

	err = s.repository.DeleteAfkirChickenSalePayment(id)
	if err != nil {
		s.log.Error("failed delete afkir chicken sale payment", zap.Error(err))
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return err
	}

	afkirChickenSale, err = s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) ConfirmationAfkirChickenSaleDraft(id uint64, request dto.CreateAfkirChickenSaleRequest, userId uuid.UUID) (dto.AfkirChickenSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	err := s.repository.DeleteAfkirChickenSaleDraft(id)
	if err != nil {
		s.log.Error("failed delete afkir chickens sale draft", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	if chickenCage.TotalChicken < request.TotalSellChicken {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("total sell chicken must be less than total chicken")
	}

	pricePerChicken, err := decimal.NewFromString(request.PricePerChicken)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid price format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment type: %s", request.PaymentType))
	}

	totalPrice := pricePerChicken.Mul(decimal.NewFromUint64(request.TotalSellChicken))

	dateNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	afkirSale := entity.AfkirChickenSale{
		AfkirChickenCustomerId: request.AfkirChickenCustomerId,
		ChickenCageId:          request.ChickenCageId,
		TotalSellChicken:       request.TotalSellChicken,
		PricePerChicken:        pricePerChicken,
		TotalPrice:             totalPrice,
		ChickenAge:             chickenCage.ChickenAge,
		PaymentStatus:          enum.PaymentStatusNotPaid,
		PaymentType:            paymentType,
		CreatedBy:              uuid.NullUUID{UUID: userId, Valid: true},
	}

	payments := make([]entity.AfkirChickenSalePayment, 0, len(request.Payments))
	totalPayment := decimal.Zero

	for _, p := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(p.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment method: %s", p.PaymentMethod))
		}

		paymentDate, err := time.Parse("02-01-2006", p.PaymentDate)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid payment date format")
		}

		nominal, err := decimal.NewFromString(p.Nominal)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest("invalid nominal format")
		}

		totalPayment = totalPayment.Add(nominal)
		payments = append(payments, entity.AfkirChickenSalePayment{
			AfkirChickenSaleId: 0,
			Nominal:            nominal,
			PaymentDate:        paymentDate,
			PaymentMethod:      paymentMethod,
			PaymentProof:       p.PaymentProof,
			CreatedBy:          uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !afkirSale.TotalPrice.Equal(totalPayment) {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}
		afkirSale.PaymentStatus = enum.PaymentStatusPaid
	} else {
		if totalPayment.Equal(totalPrice) {
			afkirSale.PaymentStatus = enum.PaymentStatusPaid
		} else if totalPayment.GreaterThan(decimal.Zero) {
			afkirSale.PaymentStatus = enum.PaymentStatusUnpaid
		} else {
			afkirSale.PaymentStatus = enum.PaymentStatusNotPaid
		}
	}

	if afkirSale.PaymentStatus != enum.PaymentStatusPaid {
		afkirSale.DeadlinePaymentDate = sql.NullTime{Time: dateNow.AddDate(0, 0, 7), Valid: true}
	}

	err = s.repository.CreateAfkirChickenSale(&afkirSale)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	if len(payments) > 0 {
		for i := range payments {
			payments[i].AfkirChickenSaleId = afkirSale.Id
		}
		err = s.repository.CreateAfkirChickenSalePaymentInBatch(&payments)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, err
		}
	}

	isUsed := false
	currentChicken := chickenCage.TotalChicken - request.TotalSellChicken
	if currentChicken > 0 {
		isUsed = true
	}

	_, err = s.cageService.UpdateCage(chickenCage.Cage.Id, dto.UpdateCageRequest{
		Name:            chickenCage.Cage.Name,
		Capacity:        chickenCage.Cage.Capacity,
		LocationId:      chickenCage.Cage.Location.Id,
		ChickenCategory: chickenCage.Cage.ChickenCategory,
		IsUsed:          &isUsed,
	}, userId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	_, err = s.cageService.UpdateChickenCage(chickenCage.Id, dto.UpdateChickenCageRequest{
		TotalChicken: currentChicken,
	}, userId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	if !isUsed {
		_, err = s.cageService.CreateChickenCage(dto.CreateChickenCageRequest{
			CageId:               chickenCage.Cage.Id,
			ChickenProcurementId: nil,
			TotalChicken:         0,
		}, userId)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, err
		}
	}

	afkirChickenCustomer, err := s.repository.GetAfkirChickenCustomer(request.AfkirChickenCustomerId)
	if err != nil {
		s.log.Error("failed get afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirChickenCustomer.LatestPrice = pricePerChicken

	err = s.repository.UpdateAfkirChickenCustomer(&afkirChickenCustomer)
	if err != nil {
		s.log.Error("failed update afkir chicken customer", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	if err = s.repository.Commit(); err != nil {
		s.log.Error("failed commit transcation", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	afkirSale, err = s.repository.GetAfkirChickenSale(afkirSale.Id)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	resPayments := make([]dto.AfkirChickenSalePaymentResponse, len(afkirSale.Payments))
	remainingPayment := afkirSale.TotalPrice
	for i, pay := range afkirSale.Payments {
		resPayments[i] = mapper.AfkirChickenSalePaymentToResponse(&pay)
		remainingPayment = remainingPayment.Sub(pay.Nominal)
		resPayments[i].Remaining = remainingPayment.String()
	}

	resp := mapper.AfkirChickenSaleToResponse(&afkirSale)
	resp.RemainingPayment = remainingPayment.String()
	resp.Payments = resPayments

	return resp, nil
}

func (s *ChickenService) GetChickenPerformances(filter dto.GetChickenPerformanceFilter) ([]dto.ChickenPerformanceResponse, error) {
	s.repository.UseTx(false)

	responses := make([]dto.ChickenPerformanceResponse, 0)

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	if filter.Date.Value().Year() == today.Year() &&
		filter.Date.Value().Month() == today.Month() &&
		filter.Date.Value().Day() == today.Day() {

		chickenCages, err := s.cageService.GetChickenCages(dto.GetChickenCageFilter{
			LocationId: filter.LocationId,
			CageId:     filter.CageId,
		})
		if err != nil {
			return nil, err
		}

		totalExpensePerMonth, err := s.cashflowService.GetTotalExpenseProductionInMonth(enum.Month(today.Month()), uint64(today.Year()))
		if err != nil {
			return nil, err
		}
		totalDayInMonth := util.TotalDaysInMonth(today.Year(), today.Month())
		totalExpensePerDay := totalExpensePerMonth.Div(decimal.NewFromUint64(totalDayInMonth))

		goodEgg, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg)
		if err != nil {
			return nil, err
		}
		itemPriceGoodEgg, err := s.itemService.GetItemPriceByItemIdAndSaleUnit(goodEgg.Id, enum.SaleUnitKg.String())
		if err != nil {
			return nil, err
		}
		price, err := decimal.NewFromString(itemPriceGoodEgg.Price)
		if err != nil {
			return nil, errx.BadRequest("invalid egg item price")
		}

		for _, chickenCage := range chickenCages {
			chickenMonitoring, err := s.repository.GetChickenMonitoringToday(chickenCage.Id, today)
			if err != nil {
				s.log.Error("failed get chicken monitoring today")
				return nil, err
			}
			eggMonitoring, err := s.eggService.GetEggMonitoringToday(chickenCage.Id, today)
			if err != nil {
				return nil, err
			}

			if chickenMonitoring.Id == 0 || eggMonitoring.Id == 0 {
				s.log.Info("skipping cage, no monitoring data for today", zap.Uint64("cageId", chickenCage.Id))
				continue
			}

			fmt.Println(chickenCage.Id)

			totalEgg := (eggMonitoring.TotalKarpetCrackedEgg * constant.TotalEggPerKarpet) +
				eggMonitoring.TotalRemainingCrackedEgg +
				(eggMonitoring.TotalKarpetGoodEgg * constant.TotalEggPerKarpet) +
				eggMonitoring.TotalRemainingGoodEgg

			var averageConsumptionPerChicken float64
			if chickenCage.TotalChicken > 0 {
				averageConsumptionPerChicken = chickenMonitoring.TotalFeed / float64(chickenCage.TotalChicken)
			} else {
				averageConsumptionPerChicken = 0
			}

			var averageWeightPerEgg float64
			if totalEgg > 0 {
				averageWeightPerEgg = eggMonitoring.TotalWeightAllEgg / float64(totalEgg)
			} else {
				averageWeightPerEgg = 0
			}

			response := dto.ChickenPerformanceResponse{
				CageName:                     chickenCage.Cage.Name,
				ChickenCategory:              chickenCage.ChickenCategory,
				ChickenAge:                   chickenCage.ChickenAge,
				TotalChicken:                 chickenCage.TotalChicken,
				TotalEgg:                     totalEgg,
				AverageConsumptionPerChicken: averageConsumptionPerChicken,
				AverageWeightPerEgg:          averageWeightPerEgg,
			}

			if chickenCage.TotalChicken > 0 {
				response.MortalityRate = float64(chickenMonitoring.TotalDeathChicken) / float64(chickenCage.TotalChicken)
			} else {
				response.MortalityRate = 0
			}

			if chickenCage.TotalChicken > 0 {
				response.HDP = float64(response.TotalEgg) / float64(chickenCage.TotalChicken) * 100.0
			} else {
				response.HDP = 0
			}

			if response.TotalEgg > 0 {
				response.FCR = float64(chickenMonitoring.TotalFeed) / float64(response.TotalEgg)
			} else {
				response.FCR = 0
			}

			if chickenCage.ChickenAge >= 90 {
				response.Productivity = enum.ChickenProductivityAfkir.String()
			} else {
				totalPrice := decimal.Zero
				if eggMonitoring.TotalWeightGoodEgg != 0.0 {
					totalPrice = price.Mul(decimal.NewFromFloat(eggMonitoring.TotalWeightGoodEgg))
				}

				if totalPrice.Sub(totalExpensePerDay).GreaterThanOrEqual(decimal.NewFromInt(constant.MinProfitForCageNotAfkir)) {
					response.Productivity = enum.ChickenProductivityProductive.String()
				} else {
					response.Productivity = enum.ChickenProductivityAfkir.String()
				}
			}

			responses = append(responses, response)
		}
	} else {
		chickenPerformances, err := s.repository.GetChickenPerformances(filter)
		if err != nil {
			return nil, fmt.Errorf("failed to get chicken performances from DB: %w", err)
		}
		for _, chickenPerformance := range chickenPerformances {
			responses = append(responses, mapper.ChickenPerformanceToResponse(&chickenPerformance))
		}
	}

	return responses, nil
}

func (s *ChickenService) GetChickenAndWarehouseOverview(filter dto.GetChickenAndWarehouseOverviewFilter) (dto.ChickenAndWarehouseOverviewResponse, error) {
	s.repository.UseTx(false)

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)

	var totalSafeItem, totalDangerItem uint64
	warehouseItems, err := s.warehouseService.GetWarehouseItems(dto.GetWarehouseItemFilter{
		LocationId:  filter.LocationId,
		WarehouseId: filter.WarehouseId,
	})
	if err != nil {
		return dto.ChickenAndWarehouseOverviewResponse{}, err
	}
	for _, warehouseItem := range warehouseItems {
		switch warehouseItem.Description {
		case constant.WarehouseItemDescriptionSafe:
			totalSafeItem += 1
		case constant.WarehouseItemDescriptionDanger:
			totalDangerItem += 1
		}
	}

	chickenPerformances, err := s.GetChickenPerformances(dto.GetChickenPerformanceFilter{
		Date:       param.DateParam(today),
		LocationId: filter.LocationId,
		CageId:     filter.CageId,
	})
	if err != nil {
		return dto.ChickenAndWarehouseOverviewResponse{}, err
	}

	var totalProductiveCage, totalAfkirCage uint64
	for _, chickenPerformance := range chickenPerformances {
		switch chickenPerformance.Productivity {
		case enum.ChickenProductivityProductive.String():
			totalProductiveCage += 1
		case enum.ChickenProductivityAfkir.String():
			totalAfkirCage += 1
		}
	}

	chickenGraphs := make([]dto.ChickenGraphResponse, 0)
	switch filter.OverviewGraphTime.Value() {
	case enum.OverviewGraphTimeThisWeek:
		chickenGraphs, err = s.buildWeeklyGraph(filter.LocationId, filter.CageId)
	case enum.OverviewGraphTimeThisMonth:
		chickenGraphs, err = s.buildMonthlyGraph(filter.LocationId, filter.CageId)
	case enum.OverviewGraphTimeThisYear:
		chickenGraphs, err = s.buildYearlyGraph(filter.LocationId, filter.CageId)
	}
	if err != nil {
		return dto.ChickenAndWarehouseOverviewResponse{}, err
	}

	var totalFeed, totalWeightEgg float64
	var totalEgg, totalDeathChicken, totalChicken uint64
	var totalDOCChicken, totalGrowerChicken, totalPreLayerChicken, totalLayerChicken, totalAfkirChicken uint64
	chickenCages, err := s.cageService.GetChickenCages(dto.GetChickenCageFilter{
		LocationId: filter.LocationId,
		CageId:     filter.CageId,
	})
	if err != nil {
		s.log.Error("failed get chicken cages", zap.Error(err))
		return dto.ChickenAndWarehouseOverviewResponse{}, err
	}

	for _, chickenCage := range chickenCages {
		chickenMonitoring, err := s.repository.GetChickenMonitoringToday(chickenCage.Id, today)
		if err != nil {
			s.log.Error("failed get chicken monitoring today", zap.Error(err))
			return dto.ChickenAndWarehouseOverviewResponse{}, err
		}

		eggMonitoring, err := s.eggService.GetEggMonitoringToday(chickenCage.Id, today)
		if err != nil {
			s.log.Error("failed get egg monitoring today")
			return dto.ChickenAndWarehouseOverviewResponse{}, err
		}

		totalChicken += chickenCage.TotalChicken
		totalFeed += chickenMonitoring.TotalFeed
		totalDeathChicken += chickenMonitoring.TotalDeathChicken
		totalEgg += (eggMonitoring.TotalKarpetCrackedEgg * constant.TotalEggPerKarpet) + eggMonitoring.TotalRemainingCrackedEgg + (eggMonitoring.TotalKarpetGoodEgg * constant.TotalEggPerKarpet) + eggMonitoring.TotalRemainingGoodEgg
		totalWeightEgg += eggMonitoring.TotalWeightAllEgg

		count := chickenCage.TotalChicken
		switch chickenCage.Cage.ChickenCategory {
		case enum.ChickenCategoryDOC.String():
			totalDOCChicken += count
		case enum.ChickenCategoryGrower.String():
			totalGrowerChicken += count
		case enum.ChickenCategoryPreLayer.String():
			totalPreLayerChicken += count
		case enum.ChickenCategoryLayer.String():
			totalLayerChicken += count
		case enum.ChickenCategoryAfkir.String():
			totalAfkirChicken += count
		}
	}

	return dto.ChickenAndWarehouseOverviewResponse{
		ChickenPerformanceSummary: dto.ChickenPerformanceSummaryResponse{
			FeedConsumption:      totalFeed,
			AverageEggWeight:     totalWeightEgg / float64(totalEgg) * 1000.0,
			AverageMortalityRate: float64(totalDeathChicken) / float64(totalChicken),
			AverageFCR:           totalFeed / float64(totalEgg),
			AverageHDP:           float64(totalEgg) / float64(totalChicken),
		},
		ChickenBarCharts: dto.ChickenBarChartResponse{
			ChickenDOC:       float64(totalDOCChicken),
			ChickenGrower:    float64(totalGrowerChicken),
			ChickentPreLayer: float64(totalPreLayerChicken),
			ChickenLayer:     float64(totalLayerChicken),
			ChickenAfkir:     float64(totalAfkirChicken),
		},
		WarehouseItemSummary: dto.WarehouseItemSummaryResponse{
			TotalSafeItem:    totalSafeItem,
			TotalNotSafeItem: totalDangerItem,
		},
		ChickenCagePerformanceSummary: dto.ChickenCagePerformanceSummaryResponse{
			TotalProductiveCage: totalProductiveCage,
			TotalAfkirCage:      totalAfkirCage,
		},
		ChickenGraphs: chickenGraphs,
	}, nil
}

func (s *ChickenService) GetChickenAndCompanyOverview(filter dto.GetChickenAndCompanyOverviewFilter) (dto.ChickenAndCompanyOverviewResponse, error) {
	// R/C Ratio = Keuntungan (Return) / Biaya Produksi (Cost)
	// Keuntungan (Return) = Penjualan Telur + Penjualan Ayam - Biaya Produksi (Pengaadan Ayam, Pengadaan Barang, Pengadaan Jagung, Gaji, Operational, Pajak)

	// Mos (%) = (Penjualan Aktual)

	s.repository.UseTx(false)

	today := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	var totalFeed, totalWeightEgg float64
	var totalEgg, totalDeathChicken, totalChicken uint64
	var totalDOCChicken, totalGrowerChicken, totalPreLayerChicken, totalLayerChicken, totalAfkirChicken uint64
	chickenCages, err := s.cageService.GetChickenCages(dto.GetChickenCageFilter{
		LocationId: filter.LocationId,
		CageId:     filter.CageId,
	})
	if err != nil {
		s.log.Error("failed get chicken cages", zap.Error(err))
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	for _, chickenCage := range chickenCages {
		chickenMonitoring, err := s.repository.GetChickenMonitoringToday(chickenCage.Id, today)
		if err != nil {
			s.log.Error("failed get chicken monitoring today", zap.Error(err))
			return dto.ChickenAndCompanyOverviewResponse{}, err
		}

		eggMonitoring, err := s.eggService.GetEggMonitoringToday(chickenCage.Id, today)
		if err != nil {
			s.log.Error("failed get egg monitoring today")
			return dto.ChickenAndCompanyOverviewResponse{}, err
		}

		totalChicken += chickenCage.TotalChicken
		totalFeed += chickenMonitoring.TotalFeed
		totalDeathChicken += chickenMonitoring.TotalDeathChicken
		totalEgg += (eggMonitoring.TotalKarpetCrackedEgg * constant.TotalEggPerKarpet) + eggMonitoring.TotalRemainingCrackedEgg + (eggMonitoring.TotalKarpetGoodEgg * constant.TotalEggPerKarpet) + eggMonitoring.TotalRemainingGoodEgg
		totalWeightEgg += eggMonitoring.TotalWeightAllEgg

		count := chickenCage.TotalChicken
		switch chickenCage.Cage.ChickenCategory {
		case enum.ChickenCategoryDOC.String():
			totalDOCChicken += count
		case enum.ChickenCategoryGrower.String():
			totalGrowerChicken += count
		case enum.ChickenCategoryPreLayer.String():
			totalPreLayerChicken += count
		case enum.ChickenCategoryLayer.String():
			totalLayerChicken += count
		case enum.ChickenCategoryAfkir.String():
			totalAfkirChicken += count
		}
	}

	totalIncomeProduction, err := s.cashflowService.GetTotalIncomeProductionInMonth(enum.Month(time.Now().Month()), uint64(time.Now().Year()))
	if err != nil {
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	totalExpenseProduction, err := s.cashflowService.GetTotalExpenseProductionInMonth(enum.Month(time.Now().Month()), uint64(time.Now().Year()))
	if err != nil {
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	goodEgg, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	goodEggItemPrice, err := s.itemService.GetItemPriceByItemIdAndSaleUnit(goodEgg.Id, enum.SaleUnitKg.String())
	if err != nil {
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	diff := totalExpenseProduction.Sub(totalIncomeProduction)
	price, err := decimal.NewFromString(goodEggItemPrice.Price)
	if err != nil {
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	bepGoodEgg := diff.Div(price)
	rcRatio := diff.Div(totalExpenseProduction).InexactFloat64() * 100.0
	mos := totalIncomeProduction.Sub(diff).Sub(totalIncomeProduction).InexactFloat64()

	incomeAndExpenseGraphs := make([]dto.IncomeAndExpenseBarChartResponse, 0)
	cashflowHistories, err := s.cashflowService.GetCashflowHistories(dto.GetCashflowHistoryFilter{
		Year:       uint64(time.Now().Year()),
		LocationId: filter.LocationId,
	})

	if err != nil {
		return dto.ChickenAndCompanyOverviewResponse{}, err
	}

	for _, cashflowHistory := range cashflowHistories {
		incomeAndExpenseGraphs = append(incomeAndExpenseGraphs, dto.IncomeAndExpenseBarChartResponse{
			Key:     cashflowHistory.CreatedAt.Format("January"),
			Income:  cashflowHistory.Income,
			Expense: cashflowHistory.Cash,
		})
	}

	return dto.ChickenAndCompanyOverviewResponse{
		ChickenPerformanceSummary: dto.ChickenPerformanceSummaryResponse{
			FeedConsumption:      totalFeed,
			AverageEggWeight:     totalWeightEgg / float64(totalEgg) * 1000.0,
			AverageMortalityRate: float64(totalDeathChicken) / float64(totalChicken),
			AverageFCR:           totalFeed / float64(totalEgg),
			AverageHDP:           float64(totalEgg) / float64(totalChicken),
		},
		ChickenBarCharts: dto.ChickenBarChartResponse{
			ChickenDOC:       float64(totalDOCChicken),
			ChickenGrower:    float64(totalGrowerChicken),
			ChickentPreLayer: float64(totalPreLayerChicken),
			ChickenLayer:     float64(totalLayerChicken),
			ChickenAfkir:     float64(totalAfkirChicken),
		},
		BEPGoodEgg:                           bepGoodEgg.InexactFloat64(),
		MarginOfSafety:                       mos,
		RCRatio:                              rcRatio,
		IncomeAndExpensePerformanceBarCharts: incomeAndExpenseGraphs,
	}, nil
}
