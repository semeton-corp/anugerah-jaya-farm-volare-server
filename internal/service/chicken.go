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
	"go.uber.org/zap"
)

type ChickenService struct {
	log         *zap.Logger
	repository  repository.IChickenRepository
	eggService  IEggService
	cageService ICageService
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
}

func NewChickenService(log *zap.Logger, repository repository.IChickenRepository, eggService IEggService, cageService ICageService) IChickenService {
	return &ChickenService{
		log:         log,
		repository:  repository,
		eggService:  eggService,
		cageService: cageService,
	}
}

func (s *ChickenService) CreateChickenMonitoring(request dto.CreateChickenMonitoringRequest, createdBy uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	s.repository.UseTx(false)

	count, err := s.repository.CountChickenMonitoringByChickenCageIdToday(request.ChickenCageId)
	if err != nil {
		s.log.Error("failed to count chicken monitoring by cage id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	if count > 0 {
		return dto.ChickenMonitoringResponse{}, errx.BadRequest("chicken monitoring already exists for today")
	}

	chickenMonitoring := entity.ChickenMonitoring{
		ChickenCageId:     request.ChickenCageId,
		TotalDeathChicken: request.TotalDeathChicken,
		TotalSickChicken:  request.TotalSickChicken,
		TotalFeed:         request.TotalFeed,
		Note:              request.Note,
		CreatedBy:         uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	// Todo : create if there are death chicken in chicken cage

	err = s.repository.CreateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		s.log.Error("failed to create chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = s.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
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

func (s *ChickenService) UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, updateBy uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	s.repository.UseTx(false)
	chickenMonitoring, err := s.repository.GetChickenMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	count, err := s.repository.CountChickenMonitoringByChickenCageIdToday(request.ChickenCageId)
	if err != nil {
		s.log.Error("failed count chicken monitoring by chicken cage id today", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	if count > 0 {
		return dto.ChickenMonitoringResponse{}, errx.BadRequest("chicken cage id is already use for another monitoring")
	}

	chickenMonitoring.TotalSickChicken = request.TotalSickChicken
	chickenMonitoring.TotalDeathChicken = request.TotalDeathChicken
	chickenMonitoring.TotalFeed = request.TotalFeed
	chickenMonitoring.Note = request.Note
	chickenMonitoring.UpdateBy = uuid.NullUUID{UUID: updateBy, Valid: true}
	chickenMonitoring.ChickenCageId = request.ChickenCageId

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
		totalLiveChicken += cm.ChickenCage.TotalChicken - cm.TotalSickChicken - cm.TotalDeathChicken
		totalSickChicken += cm.TotalSickChicken
		totalDeathChicken += cm.TotalDeathChicken

		count := cm.TotalSickChicken + cm.ChickenCage.TotalChicken - cm.ChickenCage.TotalDeathChicken
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
		chickenGraphs, err = c.buildWeeklyGraph()
	case enum.OverviewGraphTimeThisMonth:
		chickenGraphs, err = c.buildMonthlyGraph()
	case enum.OverviewGraphTimeThisYear:
		chickenGraphs, err = c.buildYearlyGraph()
	}
	if err != nil {
		return dto.ChickenOverviewResponse{}, err
	}

	denominator := float64(totalLiveChicken + totalDeathChicken + totalSickChicken)
	if denominator == 0 {
		denominator = 1
	}
	mortalityRate := float64(totalDeathChicken) / denominator
	hdpRate := float64(totalEgg) / denominator

	return dto.ChickenOverviewResponse{
		ChickenDetail: dto.ChickenDetailOverview{
			TotalLiveChicken:    totalLiveChicken,
			TotalSickChicken:    totalSickChicken,
			TotalDeathChicken:   totalDeathChicken,
			TotalKPIPerformance: (mortalityRate + hdpRate) / 2,
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

func (c *ChickenService) buildWeeklyGraph() ([]dto.ChickenGraphResponse, error) {
	endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	startDate := endDate.AddDate(0, 0, -7)

	weekMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		c.log.Error("failed to get chicken monitorings", zap.Error(err))
		return nil, err
	}

	graphs := make([]dto.ChickenGraphResponse, 0)
	for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 1) {
		var sickSum, deathSum uint64
		for _, cm := range weekMonitorings {
			if isSameDate(day, cm.CreatedAt) {
				sickSum += cm.TotalSickChicken
				deathSum += cm.TotalDeathChicken
			}
		}
		graphs = append(graphs, dto.ChickenGraphResponse{
			Key:          day.Format("2006-01-02"),
			SickChicken:  sickSum,
			DeathChicken: deathSum,
		})
	}
	return graphs, nil
}

func (c *ChickenService) buildMonthlyGraph() ([]dto.ChickenGraphResponse, error) {
	weekMaps := util.GetFourWeekRanges(time.Now().Year(), time.Now().Month())
	startDate, endDate := util.GetStartDateAndEndDateInMonth(time.Now().Year(), time.Now().Month())

	monthMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
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

func (c *ChickenService) buildYearlyGraph() ([]dto.ChickenGraphResponse, error) {
	monthMaps := util.GetTwelveMonthRanges(time.Now().Year())
	startDate, endDate := util.GetStartDateAndEndDateInYear(time.Now().Year())

	yearMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
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

func isSameDate(a, b time.Time) bool {
	return a.Year() == b.Year() && a.Month() == b.Month() && a.Day() == b.Day()
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

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price from string", zap.Error(err))
		return dto.ChickenProcurementDraftResponse{}, err
	}

	data := entity.ChickenProcurementDraft{
		CageId:     request.CageId,
		SupplierId: request.SupplierId,
		Quantity:   request.Quantity,
		Price:      price,
		TotalPrice: price.Mul(decimal.NewFromUint64(request.Quantity)),
		CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateChickenProcurementDraft(&data)
	if err != nil {
		return dto.ChickenProcurementDraftResponse{}, err
	}

	chickenProcurementDraft, err := s.repository.GetChickenProcurementDraftById(data.Id)
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

func (s *ChickenService) ConfirmationChickenProcurementDraft(id uint64, request dto.ConfirmedChickenProcurementRequest, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurementDraft, err := s.repository.GetChickenProcurementDraftById(id)
	if err != nil {
		return err
	}

	cage, err := s.cageService.GetCageById(chickenProcurementDraft.CageId)
	if err != nil {
		return err
	}

	if cage.IsUsed {
		return errx.BadRequest("cage is in used by another chicken")
	}

	chickenProcurementDraft.IsConfirmed = true

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price from string", zap.Error(err))
		return err
	}

	estimateArrivalDate, err := time.Parse("02-01-2006", request.EstimateArrivalDate)
	if err != nil {
		s.log.Error("failed to parse estimate arrival date", zap.Error(err))
		return err
	}

	chickenProcurement := entity.ChickenProcurement{
		CageId:              chickenProcurementDraft.CageId,
		SupplierId:          chickenProcurementDraft.SupplierId,
		Quantity:            request.Quantity,
		Price:               price,
		TotalPrice:          price.Mul(decimal.NewFromUint64(request.Quantity)),
		Status:              enum.ProcurementStatusSentOff,
		PaymentStatus:       enum.PaymentStatusNotPaid,
		EstimateArrivalDate: estimateArrivalDate,
		CreatedBy:           uuid.NullUUID{UUID: userId, Valid: true},
	}

	chickenProcurementPayments := make([]entity.ChickenProcurementPayment, 0)
	totalPayment := decimal.Zero
	for _, payment := range request.Payments {
		paymentDate, err := time.Parse("02-01-2006", payment.PaymentDate)
		if err != nil {
			s.log.Error("failed to parse payment date", zap.Error(err))
			return err
		}

		nominal, err := decimal.NewFromString(payment.Nominal)
		if err != nil {
			s.log.Error("failed to parse nominal from string", zap.Error(err))
			return err
		}

		paymentMethod := enum.ValueOfPaymentMethod(payment.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("invalid payment method", zap.String("paymentMethod", payment.PaymentMethod))
			return errx.BadRequest("invalid payment method")
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
	} else if totalPayment.GreaterThan(decimal.Zero) {
		chickenProcurement.PaymentStatus = enum.PaymentStatusUnpaid
	}

	err = s.repository.UpdateChickenProcurementDraft(&chickenProcurementDraft)
	if err != nil {
		s.log.Error("failed to update chicken procurement draft", zap.Error(err))
		return err
	}

	err = s.repository.CreateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed to create chicken procurement", zap.Error(err))
		return err
	}

	if chickenProcurementPayments != nil {
		for i := range chickenProcurementPayments {
			chickenProcurementPayments[i].ChickenProcurementId = chickenProcurement.Id
		}

		err = s.repository.CreateChickenProcurementPaymentInBatch(&chickenProcurementPayments)
		if err != nil {
			s.log.Error("failed to create chicken procurement in batch", zap.Error(err))
			return err
		}
	}

	err = s.repository.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (s *ChickenService) ArrivalConfirmationChickenProcurement(id uint64, request dto.ArrivalConfirmationChickenProcurementRequest, userId uuid.UUID) error {
	s.repository.UseTx(false)

	chickenProcurement, err := s.repository.GetChickenProcurementById(id)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return err
	}

	chickenProcurement.RecieveQuantity = sql.NullInt64{Int64: int64(request.Quantity), Valid: true}
	chickenProcurement.Note = request.Note
	chickenProcurement.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	chickenProcurement.TakenBy = uuid.NullUUID{UUID: userId, Valid: true}
	chickenProcurement.IsArrived = true
	chickenProcurement.Status = enum.ProcurementStatusArrived

	// Sage pattern
	_, err = s.cageService.CreateChickenCage(dto.CreateChickenCageRequest{}, userId)
	if err != nil {
		return err
	}

	err = s.repository.UpdateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed update chicken procurement", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) CreateChickenProcurementPayment(chickenProcurementId uint64, request dto.CreateChickenProcurementPaymentRequest, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurement, err := s.repository.GetChickenProcurementById(chickenProcurementId)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return err
	}

	if chickenProcurement.PaymentStatus == enum.PaymentStatusPaid {
		return errx.BadRequest("chicken procurement is already paid")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return errx.BadRequest("invalid nominal format")
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
		return errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.UpdateChickenProcurement(&chickenProcurement)
	if err != nil {
		s.log.Error("failed update chicken procurement", zap.Error(err))
		return err
	}

	err = s.repository.CreateChickenProcurementPayment(&chickenProcurementPayment)
	if err != nil {
		s.log.Error("failed create chicken procurement payment", zap.Error(err))
		return err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) UpdateChickenProcurementPayment(chickenProcurementId uint64, id uint64, request dto.UpdateChickenProcurementPaymentRequest, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurement, err := s.repository.GetChickenProcurementById(chickenProcurementId)
	if err != nil {
		s.log.Error("failed to get chicken procurement by id", zap.Error(err))
		return err
	}

	chickenProcurementPayment, err := s.repository.GetChickenProcurementPaymentById(id)
	if err != nil {
		s.log.Error("failed to get chicken procurement payment by id", zap.Error(err))
		return err
	}

	if chickenProcurement.PaymentStatus == enum.PaymentStatusPaid {
		return errx.BadRequest("chicken procurement is already paid")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return errx.BadRequest("invalid nominal format")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return errx.BadRequest("invalid payment method")
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
		return errx.BadRequest("total payment is greater than total price")
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
		return err
	}

	err = s.repository.UpdateChickenProcurementPayment(&chickenProcurementPayment)
	if err != nil {
		s.log.Error("failed create chicken procurement payment", zap.Error(err))
		return err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *ChickenService) DeleteChickenProcurement(chickenProcurementId uint64, id uint64) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenProcurement, err := s.repository.GetChickenProcurementById(chickenProcurementId)
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

	return nil
}

func (s *ChickenService) MoveChickenCage(request dto.MoveChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	cageIds := make([]uint64, 0)
	for _, e := range request.DestinationChickenCages {
		cageIds = append(cageIds, e.DestinationCageId)
	}

	chickenCages, err := s.cageService.GetChickenCagesByCageIds(cageIds)
	if err != nil {
		return nil, err
	}

	for _, chickenCage := range chickenCages {
		if chickenCage.Cage.IsUsed {
			return nil, errx.BadRequest(fmt.Sprintf("cage with id %d is used", chickenCage.Cage.Id))
		}
	}

	newChickenCages := make([]dto.CreateChickenCageRequest, 0)
	for _, destinationChickenCage := range request.DestinationChickenCages {
		newChickenCage := dto.CreateChickenCageRequest{
			CageId:       destinationChickenCage.DestinationCageId,
			TotalChicken: destinationChickenCage.TotalChicken,
		}

		for _, chickenCage := range chickenCages {
			if chickenCage.Cage.Id == destinationChickenCage.DestinationCageId {
				newChickenCage.ChickenProcurementId = chickenCage.ChickenProcurementId
			}
		}

		newChickenCages = append(newChickenCages, newChickenCage)
	}

	response, err := s.cageService.CreateChickenCageInBatch(newChickenCages, userId)
	if err != nil {
		return nil, err
	}

	return response, nil
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
		return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment type : %s", request.PaymentType))
	}

	data := entity.AfkirChickenSale{
		AfkirChickenCustomerId: request.AfkirChickenCustomerId,
		ChickenCageId:          request.ChickenCageId,
		TotalSellChicken:       request.TotalSellChicken,
		PricePerChicken:        pricePerChicken,
		TotalPrice:             pricePerChicken.Mul(decimal.NewFromUint64(request.TotalSellChicken)),
		ChickenAge:             chickenCage.ChickenAge,
		PaymentStatus:          enum.PaymentStatusNotPaid,
		PaymentType:            paymentType,
		CreatedBy:              uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.AfkirChickenSalePayment != nil {
		nominal, err := decimal.NewFromString(request.AfkirChickenSalePayment.Nominal)
		if err != nil {
			s.log.Error("failed parse nominal from string", zap.Error(err))
			return dto.AfkirChickenSaleResponse{}, err
		}

		paymentDate, err := time.Parse("02-01-2006", request.AfkirChickenSalePayment.PaymentDate)
		if err != nil {
			s.log.Error("failed parse payment date", zap.Error(err))
			return dto.AfkirChickenSaleResponse{}, err
		}

		paymentMethod := enum.ValueOfPaymentMethod(request.AfkirChickenSalePayment.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.AfkirChickenSaleResponse{}, errx.BadRequest(fmt.Sprintf("invalid payment method : %s", request.AfkirChickenSalePayment.PaymentMethod))
		}

		if data.TotalPrice.Equal(nominal) {
			data.PaymentStatus = enum.PaymentStatusPaid
		} else if nominal.GreaterThan(decimal.Zero) {
			data.PaymentStatus = enum.PaymentStatusUnpaid
		}

		err = s.repository.CreateAfkirChickenSale(&data)
		if err != nil {
			s.log.Error("failed create afkir chicken sale", zap.Error(err))
			return dto.AfkirChickenSaleResponse{}, nil
		}

		payment := entity.AfkirChickenSalePayment{
			AfkirChickenSaleId: data.AfkirChickenCustomerId,
			Nominal:            nominal,
			PaymentDate:        paymentDate,
			PaymentMethod:      paymentMethod,
			PaymentProof:       request.AfkirChickenSalePayment.PaymentProof,
			CreatedBy:          uuid.NullUUID{UUID: userId, Valid: true},
		}

		err = s.repository.CreateAfkirChickenSalePayment(&payment)
		if err != nil {
			return dto.AfkirChickenSaleResponse{}, err
		}
	} else {
		err = s.repository.CreateAfkirChickenSale(&data)
		if err != nil {
			s.log.Error("failed create afkir chicken sale", zap.Error(err))
			return dto.AfkirChickenSaleResponse{}, nil
		}
	}

	_, err = s.cageService.UpdateCage(chickenCage.Cage.Id, dto.UpdateCageRequest{
		IsUsed: false,
	}, userId)
	if err != nil {
		return dto.AfkirChickenSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	data, err = s.repository.GetAfkirChickenSale(data.Id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	return mapper.AfkirChickenSaleToResponse(&data), nil
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
		response.TotalPage = uint64(totalData) / constant.PaginationDefaultLimit
	}

	return response, nil
}

func (s *ChickenService) GetAkfirChickenSale(id uint64) (dto.AfkirChickenSaleResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetAfkirChickenSale(id)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	return mapper.AfkirChickenSaleToResponse(&data), nil
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
		return dto.AfkirChickenSaleResponse{}, err
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

	afkirChickenSale, err = s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	return mapper.AfkirChickenSaleToResponse(&afkirChickenSale), nil
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
		return dto.AfkirChickenSaleResponse{}, err
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

	if totalCurrentPrice.Add(nominal).LessThan(afkirChickenSale.TotalPrice) {
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

	afkirChickenSale, err = s.repository.GetAfkirChickenSale(afkirChickenSaleId)
	if err != nil {
		s.log.Error("failed get afkir chicken sale", zap.Error(err))
		return dto.AfkirChickenSaleResponse{}, err
	}

	return mapper.AfkirChickenSaleToResponse(&afkirChickenSale), nil
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
