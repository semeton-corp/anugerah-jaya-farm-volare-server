package service

import (
	"database/sql"
	"fmt"
	"math"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
	"go.uber.org/zap"
)

type CageService struct {
	log              *zap.Logger
	repository       repository.ICageRepository
	warehouseService IWarehouseService
}

type ICageService interface {
	GetCages(filter dto.GetCageFilter) ([]dto.CageResponse, error)
	CreateCage(request dto.CreateCageRequest, userId uuid.UUID) (dto.CageResponse, error)
	UpdateCage(id uint64, request dto.UpdateCageRequest, userId uuid.UUID) (dto.CageResponse, error)
	DeleteCage(id uint64) error
	GetCageById(id uint64) (dto.CageResponse, error)
	GetCagesByIds(ids []uint64) ([]dto.CageResponse, error)

	CreateChickenCage(request dto.CreateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error)
	GetChickenCages(filter dto.GetChickenCageFilter) ([]dto.ChickenCageResponse, error)
	GetChickenCageById(id uint64) (dto.ChickenCageResponse, error)
	UpdateChickenCage(id uint64, request dto.UpdateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error)

	GetChickenCageFeeds(filter dto.GetChickenCageFeedFilter) ([]dto.ChickenCageFeedListResponse, error)
	GetChickenCageFeed(chickenCageId uint64) (dto.ChickenCageFeedResponse, error)

	CreateCageFeed(request dto.CreateCageFeedRequest, userId uuid.UUID) (dto.CageFeedResponse, error)
	UpdateCageFeed(id uint64, request dto.UpdateCageFeedRequest, userId uuid.UUID) (dto.CageFeedResponse, error)
	GetCageFeeds() ([]dto.CageFeedResponse, error)
	GetCageFeed(id uint64) (dto.CageFeedResponse, error)

	ConfirmationChickenCageFeed(chickenCageId uint64, request dto.ConfirmationChickenCageFeedRequest, userId uuid.UUID) (dto.ChickenCageFeedResponse, error)

	MoveChickenCage(request dto.MoveChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error)

	ReduceCageFeedStocks(latestFeed float64, currentFeeds float64, cageId uint64) error
}

func NewCageService(log *zap.Logger, repository repository.ICageRepository, warehouseService IWarehouseService) ICageService {
	return &CageService{
		log:              log,
		repository:       repository,
		warehouseService: warehouseService,
	}
}

func (s *CageService) GetCages(filter dto.GetCageFilter) ([]dto.CageResponse, error) {
	s.repository.UseTx(false)

	cages, err := s.repository.GetCages(filter)
	if err != nil {
		s.log.Error("failed to get cages", zap.Error(err))
		return nil, err
	}

	cageResponses := make([]dto.CageResponse, 0)
	for _, cage := range cages {
		cageResponses = append(cageResponses, mapper.CageToResponse(&cage))
	}

	return cageResponses, nil
}

