package service

import (
	"database/sql"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
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

	count, err := s.repository.CountChickenMonitoringByCageIdToday(request.ChickenCageId)
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

	// Todo : create if there are death chicken

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

func (s *ChickenService) UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	s.repository.UseTx(false)
	chickenMonitoring, err := s.repository.GetChickenMonitoringById(id)
	if err != nil {
		s.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring.TotalSickChicken = request.TotalSickChicken
	chickenMonitoring.TotalDeathChicken = request.TotalDeathChicken
	chickenMonitoring.TotalFeed = request.TotalFeed
	chickenMonitoring.Note = request.Note
	chickenMonitoring.UpdateBy = uuid.NullUUID{UUID: accountId, Valid: true}

	// Todo : create if there are death chicken

	err = s.repository.UpdateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		s.log.Error("[UpdateChickenMonitoring] failed to update chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = s.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		s.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	return mapper.ChickenMonitoringToResponse(&chickenMonitoring), nil
}

func (c *ChickenService) DeleteChickenMonitoring(id uint64) error {
	err := c.repository.DeleteChickenMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenMonitoring] failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (c *ChickenService) GetChickenOverview(filter dto.GetChickenOverviewFilter) (dto.ChickenOverviewResponse, error) {
	c.repository.UseTx(false)

	// currentChickenMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
	// 	Date:     param.DateParam(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
	// 	Location: filter.Location,
	// })
	// if err != nil {
	// 	c.log.Error("[GetChickenOverview] failed to get chicken monitorings", zap.Error(err))
	// 	return dto.ChickenOverviewResponse{}, err
	// }

	currentEggMonitoring, err := c.eggService.GetEggMonitorings(dto.GetEggMonitoringFilter{
		Date:       param.DateParam(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
		LocationId: filter.Location,
	})
	if err != nil {
		c.log.Error("[GetChickenOverview] failed to get egg monitorings", zap.Error(err))
		return dto.ChickenOverviewResponse{}, err
	}

	totalEgg := uint64(0)
	for _, eggMonitoring := range currentEggMonitoring {
		totalEgg += eggMonitoring.TotalAllEgg
	}

	totalDOCChicken := uint64(0)
	totalGrowerChicken := uint64(0)
	totalPreLayerChicken := uint64(0)
	totalLayerChicken := uint64(0)
	totalAfkirChicken := uint64(0)

	totalLiveChicken := uint64(0)
	totalSickChicken := uint64(0)
	totalDeathChicken := uint64(0)

	chickenGraphs := make([]dto.ChickenGraphResponse, 0)

	// for _, chickenMonitoring := range currentChickenMonitorings {
	// 	totalLiveChicken += chickenMonitoring.TotalLiveChicken
	// 	totalSickChicken += chickenMonitoring.TotalSickChicken
	// 	totalDeathChicken += chickenMonitoring.TotalDeathChicken

	// 	if chickenMonitoring.ChickenCategory == enum.ChickenCategoryDOC {
	// 		totalDOCChicken += chickenMonitoring.TotalSickChicken + chickenMonitoring.TotalLiveChicken
	// 	} else if chickenMonitoring.ChickenCategory == enum.ChickenCategoryGrower {
	// 		totalGrowerChicken += chickenMonitoring.TotalSickChicken + chickenMonitoring.TotalLiveChicken
	// 	} else if chickenMonitoring.ChickenCategory == enum.ChickenCategoryPreLayer {
	// 		totalPreLayerChicken += chickenMonitoring.TotalSickChicken + chickenMonitoring.TotalLiveChicken
	// 	} else if chickenMonitoring.ChickenCategory == enum.ChickenCategoryLayer {
	// 		totalLayerChicken += chickenMonitoring.TotalSickChicken + chickenMonitoring.TotalLiveChicken
	// 	} else if chickenMonitoring.ChickenCategory == enum.ChickenCategoryAfkir {
	// 		totalAfkirChicken += chickenMonitoring.TotalSickChicken + chickenMonitoring.TotalLiveChicken
	// 	}
	// }

	if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisWeek {
		endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		startDate := endDate.AddDate(0, 0, -7)

		weekChickenMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})

		if err != nil {
			c.log.Error("[GetChickenOverview] failed to get chicken monitorings", zap.Error(err))
			return dto.ChickenOverviewResponse{}, err
		}

		for i := startDate; i.Before(endDate); i = i.AddDate(0, 0, 1) {
			for _, chickenMonitoring := range weekChickenMonitorings {
				if i.Equal(chickenMonitoring.CreatedAt) {
					chickenGraphs = append(chickenGraphs, dto.ChickenGraphResponse{
						Key:          i.Format("2006-01-02"),
						SickChicken:  chickenMonitoring.TotalSickChicken,
						DeathChicken: chickenMonitoring.TotalDeathChicken,
					})
				} else {
					chickenGraphs = append(chickenGraphs, dto.ChickenGraphResponse{
						Key:          i.Format("2006-01-02"),
						SickChicken:  0,
						DeathChicken: 0,
					})
				}
			}
		}
	} else if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisMonth {
		weekMaps := util.GetFourWeekRanges(time.Now().Year(), time.Now().Month())

		totalSickChickenGraph := make(map[int]uint64)
		totalDeathChickenGraph := make(map[int]uint64)

		startDate, endDate := util.GetStartDateAndEndDateInMonth(time.Now().Year(), time.Now().Month())

		monthChickenMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			c.log.Error("[GetChickenOverview] failed to get chicken monitorings", zap.Error(err))
			return dto.ChickenOverviewResponse{}, err
		}

		for _, chickenMonitoring := range monthChickenMonitorings {
			i := util.FindWeek(chickenMonitoring.CreatedAt, weekMaps)

			if i > 0 {
				totalSickChickenGraph[i] += chickenMonitoring.TotalSickChicken
				totalDeathChickenGraph[i] += chickenMonitoring.TotalDeathChicken
			}
		}

		keys := make([]int, 0)
		for k := range weekMaps {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, key := range keys {
			chickenGraphs = append(chickenGraphs, dto.ChickenGraphResponse{
				Key:          fmt.Sprintf("Minggu %d", key),
				SickChicken:  totalSickChickenGraph[key],
				DeathChicken: totalDeathChickenGraph[key],
			})
		}

	} else if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisYear {
		monthMaps := util.GetTwelveMonthRanges(time.Now().Year())

		totalSickChickenGraph := make(map[int]uint64)
		totalDeathChickenGraph := make(map[int]uint64)

		startDate, endDate := util.GetStartDateAndEndDateInYear(time.Now().Year())

		yearChickenMonitorings, err := c.repository.GetChickenMonitorings(&dto.GetChickenMonitoringFilter{
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			c.log.Error("[GetChickenOverview] failed to get chicken monitorings", zap.Error(err))
			return dto.ChickenOverviewResponse{}, err
		}

		for _, chickenMonitoring := range yearChickenMonitorings {
			i := util.FindMonth(chickenMonitoring.CreatedAt, monthMaps)

			if i > 0 {
				totalSickChickenGraph[i] += chickenMonitoring.TotalSickChicken
				totalDeathChickenGraph[i] += chickenMonitoring.TotalDeathChicken

			}
		}

		keys := make([]int, 0)
		for k := range monthMaps {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for _, key := range keys {
			chickenGraphs = append(chickenGraphs, dto.ChickenGraphResponse{
				Key:          time.Month(key).String(),
				SickChicken:  totalSickChickenGraph[key],
				DeathChicken: totalDeathChickenGraph[key],
			})
		}
	}

	mortalityRate := float64(totalDeathChicken) / float64(totalLiveChicken+totalDeathChicken+totalSickChicken)
	hdpRate := float64(totalEgg) / float64(totalLiveChicken+totalDeathChicken+totalSickChicken)

	chickenDetail := dto.ChickenDetailOverview{
		TotalLiveChicken:    totalLiveChicken,
		TotalSickChicken:    totalSickChicken,
		TotalDeathChicken:   totalDeathChicken,
		TotalKPIPerformance: (mortalityRate + hdpRate) / 2,
	}

	chickenPie := dto.ChickenPieResponse{
		ChickenDOCType:       float64(totalDOCChicken) / float64(totalLiveChicken+totalSickChicken),
		ChickenGrowerType:    float64(totalGrowerChicken) / float64(totalLiveChicken+totalSickChicken),
		ChickentPreLayerType: float64(totalPreLayerChicken) / float64(totalLiveChicken+totalSickChicken),
		ChickenLayer:         float64(totalLayerChicken) / float64(totalLiveChicken+totalSickChicken),
		ChickenAfkir:         float64(totalAfkirChicken) / float64(totalLiveChicken+totalSickChicken),
	}

	return dto.ChickenOverviewResponse{
		ChickenDetail: chickenDetail,
		ChickenPie:    chickenPie,
		ChickenGraphs: chickenGraphs,
	}, nil
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

	data := entity.ChickenHealthMonitoring{
		ChickenCageId:  request.ChickenCageId,
		HealthItemName: request.HealthItemName,
		Type:           chickenHealthMonitoringType,
		Dose:           request.Dose,
		Unit:           request.Unit,
		CreatedBy:      uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	if request.Disease != nil {
		data.Disease = sql.NullString{String: *request.Disease, Valid: true}
	}

	err := s.repository.CreateChickenHealthMonitoring(&data)
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
