package service

import (
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

	CreateChickenDiseaseMonitoring(chickenMonitoringId uint64, request dto.CreateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	UpdateChickenDiseaseMonitoring(id uint64, request dto.UpdateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	DeleteChickenDiseaseMonitoring(id uint64) error

	CreateChickenVaccineMonitoring(chickenMonitoringId uint64, request dto.CreateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	UpdateChickenVaccineMonitoring(id uint64, request dto.UpdateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error)
	DeleteChickenVaccineMonitoring(id uint64) error
}

func NewChickenService(log *zap.Logger, repository repository.IChickenRepository, eggService IEggService) IChickenService {
	return &ChickenService{
		log:        log,
		repository: repository,
		eggService: eggService,
	}
}

func (s *ChickenService) CreateChickenMonitoring(request dto.CreateChickenMonitoringRequest, createdBy uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	s.repository.UseTx(false)

	count, err := s.repository.CountChickenMonitoringByCageIdToday(request.ChickenCageId)
	if err != nil {
		s.log.Error("[CreateChickenMonitoring] failed to count chicken monitoring by cage id", zap.Error(err))
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
		s.log.Error("[CreateChickenMonitoring] failed to create chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = s.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		s.log.Error("[CreateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	return mapper.ChickenMonitoringToResponse(&chickenMonitoring), nil
}

func (s *ChickenService) GetChickenMonitoringById(id uint64) (dto.ChickenMonitoringResponse, error) {
	chickenMonitoring, err := s.repository.GetChickenMonitoringById(id)
	if err != nil {
		s.log.Error("[GetChickenMonitoringById] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	return mapper.ChickenMonitoringToResponse(&chickenMonitoring), nil
}

func (c *ChickenService) GetChickenMonitorings(filter dto.GetChickenMonitoringFilter) ([]dto.ChickenMonitoringListResponse, error) {
	chickenMonitorings, err := c.repository.GetChickenMonitorings(&filter)
	if err != nil {
		c.log.Error("[GetChickenMonitorings] failed to get chicken monitorings", zap.Error(err))
		return []dto.ChickenMonitoringListResponse{}, err
	}

	chickenMonitoringsResponse := make([]dto.ChickenMonitoringListResponse, len(chickenMonitorings))
	for i, chickenMonitoring := range chickenMonitorings {
		chickenMonitoringsResponse[i] = mapper.ChickenMonitoringToListResponse(&chickenMonitoring)
	}

	return chickenMonitoringsResponse, nil
}

func (c *ChickenService) UpdateChickenMonitoring(id uint64, request dto.UpdateChickenMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	c.repository.UseTx(false)
	chickenMonitoring, err := c.repository.GetChickenMonitoringById(id)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring.TotalSickChicken = request.TotalSickChicken
	chickenMonitoring.TotalDeathChicken = request.TotalDeathChicken
	chickenMonitoring.TotalFeed = request.TotalFeed
	chickenMonitoring.UpdateBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = c.repository.UpdateChickenMonitoring(&chickenMonitoring)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to update chicken monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err = c.repository.GetChickenMonitoringById(chickenMonitoring.Id)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoringResponse := mapper.ChickenMonitoringToResponse(&chickenMonitoring)

	return chickenMonitoringResponse, nil
}

func (c *ChickenService) CreateChickenDiseaseMonitoring(chickenMonitoringId uint64, request dto.CreateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenDisease := entity.ChickenDiseaseMonitoring{
		ChickenMonitoringId: chickenMonitoringId,
		Disease:             request.Disease,
		Medicine:            request.Medicine,
		Dose:                request.Dose,
		Unit:                request.Unit,
		CreatedBy:           uuid.NullUUID{UUID: accountId, Valid: true},
	}

	err := c.repository.CreateChickenDiseaseMonitoring(&[]entity.ChickenDiseaseMonitoring{chickenDisease})
	if err != nil {
		c.log.Error("[CreateChickenDiseaseMonitoring] failed to create chicken disease monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoringResponse := mapper.ChickenMonitoringToResponse(&chickenMonitoring)

	return chickenMonitoringResponse, nil
}

func (c *ChickenService) CreateChickenVaccineMonitoring(chickenMonitoringId uint64, request dto.CreateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenVaccine := entity.ChickenVaccineMonitoring{
		ChickenMonitoringId: chickenMonitoringId,
		Vaccine:             request.Vaccine,
		Dose:                request.Dose,
		Unit:                request.Unit,
		CreatedBy:           uuid.NullUUID{UUID: accountId, Valid: true},
	}

	err := c.repository.CreateChickenVaccineMonitoring(&[]entity.ChickenVaccineMonitoring{chickenVaccine})
	if err != nil {
		c.log.Error("[CreateChickenVaccineMonitoring] failed to create chicken vaccine monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoringResponse := mapper.ChickenMonitoringToResponse(&chickenMonitoring)

	return chickenMonitoringResponse, nil
}

func (c *ChickenService) UpdateChickenDiseaseMonitoring(id uint64, request dto.UpdateChickenDiseaseMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenDisease, err := c.repository.GetChickenDiseaseMonitoringById(id)
	if err != nil {
		c.log.Error("[UpdateChickenDiseaseMonitoring] failed to get chicken disease monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenDisease.Disease = request.Disease
	chickenDisease.Medicine = request.Medicine
	chickenDisease.Dose = request.Dose
	chickenDisease.Unit = request.Unit
	chickenDisease.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = c.repository.UpdateChickenDiseaseMonitoring(&chickenDisease)
	if err != nil {
		c.log.Error("[UpdateChickenDiseaseMonitoring] failed to update chicken disease monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenDisease.ChickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoringResponse := mapper.ChickenMonitoringToResponse(&chickenMonitoring)

	return chickenMonitoringResponse, nil
}

func (c *ChickenService) UpdateChickenVaccineMonitoring(id uint64, request dto.UpdateChickenVaccineMonitoringRequest, accountId uuid.UUID) (dto.ChickenMonitoringResponse, error) {
	chickenVaccine, err := c.repository.GetChickenVaccineMonitoringById(id)
	if err != nil {
		c.log.Error("[UpdateChickenVaccineMonitoring] failed to get chicken vaccine monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenVaccine.Vaccine = request.Vaccine
	chickenVaccine.Dose = request.Dose
	chickenVaccine.Unit = request.Unit
	chickenVaccine.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = c.repository.UpdateChickenVaccineMonitoring(&chickenVaccine)
	if err != nil {
		c.log.Error("[UpdateChickenVaccineMonitoring] failed to update chicken vaccine monitoring", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoring, err := c.repository.GetChickenMonitoringById(chickenVaccine.ChickenMonitoringId)
	if err != nil {
		c.log.Error("[UpdateChickenMonitoring] failed to get chicken monitoring by id", zap.Error(err))
		return dto.ChickenMonitoringResponse{}, err
	}

	chickenMonitoringResponse := mapper.ChickenMonitoringToResponse(&chickenMonitoring)

	return chickenMonitoringResponse, nil
}

func (c *ChickenService) DeleteChickenMonitoring(id uint64) error {
	err := c.repository.DeleteChickenMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenMonitoring] failed to delete chicken monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (c *ChickenService) DeleteChickenDiseaseMonitoring(id uint64) error {
	err := c.repository.DeleteChickenDiseaseMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenDiseaseMonitoring] failed to delete chicken disease monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (c *ChickenService) DeleteChickenVaccineMonitoring(id uint64) error {
	err := c.repository.DeleteChickenVaccineMonitoring(id)
	if err != nil {
		c.log.Error("[DeleteChickenVaccineMonitoring] failed to delete chicken vaccine monitoring", zap.Error(err))
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
		Date:     param.DateParam(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
		Location: filter.Location,
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