func (s *CageService) CreateCage(request dto.CreateCageRequest, userId uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("invalid chicken category")
		return dto.CageResponse{}, errx.BadRequest("invalid chicken category")
	}

	cage := entity.Cage{
		LocationId:      request.LocationId,
		Name:            request.Name,
		Capacity:        request.Capacity,
		ChickenCategory: chickenCategory,
		CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateCage(&cage)
	if err != nil {
		s.log.Error("failed to create cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	chickenCage := entity.ChickenCage{
		CageId:    cage.Id,
		CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed to create chicken cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transcation", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage, err = s.repository.GetCageById(cage.Id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) UpdateCage(id uint64, request dto.UpdateCageRequest, userId uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("invalid chicken category")
		return dto.CageResponse{}, errx.BadRequest("invalid chicken category")
	}

	cage, err := s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage.Name = request.Name
	cage.LocationId = request.LocationId
	cage.Capacity = request.Capacity
	cage.ChickenCategory = chickenCategory
	cage.IsUsed = *request.IsUsed
	cage.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateCage(&cage)
	if err != nil {
		s.log.Error("failed to update cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage, err = s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) DeleteCage(id uint64) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	cage, err := s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed get cage by id", zap.Error(err))
		return err
	}

	if cage.IsUsed {
		return errx.BadRequest("cage is in used, please make cage empty first")
	}

	// Soft delete related chicken cages first
	err = s.repository.DeleteChickenCageByCageId(id)
	if err != nil {
		s.log.Error("failed to delete chicken cages", zap.Error(err))
		return err
	}

	// Then delete the cage (this will trigger BeforeDelete hook for name update)
	err = s.repository.DeleteCage(id)
	if err != nil {
		s.log.Error("failed to delete cage", zap.Error(err))
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *CageService) GetChickenCages(filter dto.GetChickenCageFilter) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCageResponses := make([]dto.ChickenCageResponse, 0)
	chickenCages, err := s.repository.GetChickenCages(filter)
	if err != nil {
		return nil, err
	}

	for _, chickenCage := range chickenCages {
		chickenCageResponses = append(chickenCageResponses, mapper.ChickenCageToResponse(&chickenCage))
	}

	return chickenCageResponses, nil
}

func (s *CageService) GetChickenCageById(id uint64) (dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageById(id)
	if err != nil {
		return dto.ChickenCageResponse{}, err
	}

	return mapper.ChickenCageToResponse(&chickenCage), nil
}

func (s *CageService) UpdateChickenCage(id uint64, request dto.UpdateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageById(id)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	chickenCage.TotalChicken = request.TotalChicken
	chickenCage.IsNeedRoutineVaccine = request.IsNeedRoutineVaccine
	chickenCage.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if request.LatestChickenAgeVaccineRoutine != nil {
		chickenCage.LatestChickenAgeVaccineRoutine = sql.NullInt64{Int64: int64(*request.LatestChickenAgeVaccineRoutine), Valid: true}
	}

	err = s.repository.UpdateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed update chicken cage", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	return mapper.ChickenCageToResponse(&chickenCage), nil
}

func (s *CageService) GetCageById(id uint64) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	cage, err := s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) CreateChickenCage(request dto.CreateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCage := entity.ChickenCage{
		CageId:       request.CageId,
		TotalChicken: request.TotalChicken,
		CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.ChickenProcurementId != nil {
		chickenCage.ChickenProcurementId = sql.NullInt64{Int64: int64(*request.ChickenProcurementId), Valid: true}
	}

	err := s.repository.CreateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed to create chicken cage", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	chickenCage, err = s.repository.GetChickenCageById(chickenCage.Id)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	return mapper.ChickenCageToResponse(&chickenCage), err
}

func (s *CageService) GetCagesByIds(ids []uint64) ([]dto.CageResponse, error) {
	s.repository.UseTx(false)

	cages, err := s.repository.GetCagesByIds(ids)
	if err != nil {
		return nil, err
	}

	cageResponses := make([]dto.CageResponse, 0)
	for _, cage := range cages {
		cageResponses = append(cageResponses, mapper.CageToResponse(&cage))
	}

	return cageResponses, nil
}

func (s *CageService) CreateCageFeed(request dto.CreateCageFeedRequest, userId uuid.UUID) (dto.CageFeedResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		return dto.CageFeedResponse{}, errx.BadRequest("invalid chicken category")
	}

	feedType := enum.ValueOfFeedType(request.FeedType)
	if !feedType.IsValid() {
		return dto.CageFeedResponse{}, errx.BadRequest("invalid feed type")
	}

	cageFeed := entity.CageFeed{
		ChickenCategory: chickenCategory,
		FeedType:        feedType,
		TotalFeed:       request.TotalFeed,
		CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
	}

	if err := s.repository.CreateCageFeed(&cageFeed); err != nil {
		s.log.Error("failed to save cage feed", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	cageFeedDetails := make([]entity.CageFeedDetail, 0, len(request.CageFeedDetails))
	for _, e := range request.CageFeedDetails {
		cageFeedDetails = append(cageFeedDetails, entity.CageFeedDetail{
			CageFeedId: cageFeed.Id,
			ItemId:     e.ItemId,
			Percentage: e.Percentage,
			CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if len(cageFeedDetails) > 0 {
		if err := s.repository.CreateCageFeedDetails(&cageFeedDetails); err != nil {
			s.log.Error("failed to save cage feed details", zap.Error(err))
			return dto.CageFeedResponse{}, err
		}
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	data, err := s.repository.GetCageFeed(cageFeed.Id)
	if err != nil {
		s.log.Error("failed get cage feed", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	response := mapper.CageFeedToResponse(&data)

	cageFeedDetailResponses := make([]dto.CageFeedDetailResponse, 0)
	for _, e := range data.CageFeedDetails {
		cageFeedDetailResponses = append(cageFeedDetailResponses, mapper.CageFeedDetailToResponse(&e))
	}

	response.CageFeedDetails = cageFeedDetailResponses

	return response, nil
}

func (s *CageService) UpdateCageFeed(id uint64, request dto.UpdateCageFeedRequest, userId uuid.UUID) (dto.CageFeedResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		return dto.CageFeedResponse{}, errx.BadRequest("invalid chicken category")
	}

	feedType := enum.ValueOfFeedType(request.FeedType)
	if !feedType.IsValid() {
		return dto.CageFeedResponse{}, errx.BadRequest("invalid feed type")
	}

	cageFeed, err := s.repository.GetCageFeed(id)
	if err != nil {
		s.log.Error("failed to get cage feed by id", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	cageFeed.ChickenCategory = chickenCategory
	cageFeed.FeedType = feedType
	cageFeed.TotalFeed = request.TotalFeed
	cageFeed.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if err := s.repository.UpdateCageFeed(&cageFeed); err != nil {
		s.log.Error("failed to update cage feed", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	newDetails := make([]entity.CageFeedDetail, 0, len(request.CageFeedDetails))
	newItemIds := make([]uint64, 0, len(request.CageFeedDetails))
	for _, d := range request.CageFeedDetails {
		newDetails = append(newDetails, entity.CageFeedDetail{
			Id:         d.Id,
			CageFeedId: cageFeed.Id,
			ItemId:     d.ItemId,
			Percentage: d.Percentage,
			CreatedBy:  uuid.NullUUID{UUID: userId, Valid: true},
		})
		newItemIds = append(newItemIds, d.ItemId)
	}

	if err := s.repository.DeleteCageFeedDetailsNotIn(cageFeed.Id, newItemIds); err != nil {
		s.log.Error("failed to delete old cage feed details", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	if len(newDetails) > 0 {
		if err := s.repository.UpsertCageFeedDetails(&newDetails); err != nil {
			s.log.Error("failed to upsert cage feed details", zap.Error(err))
			return dto.CageFeedResponse{}, err
		}
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	data, err := s.repository.GetCageFeed(cageFeed.Id)
	if err != nil {
		s.log.Error("failed get cage feed", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	response := mapper.CageFeedToResponse(&data)

	cageFeedDetailResponses := make([]dto.CageFeedDetailResponse, 0)
	for _, e := range data.CageFeedDetails {
		cageFeedDetailResponses = append(cageFeedDetailResponses, mapper.CageFeedDetailToResponse(&e))
	}

	response.CageFeedDetails = cageFeedDetailResponses

	return response, nil
}

func (s *CageService) GetCageFeeds() ([]dto.CageFeedResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetCageFeeds()
	if err != nil {
		s.log.Error("failed get cage feed", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.CageFeedResponse, 0)
	for _, cageFeed := range data {
		response := mapper.CageFeedToResponse(&cageFeed)
		cageFeedDetailResponses := make([]dto.CageFeedDetailResponse, 0)
		for _, e := range cageFeed.CageFeedDetails {
			cageFeedDetailResponses = append(cageFeedDetailResponses, mapper.CageFeedDetailToResponse(&e))
		}

		response.CageFeedDetails = cageFeedDetailResponses

		responses = append(responses, response)
	}

	return responses, nil
}

func (s *CageService) GetCageFeed(id uint64) (dto.CageFeedResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetCageFeed(id)
	if err != nil {
		s.log.Error("failed get cage feed", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	response := mapper.CageFeedToResponse(&data)

	cageFeedDetailResponses := make([]dto.CageFeedDetailResponse, 0)
	for _, e := range data.CageFeedDetails {
		cageFeedDetailResponses = append(cageFeedDetailResponses, mapper.CageFeedDetailToResponse(&e))
	}

	response.CageFeedDetails = cageFeedDetailResponses

	return response, nil
}

func (s *CageService) GetChickenCageFeeds(filter dto.GetChickenCageFeedFilter) ([]dto.ChickenCageFeedListResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetChickenCageFeeds(filter)
	if err != nil {
		return nil, err
	}

	chickenCages, err := s.repository.GetChickenCages(dto.GetChickenCageFilter{
		LocationId: filter.LocationId,
	})
	if err != nil {
		s.log.Error("failed get chicken cages", zap.Error(err))
		return nil, err
	}

	chickenCageIds := make([]uint64, len(chickenCages))
	for i, cage := range chickenCages {
		chickenCageIds[i] = cage.Id
	}

	cageFeedStocks, err := s.repository.GetCageFeedStocks(dto.GetCageFeedStockFilter{})
	if err != nil {
		s.log.Error("failed get cage feed stocks", zap.Error(err))
		return nil, err
	}

	cageFeedStockMap := make(map[uint64]float64, 0)
	for _, m := range cageFeedStocks {
		cageFeedStockMap[m.CageId] += m.TotalFeed - m.UsedFeed
	}

	cageFeeds, err := s.repository.GetCageFeeds()
	if err != nil {
		return nil, err
	}

	cageFeedsMapByCategory := make(map[string]entity.CageFeed)
	for _, e := range cageFeeds {
		cageFeedsMapByCategory[e.ChickenCategory.String()] = e
	}

	responses := make([]dto.ChickenCageFeedListResponse, 0)
	for _, e := range data {
		if !e.Cage.IsUsed {
			continue
		}

		resp := mapper.ChickenCageFeedToListResponse(&e)
		resp.TotalFeed = (cageFeedsMapByCategory[resp.ChickenCategory].TotalFeed * float64(e.TotalChicken) / 1000.0) - cageFeedStockMap[e.CageId]
		resp.FeedPerChicken = cageFeedsMapByCategory[resp.ChickenCategory].TotalFeed
		resp.FeedType = cageFeedsMapByCategory[resp.ChickenCategory].FeedType.String()
		responses = append(responses, resp)
	}

	return responses, nil
}

func (s *CageService) GetChickenCageFeed(chickenCageId uint64) (dto.ChickenCageFeedResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageById(chickenCageId)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	chickenCategory := util.GetChickenCategoryByChickenCage(&chickenCage)
	cageFeed, err := s.repository.GetCageFeedByChickenCategory(chickenCategory)
	if err != nil {
		return dto.ChickenCageFeedResponse{}, err
	}

	needCreateFeed := (cageFeed.TotalFeed * float64(chickenCage.TotalChicken)) / 1000.0
	cageFeedStocks, err := s.repository.GetCageFeedStocks(dto.GetCageFeedStockFilter{
		CageId: chickenCage.CageId,
	})
	if err != nil {
		s.log.Error("failed get cage feed stocks", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	remainingStockFeeds := float64(0)
	for _, e := range cageFeedStocks {
		needCreateFeed -= (e.TotalFeed - e.UsedFeed)
		remainingStockFeeds += (e.TotalFeed - e.UsedFeed)
	}

	feedDetailResponse := make([]dto.FeedDetailResponse, 0)
	for _, cageFeedDetail := range cageFeed.CageFeedDetails {
		feedDetailResponse = append(feedDetailResponse, dto.FeedDetailResponse{
			Item:       mapper.ItemToResponse(&cageFeedDetail.Item),
			Percentage: cageFeedDetail.Percentage,
			Quantity:   needCreateFeed * (cageFeedDetail.Percentage / 100.0),
		})
	}

	response := mapper.ChickenCageFeedToResponse(&chickenCage)
	response.FeedType = cageFeed.FeedType.String()
	response.RemainingTotalFeed = remainingStockFeeds
	response.FeedPerChicken = cageFeed.TotalFeed

	if needCreateFeed < 0 {
		response.TotalFeed = 0
	} else {
		response.TotalFeed = needCreateFeed
	}
	response.FeedDetails = feedDetailResponse

	return response, nil
}

func (s *CageService) ConfirmationChickenCageFeed(chickenCageId uint64, request dto.ConfirmationChickenCageFeedRequest, userId uuid.UUID) (dto.ChickenCageFeedResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageById(chickenCageId)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	if !chickenCage.IsNeedFeed {
		return dto.ChickenCageFeedResponse{}, errx.BadRequest("chicken cage is already feed")
	}

	chickenCategory := util.GetChickenCategoryByChickenCage(&chickenCage)
	cageFeed, err := s.repository.GetCageFeedByChickenCategory(chickenCategory)
	if err != nil {
		return dto.ChickenCageFeedResponse{}, err
	}

	needCreateFeed := (cageFeed.TotalFeed * float64(chickenCage.TotalChicken)) / 1000.0
	cageFeedStocks, err := s.repository.GetCageFeedStocks(dto.GetCageFeedStockFilter{
		CageId: chickenCage.CageId,
	})
	if err != nil {
		s.log.Error("failed get cage feed stocks", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	for _, e := range cageFeedStocks {
		needCreateFeed -= (e.TotalFeed - e.UsedFeed)
	}

	requestToWarehouse := make([]dto.ReduceFeedRequest, 0)
	for _, cageFeedDetail := range request.ChickenCageFeedDetails {
		requestToWarehouse = append(requestToWarehouse, dto.ReduceFeedRequest{
			ItemId:       cageFeedDetail.ItemId,
			ItemCategory: cageFeedDetail.Category,
			Quantity:     cageFeedDetail.Quantity,
		})
	}

	err = s.warehouseService.ReduceWarehouseItemForFeed(request.WarehouseId, requestToWarehouse, userId, chickenCage.Cage.Name)
	if err != nil {
		return dto.ChickenCageFeedResponse{}, err
	}

	chickenCage.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	chickenCage.IsNeedFeed = false
	err = s.repository.UpdateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed update chicken cage", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	err = s.repository.CreateCageFeedStock(&entity.CageFeedStock{
		CageId:    chickenCage.CageId,
		TotalFeed: needCreateFeed,
		UsedFeed:  0,
		CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
	})
	if err != nil {
		s.log.Error("failed create cage feed history", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	return s.GetChickenCageFeed(chickenCageId)
}

func (s *CageService) MoveChickenCage(request dto.MoveChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	destinationCageIds := make([]uint64, len(request.DestinationChickenCages))
	for i, e := range request.DestinationChickenCages {
		destinationCageIds[i] = e.DestinationCageId
	}

	destinationChickenCages, err := s.repository.GetChickenCagesByCageIds(destinationCageIds)
	if err != nil {
		s.log.Error("failed to get chicken cages by ids", zap.Error(err))
		return nil, err
	}

	sourceChickenCage, err := s.repository.GetChickenCageByCageId(request.SourceCageId)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return nil, err
	}

	sourceCage, err := s.repository.GetCageById(request.SourceCageId)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return nil, err
	}

	destinationCages, err := s.repository.GetCagesByIds(destinationCageIds)
	if err != nil {
		s.log.Error("failed get cages by ids", zap.Error(err))
		return nil, err
	}

	destinationCagesMapById := make(map[uint64]entity.Cage)
	for _, dest := range destinationCages {
		destinationCagesMapById[dest.Id] = dest
	}

	destinationChickenCagesMapByCageId := make(map[uint64]entity.ChickenCage)
	for _, dest := range destinationChickenCages {
		destinationChickenCagesMapByCageId[dest.CageId] = dest
	}

	for _, chickenCage := range destinationChickenCages {
		if chickenCage.Cage.IsUsed && chickenCage.ChickenProcurementId != sourceChickenCage.ChickenProcurementId {
			return nil, errx.BadRequest(fmt.Sprintf("cage with id %d is used and not suitable batch id", chickenCage.Cage.Id))
		}
	}

	// Note : n+1 query
	for _, destinationCage := range destinationCages {
		destinationCage.IsUsed = true

		err = s.repository.UpdateCage(&destinationCage)
		if err != nil {
			s.log.Error("failed update cage", zap.Error(err))
			return nil, err
		}
	}

	newChickenCages := make([]entity.ChickenCage, 0)
	updateChickenCages := make([]entity.ChickenCage, 0)
	totalMoveChicken := uint64(0)

	for _, request := range request.DestinationChickenCages {
		if destinationChickenCagesMapByCageId[request.DestinationCageId].ChickenProcurementId == sourceChickenCage.ChickenProcurementId {
			updateChickenCage := destinationChickenCagesMapByCageId[request.DestinationCageId]
			updateChickenCage.TotalChicken += request.TotalChicken
			updateChickenCages = append(updateChickenCages, updateChickenCage)

			if updateChickenCage.TotalChicken > destinationCagesMapById[request.DestinationCageId].Capacity {
				return nil, errx.BadRequest(fmt.Sprintf("total chicken for cage id %d is more than capacity", request.DestinationCageId))
			}
		} else {
			newChickenCage := entity.ChickenCage{
				CageId:               request.DestinationCageId,
				TotalChicken:         request.TotalChicken,
				ChickenProcurementId: sourceChickenCage.ChickenProcurementId,
			}
			newChickenCages = append(newChickenCages, newChickenCage)

			if newChickenCage.TotalChicken > destinationCagesMapById[request.DestinationCageId].Capacity {
				return nil, errx.BadRequest(fmt.Sprintf("total chicken for cage id %d is more than capacity", request.DestinationCageId))
			}
		}

		totalMoveChicken += request.TotalChicken
	}

	if sourceChickenCage.TotalChicken < totalMoveChicken {
		return nil, errx.BadRequest("total move chicken more than total chicken in current cage")
	}

	sourceChickenCage.TotalChicken -= totalMoveChicken
	if sourceChickenCage.TotalChicken == 0 {
		sourceCage.IsUsed = false

		if err := s.repository.UpdateCage(&sourceCage); err != nil {
			s.log.Error("failed update cage", zap.Error(err))
			return nil, err
		}

		if err := s.repository.CreateChickenCage(&entity.ChickenCage{
			CageId:       sourceCage.Id,
			TotalChicken: 0,
			CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
		}); err != nil {
			return nil, err
		}
	}

	if err := s.repository.UpdateChickenCage(&sourceChickenCage); err != nil {
		s.log.Error("failed update chicken cage", zap.Error(err))
		return nil, err
	}

	if len(newChickenCages) > 0 {
		if err := s.repository.CreateChickenCageInBatch(&newChickenCages); err != nil {
			s.log.Error("failed create chicken cage in batch", zap.Error(err))
			return nil, err
		}
	}

	if len(updateChickenCages) > 0 {
		// Note : n+1 query problem
		for _, updateChickenCage := range updateChickenCages {
			if err := s.repository.UpdateChickenCage(&updateChickenCage); err != nil {
				s.log.Error("failed to update chicken cage", zap.Error(err))
				return nil, err
			}
		}
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return nil, err
	}

	chickenCageIds := make([]uint64, len(newChickenCages))
	for i, cc := range newChickenCages {
		chickenCageIds[i] = cc.Id
	}

	for i, cc := range updateChickenCages {
		chickenCageIds[i] = cc.Id
	}

	createdChickenCages, err := s.repository.GetChickenCageByIds(chickenCageIds)
	if err != nil {
		return nil, err
	}

	responses := make([]dto.ChickenCageResponse, len(createdChickenCages))
	for i, chickenCage := range createdChickenCages {
		responses[i] = mapper.ChickenCageToResponse(&chickenCage)
	}

	return responses, nil
}

func (s *CageService) GetTotalCageFeedHistory(cageId uint64) (float64, error) {
	s.repository.UseTx(false)

	totalCurrentFeed := float64(0)
	cageFeedStocks, err := s.repository.GetCageFeedStocks(dto.GetCageFeedStockFilter{
		CageId: cageId,
	})
	if err != nil {
		s.log.Error("failed get cage feed stocks", zap.Error(err))
		return -1, err
	}

	for _, e := range cageFeedStocks {
		totalCurrentFeed = e.TotalFeed - e.UsedFeed
	}

	return totalCurrentFeed, nil
}

func (s *CageService) ReduceCageFeedStocks(latestFeed float64, currentFeeds float64, cageId uint64) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	diff := currentFeeds - latestFeed

	if diff == 0 {
		return nil
	}

	cageFeedStocks, err := s.repository.GetCageFeedStocks(dto.GetCageFeedStockFilter{
		CageId: cageId,
	})
	if err != nil {
		s.log.Error("failed to get cage feed stocks", zap.Error(err))
		return err
	}

	for i := range cageFeedStocks {
		e := &cageFeedStocks[i]

		if diff == 0 {
			break
		}

		if diff > 0 {
			available := e.TotalFeed - e.UsedFeed
			if available <= 0 {
				continue
			}

			use := math.Min(available, diff)
			e.UsedFeed += use
			diff -= use

		} else if diff < 0 {
			reduce := math.Min(e.UsedFeed, -diff)
			e.UsedFeed -= reduce
			diff += reduce
		}

		if err := s.repository.UpdateCageFeedStock(e); err != nil {
			s.log.Error("failed to update cage feed stock", zap.Error(err))
			return err
		}
	}

	if diff != 0 {
		s.log.Error("not enough feed stock capacity to adjust completely",
			zap.Float64("remaining_diff", diff),
			zap.Uint64("cageId", cageId))
		return errx.BadRequest(fmt.Sprintf("insufficient feed stock to adjust by %.2f", diff))
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}
