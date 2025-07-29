package service

import (
	"database/sql"
	"fmt"
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
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
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
	itemService      IItemService
	customerService  ICustomerService
}

type IStoreService interface {
	CreateStore(request dto.CreateStoreRequest, createdBy uuid.UUID) (dto.StoreResponse, error)
	UpdateStore(id uint64, request dto.UpdateStoreRequest, updatedBy uuid.UUID) (dto.StoreResponse, error)
	DeleteStore(id uint64) error
	GetStores(filter dto.GetStoreFilter) ([]dto.StoreResponse, error)
	GetStoreDetailById(id uint64) (dto.StoreDetailResponse, error)
	GetStoreOverview(filter dto.GetStoreOverviewFilter) (dto.StoreOverview, error)

	CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, createdBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	CreateStoreRequestItemFromEggMonitoring(request dto.CreateStoreRequestItemRequest, createdBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) (dto.StoreRequestItemListPaginationResponse, error)
	WarehouseConfirmationStoreRequestItem(id uint64, request dto.WarehouseConfirmationStoreRequestItem, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	StoreConfirmationStoreRequestItem(id uint64, request dto.StoreConfirmationStoreRequestItem, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)
	SortingStoreRequestItem(id uint64, request dto.SortingStoreRequestItemRequest, updatedBy uuid.UUID) (dto.StoreRequestItemResponse, error)

	GetStoreItems(filter dto.GetStoreItemFilter) ([]dto.StoreItemResponse, error)
	GetStoreItemStocks(id uint64) (dto.StoreItemOverview, error)
	GetStoreItemByStoreIdAndItemId(storeId uint64, itemId uint64) (dto.StoreItemResponse, error)
	UpdateStoreItem(storeId uint64, itemId uint64, request dto.UpdateStoreItemRequest, updatedBy uuid.UUID) (dto.StoreItemResponse, error)
	GetEggStoreItemSummary(storeId uint64) ([]dto.EggStoreItemSummary, error)

	GetStoreItemHistories(filter dto.GetStoreItemHistoryFilter) (dto.StoreItemHistoryListPaginationResponse, error)
	GetStoreItemHistoryById(id uint64) (dto.StoreItemHistoryResponse, error)

	CreateStoreSale(request dto.CreateStoreSaleRequest, createdBy uuid.UUID) (dto.StoreSaleResponse, error)
	GetStoreSaleById(id uint64) (dto.StoreSaleResponse, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) (dto.StoreSaleListPaginationResponse, error)
	UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, userId uuid.UUID) (dto.StoreSaleResponse, error)
	DeleteStoreSale(id uint64, userId uuid.UUID) error
	SendStoreSale(id uint64, userId uuid.UUID) (dto.StoreSaleResponse, error)

	CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, userId uuid.UUID) (dto.StoreSaleResponse, error)
	UpdateStoreSalePayment(storeSaleId uint64, id uint64, request dto.UpdateStoreSalePaymentRequest, userId uuid.UUID) (dto.StoreSaleResponse, error)
	DeleteStoreSalePayment(storeSaleId uint64, id uint64, userId uuid.UUID) error

	CreateStoreSaleQueue(request dto.CreateStoreSaleQueueRequest, userId uuid.UUID) (dto.StoreSaleQueueResponse, error)
	GetStoreSaleQueue(id uint64) (dto.StoreSaleQueueResponse, error)
	GetStoreSaleQueues() ([]dto.StoreSaleQueueResponse, error)
	DeleteStoreSaleQueue(id uint64) error
}

func NewStoreService(log *zap.Logger, repository repository.IStoreRepository, cacheService cache.ICache, placementService IPlacementService, warehouseService IWarehouseService, userService IUserService, itemService IItemService, customerService ICustomerService) IStoreService {
	return &StoreService{
		log:              log,
		repository:       repository,
		cacheService:     cacheService,
		placementService: placementService,
		warehouseService: warehouseService,
		userService:      userService,
		itemService:      itemService,
		customerService:  customerService,
	}
}

