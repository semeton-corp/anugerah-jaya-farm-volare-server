package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
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
	UpdateCage(id uint64, request dto.UpdateCageRequest, updatedBy uuid.UUID) (dto.CageResponse, error)
	DeleteCage(id uint64) error
	GetCageById(id uint64) (dto.CageResponse, error)
	GetCagesByIds(ids []uint64) ([]dto.CageResponse, error)

	CreateChickenCage(request dto.CreateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error)
	GetChickenCages(filter dto.GetChickenCageFilter) ([]dto.ChickenCageResponse, error)
	GetChickenCageById(id uint64) (dto.ChickenCageResponse, error)
	GetChickenCagesByCageIds(cageIds []uint64) ([]dto.ChickenCageResponse, error)
	CreateChickenCageInBatch(request []dto.CreateChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error)
	UpdateChickenCage(id uint64, request dto.UpdateChickenCageRequest, updatedBy uuid.UUID) (dto.ChickenCageResponse, error)

	CreateCageFeed(request dto.CreateCageFeedRequest, userId uuid.UUID) (dto.CageFeedResponse, error)
	UpdateCageFeed(id uint64, request dto.UpdateCageFeedRequest, userId uuid.UUID) (dto.CageFeedResponse, error)
	GetCageFeeds() ([]dto.CageFeedResponse, error)
	GetCageFeed(id uint64) (dto.CageFeedResponse, error)
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

