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
) IEggService {
	return &EggService{
		log:              log,
		repository:       repository,
		warehouseService: warehouseService,
	}
}

func (e *EggService) CreateEggMonitoring(request dto.CreateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error) {
	count, err := e.repository.CountEggMonitoringByCageIdToday(request.CageId)
	if err != nil {
		e.log.Error("[CreateEggMonitoring] failed to count egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	if count > 0 {
		e.log.Error("[CreateEggMonitoring] egg monitoring already exists for today", zap.Error(errx.BadRequest("egg monitoring already exists for today")))
		return dto.EggMonitoringResponse{}, errx.BadRequest("egg monitoring already exists for today")
	}

	eggMonitoring := entity.EggMonitoring{
		CageId:          request.CageId,
		WarehouseId:     request.WarehouseId,
		TotalGoodEgg:    request.TotalGoodEgg,
		TotalCrackedEgg: request.TotalCrackedEgg,
		TotalBrokeEgg:   request.TotalBrokeEgg,
		TotalRejectEgg:  request.TotalRejectEgg,
		Weight:          request.Weight,
		CreatedBy:       accountId,
		IsArrive:        false,
	}

	if err := e.repository.CreateEggMonitoring(&eggMonitoring); err != nil {
		e.log.Error("[CreateEggMonitoring] failed to create egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = e.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		e.log.Error("[CreateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (e *EggService) GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[GetEggMonitoringById] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (e *EggService) GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error) {
	eggMonitorings, err := e.repository.GetEggMonitorings(filter)
	if err != nil {
		e.log.Error("[GetEggMonitorings] failed to get egg monitorings", zap.Error(err))
		return nil, err
	}

	eggMonitoringResponses := make([]dto.EggMonitoringListResponse, 0, len(eggMonitorings))
	for _, eggMonitoring := range eggMonitorings {
		eggMonitoringResponse := mapper.EggMonitoringToListResponse(&eggMonitoring)

		if eggMonitoringResponse.TotalAll == 0 {
			eggMonitoringResponse.AbnormalityRate = 0
			eggMonitoringResponse.Description = constant.EggMonitoringStatusSafety
		} else {
			eggMonitoringResponse.AbnormalityRate = float64(eggMonitoring.TotalCrackedEgg+eggMonitoring.TotalBrokeEgg+eggMonitoring.TotalRejectEgg) / float64(eggMonitoringResponse.TotalAll) * 100
			eggMonitoringResponse.Description = constant.EggMonitoringStatusSafety
		}

		eggMonitoringResponses = append(eggMonitoringResponses, eggMonitoringResponse)
	}

	return eggMonitoringResponses, nil
}

func (e *EggService) UpdateEggMonitoring(id uint64, request dto.UpdateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[UpdateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring.Weight = request.Weight
	eggMonitoring.CageId = request.CageId
	eggMonitoring.WarehouseId = request.WarehouseId
	eggMonitoring.TotalGoodEgg = request.TotalGoodEgg
	eggMonitoring.TotalCrackedEgg = request.TotalCrackedEgg
	eggMonitoring.TotalBrokeEgg = request.TotalBrokeEgg
	eggMonitoring.TotalRejectEgg = request.TotalRejectEgg
	eggMonitoring.UpdatedBy = accountId

	if err := e.repository.UpdateEggMonitoring(&eggMonitoring); err != nil {
		e.log.Error("[UpdateEggMonitoring] failed to update egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = e.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		e.log.Error("[UpdateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (e *EggService) DeleteEggMonitoring(id uint64) error {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[DeleteEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return err
	}

	if err := e.repository.DeleteEggMonitoring(eggMonitoring.Id); err != nil {
		e.log.Error("[DeleteEggMonitoring] failed to delete egg monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (e *EggService) TakeEggMonitoring(id uint64, accountId uuid.UUID) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[TakeEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	if eggMonitoring.IsArrive {
		e.log.Error("[TakeEggMonitoring] egg monitoring already taken", zap.Error(errx.BadRequest("egg monitoring already taken")))
		return dto.EggMonitoringResponse{}, errx.BadRequest("egg monitoring already taken")
	}

	eggMonitoring.IsArrive = true
	eggMonitoring.TakenAt = time.Now()
	eggMonitoring.TakenBy = accountId
	eggMonitoring.UpdatedBy = accountId

	// Todo : add stock into warehouse stock item

	if err := e.repository.UpdateEggMonitoring(&eggMonitoring); err != nil {
		e.log.Error("[TakeEggMonitoring] failed to update egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = e.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		e.log.Error("[TakeEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	return mapper.EggMonitoringToResponse(&eggMonitoring), nil
}

func (e *EggService) GetOverviewEggMonitoring(filter dto.GetEggOverviewFilter) (dto.EggOverviewResponse, error) {
	currentEggMonitorings, err := e.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
		Location: filter.Location,
		Date:     param.DateParam(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
	})
	if err != nil {
		e.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
		return dto.EggOverviewResponse{}, err
	}

	totalGoodEgg := uint64(0)
	totalCrackedEgg := uint64(0)
	totalBrokenEgg := uint64(0)
	totalRejectEgg := uint64(0)

	for _, eggMonitoring := range currentEggMonitorings {
		totalGoodEgg += eggMonitoring.TotalGoodEgg
		totalCrackedEgg += eggMonitoring.TotalCrackedEgg
		totalBrokenEgg += eggMonitoring.TotalBrokeEgg
		totalRejectEgg += eggMonitoring.TotalRejectEgg
	}

	eggGraphs := make([]dto.EggGraphResponse, 0)

	if filter.OverviewGraphTime.Value() == enum.OverviewGraphTimeThisWeek {
		endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
		startDate := endDate.AddDate(0, 0, -7)

		weekEggMonitorings, err := e.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
			Location:  filter.Location,
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			e.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
			return dto.EggOverviewResponse{}, err
		}

		for i := startDate; i.Before(endDate); i = i.AddDate(0, 0, 1) {
			for _, eggMonitoring := range weekEggMonitorings {
				if i.Equal(eggMonitoring.CreatedAt) {
					eggGraphs = append(eggGraphs, dto.EggGraphResponse{
						Key:        i.Format("2006-01-02"),
						GoodEgg:    eggMonitoring.TotalGoodEgg,
						CrackedEgg: eggMonitoring.TotalCrackedEgg,
						BrokenEgg:  eggMonitoring.TotalBrokeEgg,
						RejectEgg:  eggMonitoring.TotalRejectEgg,
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

		monthEggMonitorings, err := e.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
			Location:  filter.Location,
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			e.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
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
				brokenEggMaps[i] += eggMonitoring.TotalBrokeEgg
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

		yearEggMonitorings, err := e.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
			Location:  filter.Location,
			StartDate: param.DateParam(startDate),
			EndDate:   param.DateParam(endDate),
		})
		if err != nil {
			e.log.Error("[GetOverviewEggMonitoring] failed to get egg monitorings", zap.Error(err))
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
				brokenEggMaps[i] += eggMonitoring.TotalBrokeEgg
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
