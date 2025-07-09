package service

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"time"

	"github.com/google/uuid"
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
	"go.uber.org/zap"
)

type EggService struct {
	log              *zap.Logger
	repository       repository.IEggRepository
	warehouseService IWarehouseService
	cageService      ICageService
	itemService      IItemService
	cacheService     cache.ICache
	storeService     IStoreService
}

type IEggService interface {
	CreateEggMonitoring(request dto.CreateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error)
	GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error)
	GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error)
	UpdateEggMonitoring(id uint64, request dto.UpdateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error)
	DeleteEggMonitoring(id uint64, userId uuid.UUID) error

	GetOverviewEggMonitoring(filter dto.GetEggOverviewFilter) (dto.EggOverviewResponse, error)
}

func NewEggService(
	log *zap.Logger,
	repository repository.IEggRepository,
	warehouseService IWarehouseService,
	cageService ICageService,
	itemService IItemService,
	cacheService cache.ICache,
	storeService IStoreService,
) IEggService {
	return &EggService{
		log:              log,
		repository:       repository,
		warehouseService: warehouseService,
		cageService:      cageService,
		itemService:      itemService,
		cacheService:     cacheService,
		storeService:     storeService,
	}
}

// Todo : add created at in response to avoid delete after one day and saga pattern
func (s *EggService) CreateEggMonitoring(request dto.CreateEggMonitoringRequest, createdBy uuid.UUID) (dto.EggMonitoringResponse, error) {
	s.repository.UseTx(false)

	count, err := s.repository.CountEggMonitoringByChickenCageIdToday(request.ChickenCageId)
	if err != nil {
		s.log.Error("failed to count egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	if count > 0 {
		s.log.Error("egg monitoring already exists for today", zap.Error(errx.BadRequest("egg monitoring already exists for today")))
		return dto.EggMonitoringResponse{}, errx.BadRequest("egg monitoring already exists for today")
	}

	chickenCage, err := s.cageService.GetChickenCageById(request.ChickenCageId)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring := entity.EggMonitoring{
		ChickenCageId:         request.ChickenCageId,
		WarehouseId:           request.WarehouseId,
		TotalWeightCrackedEgg: request.TotalWeightCrackedEgg,
		TotalWeightGoodEgg:    request.TotalWeightGoodEgg,
		TotalGoodEgg:          (request.TotalKarpetGoodEgg * uint64(constant.TotalEggPerKarpet)) + request.TotalRemainingGoodEgg,
		TotalCrackedEgg:       (request.TotalKarpetCrackedEgg * uint64(constant.TotalEggPerKarpet)) + request.TotalRemainingCrackedEgg,
		TotalRejectEgg:        (request.TotalKarpetRejectEgg * uint64(constant.TotalEggPerKarpet)) + request.TotalRemainingRejectEgg,
		CreatedBy:             uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType("Telur OK", "Kg", enum.ItemCategoryEgg)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	goodEggWarehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(eggMonitoring.WarehouseId, goodEggItem.Id)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	jsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemId:         goodEggItem.Id,
		Source:         chickenCage.Cage.Name,
		Destination:    goodEggWarehouseItem.Warehouse.Name,
		QuantityBefore: goodEggWarehouseItem.Quantity,
		QuantityAfter:  eggMonitoring.TotalWeightGoodEgg + goodEggWarehouseItem.Quantity,
		UserId:         createdBy,
		Status:         enum.ActivityStatusIn,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.EggMonitoringResponse{}, errx.BadRequest("failed parsed struct into json")
	}
	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, string(jsonParsed))

	_, err = s.warehouseService.UpdateWarehouseItem(eggMonitoring.WarehouseId, goodEggItem.Id, dto.UpdateWarehouseItemRequest{
		Quantity: goodEggWarehouseItem.Quantity + eggMonitoring.TotalWeightGoodEgg,
	}, createdBy)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType("Telur Retak", "Kg", enum.ItemCategoryEgg)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	crackedEggWarehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(eggMonitoring.WarehouseId, crackedEggItem.Id)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	crackedEggJsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemId:         crackedEggItem.Id,
		Source:         chickenCage.Cage.Name,
		Destination:    crackedEggWarehouseItem.Warehouse.Name,
		QuantityBefore: crackedEggWarehouseItem.Quantity,
		QuantityAfter:  eggMonitoring.TotalWeightCrackedEgg + crackedEggWarehouseItem.Quantity,
		UserId:         createdBy,
		Status:         enum.ActivityStatusIn,
	})
	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.EggMonitoringResponse{}, errx.BadRequest("failed parsed struct into json")
	}
	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, string(crackedEggJsonParsed))

	_, err = s.warehouseService.UpdateWarehouseItem(eggMonitoring.WarehouseId, crackedEggItem.Id, dto.UpdateWarehouseItemRequest{
		Quantity: crackedEggWarehouseItem.Quantity + eggMonitoring.TotalWeightCrackedEgg,
	}, createdBy)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	if err := s.repository.CreateEggMonitoring(&eggMonitoring); err != nil {
		s.log.Error("failed to create egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	_, err = s.storeService.CreateStoreRequestItemFromEggMonitoring(dto.CreateStoreRequestItemRequest{
		WarehouseId: request.WarehouseId,
		Quantity:    request.TotalWeightCrackedEgg,
		ItemId:      crackedEggItem.Id,
	}, createdBy)
	if err != nil {
		s.log.Error("failed to create store request item from egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = s.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		s.log.Error("failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	return mapper.EggMonitoringToResponse(&eggMonitoring), nil
}

func (s *EggService) GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error) {
	s.repository.UseTx(false)

	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (s *EggService) GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error) {
	s.repository.UseTx(false)

	eggMonitorings, err := s.repository.GetEggMonitorings(filter)
	if err != nil {
		s.log.Error("failed to get egg monitorings", zap.Error(err))
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
	s.repository.UseTx(false)

	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get egg monitoring by id", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType("Telur OK", "Kg", enum.ItemCategoryEgg)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	goodEggWarehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(eggMonitoring.WarehouseId, goodEggItem.Id)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	goodEggJsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemId:         goodEggItem.Id,
		Source:         eggMonitoring.ChickenCage.Cage.Name,
		Destination:    goodEggWarehouseItem.Warehouse.Name,
		QuantityBefore: goodEggWarehouseItem.Quantity,
		QuantityAfter:  goodEggWarehouseItem.Quantity - eggMonitoring.TotalWeightGoodEgg + request.TotalWeightGoodEgg,
		UserId:         updatedBy,
		Status:         enum.ActivityStatusIn,
	})
	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.EggMonitoringResponse{}, errx.BadRequest("failed parsed struct into json")
	}
	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, string(goodEggJsonParsed))

	_, err = s.warehouseService.UpdateWarehouseItem(eggMonitoring.WarehouseId, goodEggItem.Id, dto.UpdateWarehouseItemRequest{
		Quantity: goodEggWarehouseItem.Quantity - eggMonitoring.TotalWeightGoodEgg + request.TotalWeightGoodEgg,
	}, updatedBy)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType("Telur Retak", "Kg", enum.ItemCategoryEgg)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	crackedEggWarehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(eggMonitoring.WarehouseId, crackedEggItem.Id)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	crackedEggJsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemId:         crackedEggItem.Id,
		Source:         eggMonitoring.ChickenCage.Cage.Name,
		Destination:    crackedEggWarehouseItem.Warehouse.Name,
		QuantityBefore: crackedEggWarehouseItem.Quantity,
		QuantityAfter:  crackedEggWarehouseItem.Quantity - eggMonitoring.TotalWeightCrackedEgg + request.TotalWeightCrackedEgg,
		UserId:         updatedBy,
		Status:         enum.ActivityStatusIn,
	})
	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.EggMonitoringResponse{}, errx.BadRequest("failed parsed struct into json")
	}
	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, string(crackedEggJsonParsed))

	_, err = s.warehouseService.UpdateWarehouseItem(eggMonitoring.WarehouseId, crackedEggItem.Id, dto.UpdateWarehouseItemRequest{
		Quantity: crackedEggWarehouseItem.Quantity - eggMonitoring.TotalWeightCrackedEgg + request.TotalWeightCrackedEgg,
	}, updatedBy)
	if err != nil {
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring.ChickenCageId = request.ChickenCageId
	eggMonitoring.WarehouseId = request.WarehouseId
	eggMonitoring.TotalGoodEgg = (request.TotalKarpetGoodEgg * uint64(constant.TotalEggPerKarpet)) + request.TotalRemainingGoodEgg
	eggMonitoring.TotalCrackedEgg = (request.TotalKarpetCrackedEgg * uint64(constant.TotalEggPerKarpet)) + request.TotalRemainingCrackedEgg
	eggMonitoring.TotalRejectEgg = (request.TotalKarpetRejectEgg * uint64(constant.TotalEggPerKarpet)) + request.TotalRemainingRejectEgg
	eggMonitoring.TotalWeightCrackedEgg = request.TotalWeightCrackedEgg
	eggMonitoring.TotalWeightGoodEgg = request.TotalWeightGoodEgg
	eggMonitoring.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	if err := s.repository.UpdateEggMonitoring(&eggMonitoring); err != nil {
		s.log.Error("failed to update egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get egg monitoring by id", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	return mapper.EggMonitoringToResponse(&eggMonitoring), nil
}

func (s *EggService) DeleteEggMonitoring(id uint64, updatedBy uuid.UUID) error {
	s.repository.UseTx(false)

	eggMonitoring, err := s.repository.GetEggMonitoringById(id)
	if err != nil {
		s.log.Error("failed to get egg monitoring by id", zap.Error(err))
		return err
	}

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType("Telur OK", "Kg", enum.ItemCategoryEgg)
	if err != nil {
		return err
	}

	goodEggWarehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(eggMonitoring.WarehouseId, goodEggItem.Id)
	if err != nil {
		return err
	}

	goodEggJsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemId:         goodEggItem.Id,
		Source:         eggMonitoring.ChickenCage.Cage.Name,
		Destination:    goodEggWarehouseItem.Warehouse.Name,
		QuantityBefore: goodEggWarehouseItem.Quantity,
		QuantityAfter:  goodEggWarehouseItem.Quantity - eggMonitoring.TotalWeightGoodEgg,
		UserId:         updatedBy,
		Status:         enum.ActivityStatusIn,
	})
	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return errx.BadRequest("failed parsed struct into json")
	}
	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, string(goodEggJsonParsed))

	_, err = s.warehouseService.UpdateWarehouseItem(eggMonitoring.WarehouseId, goodEggItem.Id, dto.UpdateWarehouseItemRequest{
		Quantity: goodEggWarehouseItem.Quantity - eggMonitoring.TotalWeightGoodEgg,
	}, updatedBy)
	if err != nil {
		return err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType("Telur Retak", "Kg", enum.ItemCategoryEgg)
	if err != nil {
		return err
	}

	crackedEggWarehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(eggMonitoring.WarehouseId, crackedEggItem.Id)
	if err != nil {
		return err
	}

	crackedEggJsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemId:         crackedEggItem.Id,
		Source:         eggMonitoring.ChickenCage.Cage.Name,
		Destination:    crackedEggWarehouseItem.Warehouse.Name,
		QuantityBefore: crackedEggWarehouseItem.Quantity,
		QuantityAfter:  crackedEggWarehouseItem.Quantity - eggMonitoring.TotalWeightCrackedEgg,
		UserId:         updatedBy,
		Status:         enum.ActivityStatusIn,
	})
	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return errx.BadRequest("failed parsed struct into json")
	}
	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, string(crackedEggJsonParsed))

	_, err = s.warehouseService.UpdateWarehouseItem(eggMonitoring.WarehouseId, crackedEggItem.Id, dto.UpdateWarehouseItemRequest{
		Quantity: crackedEggWarehouseItem.Quantity - eggMonitoring.TotalWeightCrackedEgg,
	}, updatedBy)
	if err != nil {
		return err
	}

	if err := s.repository.DeleteEggMonitoring(id); err != nil {
		s.log.Error("failed to delete egg monitoring", zap.Error(err))
		return err
	}

	return nil
}

func (s *EggService) GetOverviewEggMonitoring(filter dto.GetEggOverviewFilter) (dto.EggOverviewResponse, error) {
	currentEggMonitorings, err := s.repository.GetEggMonitorings(dto.GetEggMonitoringFilter{
		LocationId: filter.Location,
		Date:       param.DateParam(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)),
	})
	if err != nil {
		s.log.Error(" failed to get egg monitorings", zap.Error(err))
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
			LocationId: filter.Location,
			StartDate:  param.DateParam(startDate),
			EndDate:    param.DateParam(endDate),
		})
		if err != nil {
			s.log.Error("failed to get egg monitorings", zap.Error(err))
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
			LocationId: filter.Location,
			StartDate:  param.DateParam(startDate),
			EndDate:    param.DateParam(endDate),
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
			LocationId: filter.Location,
			StartDate:  param.DateParam(startDate),
			EndDate:    param.DateParam(endDate),
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
