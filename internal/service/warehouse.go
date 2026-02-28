package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"math"
	"sort"
	"strings"
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

type WarehouseService struct {
	log              *zap.Logger
	repository       repository.IWarehouseRepository
	cacheService     cache.ICache
	placementService IPlacementService
	itemService      IItemService
	customerService  ICustomerService
}

type IWarehouseService interface {
	GetWarehouses(filter dto.GetWarehouseFilter) ([]dto.WarehouseDetailResponse, error)
	CreateWarehouse(request dto.CreateWarehouseRequest, userId uuid.UUID) (dto.WarehouseResponse, error)
	DeleteWarehouse(id uint64) error
	UpdateWarehouse(id uint64, request dto.UpdateWarehouseRequest, userId uuid.UUID) (dto.WarehouseResponse, error)
	GetWarehouseWithUsersById(id uint64) (dto.WarehouseWithUsersResponse, error)
	GetWarehouseOverview(id uint64) (dto.WarehouseOverview, error)

	CreateWarehouseItem(request dto.CreateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error)
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error)
	GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) (dto.WarehouseItemResponse, error)
	UpdateWarehouseItem(warehouseId uint64, itemId uint64, request dto.UpdateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error)
	UpdateWarehouseItemCorn(id uint64, request dto.UpdateWarehouseItemCornRequest, userId uuid.UUID) (dto.WarehouseItemCornResponse, error)
	DeleteWarehouseItem(warehouseId uint64, itemId uint64) error
	GetEggWarehouseItemSummary(warehouseId uint64) ([]dto.EggWarehouseItemSummaryResponse, error)
	GetCornWarehouseItemSummary(warehouseId uint64) (dto.CornWarehouseItemSummaryResponse, error)

	GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) (dto.WarehouseItemHistoryListPaginationResponse, error)
	GetWarehouseItemHistoryById(id uint64) (dto.WarehouseItemHistoryResponse, error)

	CreateWarehouseSale(request dto.CreateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	GetWarehouseSaleById(id uint64) (dto.WarehouseSaleResponse, error)
	GetWarehouseSales(filter dto.GetWarehouseSaleFilter) (dto.WarehouseSaleListPaginationResponse, error)
	UpdateWarehouseSale(id uint64, request dto.UpdateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	DeleteWarehouseSale(id uint64, userId uuid.UUID) error

	CreateWarehouseSalePayment(warehouseSaleId uint64, request dto.CreateWarehouseSalePaymentRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	UpdateWarehouseSalePayment(warehouseSaleId uint64, id uint64, request dto.UpdateWarehouseSalePaymentRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	DeleteWarehouseSalePayment(warehouseSaleId uint64, id uint64, userId uuid.UUID) error

	SendWarehouseSale(id uint64, userId uuid.UUID) (dto.WarehouseSaleResponse, error)

	CreateWarehouseSaleQueue(request dto.CreateWarehouseSaleQueueRequest, userId uuid.UUID) (dto.WarehouseSaleQueueResponse, error)
	DeleteWarehouseSaleQueue(id uint64) error
	GetWarehouseSaleQueues(filter dto.GetWarehouseSaleQueueFilter) ([]dto.WarehouseSaleQueueResponse, error)
	GetWarehouseSaleQueue(id uint64) (dto.WarehouseSaleQueueResponse, error)
	AllocateWarehouseSaleQueue(id uint64, request dto.CreateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)

	CreateWarehouseItemProcurementDraft(request dto.CreateWarehouseItemProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemProcurementDraftResponse, error)
	GetWarehouseItemProcurementDrafts(filter dto.GetWarehouseItemProcurementDraftFilter) ([]dto.WarehouseItemProcurementDraftResponse, error)
	GetWarehouseItemProcurementDraft(id uint64) (dto.WarehouseItemProcurementDraftResponse, error)
	UpdateWarehouseItemProcurementDraft(id uint64, request dto.UpdateWarehouseItemProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemProcurementDraftResponse, error)
	DeleteWarehouseItemProcurementDraft(id uint64) error
	ConfirmationWarehouseItemProcurementDraft(id uint64, request dto.CreateWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error)

	CreateWarehouseItemProcurement(request dto.CreateWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error)
	GetWarehouseItemProcurements(filter dto.GetWarehouseItemProcurementFilter) (dto.WarehouseItemProcurementListPaginationResponse, error)
	GetWarehouseItemProcurement(id uint64) (dto.WarehouseItemProcurementResponse, error)
	ArrivalConfirmationWarehouseItemProcurement(id uint64, request dto.ArrivalConfirmationWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error)

	CreateWarehouseItemProcurementPayment(warehouseItemProcurementId uint64, request dto.CreateWarehouseItemProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error)
	UpdateWarehouseItemProcurementPayment(id uint64, warehouseItemProcurementId uint64, request dto.UpdateWarehouseItemProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error)
	DeleteWarehouseItemProcurementPayment(id uint64, warehouseItemProcurementId uint64, userId uuid.UUID) error

	CreateWarehouseItemCornProcurementDraft(request dto.CreateWarehouseItemCornProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementDraftResponse, error)
	GetWarehouseItemCornProcurementDrafts(filter dto.GetWarehouseItemCornProcurementDraftFilter) ([]dto.WarehouseItemCornProcurementDraftResponse, error)
	GetWarehouseItemCornProcurementDraft(id uint64) (dto.WarehouseItemCornProcurementDraftResponse, error)
	UpdateWarehouseItemCornProcurementDraft(id uint64, request dto.UpdateWarehouseItemCornProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementDraftResponse, error)
	DeleteWarehouseItemCornProcurementDraft(id uint64) error
	ConfirmationWarehouseItemCornProcurementDraft(id uint64, request dto.CreateWarehouseItemCornProcurementRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error)

	CreateWarehouseItemCornProcurement(request dto.CreateWarehouseItemCornProcurementRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error)
	GetWarehouseItemCornProcurements(filter dto.GetWarehouseItemCornProcurementFilter) (dto.WarehouseItemCornProcurementListPaginationResponse, error)
	GetWarehouseItemCornProcurement(id uint64) (dto.WarehouseItemCornProcurementResponse, error)
	ArrivalConfirmationWarehouseItemCornProcurement(id uint64, request dto.ArrivalConfirmationWarehouseItemCornProcurementRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error)

	CreateWarehouseItemCornProcurementPayment(warehouseItemCornProcurementId uint64, request dto.CreateWarehouseItemCornProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error)
	UpdateWarehouseItemCornProcurementPayment(id uint64, warehouseItemCornProcurementId uint64, request dto.UpdateWarehouseItemCornProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error)
	DeleteWarehouseItemCornProcurementPayment(id uint64, warehouseItemProcurementId uint64, userId uuid.UUID) error

	GetWarehouseItemCornPrices() ([]dto.WarehouseItemCornPriceResponse, error)

	CreateRawFeed(request dto.CreateRawFeedRequest, userId uuid.UUID) error
	CreateReadyToEatFeed(request dto.CreateReadyToEatFeedRequest, userId uuid.UUID) error

	ReduceWarehouseItemForFeed(warehouseId uint64, request []dto.ReduceFeedRequest, userId uuid.UUID, cageName string) error
}

func NewWarehouseService(log *zap.Logger, repository repository.IWarehouseRepository, cacheService cache.ICache, placementService IPlacementService, itemService IItemService, customerService ICustomerService) IWarehouseService {
	return &WarehouseService{
		log:              log,
		repository:       repository,
		cacheService:     cacheService,
		placementService: placementService,
		itemService:      itemService,
		customerService:  customerService,
	}
}

func (s *WarehouseService) GetWarehouseWithUsersById(id uint64) (dto.WarehouseWithUsersResponse, error) {
	s.repository.UseTx(false)

	warehouse, err := s.repository.GetWarehouseById(id)
	if err != nil {
		s.log.Error("failed to get warehouse by id")
		return dto.WarehouseWithUsersResponse{}, err
	}

	warehousePlacements, err := s.placementService.GetWarehousePlacementByWarehouseId(id)
	if err != nil {
		return dto.WarehouseWithUsersResponse{}, err
	}

	userResponses := make([]dto.UserListResponse, 0)
	for _, e := range warehousePlacements {
		userResponses = append(userResponses, e.User)
	}

	isItemsEmpty := true
	for _, e := range warehouse.WarehouseItems {
		if e.Quantity != 0 {
			isItemsEmpty = false
			break
		}
	}

	if isItemsEmpty {
		for _, e := range warehouse.WarehouseItemCorns {
			if e.Quantity != 0 {
				isItemsEmpty = false
				break
			}
		}
	}

	return dto.WarehouseWithUsersResponse{
		Id:            warehouse.Id,
		Name:          warehouse.Name,
		Location:      mapper.LocationToResponse(&warehouse.Location),
		CornCapacity:  warehouse.CornCapacity,
		TotalEmployee: uint64(len(warehouse.WarehousePlacement)),
		IsItemsEmpty:  isItemsEmpty,
		Users:         userResponses,
	}, nil
}

func (s *WarehouseService) UpdateWarehouse(id uint64, request dto.UpdateWarehouseRequest, updateBy uuid.UUID) (dto.WarehouseResponse, error) {
	s.repository.UseTx(false)

	warehouse, err := s.repository.GetWarehouseById(id)
	if err != nil {
		s.log.Error("failed to get warehouse by id", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	warehouse.Name = request.Name
	warehouse.LocationId = request.LocationId
	warehouse.CornCapacity = request.CornCapacity
	warehouse.UpdatedBy = uuid.NullUUID{UUID: updateBy, Valid: true}

	err = s.repository.UpdateWarehouse(&warehouse)
	if err != nil {
		s.log.Error("failed to get udpate warehouse", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	warehouse, err = s.repository.GetWarehouseById(warehouse.Id)
	if err != nil {
		s.log.Error("failed to get warehouse by id", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	return mapper.WarehouseToResponse(&warehouse), nil
}

func (s *WarehouseService) CreateWarehouse(request dto.CreateWarehouseRequest, userId uuid.UUID) (dto.WarehouseResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouse := entity.Warehouse{
		Name:         request.Name,
		LocationId:   request.LocationId,
		CornCapacity: request.CornCapacity,
		CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateWarehouse(&warehouse)
	if err != nil {
		s.log.Error("failed to create warehouse", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.UnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.WarehouseResponse{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.CrackedEgg, constant.UnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.WarehouseResponse{}, err

	}

	brokenEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.BrokenEgg, constant.UnitPlastik, enum.ItemCategoryEgg)
	if err != nil {
		return dto.WarehouseResponse{}, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseResponse{}, err
	}

	warehouseItems := make([]entity.WarehouseItem, 0)
	warehouseItems = append(warehouseItems, entity.WarehouseItem{
		WarehouseId: warehouse.Id,
		ItemId:      goodEggItem.Id,
		Quantity:    0,
	})

	warehouseItems = append(warehouseItems, entity.WarehouseItem{
		WarehouseId: warehouse.Id,
		ItemId:      crackedEggItem.Id,
		Quantity:    0,
	})

	warehouseItems = append(warehouseItems, entity.WarehouseItem{
		WarehouseId: warehouse.Id,
		ItemId:      brokenEggItem.Id,
		Quantity:    0,
	})

	warehouseItems = append(warehouseItems, entity.WarehouseItem{
		WarehouseId: warehouse.Id,
		ItemId:      cornItem.Id,
		Quantity:    0,
	})

	err = s.repository.CreateWarehouseItemInBatch(&warehouseItems)
	if err != nil {
		s.log.Error("failed to create warehouse item in batch", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaaction", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	warehouse, err = s.repository.GetWarehouseById(warehouse.Id)
	if err != nil {
		s.log.Error("failed to get warehouse by id", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	return mapper.WarehouseToResponse(&warehouse), nil
}

func (s *WarehouseService) GetWarehouses(filter dto.GetWarehouseFilter) ([]dto.WarehouseDetailResponse, error) {
	s.repository.UseTx(false)

	warehouses, err := s.repository.GetWarehouses(filter)
	if err != nil {
		s.log.Error("failed to get warehouses", zap.Error(err))
		return nil, err
	}

	warehouseResponses := make([]dto.WarehouseDetailResponse, 0, len(warehouses))
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, mapper.WarehouseDetailToResponse(&warehouse))
	}

	return warehouseResponses, nil
}

func (s *WarehouseService) DeleteWarehouse(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouse(id)
	if err != nil {
		s.log.Error("failed to delete warehouse", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseItem(request dto.CreateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItem := entity.WarehouseItem{
		WarehouseId: request.WarehouseId,
		ItemId:      request.ItemId,
		Quantity:    request.Quantity,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateWarehouseItem(&stockWarehouseItem)
	if err != nil {
		s.log.Error("failed to create warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	stockWarehouseItem, err = s.repository.GetWarehouseItemByWarehouseIdAndItemId(
		stockWarehouseItem.WarehouseId,
		stockWarehouseItem.ItemId,
	)
	if err != nil {
		s.log.Error("failed to get warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	return mapper.WarehouseItemToResponse(&stockWarehouseItem), nil
}

func (s *WarehouseService) GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	warehouseItems, err := s.repository.GetWarehouseItems(filter)
	if err != nil {
		s.log.Error("failed to get warehouse items", zap.Error(err))
		return nil, err
	}

	warehouseItemsResponses := make([]dto.WarehouseItemResponse, 0, len(warehouseItems))
	for _, item := range warehouseItems {
		warehouseItemsResponses = append(warehouseItemsResponses, mapper.WarehouseItemToResponse(&item))
	}

	return warehouseItemsResponses, nil
}

func (s *WarehouseService) GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item by warehouse id and item id", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	return mapper.WarehouseItemToResponse(&stockWarehouseItem), nil
}

func (s *WarehouseService) UpdateWarehouseItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	jsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemName:       warehouseItem.Item.Name,
		ItemUnit:       warehouseItem.Item.Unit,
		Source:         warehouseItem.Warehouse.Name,
		Destination:    "-",
		QuantityBefore: warehouseItem.Quantity,
		QuantityAfter:  request.Quantity,
		UserId:         userId,
		Status:         enum.ItemHistoryStockUpdated,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.WarehouseItemResponse{}, errx.BadRequest("failed parsed struct into json")
	}

	warehouseItem.Quantity = request.Quantity
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonParsed)

	warehouseStockItemResponse := mapper.WarehouseItemToResponse(&warehouseItem)

	return warehouseStockItemResponse, nil
}

func (s *WarehouseService) UpdateWarehouseItemCorn(id uint64, request dto.UpdateWarehouseItemCornRequest, userId uuid.UUID) (dto.WarehouseItemCornResponse, error) {
	s.repository.UseTx(false)

	warehouseItemCorn, err := s.repository.GetWarehouseItemCorn(id)
	if err != nil {
		s.log.Error("failed get warehouse item corn", zap.Error(err))
		return dto.WarehouseItemCornResponse{}, err
	}

	itemCorn, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornResponse{}, err
	}

	if warehouseItemCorn.Warehouse.CornCapacity < request.Quantity {
		return dto.WarehouseItemCornResponse{}, errx.BadRequest("quantity is more than max capacity")
	}

	jsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemName:       itemCorn.Name,
		ItemUnit:       itemCorn.Unit,
		Source:         warehouseItemCorn.Warehouse.Name,
		Destination:    "-",
		QuantityBefore: warehouseItemCorn.Quantity,
		QuantityAfter:  request.Quantity,
		UserId:         userId,
		Status:         enum.ItemHistoryStockUpdated,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.WarehouseItemCornResponse{}, errx.BadRequest("failed parsed struct into json")
	}

	warehouseItemCorn.Quantity = request.Quantity
	warehouseItemCorn.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItemCorn(&warehouseItemCorn)
	if err != nil {
		s.log.Error("failed update warehouse item corn", zap.Error(err))
		return dto.WarehouseItemCornResponse{}, err
	}

	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonParsed)

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornResponse{}, err
	}

	return mapper.WarehouseItemCornToResponse(&warehouseItemCorn, &cornItem), nil
}

func (s *WarehouseService) DeleteWarehouseItem(warehouseId uint64, warehouseItemId uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouseItemByWarehouseIdAndItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("failed to delete warehouse item", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) GetWarehouseOverview(id uint64) (dto.WarehouseOverview, error) {
	s.repository.UseTx(false)

	warehouseItems, err := s.repository.GetWarehouseItems(dto.GetWarehouseItemFilter{
		WarehouseId: id,
	})
	if err != nil {
		s.log.Error("failed to get warehouse items", zap.Error(err))
		return dto.WarehouseOverview{}, err
	}

	withZeroQuantity := false
	warehouseItemCorns, err := s.repository.GetWarehouseItemCorns(dto.GetWarehouseItemCornFilter{
		WarehouseId:      id,
		FromNewest:       true,
		WithZeroQuantity: &withZeroQuantity,
	})
	if err != nil {
		s.log.Error("failed get warehouse item corns", zap.Error(err))
		return dto.WarehouseOverview{}, err
	}

	warehouseItemCornResponses := make([]dto.WarehouseItemCornResponse, 0)
	warehouseItemEggResponses := make([]dto.WarehouseItemResponse, 0)
	warehouseItemEquipmentResponses := make([]dto.WarehouseItemResponse, 0)

	totalSafeStock := 0
	totalDangerStock := 0

	for _, e := range warehouseItems {
		res := mapper.WarehouseItemToResponse(&e)
		switch res.Description {
		case constant.WarehouseItemDescriptionDanger:
			totalDangerStock++
		case constant.WarehouseItemDescriptionSafe:
			totalSafeStock++
		}

		if e.Item.Category == enum.ItemCategoryEgg {
			warehouseItemEggResponses = append(warehouseItemEggResponses, res)
		} else if e.Item.Category != enum.ItemCategoryCornMaterial && e.Item.Category != enum.ItemCategoryEgg {
			warehouseItemEquipmentResponses = append(warehouseItemEquipmentResponses, res)
		}
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseOverview{}, err
	}

	for _, e := range warehouseItemCorns {
		res := mapper.WarehouseItemCornToResponse(&e, &cornItem)
		warehouseItemCornResponses = append(warehouseItemCornResponses, res)
	}

	warehouseItemProcurementSentOffCount, err := s.repository.CountWarehouseItemProcurement(dto.GetWarehouseItemProcurementFilter{
		WarehouseId: id,
		Status:      param.ProcurementStatusParam(enum.ProcurementStatusSentOff),
	})
	if err != nil {
		s.log.Error("failed to count warehouse item procurements", zap.Error(err))
		return dto.WarehouseOverview{}, err
	}

	warehouse, err := s.repository.GetWarehouseById(id)
	if err != nil {
		s.log.Error("failed get warehouse by id", zap.Error(err))
		return dto.WarehouseOverview{}, err
	}

	return dto.WarehouseOverview{
		Warehouse:        mapper.WarehouseToResponse(&warehouse),
		EggStocks:        warehouseItemEggResponses,
		CornStocks:       warehouseItemCornResponses,
		EquipmentStocks:  warehouseItemEquipmentResponses,
		TotalSafeStock:   uint64(totalSafeStock),
		TotalDangerStock: uint64(totalDangerStock),
		TotalItemInOrder: uint64(warehouseItemProcurementSentOffCount),
	}, nil
}

func (s *WarehouseService) GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) (dto.WarehouseItemHistoryListPaginationResponse, error) {
	s.repository.UseTx(false)

	warehouseItemHistories, err := s.repository.GetWarehouseItemHistories(filter)
	if err != nil {
		s.log.Error("failed to get warehouse item history", zap.Error(err))
		return dto.WarehouseItemHistoryListPaginationResponse{}, err
	}

	response := make([]dto.WarehouseItemHistoryListResponse, 0)
	for _, e := range warehouseItemHistories {
		response = append(response, mapper.WarehouseItemHistoryToListResponse(&e))
	}

	totalData, err := s.repository.CountTotalWarehouseItemHistory(filter)
	if err != nil {
		s.log.Error("failed to count warehouse item history", zap.Error(err))
		return dto.WarehouseItemHistoryListPaginationResponse{}, err
	}

	resp := dto.WarehouseItemHistoryListPaginationResponse{
		WarehouseItemHistories: response,
	}

	if filter.Page > 0 {
		resp.TotalData = uint64(totalData)
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (s *WarehouseService) GetWarehouseItemHistoryById(id uint64) (dto.WarehouseItemHistoryResponse, error) {
	s.repository.UseTx(false)

	warehouseItemHistory, err := s.repository.GetWarehouseItemHistoryById(id)
	if err != nil {
		s.log.Error("failed to get warehouse item history by id", zap.Error(err))
		return dto.WarehouseItemHistoryResponse{}, err
	}

	return mapper.WarehouseItemHistoryToResponse(&warehouseItemHistory), nil
}

func (s *WarehouseService) GetCornWarehouseItemSummary(warehouseId uint64) (dto.CornWarehouseItemSummaryResponse, error) {
	s.repository.UseTx(false)

	warehouse, err := s.repository.GetWarehouseById(warehouseId)
	if err != nil {
		s.log.Error("failed get warehouse by id", zap.Error(err))
		return dto.CornWarehouseItemSummaryResponse{}, err
	}

	withZeroQuantity := false
	warehouseItemCorns, err := s.repository.GetWarehouseItemCorns(dto.GetWarehouseItemCornFilter{
		WarehouseId:      warehouseId,
		FromNewest:       true,
		WithZeroQuantity: &withZeroQuantity,
	})
	if err != nil {
		s.log.Error("failed get warehouse item corns", zap.Error(err))
		return dto.CornWarehouseItemSummaryResponse{}, err
	}

	totalQuantity := float64(0)
	for _, e := range warehouseItemCorns {
		totalQuantity += e.Quantity
	}

	return dto.CornWarehouseItemSummaryResponse{
		Warehouse: mapper.WarehouseToResponse(&warehouse),
		Name:      constant.Corn,
		Quantity:  totalQuantity,
		Unit:      constant.UnitKg,
	}, nil
}

func (s *WarehouseService) GetEggWarehouseItemSummary(warehouseId uint64) ([]dto.EggWarehouseItemSummaryResponse, error) {
	s.repository.UseTx(false)

	response := make([]dto.EggWarehouseItemSummaryResponse, 0)
	warehouseItems, err := s.repository.GetWarehouseItems(dto.GetWarehouseItemFilter{
		WarehouseId: warehouseId,
		ItemNames:   []string{constant.GoodEgg, constant.CrackedEgg},
		Units:       []string{constant.UnitKg},
	})
	if err != nil {
		s.log.Error("failed to get warehouse items", zap.Error(err))
		return nil, err
	}

	for _, warehouseItem := range warehouseItems {
		switch warehouseItem.Item.Name {
		case constant.GoodEgg:
			response = append(response, dto.EggWarehouseItemSummaryResponse{
				Name:     constant.GoodEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.UnitKg,
			})

			response = append(response, dto.EggWarehouseItemSummaryResponse{
				Name:     constant.GoodEgg,
				Quantity: math.Floor(warehouseItem.Quantity / float64(constant.TotalEggPerIkat)),
				Unit:     constant.UnitIkat,
			})
		case constant.CrackedEgg:
			response = append(response, dto.EggWarehouseItemSummaryResponse{
				Name:     constant.CrackedEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.UnitKg,
			})

			response = append(response, dto.EggWarehouseItemSummaryResponse{
				Name:     constant.CrackedEgg,
				Quantity: math.Floor(warehouseItem.Quantity / float64(constant.TotalEggPerIkat)),
				Unit:     constant.UnitIkat,
			})
		}
	}

	return response, nil
}

func (s *WarehouseService) CreateWarehouseSale(request dto.CreateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(request.WarehouseId, request.ItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item by warehouse id and item id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	realQuantity := request.Quantity
	if saleUnit == enum.SaleUnitIkat {
		realQuantity *= float64(constant.TotalEggPerIkat)
	}

	if warehouseItem.Quantity < realQuantity {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("stock item is insuficcient")
	}

	warehouseItem.Quantity -= realQuantity
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update warehouse item", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse sent date", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid sent date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid price format")
	}

	totalPrice := price.Mul(decimal.NewFromFloat(request.Quantity))
	discountPrice := totalPrice.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	totalPrice = totalPrice.Sub(discountPrice)

	dateNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	warehouseSale := entity.WarehouseSale{
		WarehouseId:   request.WarehouseId,
		ItemId:        request.ItemId,
		Quantity:      request.Quantity,
		Price:         price,
		TotalPrice:    totalPrice,
		SendDate:      sendDate,
		Discount:      request.Discount,
		IsSend:        false,
		SaleUnit:      saleUnit,
		PaymentType:   paymentType,
		PaymentStatus: enum.PaymentStatusNotPaid,
		CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
	}

	if len(request.Payments) == 0 {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("warehouseSalePayment is required")
	}

	totalPayment := decimal.Zero
	for _, paymentReq := range request.Payments {
		nominal, err := decimal.NewFromString(paymentReq.Nominal)
		if err != nil {
			s.log.Error("failed to parse nominal", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
		}
		totalPayment = totalPayment.Add(nominal)
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !warehouseSale.TotalPrice.Equal(totalPayment) {
			s.log.Error("nominal is not equal to total price")
			return dto.WarehouseSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}
		warehouseSale.PaymentStatus = enum.PaymentStatusPaid
		warehouseSale.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else {
		if totalPayment.GreaterThan(warehouseSale.TotalPrice) {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("total payment is greater than total price")
		} else if totalPayment.Equal(warehouseSale.TotalPrice) {
			warehouseSale.PaymentStatus = enum.PaymentStatusPaid
		} else {
			warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
		}
	}

	payments := make([]entity.WarehouseSalePayment, 0, len(request.Payments))
	for _, paymentReq := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(paymentReq.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("invalid payment method", zap.String("paymentMethod", paymentReq.PaymentMethod))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment method")
		}
		paymentDate, err := time.Parse("02-01-2006", paymentReq.PaymentDate)
		if err != nil {
			s.log.Error("failed to parse payment date", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment date format")
		}
		nominal, err := decimal.NewFromString(paymentReq.Nominal)
		if err != nil {
			s.log.Error("failed to parse nominal", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
		}
		payments = append(payments, entity.WarehouseSalePayment{
			PaymentDate:     paymentDate,
			WarehouseSaleId: warehouseSale.Id,
			Nominal:         nominal,
			PaymentProof:    paymentReq.PaymentProof,
			PaymentMethod:   paymentMethod,
			CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if warehouseSale.PaymentStatus != enum.PaymentStatusPaid {
		warehouseSale.DeadlinePaymentDate = sql.NullTime{Time: dateNow.AddDate(0, 0, 7), Valid: true}
	}

	if request.CustomerType == constant.OldCustomerType {
		if request.CustomerId < 1 {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer id is required")
		}

		warehouseSale.CustomerId = request.CustomerId
	} else {
		customer := dto.CreateCustomerRequest{
			Name:        request.CustomerName,
			PhoneNumber: request.CustomerPhoneNumber,
		}

		if request.CustomerName == "" || request.CustomerPhoneNumber == "" {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer name and customer phone number is required")
		}

		if len(request.CustomerPhoneNumber) < 2 || request.CustomerPhoneNumber[:2] != "08" {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer phone number must be in valid format 08")
		}

		resp, err := s.customerService.CreateCustomer(customer, userId)
		if err != nil {
			return dto.WarehouseSaleResponse{}, err
		}

		warehouseSale.CustomerId = resp.Id
	}

	err = s.repository.CreateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to create warehouse sale", zap.Error(err))
		if err := s.customerService.DeleteCustomer(warehouseSale.CustomerId); err != nil {
			s.log.Error("failed to delete customer", zap.Error(err))
		}
		return dto.WarehouseSaleResponse{}, err
	}

	for i := range payments {
		payments[i].WarehouseSaleId = warehouseSale.Id
	}

	if len(payments) > 0 {
		err = s.repository.CreateWarehouseSalePaymentInBatch(&payments)
		if err != nil {
			s.log.Error("failed to create warehouse sale payment in batch", zap.Error(err))
			if err := s.customerService.DeleteCustomer(warehouseSale.CustomerId); err != nil {
				s.log.Error("failed to delete customer", zap.Error(err))
			}
			return dto.WarehouseSaleResponse{}, err
		}
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale, err = s.repository.GetWarehouseSaleById(warehouseSale.Id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	jsonWarehouseHistoryParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemName:       warehouseSale.Item.Name,
		ItemUnit:       warehouseSale.Item.Unit,
		Source:         warehouseSale.Warehouse.Name,
		Destination:    warehouseSale.Customer.Name,
		QuantityBefore: warehouseSale.Quantity,
		QuantityAfter:  warehouseSale.Quantity - request.Quantity,
		UserId:         userId,
		Status:         enum.ItemHistoryStatusOut,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("failed parsed struct into json")
	}

	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonWarehouseHistoryParsed)

	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))
	remainingPayment := warehouseSale.TotalPrice
	for i, warehouseSalePayment := range warehouseSale.Payments {
		warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&warehouseSalePayment)
		remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
		warehouseSalePayments[i].Remaining = remainingPayment.String()
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) GetWarehouseSaleById(id uint64) (dto.WarehouseSaleResponse, error) {
	warehouseSale, err := s.repository.GetWarehouseSaleById(id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))

	remainingPayment := warehouseSale.TotalPrice
	for i, warehouseSalePayment := range warehouseSale.Payments {
		warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&warehouseSalePayment)
		remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
		warehouseSalePayments[i].Remaining = remainingPayment.String()
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) GetWarehouseSales(filter dto.GetWarehouseSaleFilter) (dto.WarehouseSaleListPaginationResponse, error) {
	warehouseSales, err := s.repository.GetWarehouseSales(filter)
	if err != nil {
		s.log.Error("failed to get warehouse sales", zap.Error(err))
		return dto.WarehouseSaleListPaginationResponse{}, err
	}

	warehouseSaleResponses := make([]dto.WarehouseSaleListResponse, len(warehouseSales))
	for i, warehouseSale := range warehouseSales {
		warehouseSaleResponses[i] = mapper.WarehouseSaleToListResponse(&warehouseSale)
	}

	totalData, err := s.repository.CountTotalWarehouseSale(
		dto.GetWarehouseSaleFilter{
			Date:          filter.Date,
			PaymentStatus: filter.PaymentStatus,
			WarehouseId:   filter.WarehouseId,
		},
	)
	if err != nil {
		s.log.Error("failed to get warehouse sales", zap.Error(err))
		return dto.WarehouseSaleListPaginationResponse{}, err
	}

	resp := dto.WarehouseSaleListPaginationResponse{
		WarehouseSales: warehouseSaleResponses,
	}

	if filter.Page > 0 {
		resp.TotalData = totalData
		resp.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return resp, nil
}

func (s *WarehouseService) CreateWarehouseSalePayment(warehouseSaleId uint64, request dto.CreateWarehouseSalePaymentRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	warehouseSalePayment := entity.WarehouseSalePayment{
		WarehouseSaleId: warehouseSaleId,
		PaymentDate:     paymentDate,
		PaymentMethod:   paymentMethod,
		Nominal:         nominal,
		PaymentProof:    request.PaymentProof,
		CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
	}

	warehouseSale, err := s.repository.GetWarehouseSaleById(warehouseSaleId)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	if warehouseSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("warehouse sale is already paid", zap.Uint64("id", warehouseSaleId))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("warehouse sale is already paid")
	}

	totalPayment := nominal
	for _, payment := range warehouseSale.Payments {
		totalPayment = totalPayment.Add(payment.Nominal)
	}

	if totalPayment.Equal(warehouseSale.TotalPrice) {
		warehouseSale.PaymentStatus = enum.PaymentStatusPaid
		warehouseSale.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.GreaterThan(warehouseSale.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.CreateWarehouseSalePayment(&warehouseSalePayment)
	if err != nil {
		s.log.Error("failed to create warehouse sale payment", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	err = s.repository.UpdateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to update warehouse sale", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale.Payments = append(warehouseSale.Payments, warehouseSalePayment)
	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))

	remainingPayment := warehouseSale.TotalPrice
	for i, warehouseSalePayment := range warehouseSale.Payments {
		warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&warehouseSalePayment)
		remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
		warehouseSalePayments[i].Remaining = remainingPayment.String()
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) UpdateWarehouseSale(id uint64, request dto.UpdateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error) {
	warehouseSale, err := s.repository.GetWarehouseSaleById(id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseSale.WarehouseId, warehouseSale.ItemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	realQuantity := request.Quantity
	tempStoreSaleQuantity := warehouseSale.Quantity
	if saleUnit == enum.SaleUnitIkat {
		realQuantity *= float64(constant.TotalEggPerIkat)
		tempStoreSaleQuantity *= float64(constant.TotalEggPerIkat)
	}

	if warehouseItem.Quantity+tempStoreSaleQuantity < realQuantity {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("stock item is insuficcient")
	}

	warehouseItem.Quantity += warehouseSale.Quantity - realQuantity
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price from string", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale.Quantity = request.Quantity
	totalPrice := price.Mul(decimal.NewFromFloat(request.Quantity))
	discountPrice := totalPrice.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	warehouseSale.TotalPrice = totalPrice.Sub(discountPrice)
	warehouseSale.Price = price
	warehouseSale.Discount = request.Discount

	totalPayment := decimal.Zero
	for _, payment := range warehouseSale.Payments {
		totalPayment = totalPayment.Add(payment.Nominal)
	}

	if totalPayment.Equal(warehouseSale.TotalPrice) {
		warehouseSale.PaymentStatus = enum.PaymentStatusPaid
		warehouseSale.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.LessThan(warehouseSale.TotalPrice) {
		warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseSale.PaidDate = sql.NullTime{Valid: false}
	} else if totalPayment.GreaterThan(warehouseSale.TotalPrice) {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("quantity can't be lower")
	}

	warehouseSale.SendDate, err = time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse send date", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid send date format")
	}

	warehouseSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to update warehouse sale", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale, err = s.repository.GetWarehouseSaleById(warehouseSale.Id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))

	remainingPayment := warehouseSale.TotalPrice
	for i, warehouseSalePayment := range warehouseSale.Payments {
		warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&warehouseSalePayment)
		remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
		warehouseSalePayments[i].Remaining = remainingPayment.String()
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) UpdateWarehouseSalePayment(warehouseSaleId uint64, id uint64, request dto.UpdateWarehouseSalePaymentRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseSalePayment, err := s.repository.GetWarehouseSalePaymentById(id)
	if err != nil {
		s.log.Error("failed to get warehouse sale payment by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale, err := s.repository.GetWarehouseSaleById(warehouseSaleId)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed to parse payment date", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	totalPayment := nominal
	for _, payment := range warehouseSale.Payments {
		if payment.Id != warehouseSalePayment.Id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.Equal(warehouseSale.TotalPrice) {
		warehouseSale.PaymentStatus = enum.PaymentStatusPaid
		warehouseSale.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.GreaterThan(warehouseSale.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	} else if totalPayment.LessThan(warehouseSale.TotalPrice) {
		warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseSale.PaidDate = sql.NullTime{Valid: false}
	}

	warehouseSalePayment.Nominal = nominal
	warehouseSalePayment.PaymentProof = request.PaymentProof
	warehouseSalePayment.PaymentDate = paymentDate
	warehouseSalePayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to update warehouse sale", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	err = s.repository.UpdateWarehouseSalePayment(&warehouseSalePayment)
	if err != nil {
		s.log.Error("failed to update warehouse sale payment", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))

	remainingPayment := warehouseSale.TotalPrice
	for i, payment := range warehouseSale.Payments {
		if payment.Id == id {
			warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&payment)
			remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
			warehouseSalePayments[i].Remaining = remainingPayment.String()
		} else {
			warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&payment)
			remainingPayment = remainingPayment.Sub(payment.Nominal)
			warehouseSalePayments[i].Remaining = remainingPayment.String()
		}
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) SendWarehouseSale(id uint64, userId uuid.UUID) (dto.WarehouseSaleResponse, error) {
	warehouseSale, err := s.repository.GetWarehouseSaleById(id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	if warehouseSale.IsSend {
		s.log.Error("warehouse sale is already sent", zap.Uint64("id", id))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale.IsSend = true
	warehouseSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to update wareehouse sale", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale, err = s.repository.GetWarehouseSaleById(warehouseSale.Id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))

	remainingPayment := warehouseSale.TotalPrice
	for i, warehouseSalePayment := range warehouseSale.Payments {
		warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&warehouseSalePayment)
		remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
		warehouseSalePayments[i].Remaining = remainingPayment.String()
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) DeleteWarehouseSale(id uint64, userId uuid.UUID) error {
	warehouseSale, err := s.repository.GetWarehouseSaleById(id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return err
	}

	if warehouseSale.IsSend {
		s.log.Error("warehouse sale is already sent", zap.Uint64("id", id))
		return errx.BadRequest("warehouse sale already send")
	}

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseSale.WarehouseId, warehouseSale.ItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item by store id and item id", zap.Error(err))
		return err
	}

	realQuantity := warehouseSale.Quantity
	if warehouseSale.SaleUnit == enum.SaleUnitIkat {
		realQuantity *= float64(constant.TotalEggPerIkat)
	}

	warehouseItem.Quantity += realQuantity
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return err
	}

	err = s.repository.DeleteWarehouseSale(id)
	if err != nil {
		s.log.Error("failed to delete warehouse sale", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) DeleteWarehouseSalePayment(warehouseSaleId uint64, id uint64, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseSale, err := s.repository.GetWarehouseSaleById(warehouseSaleId)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return err
	}

	totalPayment := decimal.Zero
	for _, payment := range warehouseSale.Payments {
		if payment.Id != id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.LessThan(warehouseSale.TotalPrice) && totalPayment.GreaterThan(decimal.Zero) {
		warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseSale.PaidDate = sql.NullTime{Valid: false}
		warehouseSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	} else if totalPayment.LessThan(decimal.Zero) {
		s.log.Error("delete this payment make minus", zap.Error(err))
		return errx.BadRequest("payment minus")
	}

	err = s.repository.UpdateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to update warehouse sale", zap.Error(err))
		return err
	}

	err = s.repository.DeleteWarehouseSalePayment(id)
	if err != nil {
		s.log.Error("failed to update warehouse sale", zap.Error(err))
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseSaleQueue(request dto.CreateWarehouseSaleQueueRequest, userId uuid.UUID) (dto.WarehouseSaleQueueResponse, error) {
	s.repository.UseTx(false)

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		return dto.WarehouseSaleQueueResponse{}, errx.BadRequest("invalid sale unit")
	}

	customerType := enum.ValueOfCustomerType(request.CustomerType)
	if !customerType.IsValid() {
		return dto.WarehouseSaleQueueResponse{}, errx.BadRequest("invalid customer type")
	}

	data := entity.WarehouseSaleQueue{
		ItemId:       request.ItemId,
		WarehouseId:  request.WarehouseId,
		SaleUnit:     saleUnit,
		CustomerType: customerType,
		Quantity:     request.Quantity,
		CreatedBy:    uuid.NullUUID{UUID: userId, Valid: true},
	}

	if customerType == enum.CustomerTypeNew {
		if request.CustomerName == "" || request.CustomerPhoneNumber == "" {
			return dto.WarehouseSaleQueueResponse{}, errx.BadRequest("customer name and phone number is required")
		}

		data.CustomerName = sql.NullString{String: request.CustomerName, Valid: true}
		data.CustomerPhoneNumber = sql.NullString{String: request.CustomerPhoneNumber, Valid: true}

	} else {
		if request.CustomerId < 1 {
			return dto.WarehouseSaleQueueResponse{}, errx.BadRequest("customer id is required")
		}

		data.CustomerId = sql.NullInt64{Int64: int64(request.CustomerId), Valid: true}
	}

	err := s.repository.CreateWarehouseSaleQueue(&data)
	if err != nil {
		s.log.Error("failed create Warehouse sale queue", zap.Error(err))
		return dto.WarehouseSaleQueueResponse{}, err
	}

	data, err = s.repository.GetWarehouseSaleQueueById(data.Id)
	if err != nil {
		return dto.WarehouseSaleQueueResponse{}, err
	}

	return mapper.WarehouseSaleQueueToResponse(&data), nil
}

func (s *WarehouseService) GetWarehouseSaleQueue(id uint64) (dto.WarehouseSaleQueueResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseSaleQueueById(id)
	if err != nil {
		s.log.Error("failed get Warehouse sale queue by id", zap.Error(err))
		return dto.WarehouseSaleQueueResponse{}, err
	}

	return mapper.WarehouseSaleQueueToResponse(&data), nil
}

func (s *WarehouseService) GetWarehouseSaleQueues(filter dto.GetWarehouseSaleQueueFilter) ([]dto.WarehouseSaleQueueResponse, error) {
	s.repository.UseTx(false)

	warehouseSaleQueues, err := s.repository.GetWarehouseSaleQueues(filter)
	if err != nil {
		s.log.Error("failed get Warehouse sale queues", zap.Error(err))
		return nil, err
	}

	warehouseIds := make([]uint64, 0)
	warehouseQueueMap := make(map[uint64]map[uint64][]entity.WarehouseSaleQueue)
	warehouseItemMap := make(map[uint64]map[uint64]entity.WarehouseItem)

	for _, warehouseSaleQueue := range warehouseSaleQueues {
		if _, ok := warehouseQueueMap[warehouseSaleQueue.WarehouseId]; !ok {
			warehouseQueueMap[warehouseSaleQueue.WarehouseId] = make(map[uint64][]entity.WarehouseSaleQueue)
		}
		warehouseQueueMap[warehouseSaleQueue.WarehouseId][warehouseSaleQueue.ItemId] =
			append(warehouseQueueMap[warehouseSaleQueue.WarehouseId][warehouseSaleQueue.ItemId], warehouseSaleQueue)

		warehouseIds = append(warehouseIds, warehouseSaleQueue.WarehouseId)
	}

	warehouseItems, err := s.repository.GetWarehouseItems(dto.GetWarehouseItemFilter{
		WarehouseIds: warehouseIds,
		Category:     param.ItemCategoryParam(enum.ItemCategoryEgg),
	})
	if err != nil {
		s.log.Error("failed get warehouse items", zap.Error(err))
		return nil, err
	}
	for _, warehouseItem := range warehouseItems {
		if _, ok := warehouseItemMap[warehouseItem.WarehouseId]; !ok {
			warehouseItemMap[warehouseItem.WarehouseId] = make(map[uint64]entity.WarehouseItem)
		}
		warehouseItemMap[warehouseItem.WarehouseId][warehouseItem.ItemId] = warehouseItem
	}

	weightPerWarehouseSaleQueueMap := make(map[uint64]float64)
	startAllocationWarehouseSaleQueueMap := make(map[uint64]float64)
	startBacklogQueueMap := make(map[uint64]float64)
	additionalAllocationWarehouseSaleQueueMap := make(map[uint64]float64)

	for storeId, warehouseSaleQueueItemMap := range warehouseQueueMap {
		for warehouseSaleQueueItemId, warehouseSaleQueues := range warehouseSaleQueueItemMap {
			storeItem := warehouseItemMap[storeId][warehouseSaleQueueItemId]

			totalDemand := 0.0
			for _, q := range warehouseSaleQueues {
				if q.SaleUnit == enum.SaleUnitIkat {
					totalDemand += q.Quantity * float64(constant.TotalEggPerIkat)
				} else {
					totalDemand += q.Quantity
				}
			}

			totalWeight := 0.0
			for _, q := range warehouseSaleQueues {
				demandRatio := 0.0
				if totalDemand > 0 {
					demandRatio = q.Quantity / totalDemand
				}

				weight := 0.0
				switch q.CustomerType {
				case enum.CustomerTypeNew:
					weight += constant.CustomerTypeNewWeight * constant.CustomerIndex
				case enum.CustomerTypeOld:
					weight += constant.CustomerTypeOldWeight * constant.CustomerIndex
				}
				weight += constant.DemandIndex * demandRatio

				totalWeight += weight
				weightPerWarehouseSaleQueueMap[q.Id] = weight
			}

			totalStartAllocation := 0.0
			for _, q := range warehouseSaleQueues {
				currDemand := 0.0
				if q.SaleUnit == enum.SaleUnitIkat {
					currDemand += q.Quantity * float64(constant.TotalEggPerIkat)
				} else {
					currDemand += q.Quantity
				}

				allocation := 0.0
				denom := totalWeight * storeItem.Quantity
				if denom > 0 {
					allocation = weightPerWarehouseSaleQueueMap[q.Id] / denom
				}

				startAllocation := math.Min(currDemand, allocation)
				totalStartAllocation += startAllocation
				startAllocationWarehouseSaleQueueMap[q.Id] = startAllocation
			}

			totalStartBacklog := 0.0
			for _, q := range warehouseSaleQueues {
				currDemand := 0.0
				if q.SaleUnit == enum.SaleUnitIkat {
					currDemand += q.Quantity * float64(constant.TotalEggPerIkat)
				} else {
					currDemand += q.Quantity
				}

				backlog := math.Max(0, currDemand-startAllocationWarehouseSaleQueueMap[q.Id])
				totalStartBacklog += backlog
				startBacklogQueueMap[q.Id] = backlog
			}

			remainingQuantity := storeItem.Quantity - totalStartAllocation
			for _, q := range warehouseSaleQueues {
				additional := 0.0
				if totalStartBacklog > 0 {
					ratio := startAllocationWarehouseSaleQueueMap[q.Id] / totalStartBacklog
					additional = math.Min(startAllocationWarehouseSaleQueueMap[q.Id], remainingQuantity*ratio)
				}
				additionalAllocationWarehouseSaleQueueMap[q.Id] = additional
			}
		}
	}

	responses := make([]dto.WarehouseSaleQueueResponse, 0, len(warehouseSaleQueues))
	for _, q := range warehouseSaleQueues {
		resp := mapper.WarehouseSaleQueueToResponse(&q)
		resp.TotalAllocation = startAllocationWarehouseSaleQueueMap[q.Id] + additionalAllocationWarehouseSaleQueueMap[q.Id]
		responses = append(responses, resp)
	}

	sort.Slice(responses, func(i, j int) bool {
		return responses[i].TotalAllocation > responses[j].TotalAllocation
	})

	for i := range responses {
		responses[i].OrderPriority = uint64(i + 1)
	}

	return responses, nil
}

func (s *WarehouseService) DeleteWarehouseSaleQueue(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouseSaleQueue(id)
	if err != nil {
		s.log.Error("failed delete Warehouse sale queue", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) AllocateWarehouseSaleQueue(id uint64, request dto.CreateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	err := s.repository.DeleteWarehouseSaleQueue(id)
	if err != nil {
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(request.WarehouseId, request.ItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item by warehouse id and item id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	realQuantity := request.Quantity
	if saleUnit == enum.SaleUnitIkat {
		realQuantity *= float64(constant.TotalEggPerIkat)
	}

	if warehouseItem.Quantity < realQuantity {
		return dto.WarehouseSaleResponse{}, errx.BadRequest("stock item is insuficcient")
	}

	warehouseItem.Quantity -= realQuantity
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update warehouse item", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse sent date", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid sent date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid price format")
	}

	totalPrice := price.Mul(decimal.NewFromFloat(request.Quantity))
	discountPrice := totalPrice.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	totalPrice = totalPrice.Sub(discountPrice)

	dateNow := time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)
	warehouseSale := entity.WarehouseSale{
		WarehouseId:   request.WarehouseId,
		ItemId:        request.ItemId,
		Quantity:      request.Quantity,
		Price:         price,
		TotalPrice:    totalPrice,
		SendDate:      sendDate,
		Discount:      request.Discount,
		IsSend:        false,
		SaleUnit:      saleUnit,
		PaymentType:   paymentType,
		PaymentStatus: enum.PaymentStatusNotPaid,
		CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
	}

	totalPayment := decimal.Zero
	for _, paymentReq := range request.Payments {
		nominal, err := decimal.NewFromString(paymentReq.Nominal)
		if err != nil {
			s.log.Error("failed to parse nominal", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
		}
		totalPayment = totalPayment.Add(nominal)
	}

	payments := make([]entity.WarehouseSalePayment, 0, len(request.Payments))
	for _, paymentReq := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(paymentReq.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("invalid payment method", zap.String("paymentMethod", paymentReq.PaymentMethod))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment method")
		}
		paymentDate, err := time.Parse("02-01-2006", paymentReq.PaymentDate)
		if err != nil {
			s.log.Error("failed to parse payment date", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment date format")
		}
		nominal, err := decimal.NewFromString(paymentReq.Nominal)
		if err != nil {
			s.log.Error("failed to parse nominal", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
		}
		payments = append(payments, entity.WarehouseSalePayment{
			PaymentDate:     paymentDate,
			WarehouseSaleId: warehouseSale.Id,
			Nominal:         nominal,
			PaymentProof:    paymentReq.PaymentProof,
			PaymentMethod:   paymentMethod,
			CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !warehouseSale.TotalPrice.Equal(totalPayment) {
			s.log.Error("nominal is not equal to total price")
			return dto.WarehouseSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}
		warehouseSale.PaymentStatus = enum.PaymentStatusPaid
		warehouseSale.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else {
		if totalPayment.GreaterThan(warehouseSale.TotalPrice) {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("total payment is greater than total price")
		} else if totalPayment.Equal(warehouseSale.TotalPrice) {
			warehouseSale.PaymentStatus = enum.PaymentStatusPaid
			warehouseSale.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
		} else {
			warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
		}
	}

	if warehouseSale.PaymentStatus != enum.PaymentStatusPaid {
		warehouseSale.DeadlinePaymentDate = sql.NullTime{Time: dateNow.AddDate(0, 0, 7), Valid: true}
	}

	if request.CustomerType == constant.OldCustomerType {
		if request.CustomerId < 1 {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer id is required")
		}

		warehouseSale.CustomerId = request.CustomerId
	} else {
		customer := dto.CreateCustomerRequest{
			Name:        request.CustomerName,
			PhoneNumber: request.CustomerPhoneNumber,
		}

		if request.CustomerName == "" || request.CustomerPhoneNumber == "" {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer name and customer phone number is required")
		}

		if len(request.CustomerPhoneNumber) < 2 || request.CustomerPhoneNumber[:2] != "08" {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer phone number must be in valid format 08")
		}

		resp, err := s.customerService.CreateCustomer(customer, userId)
		if err != nil {
			return dto.WarehouseSaleResponse{}, err
		}

		warehouseSale.CustomerId = resp.Id
	}

	err = s.repository.CreateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to create warehouse sale", zap.Error(err))
		if err := s.customerService.DeleteCustomer(warehouseSale.CustomerId); err != nil {
			s.log.Error("failed to delete customer", zap.Error(err))
		}
		return dto.WarehouseSaleResponse{}, err
	}

	for i := range payments {
		payments[i].WarehouseSaleId = warehouseSale.Id
	}

	if len(payments) > 0 {
		err = s.repository.CreateWarehouseSalePaymentInBatch(&payments)
		if err != nil {
			s.log.Error("failed to create warehouse sale payment in batch", zap.Error(err))
			if err := s.customerService.DeleteCustomer(warehouseSale.CustomerId); err != nil {
				s.log.Error("failed to delete customer", zap.Error(err))
			}
			return dto.WarehouseSaleResponse{}, err
		}
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale, err = s.repository.GetWarehouseSaleById(warehouseSale.Id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	jsonWarehouseHistoryParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemName:       warehouseSale.Item.Name,
		ItemUnit:       warehouseSale.Item.Unit,
		Source:         warehouseSale.Warehouse.Name,
		Destination:    warehouseSale.Customer.Name,
		QuantityBefore: warehouseSale.Quantity,
		QuantityAfter:  warehouseSale.Quantity - request.Quantity,
		UserId:         userId,
		Status:         enum.ItemHistoryStatusOut,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("failed parsed struct into json")
	}

	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonWarehouseHistoryParsed)

	warehouseSalePayments := make([]dto.WarehouseSalePaymentResponse, len(warehouseSale.Payments))
	remainingPayment := warehouseSale.TotalPrice
	for i, warehouseSalePayment := range warehouseSale.Payments {
		warehouseSalePayments[i] = mapper.WarehouseSalePaymentToResponse(&warehouseSalePayment)
		remainingPayment = remainingPayment.Sub(warehouseSalePayment.Nominal)
		warehouseSalePayments[i].Remaining = remainingPayment.String()
	}

	warehouseSaleResponse := mapper.WarehouseSaleToResponse(&warehouseSale)
	warehouseSaleResponse.Payments = warehouseSalePayments
	warehouseSaleResponse.RemainingPayment = remainingPayment.String()

	return warehouseSaleResponse, nil
}

func (s *WarehouseService) CreateWarehouseItemProcurementDraft(request dto.CreateWarehouseItemProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	data := entity.WarehouseItemProcurementDraft{
		WarehouseId:   request.WarehouseId,
		ItemId:        request.ItemId,
		SupplierId:    sql.NullInt64{Int64: int64(request.SupplierId), Valid: true},
		DailySpending: request.DailySpending,
		DaysNeed:      request.DaysNeed,
		Price:         price,
		CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateWarehouseItemProcurementDraft(&data)
	if err != nil {
		s.log.Error("failed create warehouse item procurement draft", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemProcurementDraft(data.Id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement draft", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	return mapper.WarehouseItemProcurementDraftToResponse(&data), nil
}

func (s *WarehouseService) GetWarehouseItemProcurementDrafts(filter dto.GetWarehouseItemProcurementDraftFilter) ([]dto.WarehouseItemProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemProcurementDrafts(filter)
	if err != nil {
		s.log.Error("failed get warehouse item procurement drafts", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehouseItemProcurementDraftResponse, 0)
	for _, e := range data {
		response = append(response, mapper.WarehouseItemProcurementDraftToResponse(&e))
	}

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemProcurementDraft(id uint64) (dto.WarehouseItemProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemProcurementDraft(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement draft", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	return mapper.WarehouseItemProcurementDraftToResponse(&data), nil
}

func (s *WarehouseService) UpdateWarehouseItemProcurementDraft(id uint64, request dto.UpdateWarehouseItemProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	warehouseItemProcurementDraft, err := s.repository.GetWarehouseItemProcurementDraft(id)
	if err != nil {
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	warehouseItemProcurementDraft.WarehouseId = request.WarehouseId
	warehouseItemProcurementDraft.ItemId = request.ItemId
	warehouseItemProcurementDraft.DailySpending = request.DailySpending
	warehouseItemProcurementDraft.DaysNeed = request.DaysNeed
	warehouseItemProcurementDraft.Price = price
	warehouseItemProcurementDraft.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if request.SupplierId != nil {
		warehouseItemProcurementDraft.SupplierId = sql.NullInt64{Int64: int64(*request.SupplierId), Valid: true}
	} else {
		warehouseItemProcurementDraft.SupplierId = sql.NullInt64{Valid: false}
	}

	err = s.repository.UpdateWarehouseItemProcurementDraft(&warehouseItemProcurementDraft)
	if err != nil {
		s.log.Error("failed udpate warehouse item procurement draft", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	warehouseItemProcurementDraft, err = s.repository.GetWarehouseItemProcurementDraft(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement draft", zap.Error(err))
		return dto.WarehouseItemProcurementDraftResponse{}, err
	}

	return mapper.WarehouseItemProcurementDraftToResponse(&warehouseItemProcurementDraft), nil
}

func (s *WarehouseService) ConfirmationWarehouseItemProcurementDraft(id uint64, request dto.CreateWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	err := s.repository.DeleteWarehouseItemProcurementDraft(id)
	if err != nil {
		s.log.Error("failed delete warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	estimationArrivalDate, err := time.Parse("02-01-2006", request.EstimationArrivalDate)
	if err != nil {
		s.log.Error("failed parse time", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid estimation arrival date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type")
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment type")
	}

	data := entity.WarehouseItemProcurement{
		WarehouseId:           request.WarehouseId,
		SupplierId:            request.SupplierId,
		ItemId:                request.ItemId,
		DailySpending:         request.DailySpending,
		DaysNeed:              request.DaysNeed,
		Price:                 price,
		TotalPrice:            price.Mul(decimal.NewFromFloat(request.DailySpending * float64(request.DaysNeed))),
		Quantity:              request.DailySpending * float64(request.DaysNeed),
		EstimationArrivalDate: estimationArrivalDate,
		Status:                enum.ProcurementStatusSentOff,
		PaymentStatus:         enum.PaymentStatusNotPaid,
		PaymentType:           paymentType,
		CreatedBy:             uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.DeadlinePaymentDate != nil {
		deadlinePaymentDate, err := time.Parse("02-01-2006", *request.DeadlinePaymentDate)
		if err != nil {
			s.log.Error("failed parse deadline payment date", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid deadline payment date format")
		}

		data.DeadlinePaymentDate = sql.NullTime{Time: deadlinePaymentDate, Valid: true}

	}

	if request.ExpiredAt != nil {
		expiredAt, err := time.Parse("02-01-2006", *request.ExpiredAt)
		if err != nil {
			s.log.Error("failed parse expired at", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid expired at format")
		}

		data.ExpiredAt = sql.NullTime{Time: expiredAt, Valid: true}
	}

	if len(request.Payments) == 0 {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("payments are required")
	}

	payments := make([]entity.WarehouseItemProcurementPayment, 0, len(request.Payments))
	totalPayment := decimal.Zero
	for _, p := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(p.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment method")
		}
		nominal, err := decimal.NewFromString(p.Nominal)
		if err != nil {
			s.log.Error("failed parse payment nominal", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, err
		}
		paymentDate, err := time.Parse("02-01-2006", p.PaymentDate)
		if err != nil {
			s.log.Error("failed parse payment date", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment date format")
		}
		totalPayment = totalPayment.Add(nominal)
		payments = append(payments, entity.WarehouseItemProcurementPayment{
			PaymentDate:                paymentDate,
			Nominal:                    nominal,
			PaymentProof:               p.PaymentProof,
			PaymentMethod:              paymentMethod,
			WarehouseItemProcurementId: data.Id,
			CreatedBy:                  uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff && !totalPayment.Equal(data.TotalPrice) {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("need more payment for paid off")
	}

	if totalPayment.Equal(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusPaid
		data.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.LessThan(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("total payment more than total price")
	}

	err = s.repository.CreateWarehouseItemProcurement(&data)
	if err != nil {
		s.log.Error("failed create warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	for i := range payments {
		payments[i].WarehouseItemProcurementId = data.Id
	}

	err = s.repository.CreateWarehouseItemProcurementPaymentInBatch(&payments)
	if err != nil {
		s.log.Error("failed create warehouse procurement payments in batch", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemProcurement(data.Id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	paymentResponses := make([]dto.WarehouseItemProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		paymentResponses = append(paymentResponses, payment)
	}

	response := mapper.WarehouseItemProcurementToResponse(&data)
	response.Payments = paymentResponses
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) DeleteWarehouseItemProcurementDraft(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouseItemProcurementDraft(id)
	if err != nil {
		s.log.Error("failed delete warehouse item procurement draft", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseItemProcurement(request dto.CreateWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	estimationArrivalDate, err := time.Parse("02-06-2006", request.EstimationArrivalDate)
	if err != nil {
		s.log.Error("failed parse time", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid estimation arrival date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type")
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment type")
	}

	data := entity.WarehouseItemProcurement{
		WarehouseId:           request.WarehouseId,
		SupplierId:            request.SupplierId,
		ItemId:                request.ItemId,
		DailySpending:         request.DailySpending,
		DaysNeed:              request.DaysNeed,
		Price:                 price,
		TotalPrice:            price.Mul(decimal.NewFromFloat(request.DailySpending * float64(request.DaysNeed))),
		Quantity:              request.DailySpending * float64(request.DaysNeed),
		EstimationArrivalDate: estimationArrivalDate,
		Status:                enum.ProcurementStatusSentOff,
		PaymentStatus:         enum.PaymentStatusNotPaid,
		PaymentType:           paymentType,
		CreatedBy:             uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.DeadlinePaymentDate != nil {
		deadlinePaymentDate, err := time.Parse("02-01-2006", *request.DeadlinePaymentDate)
		if err != nil {
			s.log.Error("failed parse deadline payment date", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid deadline payment date format")
		}

		data.DeadlinePaymentDate = sql.NullTime{Time: deadlinePaymentDate, Valid: true}

	}

	if request.ExpiredAt != nil {
		expiredAt, err := time.Parse("02-01-2006", *request.ExpiredAt)
		if err != nil {
			s.log.Error("failed parse expired at", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid expired at format")
		}

		data.ExpiredAt = sql.NullTime{Time: expiredAt, Valid: true}
	}

	if len(request.Payments) == 0 {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("payments are required")
	}

	payments := make([]entity.WarehouseItemProcurementPayment, 0, len(request.Payments))
	totalPayment := decimal.Zero
	for _, p := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(p.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment method")
		}
		nominal, err := decimal.NewFromString(p.Nominal)
		if err != nil {
			s.log.Error("failed parse payment nominal", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, err
		}
		paymentDate, err := time.Parse("02-01-2006", p.PaymentDate)
		if err != nil {
			s.log.Error("failed parse payment date", zap.Error(err))
			return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment date")
		}
		totalPayment = totalPayment.Add(nominal)
		payments = append(payments, entity.WarehouseItemProcurementPayment{
			PaymentDate:                paymentDate,
			Nominal:                    nominal,
			PaymentProof:               p.PaymentProof,
			PaymentMethod:              paymentMethod,
			WarehouseItemProcurementId: data.Id,
			CreatedBy:                  uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff && !totalPayment.Equal(data.TotalPrice) {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("need more payment for paid off")
	}

	if totalPayment.Equal(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusPaid
		data.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.LessThan(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("total payment more than total price")
	}

	err = s.repository.CreateWarehouseItemProcurement(&data)
	if err != nil {
		s.log.Error("failed create warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	for i := range payments {
		payments[i].WarehouseItemProcurementId = data.Id
	}

	err = s.repository.CreateWarehouseItemProcurementPaymentInBatch(&payments)
	if err != nil {
		s.log.Error("failed create warehouse procurement payments in batch", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemProcurement(data.Id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	paymentResponses := make([]dto.WarehouseItemProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		paymentResponses = append(paymentResponses, payment)
	}

	response := mapper.WarehouseItemProcurementToResponse(&data)
	response.Payments = paymentResponses
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemProcurements(filter dto.GetWarehouseItemProcurementFilter) (dto.WarehouseItemProcurementListPaginationResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemProcurements(filter)
	if err != nil {
		s.log.Error("failed get warehouse item procurements", zap.Error(err))
		return dto.WarehouseItemProcurementListPaginationResponse{}, err
	}

	totalData, err := s.repository.CountWarehouseItemProcurement(filter)
	if err != nil {
		s.log.Error("failed count warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementListPaginationResponse{}, err
	}

	warehouseItemProcurementResponses := make([]dto.WarehouseItemProcurementListResponse, 0)
	for _, e := range data {
		warehouseItemProcurementResponses = append(warehouseItemProcurementResponses, mapper.WarehouseItemProcurementToListResponse(&e))
	}

	response := dto.WarehouseItemProcurementListPaginationResponse{
		WarehouseItemProcurementes: warehouseItemProcurementResponses,
	}

	if filter.Page > 0 {
		response.TotalData = uint64(totalData)
		response.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemProcurement(id uint64) (dto.WarehouseItemProcurementResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemProcurement(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	response := mapper.WarehouseItemProcurementToResponse(&data)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) CreateWarehouseItemProcurementPayment(warehouseItemProcurementId uint64, request dto.CreateWarehouseItemProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemProcurement, err := s.repository.GetWarehouseItemProcurement(warehouseItemProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed parse payment date", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed parse nominal", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid nominal")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment method")
	}

	payment := entity.WarehouseItemProcurementPayment{
		WarehouseItemProcurementId: warehouseItemProcurement.Id,
		PaymentDate:                paymentDate,
		PaymentProof:               request.PaymentProof,
		Nominal:                    nominal,
		PaymentMethod:              paymentMethod,
		CreatedBy:                  uuid.NullUUID{UUID: userId, Valid: true},
	}

	totalPrice := nominal
	for _, e := range warehouseItemProcurement.Payments {
		totalPrice = totalPrice.Add(e.Nominal)
	}

	warehouseItemProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if totalPrice.Equal(warehouseItemProcurement.TotalPrice) {
		warehouseItemProcurement.PaymentStatus = enum.PaymentStatusPaid
		warehouseItemProcurement.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPrice.LessThan(warehouseItemProcurement.TotalPrice) {
		warehouseItemProcurement.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("nominal is to high")
	}

	err = s.repository.CreateWarehouseItemProcurementPayment(&payment)
	if err != nil {
		s.log.Error("failed create warehouse procurement payment", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.UpdateWarehouseItemProcurement(&warehouseItemProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemProcurement(warehouseItemProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	response := mapper.WarehouseItemProcurementToResponse(&data)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) UpdateWarehouseItemProcurementPayment(id uint64, warehouseItemProcurementId uint64, request dto.UpdateWarehouseItemProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemProcurement, err := s.repository.GetWarehouseItemProcurement(warehouseItemProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	warehouseItemProcurementPayment, err := s.repository.GetWarehouseItemProcurementPayment(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement payment", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed parse payment date", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment date")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed parse nominal", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid nominal")
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("invalid payment method")
	}

	warehouseItemProcurementPayment.Nominal = nominal
	warehouseItemProcurementPayment.PaymentDate = paymentDate
	warehouseItemProcurementPayment.PaymentMethod = paymentMethod
	warehouseItemProcurementPayment.PaymentProof = request.PaymentProof
	warehouseItemProcurementPayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	totalPrice := nominal
	for _, e := range warehouseItemProcurement.Payments {
		if e.Id != id {
			totalPrice = totalPrice.Add(e.Nominal)
		}
	}

	warehouseItemProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if totalPrice.Equal(warehouseItemProcurement.TotalPrice) {
		warehouseItemProcurement.PaymentStatus = enum.PaymentStatusPaid
		warehouseItemProcurement.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPrice.LessThan(warehouseItemProcurement.TotalPrice) {
		warehouseItemProcurement.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseItemProcurement.PaidDate = sql.NullTime{Valid: false}
	} else {
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("nominal is to high")
	}

	err = s.repository.UpdateWarehouseItemProcurementPayment(&warehouseItemProcurementPayment)
	if err != nil {
		s.log.Error("failed update warehouse item procurement payment", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.UpdateWarehouseItemProcurement(&warehouseItemProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemProcurement(warehouseItemProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	response := mapper.WarehouseItemProcurementToResponse(&data)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) DeleteWarehouseItemProcurementPayment(id uint64, warehouseItemProcurementId uint64, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemProcurement, err := s.repository.GetWarehouseItemProcurement(warehouseItemProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item procurement payment", zap.Error(err))
		return err
	}

	totalPrice := decimal.Zero
	for _, e := range warehouseItemProcurement.Payments {
		if e.Id != id {
			totalPrice = totalPrice.Add(e.Nominal)
		}
	}

	warehouseItemProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if totalPrice.LessThan(decimal.Zero) {
		return errx.BadRequest("delete this payment make minus")
	} else if totalPrice.LessThan(warehouseItemProcurement.TotalPrice) {
		warehouseItemProcurement.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseItemProcurement.PaidDate = sql.NullTime{Valid: false}
	}

	err = s.repository.DeleteWarehouseItemProcurementPayment(id)
	if err != nil {
		s.log.Error("failed delete warehouse item procurement payment", zap.Error(err))
		return err
	}

	err = s.repository.UpdateWarehouseItemProcurement(&warehouseItemProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item procurement", zap.Error(err))
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) ArrivalConfirmationWarehouseItemProcurement(id uint64, request dto.ArrivalConfirmationWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItemProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemProcurement, err := s.repository.GetWarehouseItemProcurement(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	warehouseItemProcurement.ReceiveQuantity = sql.NullFloat64{Float64: request.Quantity, Valid: true}
	warehouseItemProcurement.Note = request.Note
	warehouseItemProcurement.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	warehouseItemProcurement.TakenBy = uuid.NullUUID{UUID: userId, Valid: true}
	warehouseItemProcurement.IsArrived = true

	if warehouseItemProcurement.Quantity != request.Quantity {
		warehouseItemProcurement.Status = enum.ProcurementStatusArrivedNotOk
	} else {
		warehouseItemProcurement.Status = enum.ProcurementStatusArrivedOk
	}

	warehouseItem := entity.WarehouseItem{
		ItemId:      warehouseItemProcurement.ItemId,
		WarehouseId: warehouseItemProcurement.WarehouseId,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	warehouseItem, err = s.repository.FirstOrCreateWarehouseItem(warehouseItem)
	if err != nil {
		s.log.Error("failed first or create warehouse item", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	warehouseItem.Quantity = warehouseItem.Quantity + request.Quantity

	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	err = s.repository.UpdateWarehouseItemProcurement(&warehouseItemProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed update warehouse item", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	jsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemName:       warehouseItem.Item.Name,
		ItemUnit:       warehouseItem.Item.Unit,
		Source:         warehouseItem.Item.Name,
		Destination:    warehouseItemProcurement.Warehouse.Name,
		QuantityBefore: warehouseItem.Quantity,
		QuantityAfter:  request.Quantity + warehouseItem.Quantity,
		UserId:         userId,
		Status:         enum.ItemHistoryStatusIn,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, errx.BadRequest("failed parsed struct into json")
	}

	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonParsed)

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemProcurement(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	response := mapper.WarehouseItemProcurementToResponse(&data)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) CreateWarehouseItemCornProcurementDraft(request dto.CreateWarehouseItemCornProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	ovenCondition := enum.OvenConditionNotInput
	if request.OvenCondition != nil {
		ovenCondition := enum.ValueOfOvenCondition(*request.OvenCondition)
		if !ovenCondition.IsValid() {
			return dto.WarehouseItemCornProcurementDraftResponse{}, errx.BadRequest("invalid oven condition")
		}
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	data := entity.WarehouseItemCornProcurementDraft{
		WarehouseId:    request.WarehouseId,
		SupplierId:     sql.NullInt64{Int64: int64(request.SupplierId), Valid: true},
		CornWaterLevel: sql.NullFloat64{Float64: request.CornWaterLevel, Valid: true},
		OvenCondition:  ovenCondition,
		Quantity:       request.Quantity,
		Price:          price,
		Discount:       sql.NullFloat64{Float64: request.Discount, Valid: true},
		CreatedBy:      uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.IsOvenCanOperateInNearDay != nil {
		data.IsOvenCanOperateInNearDay = sql.NullBool{Bool: *request.IsOvenCanOperateInNearDay, Valid: true}
	}

	err = s.repository.CreateWarehouseItemCornProcurementDraft(&data)
	if err != nil {
		s.log.Error("failed create warehouse item corn procurement draft", zap.Error(err))
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemCornProcurementDraft(data.Id)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement draft", zap.Error(err))
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	return mapper.WarehouseItemCornProcurementDraftToResponse(&data, cornItem), nil
}

func (s *WarehouseService) GetWarehouseItemCornProcurementDrafts(filter dto.GetWarehouseItemCornProcurementDraftFilter) ([]dto.WarehouseItemCornProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemCornProcurementDrafts(filter)
	if err != nil {
		return nil, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return nil, err
	}

	response := make([]dto.WarehouseItemCornProcurementDraftResponse, 0)
	for _, e := range data {
		response = append(response, mapper.WarehouseItemCornProcurementDraftToResponse(&e, cornItem))
	}

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemCornProcurementDraft(id uint64) (dto.WarehouseItemCornProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemCornProcurementDraft(id)
	if err != nil {
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	return mapper.WarehouseItemCornProcurementDraftToResponse(&data, cornItem), nil
}

func (s *WarehouseService) UpdateWarehouseItemCornProcurementDraft(id uint64, request dto.UpdateWarehouseItemCornProcurementDraftRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementDraftResponse, error) {
	s.repository.UseTx(false)

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemCornProcurementDraft(id)
	if err != nil {
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	data.WarehouseId = request.WarehouseId
	data.SupplierId = sql.NullInt64{Int64: int64(request.SupplierId), Valid: true}
	data.CornWaterLevel = sql.NullFloat64{Float64: request.CornWaterLevel, Valid: true}
	data.Price = price
	data.Quantity = request.Quantity
	data.Discount = sql.NullFloat64{Float64: request.Discount, Valid: true}
	data.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if request.IsOvenCanOperateInNearDay != nil {
		data.IsOvenCanOperateInNearDay = sql.NullBool{Bool: *request.IsOvenCanOperateInNearDay, Valid: true}
	}

	if request.OvenCondition != nil {
		ovenCondition := enum.ValueOfOvenCondition(*request.OvenCondition)
		if !ovenCondition.IsValid() {
			return dto.WarehouseItemCornProcurementDraftResponse{}, errx.BadRequest("invalid oven condition")

		}
		data.OvenCondition = ovenCondition

	}

	err = s.repository.UpdateWarehouseItemCornProcurementDraft(&data)
	if err != nil {
		s.log.Error("failed update warehouse item corn procurement draft", zap.Error(err))
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemCornProcurementDraft(id)
	if err != nil {
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementDraftResponse{}, err
	}

	return mapper.WarehouseItemCornProcurementDraftToResponse(&data, cornItem), nil
}

func (s *WarehouseService) ConfirmationWarehouseItemCornProcurementDraft(id uint64, request dto.CreateWarehouseItemCornProcurementRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	err := s.repository.DeleteWarehouseItemCornProcurementDraft(id)
	if err != nil {
		s.log.Error("failed delete warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	ovenCondition := enum.OvenConditionNotInput
	if request.OvenCondition != nil {
		ovenCondition = enum.ValueOfOvenCondition(*request.OvenCondition)
		if !ovenCondition.IsValid() {
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid oven condition")
		}
	}

	expiredAt, err := time.Parse("02-01-2006", request.ExpiredAt)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid expired at")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type")
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment type")
	}

	data := entity.WarehouseItemCornProcurement{
		WarehouseId:               request.WarehouseId,
		SupplierId:                request.SupplierId,
		ExpiredAt:                 expiredAt,
		Price:                     price,
		Quantity:                  request.Quantity,
		Status:                    enum.ProcurementStatusSentOff,
		Discount:                  request.Discount,
		PaymentStatus:             enum.PaymentStatusNotPaid,
		CornWaterLevel:            request.Quantity,
		OvenCondition:             ovenCondition,
		IsOvenCanOperateInNearDay: *request.IsOvenCanOperateInNearDay,
		PaymentType:               paymentType,
		CreatedBy:                 uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.DeadlinePaymentDate != nil {
		deadlinePaymentDate, err := time.Parse("02-01-2006", *request.DeadlinePaymentDate)
		if err != nil {
			s.log.Error("failed parse deadline payment date", zap.Error(err))
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid deadline payment date format")
		}

		data.DeadlinePaymentDate = sql.NullTime{Time: deadlinePaymentDate, Valid: true}
	}

	discountPrice := price.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	data.TotalPrice = price.Sub(discountPrice).Mul(decimal.NewFromFloat(request.Quantity))

	if len(request.Payments) == 0 {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("payments are required")
	}

	payments := make([]entity.WarehouseItemCornProcurementPayment, 0, len(request.Payments))
	totalPayment := decimal.Zero
	for _, p := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(p.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment method")
		}
		nominal, err := decimal.NewFromString(p.Nominal)
		if err != nil {
			s.log.Error("failed parse payment nominal", zap.Error(err))
			return dto.WarehouseItemCornProcurementResponse{}, err
		}
		paymentDate, err := time.Parse("02-01-2006", p.PaymentDate)
		if err != nil {
			s.log.Error("failed parse payment date", zap.Error(err))
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment date")
		}
		totalPayment = totalPayment.Add(nominal)
		payments = append(payments, entity.WarehouseItemCornProcurementPayment{
			PaymentDate:                    paymentDate,
			Nominal:                        nominal,
			PaymentProof:                   p.PaymentProof,
			PaymentMethod:                  paymentMethod,
			WarehouseItemCornProcurementId: data.Id,
			CreatedBy:                      uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff && !totalPayment.Equal(data.TotalPrice) {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("need payment to paid off")
	}

	if totalPayment.Equal(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusPaid
		data.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.LessThan(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("total payment more than total price")
	}

	err = s.repository.CreateWarehouseItemCornProcurement(&data)
	if err != nil {
		s.log.Error("failed create warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	for i := range payments {
		payments[i].WarehouseItemCornProcurementId = data.Id
	}

	err = s.repository.CreateWarehouseItemCornProcurementPaymentInBatch(&payments)
	if err != nil {
		s.log.Error("failed create warehouse item corn procurement payments in batch", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemCornProcurement(data.Id)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentResponses := make([]dto.WarehouseItemCornProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemCornProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		paymentResponses = append(paymentResponses, payment)
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	response := mapper.WarehouseItemCornProcurementToResponse(&data, &cornItem)
	response.Payments = paymentResponses
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) DeleteWarehouseItemCornProcurementDraft(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouseItemCornProcurementDraft(id)
	if err != nil {
		s.log.Error("failed delete warehouse item corn procurement draft", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseItemCornProcurement(request dto.CreateWarehouseItemCornProcurementRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed parse price", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	ovenCondition := enum.OvenConditionNotInput
	if request.OvenCondition != nil {
		ovenCondition = enum.ValueOfOvenCondition(*request.OvenCondition)
		if !ovenCondition.IsValid() {
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid oven condition")
		}
	}

	expiredAt, err := time.Parse("02-01-2006", request.ExpiredAt)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid expired at")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("invalid payment type")
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment type")
	}

	data := entity.WarehouseItemCornProcurement{
		WarehouseId:               request.WarehouseId,
		SupplierId:                request.SupplierId,
		ExpiredAt:                 expiredAt,
		Price:                     price,
		Quantity:                  request.Quantity,
		Discount:                  request.Discount,
		Status:                    enum.ProcurementStatusSentOff,
		PaymentStatus:             enum.PaymentStatusNotPaid,
		CornWaterLevel:            request.CornWaterLevel,
		OvenCondition:             ovenCondition,
		IsOvenCanOperateInNearDay: *request.IsOvenCanOperateInNearDay,
		PaymentType:               paymentType,
	}

	if request.DeadlinePaymentDate != nil {
		deadlinePaymentDate, err := time.Parse("02-01-2006", *request.DeadlinePaymentDate)
		if err != nil {
			s.log.Error("failed parse deadline payment date", zap.Error(err))
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid deadline payment date format")
		}

		data.DeadlinePaymentDate = sql.NullTime{Time: deadlinePaymentDate, Valid: true}

	}

	discountPrice := price.Mul(decimal.NewFromFloat(request.Discount / 100.0))
	data.TotalPrice = price.Sub(discountPrice).Mul(decimal.NewFromFloat(request.Quantity))

	if len(request.Payments) == 0 {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("payments are required")
	}

	payments := make([]entity.WarehouseItemCornProcurementPayment, 0, len(request.Payments))
	totalPayment := decimal.Zero
	for _, p := range request.Payments {
		paymentMethod := enum.ValueOfPaymentMethod(p.PaymentMethod)
		if !paymentMethod.IsValid() {
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment method")
		}
		nominal, err := decimal.NewFromString(p.Nominal)
		if err != nil {
			s.log.Error("failed parse payment nominal", zap.Error(err))
			return dto.WarehouseItemCornProcurementResponse{}, err
		}
		paymentDate, err := time.Parse("02-01-2006", p.PaymentDate)
		if err != nil {
			s.log.Error("failed parse payment date", zap.Error(err))
			return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment date")
		}
		totalPayment = totalPayment.Add(nominal)
		payments = append(payments, entity.WarehouseItemCornProcurementPayment{
			PaymentDate:                    paymentDate,
			Nominal:                        nominal,
			PaymentProof:                   p.PaymentProof,
			PaymentMethod:                  paymentMethod,
			WarehouseItemCornProcurementId: data.Id,
			CreatedBy:                      uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	if paymentType == enum.PaymentTypePaidOff && !totalPayment.Equal(data.TotalPrice) {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("need payment to paid off")
	}

	if totalPayment.Equal(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusPaid
		data.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPayment.LessThan(data.TotalPrice) {
		data.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("total payment more than total price")
	}

	err = s.repository.CreateWarehouseItemCornProcurement(&data)
	if err != nil {
		s.log.Error("failed create warehouse item procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	for i := range payments {
		payments[i].WarehouseItemCornProcurementId = data.Id
	}

	err = s.repository.CreateWarehouseItemCornProcurementPaymentInBatch(&payments)
	if err != nil {
		s.log.Error("failed create warehouse item corn procurement payments in batch", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	data, err = s.repository.GetWarehouseItemCornProcurement(data.Id)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentResponses := make([]dto.WarehouseItemCornProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemCornProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		paymentResponses = append(paymentResponses, payment)
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	response := mapper.WarehouseItemCornProcurementToResponse(&data, &cornItem)
	response.Payments = paymentResponses
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemCornProcurements(filter dto.GetWarehouseItemCornProcurementFilter) (dto.WarehouseItemCornProcurementListPaginationResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemCornProcurements(filter)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurements", zap.Error(err))
		return dto.WarehouseItemCornProcurementListPaginationResponse{}, err
	}

	totalData, err := s.repository.CountWarehouseItemCornProcurement(filter)
	if err != nil {
		s.log.Error("failed count warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementListPaginationResponse{}, err
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementListPaginationResponse{}, err
	}

	warehouseItemProcurementResponses := make([]dto.WarehouseItemCornProcurementListResponse, 0)
	for _, e := range data {
		warehouseItemProcurementResponses = append(warehouseItemProcurementResponses, mapper.WarehouseItemCornProcurementToListResponse(&e, &cornItem))
	}

	response := dto.WarehouseItemCornProcurementListPaginationResponse{
		WarehouseItemCornProcurements: warehouseItemProcurementResponses,
	}

	if filter.Page > 0 {
		response.TotalData = uint64(totalData)
		response.TotalPage = uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit)))
	}

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemCornProcurement(id uint64) (dto.WarehouseItemCornProcurementResponse, error) {
	s.repository.UseTx(false)

	data, err := s.repository.GetWarehouseItemCornProcurement(id)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemCornProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemCornProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	response := mapper.WarehouseItemCornProcurementToResponse(&data, &cornItem)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) CreateWarehouseItemCornProcurementPayment(warehouseItemCornProcurementId uint64, request dto.CreateWarehouseItemCornProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemCornProcurement, err := s.repository.GetWarehouseItemCornProcurement(warehouseItemCornProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed parse payment date", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment date")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed parse nominal", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment method")
	}

	payment := entity.WarehouseItemCornProcurementPayment{
		WarehouseItemCornProcurementId: warehouseItemCornProcurement.Id,
		PaymentDate:                    paymentDate,
		PaymentProof:                   request.PaymentProof,
		Nominal:                        nominal,
		PaymentMethod:                  paymentMethod,
		CreatedBy:                      uuid.NullUUID{UUID: userId, Valid: true},
	}

	totalPrice := nominal
	for _, e := range warehouseItemCornProcurement.Payments {
		totalPrice = totalPrice.Add(e.Nominal)
	}

	warehouseItemCornProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if totalPrice.Equal(warehouseItemCornProcurement.TotalPrice) {
		warehouseItemCornProcurement.PaymentStatus = enum.PaymentStatusPaid
		warehouseItemCornProcurement.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPrice.LessThan(warehouseItemCornProcurement.TotalPrice) {
		warehouseItemCornProcurement.PaymentStatus = enum.PaymentStatusUnpaid
	} else {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("nominal is to high")
	}

	err = s.repository.CreateWarehouseItemCornProcurementPayment(&payment)
	if err != nil {
		s.log.Error("failed create warehouse item corn procurement payment", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.UpdateWarehouseItemCornProcurement(&warehouseItemCornProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemCornProcurement(warehouseItemCornProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemCornProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemCornProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	response := mapper.WarehouseItemCornProcurementToResponse(&data, &cornItem)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) UpdateWarehouseItemCornProcurementPayment(id uint64, warehouseItemCornProcurementId uint64, request dto.UpdateWarehouseItemCornProcurementPaymentRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemCornProcurement, err := s.repository.GetWarehouseItemCornProcurement(warehouseItemCornProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	warehouseItemCornProcurementPayment, err := s.repository.GetWarehouseItemCornProcurementPayment(id)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement payment", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("failed parse payment date", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment date")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("failed parse nominal", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("invalid payment method")
	}

	warehouseItemCornProcurementPayment.Nominal = nominal
	warehouseItemCornProcurementPayment.PaymentDate = paymentDate
	warehouseItemCornProcurementPayment.PaymentMethod = paymentMethod
	warehouseItemCornProcurementPayment.PaymentProof = request.PaymentProof
	warehouseItemCornProcurementPayment.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	totalPrice := nominal
	for _, e := range warehouseItemCornProcurement.Payments {
		if e.Id != id {
			totalPrice = totalPrice.Add(e.Nominal)
		}
	}

	warehouseItemCornProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if totalPrice.Equal(warehouseItemCornProcurement.TotalPrice) {
		warehouseItemCornProcurement.PaymentStatus = enum.PaymentStatusPaid
		warehouseItemCornProcurement.PaidDate = sql.NullTime{Time: time.Now(), Valid: true}
	} else if totalPrice.LessThan(warehouseItemCornProcurement.TotalPrice) {
		warehouseItemCornProcurement.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseItemCornProcurement.PaidDate = sql.NullTime{Valid: false}
	} else {
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("nominal is to high")
	}

	err = s.repository.UpdateWarehouseItemCornProcurementPayment(&warehouseItemCornProcurementPayment)
	if err != nil {
		s.log.Error("failed update warehouse item corn procurement payment", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.UpdateWarehouseItemCornProcurement(&warehouseItemCornProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemCornProcurement(warehouseItemCornProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	payments := make([]dto.WarehouseItemCornProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemCornProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		payments = append(payments, payment)
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	response := mapper.WarehouseItemCornProcurementToResponse(&data, &cornItem)
	response.Payments = payments
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) DeleteWarehouseItemCornProcurementPayment(id uint64, warehouseItemProcurementId uint64, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemCornProcurement, err := s.repository.GetWarehouseItemCornProcurement(warehouseItemProcurementId)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement payment", zap.Error(err))
		return err
	}

	totalPrice := decimal.Zero
	for _, e := range warehouseItemCornProcurement.Payments {
		if e.Id != id {
			totalPrice = totalPrice.Add(e.Nominal)
		}
	}

	warehouseItemCornProcurement.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	if totalPrice.LessThan(decimal.Zero) {
		return errx.BadRequest("delete this payment make minus")
	} else if totalPrice.LessThan(warehouseItemCornProcurement.TotalPrice) {
		warehouseItemCornProcurement.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseItemCornProcurement.PaidDate = sql.NullTime{Valid: false}
	}

	err = s.repository.DeleteWarehouseItemCornProcurementPayment(id)
	if err != nil {
		s.log.Error("failed delete warehouse item corn procurement payment", zap.Error(err))
		return err
	}

	err = s.repository.UpdateWarehouseItemCornProcurement(&warehouseItemCornProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item corn procurement", zap.Error(err))
		return err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) ArrivalConfirmationWarehouseItemCornProcurement(id uint64, request dto.ArrivalConfirmationWarehouseItemCornProcurementRequest, userId uuid.UUID) (dto.WarehouseItemCornProcurementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemCornProcurement, err := s.repository.GetWarehouseItemCornProcurement(id)
	if err != nil {
		s.log.Error("failed get warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	warehouseItemCornProcurement.ReceiveQuantity = sql.NullFloat64{Float64: request.Quantity, Valid: true}
	warehouseItemCornProcurement.Note = request.Note
	warehouseItemCornProcurement.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	warehouseItemCornProcurement.TakenBy = uuid.NullUUID{UUID: userId, Valid: true}
	warehouseItemCornProcurement.IsArrived = true

	if warehouseItemCornProcurement.Quantity != request.Quantity {
		warehouseItemCornProcurement.Status = enum.ProcurementStatusArrivedNotOk
	} else {
		warehouseItemCornProcurement.Status = enum.ProcurementStatusArrivedOk
	}

	warehouseItemCorn := entity.WarehouseItemCorn{
		WarehouseId: warehouseItemCornProcurement.WarehouseId,
		SupplierId:  warehouseItemCornProcurement.SupplierId,
		Quantity:    warehouseItemCornProcurement.Quantity,
		OrderDate:   warehouseItemCornProcurement.CreatedAt,
		ExpiredAt:   warehouseItemCornProcurement.ExpiredAt,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	totalItemCornQuantity, err := s.repository.CountQuantityWarehouseItemCornByWarehouseId(warehouseItemCornProcurement.WarehouseId)
	if err != nil {
		s.log.Error("failed to count total item corn quantity", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.CreateWarehouseItemCorn(&warehouseItemCorn)
	if err != nil {
		s.log.Error("failed create warehouse item corn", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	err = s.repository.UpdateWarehouseItemCornProcurement(&warehouseItemCornProcurement)
	if err != nil {
		s.log.Error("failed update warehouse item corn procurement", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	itemCorn, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	jsonParsed, err := json.Marshal(entity.WarehouseItemHistory{
		ItemName:       itemCorn.Name,
		ItemUnit:       itemCorn.Unit,
		Source:         warehouseItemCornProcurement.Warehouse.Name,
		Destination:    warehouseItemCornProcurement.Warehouse.Name,
		QuantityBefore: totalItemCornQuantity,
		QuantityAfter:  totalItemCornQuantity + request.Quantity,
		UserId:         userId,
		Status:         enum.ItemHistoryStatusIn,
	})

	if err != nil {
		s.log.Error("failed to parse struct into json", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, errx.BadRequest("failed parsed struct into json")
	}

	s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonParsed)

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	data, err := s.repository.GetWarehouseItemCornProcurement(id)
	if err != nil {
		s.log.Error("failed get warehouse item procurement corn", zap.Error(err))
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	paymentResponses := make([]dto.WarehouseItemCornProcurementPaymentResponse, 0)
	remainingPayment := data.TotalPrice
	for _, e := range data.Payments {
		payment := mapper.WarehouseItemCornProcurementPaymentToResponse(&e)
		remainingPayment = remainingPayment.Sub(e.Nominal)
		payment.Remaining = remainingPayment.String()
		paymentResponses = append(paymentResponses, payment)
	}

	cornItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.Corn, constant.UnitKg, enum.ItemCategoryCornMaterial)
	if err != nil {
		return dto.WarehouseItemCornProcurementResponse{}, err
	}

	response := mapper.WarehouseItemCornProcurementToResponse(&data, &cornItem)
	response.Payments = paymentResponses
	response.RemainingPayment = remainingPayment.String()

	return response, nil
}

func (s *WarehouseService) GetWarehouseItemCornPrices() ([]dto.WarehouseItemCornPriceResponse, error) {
	data, err := s.repository.GetWarehouseItemCornPrices()
	if err != nil {
		s.log.Error("failed get warehouse item corn prices", zap.Error(err))
		return nil, err
	}

	responses := make([]dto.WarehouseItemCornPriceResponse, 0)
	for _, e := range data {
		responses = append(responses, mapper.WarehouseItemCornPriceToResponse(&e))
	}

	return responses, nil
}

func (s *WarehouseService) CreateRawFeed(request dto.CreateRawFeedRequest, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	cornPrice, err := decimal.NewFromString(request.CornPrice)
	if err != nil {
		s.log.Error("failed to parse corn price", zap.Error(err))
		return err
	}

	cornDraft := entity.WarehouseItemCornProcurementDraft{
		WarehouseId:   request.WarehouseId,
		OvenCondition: enum.OvenConditionNotInput,
		Quantity:      request.CornQuantity,
		Price:         cornPrice,
		CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
	}
	err = s.repository.CreateWarehouseItemCornProcurementDraft(&cornDraft)
	if err != nil {
		s.log.Error("failed to create corn procurement draft", zap.Error(err))
		return err
	}

	drafts := make([]entity.WarehouseItemProcurementDraft, 0, len(request.RawMaterials))
	for _, rm := range request.RawMaterials {
		price, err := decimal.NewFromString(rm.Price)
		if err != nil {
			s.log.Error("failed to parse raw material price", zap.Error(err))
			return err
		}
		draft := entity.WarehouseItemProcurementDraft{
			WarehouseId:   request.WarehouseId,
			ItemId:        rm.ItemId,
			DailySpending: rm.DailySpending,
			DaysNeed:      request.DaysNeed,
			Price:         price,
			CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
		}
		drafts = append(drafts, draft)
	}

	if len(drafts) > 0 {
		err := s.repository.CreateWarehouseItemProcurementDraftInBatch(&drafts)
		if err != nil {
			s.log.Error("failed to create warehouse item procurement drafts in batch", zap.Error(err))
			return err
		}
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateReadyToEatFeed(request dto.CreateReadyToEatFeedRequest, userId uuid.UUID) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price", zap.Error(err))
		return err
	}

	draft := entity.WarehouseItemProcurementDraft{
		WarehouseId:   request.WarehouseId,
		ItemId:        request.ItemId,
		DailySpending: request.DailySpending,
		DaysNeed:      request.DaysNeed,
		Price:         price,
		CreatedBy:     uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateWarehouseItemProcurementDraft(&draft)
	if err != nil {
		s.log.Error("failed to create warehouse item procurement draft", zap.Error(err))
		return err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) ReduceWarehouseItemForFeed(warehouseId uint64, request []dto.ReduceFeedRequest, userId uuid.UUID, cageName string) error {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouseItemHistories := make([]entity.WarehouseItemHistory, 0)
	warehouse, err := s.repository.GetWarehouseById(warehouseId)
	if err != nil {
		s.log.Error("failed get warehouse by id", zap.Error(err))
		return err
	}

	itemIds := make([]uint64, 0)
	for _, r := range request {
		if r.ItemCategory != enum.ItemCategoryCornMaterial.String() {
			itemIds = append(itemIds, r.ItemId)
		}
	}

	itemMap := make(map[uint64]entity.WarehouseItem)
	items, err := s.repository.GetWarehouseItemByWarehouseIdAndItemIds(warehouseId, itemIds)
	if err != nil {
		s.log.Error("failed get warehouse item by warehouse id and item id", zap.Error(err))
		return err
	}

	for _, item := range items {
		itemMap[item.ItemId] = item
	}

	errorStockNotEnough := make([]string, 0)
	itemAndReduceQuantityMap := make(map[uint64]float64)

	for _, r := range request {
		withZeroQuantity := false
		if r.ItemCategory == enum.ItemCategoryCornMaterial.String() {
			corns, err := s.repository.GetWarehouseItemCorns(dto.GetWarehouseItemCornFilter{
				WarehouseId:      warehouseId,
				WithZeroQuantity: &withZeroQuantity,
				FromNewest:       false,
			})
			if err != nil {
				s.log.Error("failed get corn items", zap.Error(err))
				return err
			}

			totalCornQuantity, err := s.repository.CountQuantityWarehouseItemCornByWarehouseId(warehouseId)
			if err != nil {
				s.log.Error("failed to count total corn quantity", zap.Error(err))
				return err
			}

			qtyNeeded := r.Quantity
			for _, corn := range corns {
				if qtyNeeded <= 0 {
					break
				}

				if corn.Quantity >= qtyNeeded {
					corn.Quantity -= qtyNeeded
					qtyNeeded = 0
				} else {
					qtyNeeded -= corn.Quantity
					corn.Quantity = 0
				}

				corn.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
				if err := s.repository.UpdateWarehouseItemCorn(&corn); err != nil {
					s.log.Error("failed update corn item", zap.Error(err))
					return err
				}
			}
			if qtyNeeded > 0 {
				return errx.BadRequest(fmt.Sprintf(
					"stok jagung tidak mencukupi, memerlukan %.2f", qtyNeeded,
				))
			}

			warehouseItemHistories = append(warehouseItemHistories, entity.WarehouseItemHistory{
				ItemName:       constant.Corn,
				ItemUnit:       constant.UnitKg,
				Source:         warehouse.Name,
				Destination:    cageName,
				QuantityBefore: totalCornQuantity,
				QuantityAfter:  totalCornQuantity - r.Quantity,
				UserId:         userId,
				Status:         enum.ItemHistoryStockUpdated,
			})
		} else {
			item, ok := itemMap[r.ItemId]
			if !ok {
				s.log.Error("failed get warehouse item by id and warehouse id", zap.Error(err))
				return err
			}

			if item.Quantity < r.Quantity {
				errorStockNotEnough = append(errorStockNotEnough, fmt.Sprintf(
					"stok barang %s tidak cukup memiliki %.2f, sedangkan memerlukan %.2f",
					item.Item.Name, item.Quantity, r.Quantity,
				))
			} else {
				warehouseItemHistories = append(warehouseItemHistories, entity.WarehouseItemHistory{
					ItemName:       item.Item.Name,
					ItemUnit:       item.Item.Unit,
					Source:         warehouse.Name,
					Destination:    cageName,
					QuantityBefore: item.Quantity,
					QuantityAfter:  item.Quantity - r.Quantity,
					UserId:         userId,
					Status:         enum.ItemHistoryStockUpdated,
				})
			}

			itemAndReduceQuantityMap[item.ItemId] = r.Quantity
		}
	}

	if len(errorStockNotEnough) > 0 {
		return fmt.Errorf("stock not enough: %s", strings.Join(errorStockNotEnough, ", "))
	}

	for _, item := range items {
		item.Quantity -= itemAndReduceQuantityMap[item.ItemId]
		item.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
		if err := s.repository.UpdateWarehouseItem(&item); err != nil {
			s.log.Error("failed update warehouse item", zap.Error(err))
			return err
		}
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("failed commit transaction", zap.Error(err))
		return err
	}

	// Todo : it can be optimized by just sending in array to the queue
	for _, history := range warehouseItemHistories {
		jsonParsed, err := json.Marshal(history)

		if err != nil {
			s.log.Error("failed to parse struct into json", zap.Error(err))
			return errx.BadRequest("failed parsed struct into json")
		}

		s.cacheService.Publish(context.Background(), constant.WarehouseItemHistoryTopic, jsonParsed)
	}

	return nil
}