func (s *CageService) CreateCage(request dto.CreateCageRequest, createdBy uuid.UUID) (dto.CageResponse, error) {
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
		CreatedBy:       uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateCage(&cage)
	if err != nil {
		s.log.Error("failed to create cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	chickenCage := entity.ChickenCage{
		CageId:    cage.Id,
		CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
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

func (s *CageService) UpdateCage(id uint64, request dto.UpdateCageRequest, updatedBy uuid.UUID) (dto.CageResponse, error) {
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
	cage.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

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
	s.repository.UseTx(false)

	err := s.repository.DeleteCage(id)
	if err != nil {
		s.log.Error("failed to delete cage", zap.Error(err))
		return err
	}

	return nil
}

func (s *CageService) GetChickenCages(filter dto.GetChickenCageFilter) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCageResponses := make([]dto.ChickenCageResponse, 0)
	chickenCages, err := s.repository.GetChickenCages(filter)
	if err != nil {
		return chickenCageResponses, err
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

func (s *CageService) UpdateChickenCage(id uint64, request dto.UpdateChickenCageRequest, updatedBy uuid.UUID) (dto.ChickenCageResponse, error) {
	return dto.ChickenCageResponse{}, nil
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
		CageId:               request.CageId,
		ChickenProcurementId: request.ChickenProcurementId,
		TotalChicken:         request.TotalChicken,
		CreatedBy:            uuid.NullUUID{UUID: userId, Valid: true},
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

func (s *CageService) GetChickenCagesByCageIds(ids []uint64) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCages, err := s.repository.GetChickenCagesByCageIds(ids)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ChickenCageResponse, 0)
	for _, chickenCage := range chickenCages {
		response = append(response, mapper.ChickenCageToResponse(&chickenCage))
	}

	return response, nil
}

func (s *CageService) CreateChickenCageInBatch(request []dto.CreateChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCages := make([]entity.ChickenCage, 0)
	for _, req := range request {
		chickenCages = append(chickenCages, entity.ChickenCage{
			CageId:               req.CageId,
			ChickenProcurementId: req.ChickenProcurementId,
			TotalChicken:         req.TotalChicken,
		})
	}

	err := s.repository.CreateChickenCageInBatch(&chickenCages)
	if err != nil {
		s.log.Error("failed create chicken cage in batch", zap.Error(err))
		return nil, err
	}

	chickenCageIds := make([]uint64, 0)
	for _, chichickenCage := range chickenCages {
		chickenCageIds = append(chickenCageIds, chichickenCage.Id)
	}

	chickenCages, err = s.repository.GetChickenCageByIds(chickenCageIds)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ChickenCageResponse, 0)
	for _, chickenCage := range chickenCages {
		response = append(response, mapper.ChickenCageToResponse(&chickenCage))
	}

	return response, nil
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

	// Get existing cage feed
	cageFeed, err := s.repository.GetCageFeed(id)
	if err != nil {
		s.log.Error("failed to get cage feed by id", zap.Error(err))
		return dto.CageFeedResponse{}, err
	}

	// Update main cage feed fields
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

	now := time.Now()
	yesterday := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)

	chickenMonitoringsYesterday, err := s.repository.GetChickenMonitoringYesterdayByChickenCageIds(chickenCageIds, yesterday)
	if err != nil {
		s.log.Error("failed get chicken monitoring yesterday by chicken cage ids", zap.Error(err))
		return nil, err
	}

	monitoringMap := make(map[uint64]float64, len(chickenMonitoringsYesterday))
	for _, m := range chickenMonitoringsYesterday {
		monitoringMap[m.ChickenCageId] = m.TotalFeed
	}

	responses := make([]dto.ChickenCageFeedListResponse, len(data))
	for i, e := range data {
		resp := mapper.ChickenCageFeedToListResponse(&e)
		if yesterdayFeed, ok := monitoringMap[resp.Id]; ok {
			resp.YesterdayTotalFeed = yesterdayFeed
			resp.TotalFeed = resp.ExpectedTotalFeed - yesterdayFeed
		}
		responses[i] = resp
	}

	return responses, nil
}

func (s *CageService) GetChickenCageFeed(chickenCageId uint64) (dto.ChickenCageFeedResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageFeedById(chickenCageId)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	yesterdayDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	chickenMonitoring, err := s.repository.GetChickenMonitoringYesterday(chickenCageId, yesterdayDate)
	if err != nil {
		s.log.Error("failed get chicken monitoring yesterday", zap.Error(err))
		return dto.ChickenCageFeedResponse{}, err
	}

	needCreateFeed := chickenMonitoring.TotalFeed
	feedDetailResponse := make([]dto.FeedDetailResponse, 0)
	for _, cageFeedDetail := range chickenCage.Cage.CageFeed.CageFeedDetails {
		feedDetailResponse = append(feedDetailResponse, dto.FeedDetailResponse{
			Item:       mapper.ItemToResponse(&cageFeedDetail.Item),
			Percentage: cageFeedDetail.Percentage,
			Quantity:   needCreateFeed * (cageFeedDetail.Percentage / 100.0),
		})
	}

	response := mapper.ChickenCageFeedToResponse(&chickenCage)
	response.YesterdayTotalFeed = chickenMonitoring.TotalFeed
	response.TotalFeed = response.ExpectedTotalFeed - response.YesterdayTotalFeed
	response.FeedDetails = feedDetailResponse

	return response, nil
}

func (s *CageService) ConfirmationChickenCageFeed(chickenCageId uint64, request dto.ConfirmationChickenCageFeedRequest) error {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageById(chickenCageId)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return err
	}

	yesterdayDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local).AddDate(0, 0, -1)
	chickenMonitoring, err := s.repository.GetChickenMonitoringYesterday(chickenCageId, yesterdayDate)
	if err != nil {
		s.log.Error("failed get chicken monitoring yesterday", zap.Error(err))
		return err
	}

	chickenCage.IsNeedFeed = false

	needCreateFeed := chickenMonitoring.TotalFeed
	requestToWarehouse := make([]dto.ReduceFeedRequest, 0)
	for _, cageFeedDetail := range chickenCage.Cage.CageFeed.CageFeedDetails {
		requestToWarehouse = append(requestToWarehouse, dto.ReduceFeedRequest{
			ItemId:       cageFeedDetail.ItemId,
			ItemCategory: cageFeedDetail.Item.Category.String(),
			Quantity:     needCreateFeed * (cageFeedDetail.Percentage / 100.0),
		})
	}

	err = s.warehouseService.ReduceWarehouseItemForFeed(request.WarehouseId, requestToWarehouse)
	if err != nil {
		return err
	}

	err = s.repository.UpdateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed update chicken cage", zap.Error(err))
		return err
	}

	return nil
}
