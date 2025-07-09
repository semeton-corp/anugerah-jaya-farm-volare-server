package service

import (
	"database/sql"
	"math"
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
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type StoreService struct {
	log              *zap.Logger
	repository       repository.IStoreRepository
	cacheService     cache.ICache
	placementService IPlacementService
	warehouseService IWarehouseService
	userService      IUserService
}

type IStoreService interface {
	CreateStore(request dto.CreateStoreRequest, createdBy uuid.UUID) (dto.StoreResponse, error)
	UpdateStore(id uint64, request dto.UpdateStoreRequest, updatedBy uuid.UUID) (dto.StoreResponse, error)
	DeleteStore(id uint64) error
	GetStores(filter dto.GetStoreFilter) ([]dto.StoreResponse, error)
	GetStoreDetailById(id uint64) (dto.StoreDetailResponse, error)

	CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, createdBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	CreateStoreRequestItemFromEggMonitoring(request dto.CreateStoreRequestItemRequest, createdBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) (dto.StoreRequestItemListPaginationResponse, error)
	WarehouseConfirmationStoreRequestItem(id uint64, request dto.WarehouseConfirmationStoreRequestItem, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	StoreConfirmationStoreRequestItem(id uint64, request dto.StoreConfirmationStoreRequestItem, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	SortingStoreRequestItem(id uint64, request dto.SortingStoreRequestItemRequest, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)

	GetStoreItems(filter dto.GetStoreItemFilter) ([]dto.StoreItemResponse, error)
	GetStoreOverview(id uint64) (dto.StoreItemOverview, error)
	GetStoreItemByStoreIdAndItemId(storeId uint64, itemId uint64) (dto.StoreItemResponse, error)
	UpdateStoreItem(storeId uint64, itemId uint64, request dto.UpdateStoreItemRequest, updatedBy uuid.UUID) (dto.StoreItemResponse, error)
	GetEggStoreItemSummary(storeId uint64) ([]dto.EggStoreItemSummary, error)

	GetStoreItemHistories(filter dto.GetStoreItemHistoryFilter) (dto.StoreItemHistoryListPaginationResponse, error)
	GetStoreItemHistoryById(id uint64) (dto.StoreItemHistoryResponse, error)

	CreateStoreSale(request dto.CreateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)
	GetStoreSaleById(id uint64) (dto.StoreSaleResponse, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) (dto.StoreSaleListPaginationResponse, error)
	UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)

	CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)
	UpdateStoreSalePayment(id uint64, request dto.UpdateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)

	SendStoreSale(id uint64, accountId uuid.UUID) (dto.StoreSaleResponse, error)
}

func NewStoreService(log *zap.Logger, repository repository.IStoreRepository, cacheService cache.ICache, placementService IPlacementService, warehouseService IWarehouseService, userService IUserService) IStoreService {
	return &StoreService{
		log:              log,
		repository:       repository,
		cacheService:     cacheService,
		placementService: placementService,
		warehouseService: warehouseService,
		userService:      userService,
	}
}

// Todo : When created store auto create 3 egg object
func (s *StoreService) CreateStore(request dto.CreateStoreRequest, createdBy uuid.UUID) (dto.StoreResponse, error) {
	s.repository.UseTx(false)

	store := entity.Store{
		Name:       request.Name,
		LocationId: request.LocationId,
		CreatedBy:  uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateStore(&store)
	if err != nil {
		s.log.Error("failed to create store", zap.Error(err))
		return dto.StoreResponse{}, err
	}

	store, err = s.repository.GetStoreById(store.Id)
	if err != nil {
		s.log.Error("failed to get store by id", zap.Error(err))
		return dto.StoreResponse{}, err
	}

	return mapper.StoreToResponse(&store), nil
}

func (s *StoreService) UpdateStore(id uint64, request dto.UpdateStoreRequest, updatedBy uuid.UUID) (dto.StoreResponse, error) {
	s.repository.UseTx(false)

	store, err := s.repository.GetStoreById(id)
	if err != nil {
		s.log.Error("failed to get store by id", zap.Error(err))
		return dto.StoreResponse{}, err
	}

	store.Name = request.Name
	store.LocationId = request.LocationId
	store.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateStore(&store)
	if err != nil {
		s.log.Error("failed to update store", zap.Error(err))
		return dto.StoreResponse{}, err
	}

	store, err = s.repository.GetStoreById(id)
	if err != nil {
		s.log.Error("failed to get store by id", zap.Error(err))
		return dto.StoreResponse{}, err
	}

	return mapper.StoreToResponse(&store), nil
}

func (s *StoreService) DeleteStore(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteStore(id)
	if err != nil {
		s.log.Error("failed to delete store", zap.Error(err))
		return err
	}

	return nil
}

func (s *StoreService) GetStores(filter dto.GetStoreFilter) ([]dto.StoreResponse, error) {
	s.repository.UseTx(false)

	stores, err := s.repository.GetStores(filter)
	if err != nil {
		s.log.Error("failed to get stores", zap.Error(err))
		return nil, err
	}

	storeResponses := make([]dto.StoreResponse, len(stores))
	for i, store := range stores {
		storeResponses[i] = mapper.StoreToResponse(&store)
	}

	return storeResponses, nil
}

func (s *StoreService) GetStoreDetailById(id uint64) (dto.StoreDetailResponse, error) {
	s.repository.UseTx(false)

	store, err := s.repository.GetStoreById(id)
	if err != nil {
		s.log.Error("failed to get store by id", zap.Error(err))
		return dto.StoreDetailResponse{}, err
	}

	storePlacements, err := s.placementService.GetStorePlacementByStoreId(id)
	if err != nil {
		return dto.StoreDetailResponse{}, err
	}

	userResponses := make([]dto.UserResponse, 0)
	for _, e := range storePlacements {
		userResponses = append(userResponses, e.User)
	}

	return dto.StoreDetailResponse{
		Id:       store.Id,
		Name:     store.Name,
		Location: mapper.LocationToResponse(&store.Location),
		Users:    userResponses,
	}, nil
}

// Todo : pubsub redis store item history
func (s *StoreService) CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, createdBy uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(false)

	storePlacement, err := s.placementService.GetStorePlacementByUserId(createdBy)
	if err != nil {
		return dto.StoreRequestItemResponse{}, err
	}

	warehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(request.WarehouseId, request.ItemId)
	if err != nil {
		return dto.StoreRequestItemResponse{}, err
	}

	// Note & Todo : the stock must be in Kg && fix
	if warehouseItem.Quantity < request.Quantity*float64(constant.TotalEggPerIkat) && warehouseItem.Item.Unit == "Kg" {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("insuficcient stock for request item")
	}

	storeRequestItem := entity.StoreRequestItem{
		WarehouseId: request.WarehouseId,
		ItemId:      request.ItemId,
		StoreId:     sql.NullInt64{Int64: int64(storePlacement.Store.Id), Valid: true},
		Quantity:    request.Quantity,
		Status:      enum.RequestItemStatusPending,
		CreatedBy:   uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err = s.repository.CreateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("failed to create store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return s.GetStoreRequestItemById(storeRequestItem.Id)
}

func (s *StoreService) CreateStoreRequestItemFromEggMonitoring(request dto.CreateStoreRequestItemRequest, createdBy uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(false)

	storeRequestItem := entity.StoreRequestItem{
		WarehouseId: request.WarehouseId,
		ItemId:      request.ItemId,
		Quantity:    request.Quantity,
		Status:      enum.RequestItemStatusPending,
		CreatedBy:   uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("failed to create store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return s.GetStoreRequestItemById(storeRequestItem.Id)
}

func (s *StoreService) GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error) {
	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	user, err := s.userService.GetUserById(storeRequestItem.CreatedBy.UUID)
	if err != nil {
		return dto.StoreRequestItemResponse{}, err
	}

	response := mapper.StoreRequestItemToResponse(&storeRequestItem)
	response.CreatedBy = user.Name

	return response, nil
}

func (s *StoreService) GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) (dto.StoreRequestItemListPaginationResponse, error) {
	s.repository.UseTx(false)

	storeRequestItems, err := s.repository.GetStoreRequestItems(filter)
	if err != nil {
		s.log.Error("failed to get store request items", zap.Error(err))
		return dto.StoreRequestItemListPaginationResponse{}, err
	}

	storeRequestItemResponses := make([]dto.StoreRequestItemResponse, len(storeRequestItems))
	for i, storeRequestItem := range storeRequestItems {
		storeRequestItemResponses[i] = mapper.StoreRequestItemToResponse(&storeRequestItem)
	}

	totalData, err := s.repository.CountTotalStoreRequestItem(dto.GetStoreRequestItemFilter{
		Date: filter.Date,
	})
	if err != nil {
		s.log.Error("failed to count request items", zap.Error(err))
		return dto.StoreRequestItemListPaginationResponse{}, err
	}

	resp := dto.StoreRequestItemListPaginationResponse{
		TotalPage:         uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData:         totalData,
		StoreRequestItems: storeRequestItemResponses,
	}

	return resp, nil
}

func (s *StoreService) WarehouseConfirmationStoreRequestItem(id uint64, request dto.WarehouseConfirmationStoreRequestItem, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	if storeRequestItem.Status == enum.RequestItemStatusCanceled || storeRequestItem.Status == enum.RequestItemStatusRejected || storeRequestItem.Status == enum.RequestItemStatusArrivedNotOk || storeRequestItem.Status == enum.RequestItemStatusArrivedOk {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is in another status")
	} else if storeRequestItem.Status != enum.RequestItemStatusPending {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item not pending")
	}

	if storeRequestItem.Quantity > request.Quantity {
		remainingQuantity := storeRequestItem.Quantity - request.Quantity

		storeRequestItem.WarehouseFulfillment = request.Quantity
		storeRequestItem.WarehouseNote = request.WarehouseNote
		storeRequestItem.Status = enum.RequestItemStatusSentOff
		storeRequestItem.StoreId = sql.NullInt64{Int64: int64(request.StoreId), Valid: true}
		storeRequestItem.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

		newStoreRequestItem := entity.StoreRequestItem{
			WarehouseId: storeRequestItem.WarehouseId,
			ItemId:      storeRequestItem.ItemId,
			StoreId:     storeRequestItem.StoreId,
			Quantity:    remainingQuantity,
			Status:      enum.RequestItemStatusPending,
			CreatedBy:   uuid.NullUUID{UUID: updatedBy, Valid: true},
		}

		err = s.repository.CreateStoreRequestItem(&newStoreRequestItem)
		if err != nil {
			s.log.Error("failed to create store request item", zap.Error(err))
			return dto.StoreRequestItemResponse{}, err
		}
	} else {
		storeRequestItem.WarehouseFulfillment = request.Quantity
		storeRequestItem.Status = enum.RequestItemStatusSentOff
		storeRequestItem.StoreId = sql.NullInt64{Int64: int64(request.StoreId), Valid: true}
		storeRequestItem.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}
	}

	err = s.repository.UpdateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("failed to update store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return s.GetStoreRequestItemById(id)
}

func (s *StoreService) StoreConfirmationStoreRequestItem(id uint64, request dto.StoreConfirmationStoreRequestItem, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	if storeRequestItem.Status == enum.RequestItemStatusCanceled || storeRequestItem.Status == enum.RequestItemStatusRejected || storeRequestItem.Status == enum.RequestItemStatusArrivedNotOk || storeRequestItem.Status == enum.RequestItemStatusArrivedOk {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is in another status")
	} else if storeRequestItem.Status != enum.RequestItemStatusSentOff {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item not sent off")
	}

	storeRequestItem.RecieveQuantity = request.Quantity
	storeRequestItem.StoreNote = request.StoreNote
	storeRequestItem.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}
	storeRequestItem.RecieveDate = sql.NullTime{Time: time.Now(), Valid: true}

	if storeRequestItem.RecieveQuantity != storeRequestItem.Quantity {
		storeRequestItem.Status = enum.RequestItemStatusArrivedNotOk
	} else {
		storeRequestItem.Status = enum.RequestItemStatusArrivedOk
	}

	storeItem, err := s.repository.GetStoreItemByStoreIdAndItemId(uint64(storeRequestItem.StoreId.Int64), storeRequestItem.ItemId)
	if err != nil {
		s.log.Error("failed to get store item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	storeItem.Quantity += request.Quantity

	err = s.repository.UpdateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("failed to update store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	err = s.repository.UpdateStoreItem(&storeItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return s.GetStoreRequestItemById(id)
}

func (s *StoreService) UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(false)

	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	if storeRequestItem.Status == enum.RequestItemStatusCanceled || storeRequestItem.Status == enum.RequestItemStatusRejected || storeRequestItem.Status == enum.RequestItemStatusArrivedNotOk || storeRequestItem.Status == enum.RequestItemStatusArrivedOk {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is in another status")
	}

	status := enum.ValueOfRequestItemStatus(request.Status)
	if !status.IsValid() {
		s.log.Warn("invalid status request item", zap.String("status", request.Status))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("invalid request item status")
	}

	storeRequestItem.Status = status
	storeRequestItem.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("failed to update store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return s.GetStoreRequestItemById(id)
}

func (s *StoreService) SortingStoreRequestItem(id uint64, request dto.SortingStoreRequestItemRequest, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	storeItems, err := s.repository.GetStoreItems(dto.GetStoreItemFilter{
		StoreId:   uint64(storeRequestItem.StoreId.Int64),
		Category:  param.ItemCategoryParam(enum.ItemCategoryEgg),
		ItemNames: []string{constant.CrackedEgg, constant.BrokenEgg},
		Units:     []string{constant.EggUnitPlastik, constant.EggUnitKg},
	})

	if err != nil {
		s.log.Error("failed to get store items", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	var crackedEgg, brokenEgg entity.StoreItem
	for _, storeItem := range storeItems {
		switch storeItem.Item.Name {
		case constant.CrackedEgg:
			crackedEgg = storeItem
		case constant.BrokenEgg:
			brokenEgg = storeItem
		}
	}

	if crackedEgg.Item.Id == 0 || brokenEgg.Item.Id == 0 {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("cracked egg or broken egg not found")
	}

	crackedEgg.Quantity -= request.BrokenEggInKg
	crackedEgg.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}
	brokenEgg.Quantity += math.Ceil(float64(request.BrokenEggInButir) / 4)
	brokenEgg.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	storeRequestItem.IsSorted = true
	storeRequestItem.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("failed to update store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	err = s.repository.UpdateStoreItem(&crackedEgg)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	err = s.repository.UpdateStoreItem(&brokenEgg)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return mapper.StoreRequestItemToResponse(&storeRequestItem), nil
}

func (s *StoreService) GetStoreItems(filter dto.GetStoreItemFilter) ([]dto.StoreItemResponse, error) {
	s.repository.UseTx(false)

	storeItems, err := s.repository.GetStoreItems(filter)
	if err != nil {
		s.log.Error("failed to get store items", zap.Error(err))
		return nil, err
	}

	storeItemResponses := make([]dto.StoreItemResponse, len(storeItems))
	for i, storeItem := range storeItems {
		storeItemResponses[i] = mapper.StoreItemToResponse(&storeItem)
	}

	return storeItemResponses, nil
}

func (s *StoreService) GetStoreOverview(id uint64) (dto.StoreItemOverview, error) {
	s.repository.UseTx(false)

	storeItems, err := s.repository.GetStoreItems(dto.GetStoreItemFilter{
		StoreId: id,
	})
	if err != nil {
		s.log.Error("failed to get store items", zap.Error(err))
		return dto.StoreItemOverview{}, err
	}

	eggStoreItemSummaries := make([]dto.EggStoreItemSummary, 0)
	for _, warehouseItem := range storeItems {
		switch warehouseItem.Item.Name {
		case constant.GoodEgg:
			eggStoreItemSummaries = append(eggStoreItemSummaries, dto.EggStoreItemSummary{
				Name:     constant.GoodEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.EggUnitKg,
			})

			eggStoreItemSummaries = append(eggStoreItemSummaries, dto.EggStoreItemSummary{
				Name:     constant.GoodEgg,
				Quantity: warehouseItem.Quantity / float64(constant.TotalEggPerIkat),
				Unit:     constant.EggUnitIkat,
			})
		case constant.CrackedEgg:
			eggStoreItemSummaries = append(eggStoreItemSummaries, dto.EggStoreItemSummary{
				Name:     constant.CrackedEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.EggUnitKg,
			})

			eggStoreItemSummaries = append(eggStoreItemSummaries, dto.EggStoreItemSummary{
				Name:     constant.CrackedEgg,
				Quantity: warehouseItem.Quantity / float64(constant.TotalEggPerIkat),
				Unit:     constant.EggUnitIkat,
			})
		}
	}

	storeItemResponses := make([]dto.StoreItemResponse, len(storeItems))
	for i, storeItem := range storeItems {
		storeItemResponses[i] = mapper.StoreItemToResponse(&storeItem)
	}

	return dto.StoreItemOverview{
		StoreItems:            storeItemResponses,
		EggStoreItemSummaries: eggStoreItemSummaries,
	}, nil
}

func (s *StoreService) GetStoreItemByStoreIdAndItemId(storeId uint64, itemId uint64) (dto.StoreItemResponse, error) {
	s.repository.UseTx(false)

	storeItem, err := s.repository.GetStoreItemByStoreIdAndItemId(storeId, itemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.StoreItemResponse{}, err
	}

	return mapper.StoreItemToResponse(&storeItem), nil
}

func (s *StoreService) UpdateStoreItem(storeId uint64, itemId uint64, request dto.UpdateStoreItemRequest, updatedBy uuid.UUID) (dto.StoreItemResponse, error) {
	s.repository.UseTx(false)

	storeItem, err := s.repository.GetStoreItemByStoreIdAndItemId(storeId, itemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.StoreItemResponse{}, err
	}

	storeItem.Quantity = request.Quantity
	storeItem.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateStoreItem(&storeItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.StoreItemResponse{}, err
	}

	storeItem, err = s.repository.GetStoreItemByStoreIdAndItemId(storeId, itemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.StoreItemResponse{}, err
	}

	return mapper.StoreItemToResponse(&storeItem), nil
}

func (s *StoreService) GetEggStoreItemSummary(storeId uint64) ([]dto.EggStoreItemSummary, error) {
	s.repository.UseTx(false)

	response := make([]dto.EggStoreItemSummary, 0)
	storeItems, err := s.repository.GetStoreItems(dto.GetStoreItemFilter{
		StoreId:   storeId,
		ItemNames: []string{constant.GoodEgg, constant.CrackedEgg, constant.BrokenEgg},
		Units:     []string{constant.EggUnitKg, constant.EggUnitPlastik},
	})
	if err != nil {
		s.log.Error("failed to get store items", zap.Error(err))
		return nil, err
	}

	for _, storeItem := range storeItems {
		switch storeItem.Item.Name {
		case constant.GoodEgg:
			response = append(response, dto.EggStoreItemSummary{
				Name:     constant.GoodEgg,
				Quantity: storeItem.Quantity,
				Unit:     constant.EggUnitKg,
			})

			response = append(response, dto.EggStoreItemSummary{
				Name:     constant.GoodEgg,
				Quantity: storeItem.Quantity / float64(constant.TotalEggPerIkat),
				Unit:     constant.EggUnitIkat,
			})
		case constant.CrackedEgg:
			response = append(response, dto.EggStoreItemSummary{
				Name:     constant.CrackedEgg,
				Quantity: storeItem.Quantity,
				Unit:     constant.EggUnitKg,
			})

			response = append(response, dto.EggStoreItemSummary{
				Name:     constant.CrackedEgg,
				Quantity: storeItem.Quantity / float64(constant.TotalEggPerIkat),
				Unit:     constant.EggUnitIkat,
			})
		case constant.BrokenEgg:
			response = append(response, dto.EggStoreItemSummary{
				Name:     constant.BrokenEgg,
				Quantity: storeItem.Quantity,
				Unit:     constant.EggUnitPlastik,
			})
		}
	}

	return response, nil
}

func (s *StoreService) GetStoreItemHistories(filter dto.GetStoreItemHistoryFilter) (dto.StoreItemHistoryListPaginationResponse, error) {
	s.repository.UseTx(false)

	storeItemHistories, err := s.repository.GetStoreItemHistories(filter)
	if err != nil {
		s.log.Error("failed to get Store item history", zap.Error(err))
		return dto.StoreItemHistoryListPaginationResponse{}, err
	}

	response := make([]dto.StoreItemHistoryListResponse, 0)
	for _, e := range storeItemHistories {
		response = append(response, mapper.StoreItemHistoryToListResponse(&e))
	}

	totalData, err := s.repository.CountTotalStoreItemHistory(filter)
	if err != nil {
		s.log.Error("failed to count Store item history", zap.Error(err))
		return dto.StoreItemHistoryListPaginationResponse{}, err
	}

	return dto.StoreItemHistoryListPaginationResponse{
		TotalPage:          uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData:          uint64(totalData),
		StoreItemHistories: response,
	}, nil
}

func (s *StoreService) GetStoreItemHistoryById(id uint64) (dto.StoreItemHistoryResponse, error) {
	s.repository.UseTx(false)

	storeItemHistory, err := s.repository.GetStoreItemHistoryById(id)
	if err != nil {
		s.log.Error("failed to get Store item history by id", zap.Error(err))
		return dto.StoreItemHistoryResponse{}, err
	}

	return mapper.StoreItemHistoryToResponse(&storeItemHistory), nil
}

func (s *StoreService) CreateStoreSale(request dto.CreateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	// Todo : reduce stock in warehouse item id in store id

	s.repository.UseTx(true)
	defer s.repository.Rollback()

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse sent date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sent date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("[CreateStoreSale] invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid price format")
	}

	totalPrice := price.Mul(decimal.NewFromInt(int64(request.Quantity)))

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		s.log.Error("[CreateStoreSale] invalid sale unit", zap.String("saleUnit", request.SaleUnit))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	storeSale := entity.StoreSale{
		Customer:    request.Customer,
		Phone:       request.Phone,
		StoreId:     request.StoreId,
		ItemId:      request.WarehouseItemId,
		Quantity:    request.Quantity,
		Price:       price,
		TotalPrice:  totalPrice,
		SendDate:    sendDate,
		IsSend:      false,
		SaleUnit:    saleUnit,
		PaymentType: paymentType,
		CreatedBy:   uuid.NullUUID{UUID: accountId, Valid: true},
	}

	nominal, err := decimal.NewFromString(request.StoreSalePayment.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !storeSale.TotalPrice.Equal(nominal) {
			s.log.Error("[CreateStoreSale] nominal is not equal to total price", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}

		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	err = s.repository.CreateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to create store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if request.StoreSalePayment.Nominal != "" &&
		request.StoreSalePayment.PaymentDate != "" &&
		request.StoreSalePayment.PaymentProof != "" &&
		request.StoreSalePayment.PaymentMethod != "" {
		paymentMethod := enum.ValueOfPaymentMethod(request.StoreSalePayment.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("[CreateStoreSale] invalid payment method", zap.String("paymentMethod", request.StoreSalePayment.PaymentMethod))
			return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
		}

		paymentDate, err := time.Parse("02-01-2006", request.StoreSalePayment.PaymentDate)
		if err != nil {
			s.log.Error("[CreateStoreSale] failed to parse payment date", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
		}

		storeSalePayment := entity.StoreSalePayment{
			PaymentDate:   paymentDate,
			StoreSaleId:   storeSale.Id,
			Nominal:       nominal,
			PaymentProof:  request.StoreSalePayment.PaymentProof,
			PaymentMethod: paymentMethod,
			CreatedBy:     uuid.NullUUID{UUID: accountId, Valid: true},
		}

		err = s.repository.CreateStoreSalePayment(&storeSalePayment)
		if err != nil {
			s.log.Error("[CreateStoreSale] failed to create store sale payment", zap.Error(err))
			return dto.StoreSaleResponse{}, err
		}
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))
	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) GetStoreSaleById(id uint64) (dto.StoreSaleResponse, error) {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) GetStoreSales(filter dto.GetStoreSaleFilter) (dto.StoreSaleListPaginationResponse, error) {
	storeSales, err := s.repository.GetStoreSales(filter)
	if err != nil {
		s.log.Error("[GetStoreSales] failed to get store sales", zap.Error(err))
		return dto.StoreSaleListPaginationResponse{}, err
	}

	storeSaleResponses := make([]dto.StoreSaleListResponse, len(storeSales))
	for i, storeSale := range storeSales {
		storeSaleResponses[i] = mapper.StoreSaleToListResponse(&storeSale)
	}

	totalData, err := s.repository.CountTotalStoreSale(
		dto.GetStoreSaleFilter{
			Date:          filter.Date,
			PaymentMethod: filter.PaymentMethod,
		},
	)
	if err != nil {
		s.log.Error("[GetStoreSales] failed to get store sales", zap.Error(err))
		return dto.StoreSaleListPaginationResponse{}, err
	}

	resp := dto.StoreSaleListPaginationResponse{
		TotalPage:  uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData:  totalData,
		StoreSales: storeSaleResponses,
	}

	return resp, nil
}

func (s *StoreService) CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("[CreateStoreSalePayment] invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse payment date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	storeSalePayment := entity.StoreSalePayment{
		StoreSaleId:   storeSaleId,
		PaymentDate:   paymentDate,
		PaymentMethod: paymentMethod,
		Nominal:       nominal,
		PaymentProof:  request.PaymentProof,
		CreatedBy:     uuid.NullUUID{UUID: accountId, Valid: true},
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSaleId)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("[CreateStoreSalePayment] store sale is already paid", zap.Uint64("id", storeSaleId))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already paid")
	}

	totalPayment := nominal
	for _, payment := range storeSale.Payments {
		totalPayment = totalPayment.Add(payment.Nominal)
	}

	if totalPayment.Equal(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(storeSale.TotalPrice) {
		s.log.Error("[CreateStoreSalePayment] total payment is greater than total price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.CreateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to create store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale.Payments = append(storeSale.Payments, storeSalePayment)
	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("[UpdateStoreSale] store sale is already sent", zap.Uint64("id", id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already sent")
	}

	storeSale.Customer = request.Customer
	storeSale.Phone = request.Phone
	storeSale.StoreId = request.StoreId
	storeSale.ItemId = request.WarehouseItemId
	storeSale.Quantity = request.Quantity
	storeSale.Price, err = decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to parse price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid price format")
	}

	storeSale.TotalPrice = storeSale.Price.Mul(decimal.NewFromInt(int64(request.Quantity)))

	storeSale.SaleUnit = enum.ValueOfSaleUnit(request.SaleUnit)
	if !storeSale.SaleUnit.IsValid() {
		s.log.Error("[UpdateStoreSale] invalid sale unit", zap.String("saleUnit", request.SaleUnit))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	storeSale.PaymentType = enum.ValueOfPaymentType(request.PaymentType)
	if !storeSale.PaymentType.IsValid() {
		s.log.Error("[UpdateStoreSale] invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	storeSale.SendDate, err = time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to parse send date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid send date format")
	}

	storeSale.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) UpdateStoreSalePayment(id uint64, request dto.UpdateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeSalePayment, err := s.repository.GetStoreSalePaymentById(id)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to get store sale payment by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSalePayment.StoreSaleId)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("[UpdateStoreSalePayment] store sale is already sent", zap.Uint64("id", storeSale.Id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already sent")
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("[CreateStoreSalePayment] store sale is already paid", zap.Uint64("id", storeSale.Id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already paid")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse payment date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	totalPayment := nominal
	for _, payment := range storeSale.Payments {
		if payment.Id != storeSalePayment.Id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.Equal(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(storeSale.TotalPrice) {
		s.log.Error("[CreateStoreSalePayment] total payment is greater than total price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	} else if totalPayment.LessThan(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	storeSalePayment.Nominal = nominal
	storeSalePayment.PaymentProof = request.PaymentProof
	storeSalePayment.PaymentDate = paymentDate
	storeSalePayment.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.UpdateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to update store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, payment := range storeSale.Payments {
		if payment.Id == id {
			storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
			remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
			storeSalePayments[i].Remaining = remainingPayment.String()
		} else {
			storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&payment)
			remainingPayment = remainingPayment.Sub(payment.Nominal)
			storeSalePayments[i].Remaining = remainingPayment.String()
		}
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) SendStoreSale(id uint64, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("[SendStoreSale] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("[SendStoreSale] store sale is already sent", zap.Uint64("id", id))
		return dto.StoreSaleResponse{}, err
	}

	storeSale.IsSend = true
	storeSale.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[SendStoreSale] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}
