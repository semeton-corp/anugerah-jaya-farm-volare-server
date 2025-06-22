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
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
)

type EggService struct {
	log              *zap.Logger
	repository       repository.IEggRepository
	warehouseService IWarehouseService
	cageService      ICageService
}

type IEggService interface {
	CreateEggMonitoring(request dto.CreateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error)
	GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error)
	GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error)
	UpdateEggMonitoring(id uint64, request dto.UpdateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error)
	DeleteEggMonitoring(id uint64) error

	TakeEggMonitoring(id uint64, accountId uuid.UUID) (dto.EggMonitoringResponse, error)

	GetOverviewEggMonitoring(filter dto.GetEggOverviewFilter) (dto.EggOverviewResponse, error)
}

func NewEggService(
	log *zap.Logger,
	repository repository.IEggRepository,
	warehouseService IWarehouseService,
	cageService ICageService,
) IEggService {
	return &EggService{
		log:              log,
		repository:       repository,
		warehouseService: warehouseService,
		cageService:      cageService,
	}
}

func (s *EggService) CreateEggMonitoring(request dto.CreateEggMonitoringRequest, createdBy uuid.UUID) (dto.EggMonitoringResponse, error) {
	chickenCage, err := s.cageService.GetChickenCageByCageId(request.CageId)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	if chickenCage.TotalChicken == 0 {
		s.log.Warn("[CreateEggMonitoring] cage is empty")
		return dto.EggMonitoringResponse{}, errx.BadRequest("cage is empty or no chicken in there")
	}

	count, err := s.repository.CountEggMonitoringByChickenCageIdToday(chickenCage.Id)
	if err != nil {
		s.log.Error("[CreateEggMonitoring] failed to count egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	if count > 0 {
		s.log.Error("[CreateEggMonitoring] egg monitoring already exists for today", zap.Error(errx.BadRequest("egg monitoring already exists for today")))
		return dto.EggMonitoringResponse{}, errx.BadRequest("egg monitoring already exists for today")
	}

	eggMonitoring := entity.EggMonitoring{
		ChickenCageId:         chickenCage.Id,
		WarehouseId:           request.WarehouseId,
		TotalWeightCrackedEgg: request.TotalWeightCrackedEgg,
		TotalWeightGoodEgg:    request.TotalWeightGoodEgg,
		TotalGoodEgg:          (request.TotalKarpetGoodEgg * uint64(constant.TotalEggKarpet)) + request.TotalRemainingGoodEgg,
		TotalCrackedEgg:       (request.TotalKarpetCrackedEgg * uint64(constant.TotalEggKarpet)) + request.TotalRemainingCrackedEgg,
		TotalRejectEgg:        (request.TotalKarpetRejectEgg * uint64(constant.TotalEggKarpet)) + request.TotalRemainingRejectEgg,
		IsTaken:               false,
		CreatedBy:             uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	if err := s.repository.CreateEggMonitoring(&eggMonitoring); err != nil {
		s.log.Error("[CreateEggMonitoring] failed to create egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = s.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		s.log.Error("[CreateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (s *EggService) GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("[GetEggMonitoringById] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (s *EggService) GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error) {
	eggMonitorings, err := s.repository.GetEggMonitorings(filter)
	if err != nil {
		s.log.Error("[GetEggMonitorings] failed to get egg monitorings", zap.Error(err))
		return nil, err
	}

	eggMonitoringResponses := make([]dto.EggMonitoringListResponse, 0, len(eggMonitorings))
	for _, eggMonitoring := range eggMonitorings {
		// Todo : create notification (?)
		eggMonitoringResponses = append(eggMonitoringResponses, mapper.EggMonitoringToListResponse(&eggMonitoring))
	}

	return eggMonitoringResponses, nil
}

func (s *EggService) UpdateEggMonitoring(id uint64, request dto.UpdateEggMonitoringRequest, updatedBy uuid.UUID) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("[UpdateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	chickenCage, err := s.cageService.GetChickenCageByCageId(request.CageId)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring.ChickenCageId = chickenCage.Id
	eggMonitoring.WarehouseId = request.WarehouseId
	eggMonitoring.TotalGoodEgg = (request.TotalKarpetGoodEgg * uint64(constant.TotalEggKarpet)) + request.TotalRemainingGoodEgg
	eggMonitoring.TotalCrackedEgg = (request.TotalKarpetCrackedEgg * uint64(constant.TotalEggKarpet)) + request.TotalRemainingCrackedEgg
	eggMonitoring.TotalRejectEgg = (request.TotalKarpetRejectEgg * uint64(constant.TotalEggKarpet)) + request.TotalRemainingRejectEgg
	eggMonitoring.TotalWeightCrackedEgg = request.TotalWeightCrackedEgg
	eggMonitoring.TotalWeightGoodEgg = request.TotalWeightGoodEgg
	eggMonitoring.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	if err := s.repository.UpdateEggMonitoring(&eggMonitoring); err != nil {
		s.log.Error("[UpdateEggMonitoring] failed to update egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = s.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		s.log.Error("[UpdateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (s *EggService) DeleteEggMonitoring(id uint64) error {
	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("[DeleteEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return err
	}

	if err := s.repository.DeleteEggMonitoring(eggMonitoring.Id); err != nil {
		s.log.Error("[DeleteEggMonitoring] failed to delete egg monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (s *EggService) TakeEggMonitoring(id uint64, accountId uuid.UUID) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("[TakeEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	if eggMonitoring.IsTaken {
		s.log.Error("[TakeEggMonitoring] egg monitoring already taken", zap.Error(errx.BadRequest("egg monitoring already taken")))
		return dto.EggMonitoringResponse{}, errx.BadRequest("egg monitoring already taken")
	}

	eggMonitoring.IsTaken = true
	eggMonitoring.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	eggMonitoring.TakenBy = uuid.NullUUID{UUID: accountId, Valid: true}
	eggMonitoring.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	// Todo : add stock into warehouse stock item

	if err := s.repository.UpdateEggMonitoring(&eggMonitoring); err != nil {
		s.log.Error("[TakeEggMonitoring] failed to update egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = s.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		s.log.Error("[TakeEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	return mapper.EggMonitoringToResponse(&eggMonitoring), nil
}

func (s *EggService) GetOverviewEggMonitoring(filter dto.GetEggOverviewFilter) (dto.EggOverviewResponse, error) {
	currentEggMonitorings, err := s.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
		Location: filter.Location,
		Date:     param.DateParam(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
	})
	if err != nil {
		s.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
		return dto.EggOverviewResponse{}, err
	}

	totalGoodEgg := uint64(0)
	totalCrackedEgg := uint64(0)
	totalBrokenEgg := uint64(0)
	totalRejectEgg := uint64(0)

	for _, eggMonitoring := range currentEggMonitorings {
		totalGoodEgg += eggMonitoring.TotalGoodEgg
		totalCrackedEgg += eggMonitoring.TotalCrackedEgg
		// totalBrokenEgg += eggMonitoring.TotalBrokeEgg
		totalRejectEgg += eggMonitoring.TotalRejectEgg
	}

	eggGraphs := make([]dto.EggGraphResponse, 0)

	if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisWeek {
		endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		startDate := endDate.AddDate(0, 0, -7)

		weekEggMonitorings, err := s.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
			Location:  filter.Location,
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			s.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
			return dto.EggOverviewResponse{}, err
		}

		for i := startDate; i.Before(endDate); i = i.AddDate(0, 0, 1) {
			for _, eggMonitoring := range weekEggMonitorings {
				if i.Equal(eggMonitoring.CreatedAt) {
					eggGraphs = append(eggGraphs, dto.EggGraphResponse{
						Key:        i.Format("2006-01-02"),
						GoodEgg:    eggMonitoring.TotalGoodEgg,
						CrackedEgg: eggMonitoring.TotalCrackedEgg,
						// BrokenEgg:  eggMonitoring.TotalBrokeEgg,
						RejectEgg: eggMonitoring.TotalRejectEgg,
					})
				} else {
					eggGraphs = append(eggGraphs, dto.EggGraphResponse{
						Key:        i.Format("2006-01-02"),
						GoodEgg:    0,
						CrackedEgg: 0,
						BrokenEgg:  0,
						RejectEgg:  0,
					})
				}
			}
		}
	} else if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisMonth {
		weekMaps := util.GetFourWeekRanges(time.Now().Year(), time.Now().Month())
		startDate, endDate := util.GetStartDateAndEndDateInMonth(time.Now().Year(), time.Now().Month())

		monthEggMonitorings, err := s.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
			Location:  filter.Location,
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			s.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
			return dto.EggOverviewResponse{}, err
		}

		goodEggMaps := make(map[int]uint64)
		crackedEggMaps := make(map[int]uint64)
		brokenEggMaps := make(map[int]uint64)
		rejectEggMaps := make(map[int]uint64)

		for _, eggMonitoring := range monthEggMonitorings {
			i := util.FindWeek(eggMonitoring.CreatedAt, weekMaps)
			if i != 0 {
				goodEggMaps[i] += eggMonitoring.TotalGoodEgg
				crackedEggMaps[i] += eggMonitoring.TotalCrackedEgg
				// brokenEggMaps[i] += eggMonitoring.TotalBrokeEgg
				rejectEggMaps[i] += eggMonitoring.TotalRejectEgg
			}
		}

		keys := make([]int, 0)
		for k := range weekMaps {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for i := range keys {
			eggGraphs = append(eggGraphs, dto.EggGraphResponse{
				Key:        fmt.Sprintf("Minggu %d", i),
				GoodEgg:    goodEggMaps[i],
				CrackedEgg: crackedEggMaps[i],
				BrokenEgg:  brokenEggMaps[i],
				RejectEgg:  rejectEggMaps[i],
			})
		}

	} else if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisYear {
		monthMaps := util.GetTwelveMonthRanges(time.Now().Year())
		startDate, endDate := util.GetStartDateAndEndDateInYear(time.Now().Year())

		yearEggMonitorings, err := s.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
			Location:  filter.Location,
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			s.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
			return dto.EggOverviewResponse{}, err
		}

		goodEggMaps := make(map[int]uint64)
		crackedEggMaps := make(map[int]uint64)
		brokenEggMaps := make(map[int]uint64)
		rejectEggMaps := make(map[int]uint64)

		for _, eggMonitoring := range yearEggMonitorings {
			i := util.FindMonth(eggMonitoring.CreatedAt, monthMaps)
			if i != 0 {
				goodEggMaps[i] += eggMonitoring.TotalGoodEgg
				crackedEggMaps[i] += eggMonitoring.TotalCrackedEgg
				// brokenEggMaps[i] += eggMonitoring.TotalBrokeEgg
				rejectEggMaps[i] += eggMonitoring.TotalRejectEgg
			}
		}

		keys := make([]int, 0)
		for k := range monthMaps {
			keys = append(keys, k)
		}
		sort.Ints(keys)

		for i := range keys {
			eggGraphs = append(eggGraphs, dto.EggGraphResponse{
				Key:        time.Month(i).String(),
				GoodEgg:    goodEggMaps[i],
				CrackedEgg: crackedEggMaps[i],
				BrokenEgg:  brokenEggMaps[i],
				RejectEgg:  rejectEggMaps[i],
			})
		}
	}

	eggOverviewDetail := dto.EggOverviewDetailResponse{
		TotalGoodEgg:    totalGoodEgg,
		TotalCrackedEgg: totalCrackedEgg,
		TotalBrokenEgg:  totalBrokenEgg,
		TotalRejectEgg:  totalRejectEgg,
	}

	eggOverview := dto.EggOverviewResponse{
		EggOverviewDetail: eggOverviewDetail,
		EggGraphs:         eggGraphs,
	}

	return eggOverview, nil
}