// Todo : When created store auto create 3 egg object
func (s *StoreService) CreateStore(request dto.CreateStoreRequest, createdBy uuid.UUID) (dto.StoreResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

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

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.StoreResponse{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.CrackedEgg, constant.EggUnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.StoreResponse{}, err

	}

	brokenEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.BrokenEgg, constant.EggUnitPlastik, enum.ItemCategoryEgg)
	if err != nil {
		return dto.StoreResponse{}, err
	}

	storeItems := make([]entity.StoreItem, 0)
	storeItems = append(storeItems, entity.StoreItem{
		StoreId:  store.Id,
		ItemId:   goodEggItem.Id,
		Quantity: 0,
	})

	storeItems = append(storeItems, entity.StoreItem{
		StoreId:  store.Id,
		ItemId:   crackedEggItem.Id,
		Quantity: 0,
	})

	storeItems = append(storeItems, entity.StoreItem{
		StoreId:  store.Id,
		ItemId:   brokenEggItem.Id,
		Quantity: 0,
	})

	err = s.repository.CreateStoreItemsInBatch(&storeItems)
	if err != nil {
		s.log.Error("failed to create store items in batch", zap.Error(err))
		return dto.StoreResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
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

	warehouseItem, err := s.warehouseService.GetWarehouseItemByWarehouseIdAndItemId(request.WarehouseId, request.ItemId)
	if err != nil {
		return dto.StoreRequestItemResponse{}, err
	}

	if warehouseItem.Quantity < request.Quantity*float64(constant.TotalEggPerIkat) && warehouseItem.Item.Unit == constant.EggUnitKg {
		return dto.StoreRequestItemResponse{}, errx.BadRequest("insuficcient stock for request item")
	}

	storeRequestItem := entity.StoreRequestItem{
		WarehouseId: request.WarehouseId,
		ItemId:      request.ItemId,
		StoreId:     sql.NullInt64{Int64: int64(request.StoreId), Valid: true},
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

	totalData, err := s.repository.CountTotalStoreRequestItem(filter)
	if err != nil {
		s.log.Error("failed to count request items", zap.Error(err))
		return dto.StoreRequestItemListPaginationResponse{}, err
	}

	resp := dto.StoreRequestItemListPaginationResponse{
		StoreRequestItems: storeRequestItemResponses,
	}

	if filter.Page > 0 {
		resp.TotalData = totalData
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
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

func (s *StoreService) GetStoreItemStocks(id uint64) (dto.StoreItemOverview, error) {
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
		case constant.BrokenEgg:
			eggStoreItemSummaries = append(eggStoreItemSummaries, dto.EggStoreItemSummary{
				Name:     constant.BrokenEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.EggUnitPlastik,
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

	resp := dto.StoreItemHistoryListPaginationResponse{
		StoreItemHistories: response,
	}

	if filter.Page > 0 {
		resp.TotalData = uint64(totalData)
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
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

func (s *StoreService) CreateStoreSale(request dto.CreateStoreSaleRequest, userId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeItem, err := s.repository.GetStoreItemByStoreIdAndItemId(request.StoreId, request.ItemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeItem.Quantity -= request.Quantity
	storeItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateStoreItem(&storeItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse sent date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sent date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid price format")
	}

	totalPrice := price.Mul(decimal.NewFromFloat(request.Quantity))
	discountPrice := totalPrice.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	totalPrice = totalPrice.Sub(discountPrice)

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		s.log.Error("invalid sale unit", zap.String("saleUnit", request.SaleUnit))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	storeSale := entity.StoreSale{
		StoreId:     request.StoreId,
		ItemId:      request.ItemId,
		Quantity:    request.Quantity,
		Price:       price,
		TotalPrice:  totalPrice,
		SendDate:    sendDate,
		Discount:    request.Discount,
		IsSend:      false,
		SaleUnit:    saleUnit,
		PaymentType: paymentType,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.CustomerType == constant.OldCustomerType {
		if request.CustomerId < 1 {
			return dto.StoreSaleResponse{}, errx.BadRequest("customer id is required")
		}

		storeSale.CustomerId = request.CustomerId
	} else {
		customer := dto.CreateCustomerRequest{
			Name:        request.CustomerName,
			PhoneNumber: request.CustomerPhoneNumber,
		}

		if request.CustomerName == "" || request.CustomerPhoneNumber == "" {
			return dto.StoreSaleResponse{}, errx.BadRequest("customer name and customer phone number is required")
		}

		if request.CustomerPhoneNumber[:2] != "08" {
			return dto.StoreSaleResponse{}, errx.BadRequest("customer phone number must be in valid format 08")
		}

		// Saga pattern
		resp, err := s.customerService.CreateCustomer(customer)
		if err != nil {
			return dto.StoreSaleResponse{}, err
		}

		storeSale.CustomerId = resp.Id
	}

	nominal, err := decimal.NewFromString(request.StoreSalePayment.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !storeSale.TotalPrice.Equal(nominal) {
			s.log.Error("nominal is not equal to total price", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}

		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	err = s.repository.CreateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("failed to create store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if request.StoreSalePayment.Nominal != "" &&
		request.StoreSalePayment.PaymentDate != "" &&
		request.StoreSalePayment.PaymentProof != "" &&
		request.StoreSalePayment.PaymentMethod != "" {
		paymentMethod := enum.ValueOfPaymentMethod(request.StoreSalePayment.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("invalid payment method", zap.String("paymentMethod", request.StoreSalePayment.PaymentMethod))
			return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
		}

		paymentDate, err := time.Parse("02-01-2006", request.StoreSalePayment.PaymentDate)
		if err != nil {
			s.log.Error("failed to parse payment date", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
		}

		storeSalePayment := entity.StoreSalePayment{
			PaymentDate:   paymentDate,
			StoreSaleId:   storeSale.Id,
			Nominal:       nominal,
			PaymentProof:  request.StoreSalePayment.PaymentProof,
			PaymentMethod: paymentMethod,
			CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
		}

		err = s.repository.CreateStoreSalePayment(&storeSalePayment)
		if err != nil {
			s.log.Error("failed to create store sale payment", zap.Error(err))
			return dto.StoreSaleResponse{}, err
		}
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
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
		s.log.Error("failed to get store sale by id", zap.Error(err))
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
		s.log.Error("failed to get store sales", zap.Error(err))
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
		s.log.Error("failed to get store sales", zap.Error(err))
		return dto.StoreSaleListPaginationResponse{}, err
	}

	resp := dto.StoreSaleListPaginationResponse{
		StoreSales: storeSaleResponses,
	}

	if filter.Page > 0 {
		resp.TotalData = totalData
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (s *StoreService) CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, userId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	storeSalePayment := entity.StoreSalePayment{
		StoreSaleId:   storeSaleId,
		PaymentDate:   paymentDate,
		PaymentMethod: paymentMethod,
		Nominal:       nominal,
		PaymentProof:  request.PaymentProof,
		CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSaleId)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("store sale is already paid", zap.Uint64("id", storeSaleId))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already paid")
	}

	totalPayment := nominal
	for _, payment := range storeSale.Payments {
		totalPayment = totalPayment.Add(payment.Nominal)
	}

	if totalPayment.Equal(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(storeSale.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.CreateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("failed to create store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
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

func (s *StoreService) UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, userId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("store sale is already sent", zap.Uint64("id", id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already sent")
	}

	storeItem, err := s.repository.GetStoreItemByStoreIdAndItemId(storeSale.StoreId, storeSale.ItemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeItem.Quantity += storeSale.Quantity - request.Quantity
	storeItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateStoreItem(&storeItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price from string", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale.Quantity = request.Quantity
	totalPrice := price.Mul(decimal.NewFromFloat(request.Quantity))
	discountPrice := totalPrice.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	storeSale.TotalPrice = totalPrice.Sub(discountPrice)
	storeSale.Price = price
	storeSale.Discount = request.Discount

	storeSale.SendDate, err = time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse send date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid send date format")
	}

	storeSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
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

func (s *StoreService) UpdateStoreSalePayment(storeSaleId uint64, id uint64, request dto.UpdateStoreSalePaymentRequest, userId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeSalePayment, err := s.repository.GetStoreSalePaymentById(id)
	if err != nil {
		s.log.Error("failed to get store sale payment by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSaleId)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("store sale is already sent", zap.Uint64("id", storeSale.Id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already sent")
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("store sale is already paid", zap.Uint64("id", storeSale.Id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already paid")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
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
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	} else if totalPayment.LessThan(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	storeSalePayment.PaymentMethod = paymentMethod
	storeSalePayment.Nominal = nominal
	storeSalePayment.PaymentProof = request.PaymentProof
	storeSalePayment.PaymentDate = paymentDate
	storeSalePayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.UpdateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("failed to update store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
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

func (s *StoreService) DeleteStoreSalePayment(storeSaleId uint64, id uint64, userId uuid.UUID) error {
	s.repository.UseTx(false)

	storeSale, err := s.repository.GetStoreSaleById(storeSaleId)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return err
	}

	if storeSale.IsSend {
		s.log.Error("store sale is already sent", zap.Uint64("id", storeSale.Id))
		return errx.BadRequest("store sale is already sent")
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("store sale is already paid", zap.Uint64("id", storeSale.Id))
		return errx.BadRequest("store sale is already paid")
	}

	totalPayment := decimal.Zero
	for _, payment := range storeSale.Payments {
		if id != payment.Id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.LessThan(storeSale.TotalPrice) && totalPayment.GreaterThan(decimal.Zero) {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
		storeSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	} else if totalPayment.LessThan(decimal.Zero) {
		s.log.Error("delete this payment make minus", zap.Error(err))
		return errx.BadRequest("delete this payment make minus")
	}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("failed to update store sale", zap.Error(err))
		return err
	}

	err = s.repository.DeleteStoreSalePayment(id)
	if err != nil {
		s.log.Error("failed to update store sale", zap.Error(err))
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *StoreService) SendStoreSale(id uint64, userId uuid.UUID) (dto.StoreSaleResponse, error) {
	// Todo : emit event for store item histories

	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("store sale is already sent", zap.Uint64("id", id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale already send")
	}

	storeSale.IsSend = true
	storeSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
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

func (s *StoreService) DeleteStoreSale(id uint64, userId uuid.UUID) error {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("failed to get store sale by id", zap.Error(err))
		return err
	}

	if storeSale.IsSend {
		s.log.Error("store sale is already sent", zap.Uint64("id", id))
		return errx.BadRequest("store sale already send")
	}

	storeItem, err := s.repository.GetStoreItemByStoreIdAndItemId(storeSale.StoreId, storeSale.ItemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return err
	}

	storeItem.Quantity += storeSale.Quantity
	storeItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateStoreItem(&storeItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return err
	}

	err = s.repository.DeleteStoreSale(id)
	if err != nil {
		s.log.Error("failed to delete store sale", zap.Error(err))
		return err
	}

	return nil
}

func (s *StoreService) GetStoreOverview(filter dto.GetStoreOverviewFilter) (dto.StoreOverview, error) {
	s.repository.UseTx(false)

	startDate, endDate := util.GetStartDateAndEndDateInMonth(int(filter.Year), time.Month(filter.Month.Value()))

	storeSales, err := s.repository.GetStoreSales(dto.GetStoreSaleFilter{
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
	})
	if err != nil {
		s.log.Error("failed to get store sales", zap.Error(err))
		return dto.StoreOverview{}, err
	}

	goodEggInKg := float64(0)
	crackedEggInKg := float64(0)
	brokenEggInPlastik := float64(0)

	income := decimal.Zero
	receivables := decimal.Zero

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.StoreOverview{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.CrackedEgg, constant.EggUnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.StoreOverview{}, err
	}

	brokenEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.BrokenEgg, constant.EggUnitPlastik, enum.ItemCategoryEgg)
	if err != nil {
		return dto.StoreOverview{}, err
	}

	for _, storeSale := range storeSales {
		if storeSale.SaleUnit.String() == constant.EggUnitKg && goodEggItem.Id == storeSale.ItemId {
			goodEggInKg += storeSale.Quantity
		} else if storeSale.SaleUnit.String() == constant.EggUnitIkat && goodEggItem.Id == storeSale.ItemId {
			goodEggInKg += storeSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if storeSale.SaleUnit.String() == constant.EggUnitKg && crackedEggItem.Id == storeSale.ItemId {
			crackedEggInKg += storeSale.Quantity
		} else if storeSale.SaleUnit.String() == constant.EggUnitIkat && crackedEggItem.Id == storeSale.ItemId {
			crackedEggInKg += storeSale.Quantity * float64(constant.TotalEggPerIkat)
		} else if storeSale.SaleUnit.String() == constant.EggUnitPlastik && brokenEggItem.Id == storeSale.ItemId {
			brokenEggInPlastik += storeSale.Quantity
		}

		payment := decimal.Zero
		for _, storeSalePayment := range storeSale.Payments {
			payment = payment.Add(storeSalePayment.Nominal)
		}

		income = income.Add(payment)
		receivables = receivables.Add(storeSale.TotalPrice.Sub(payment))
	}

	storeOverviewDetail := dto.StoreOverviewDetail{
		TotalIncome:        income.String(),
		TotalReceivables:   receivables.String(),
		GoodEggInKg:        goodEggInKg,
		GoodEggInIkat:      math.Floor(goodEggInKg / float64(constant.TotalEggPerIkat)),
		CrackedEggInKg:     crackedEggInKg,
		CrackedEggInIkat:   math.Floor(crackedEggInKg / float64(constant.TotalEggPerIkat)),
		BrokenEggInPlastik: brokenEggInPlastik,
	}

	storeGraphs := make([]dto.StoreGraphResponse, 0)
	switch filter.OverviewGraphTime.Value() {
	case enum.OverviewGraphTimeThisWeek:
		storeGraphs, err = s.buildStoreOverviewWeeklyGraph(filter.StoreId, filter.ItemId)
	case enum.OverviewGraphTimeThisMonth:
		storeGraphs, err = s.buildStoreOverviewMonthlyGraph(filter.StoreId, filter.ItemId)
	case enum.OverviewGraphTimeThisYear:
		storeGraphs, err = s.buildStoreOverviewYearlyGraph(filter.StoreId, filter.ItemId)
	}
	if err != nil {
		return dto.StoreOverview{}, err
	}

	return dto.StoreOverview{
		StoreOverviewDetail: storeOverviewDetail,
		StoreGraphs:         storeGraphs,
	}, nil
}

func (s *StoreService) buildStoreOverviewWeeklyGraph(storeId uint64, itemId uint64) ([]dto.StoreGraphResponse, error) {
	endDate := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	startDate := endDate.AddDate(0, 0, -7)

	weekStoreSales, err := s.repository.GetStoreSales(dto.GetStoreSaleFilter{
		StoreId:   storeId,
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
		ItemId:    itemId,
	})
	if err != nil {
		s.log.Error("failed to get store sales weekly", zap.Error(err))
		return nil, err
	}

	graphs := make([]dto.StoreGraphResponse, 0)
	for day := startDate; day.Before(endDate); day = day.AddDate(0, 0, 1) {
		var itemSale float64
		for _, storeSale := range weekStoreSales {
			if isSameDate(day, storeSale.CreatedAt) {
				if storeSale.SaleUnit.String() == constant.EggUnitKg {
					itemSale += storeSale.Quantity
				} else if storeSale.SaleUnit.String() == constant.EggUnitIkat {
					itemSale += storeSale.Quantity * float64(constant.TotalEggPerIkat)
				} else if storeSale.SaleUnit.String() == constant.EggUnitPlastik {
					itemSale += storeSale.Quantity
				}
			}
		}
		graphs = append(graphs, dto.StoreGraphResponse{
			Key:   day.Format("2006-01-02"),
			Value: itemSale,
		})
	}
	return graphs, nil
}

func (s *StoreService) buildStoreOverviewMonthlyGraph(storeId uint64, itemId uint64) ([]dto.StoreGraphResponse, error) {
	weekMaps := util.GetFourWeekRanges(time.Now().Year(), time.Now().Month())
	startDate, endDate := util.GetStartDateAndEndDateInMonth(time.Now().Year(), time.Now().Month())

	monthStoreSales, err := s.repository.GetStoreSales(dto.GetStoreSaleFilter{
		StoreId:   storeId,
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
		ItemId:    itemId,
	})
	if err != nil {
		s.log.Error("failed to get store sale monthly", zap.Error(err))
		return nil, err
	}

	itemSales := make(map[int]float64)
	for _, storeSale := range monthStoreSales {
		week := util.FindWeek(storeSale.CreatedAt, weekMaps)
		if week > 0 {
			if storeSale.SaleUnit.String() == constant.EggUnitKg {
				itemSales[week] += storeSale.Quantity
			} else if storeSale.SaleUnit.String() == constant.EggUnitIkat {
				itemSales[week] += storeSale.Quantity * float64(constant.TotalEggPerIkat)
			} else if storeSale.SaleUnit.String() == constant.EggUnitPlastik {
				itemSales[week] += storeSale.Quantity
			}
		}
	}

	keys := util.GetSortedKeys(weekMaps)
	graphs := make([]dto.StoreGraphResponse, 0)
	for _, k := range keys {
		graphs = append(graphs, dto.StoreGraphResponse{
			Key:   fmt.Sprintf("Minggu %d", k),
			Value: itemSales[k],
		})
	}

	return graphs, nil
}

func (s *StoreService) buildStoreOverviewYearlyGraph(storeId uint64, itemId uint64) ([]dto.StoreGraphResponse, error) {
	monthMaps := util.GetTwelveMonthRanges(time.Now().Year())
	startDate, endDate := util.GetStartDateAndEndDateInYear(time.Now().Year())

	yearStoreSales, err := s.repository.GetStoreSales(dto.GetStoreSaleFilter{
		StoreId:   storeId,
		StartDate: param.DateParam(startDate),
		EndDate:   param.DateParam(endDate),
		ItemId:    itemId,
	})
	if err != nil {
		s.log.Error("failed to get store items yearly", zap.Error(err))
		return nil, err
	}

	itemSales := make(map[int]float64)
	for _, storeSale := range yearStoreSales {
		month := util.FindMonth(storeSale.CreatedAt, monthMaps)
		if month > 0 {
			if storeSale.SaleUnit.String() == constant.EggUnitKg {
				itemSales[month] += storeSale.Quantity
			} else if storeSale.SaleUnit.String() == constant.EggUnitIkat {
				itemSales[month] += storeSale.Quantity * float64(constant.TotalEggPerIkat)
			} else if storeSale.SaleUnit.String() == constant.EggUnitPlastik {
				itemSales[month] += storeSale.Quantity
			}
		}
	}

	keys := util.GetSortedKeys(monthMaps)
	graphs := make([]dto.StoreGraphResponse, 0)
	for _, k := range keys {
		graphs = append(graphs, dto.StoreGraphResponse{
			Key:   util.IndoMonthName(k),
			Value: itemSales[k],
		})
	}
	return graphs, nil
}

func (s *StoreService) CreateStoreSaleQueue(request dto.CreateStoreSaleQueueRequest, userId uuid.UUID) (dto.StoreSaleQueueResponse, error) {
	s.repository.UseTx(false)

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse send date", zap.Error(err))
		return dto.StoreSaleQueueResponse{}, err
	}

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		return dto.StoreSaleQueueResponse{}, errx.BadRequest("invalid sale unit")
	}

	customerType := enum.ValueOfCustomerType(request.CustomerType)
	if !customerType.IsValid() {
		return dto.StoreSaleQueueResponse{}, errx.BadRequest("invalid customer type")
	}

	data := entity.StoreSaleQueue{
		ItemId:       request.ItemId,
		StoreId:      request.StoreId,
		Quantity:     request.Quantity,
		SaleUnit:     saleUnit,
		SendDate:     sendDate,
		CustomerType: customerType,
		CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	if customerType == enum.CustomerTypeNew {
		if request.CustomerName == "" || request.CustomerPhoneNumber == "" {
			return dto.StoreSaleQueueResponse{}, errx.BadRequest("customer name and phone number is required")
		}

		data.CustomerName = sql.NullString{String: request.CustomerName, Valid: true}
		data.CustomerPhoneNumber = sql.NullString{String: request.CustomerPhoneNumber, Valid: true}

	} else {
		if request.CustomerId < 1 {
			return dto.StoreSaleQueueResponse{}, errx.BadRequest("customer id is required")
		}

		data.CustomerId = sql.NullInt64{Int64: int64(request.CustomerId), Valid: true}
	}

	err = s.repository.CreateStoreSaleQueue(&data)
	if err != nil {
		s.log.Error("failed create store sale queue", zap.Error(err))
		return dto.StoreSaleQueueResponse{}, err
	}

	data, err = s.repository.GetStoreSaleQueueById(data.Id)
	if err != nil {
		return dto.StoreSaleQueueResponse{}, err
	}

	return mapper.StoreSaleQueueToResponse(&data), nil
}

func (s *StoreService) GetStoreSaleQueue(id uint64) (dto.StoreSaleQueueResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetStoreSaleQueueById(id)
	if err != nil {
		s.log.Error("failed get store sale queue by id", zap.Error(err))
		return dto.StoreSaleQueueResponse{}, err
	}

	return mapper.StoreSaleQueueToResponse(&data), nil
}

func (s *StoreService) GetStoreSaleQueues() ([]dto.StoreSaleQueueResponse, error) {
	s.repository.UseTx(false)

	// Todo : formula for integrated planning
	storeSaleQueues, err := s.repository.GetStoreSaleQueues()
	if err != nil {
		s.log.Error("failed get store sale queues", zap.Error(err))
		return nil, err
	}

	response := make([]dto.StoreSaleQueueResponse, 0)
	for _, storeSaleQueue := range storeSaleQueues {
		response = append(response, mapper.StoreSaleQueueToResponse(&storeSaleQueue))
	}

	return response, nil
}

func (s *StoreService) DeleteStoreSaleQueue(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteStoreSaleQueue(id)
	if err != nil {
		s.log.Error("failed delete store sale queue", zap.Error(err))
		return err
	}

	return nil
}
