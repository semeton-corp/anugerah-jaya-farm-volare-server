package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type WarehouseService struct {
	log          *zap.Logger
	repository   repository.IWarehouseRepository
	storeService IStoreService
}

type IWarehouseService interface {
	CreateWarehouseItem(request dto.CreateWarehouseItemRequest, accountId uuid.UUID) (dto.WarehouseItemResponse, error)
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error)
	UpdateWarehouseItem(warehouseItemId uint64, request dto.UpdateWarehouseItemRequest, accountId uuid.UUID) (dto.WarehouseItemResponse, error)
	GetWarehouseItemById(id uint64) (dto.WarehouseItemResponse, error)

	GetWarehouses() ([]dto.WarehouseResponse, error)

	CreateWarehouseStockItem(request *dto.CreateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error)
	GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]dto.WarehouseStockItemResponse, error)
	GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseStockItemResponse, error)
	UpdateWarehouseStockItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error)
	DeleteWarehouseStockItem(warehouseId uint64, warehouseItemId uint64) error

	CreateWarehouseOrderItem(request dto.CreateWarehouseOrderItemRequest, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error)
	GetWarehouseOrderItemById(id uint64) (dto.WarehouseOrderItemResponse, error)
	GetWarehouseOrderItems() ([]dto.WarehouseOrderItemResponse, error)
	DeleteWarehouseOrderItem(id uint64) error
	TakeWarehouseOrderItem(id uint64, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error)
}

func NewWarehouseService(log *zap.Logger, repository repository.IWarehouseRepository, storeService IStoreService) IWarehouseService {
	return &WarehouseService{
		log:          log,
		repository:   repository,
		storeService: storeService,
	}
}

func (w *WarehouseService) CreateWarehouseItem(request dto.CreateWarehouseItemRequest, accountId uuid.UUID) (dto.WarehouseItemResponse, error) {
	warehouseItemCategory := enum.ValueOfWarehouseItemCategory(request.Category)
	if !warehouseItemCategory.IsValid() {
		w.log.Error("[CreateWarehouseItem] invalid warehouse item category", zap.String("category", request.Category))
		return dto.WarehouseItemResponse{}, errx.BadRequest("invalid warehouse item category")
	}

	if !warehouseItemCategory.IsValid() {
		w.log.Error("[CreateWarehouseItem] invalid warehouse item category", zap.String("category", request.Category))
		return dto.WarehouseItemResponse{}, errx.BadRequest("invalid warehouse item category")
	}

	warehouseItem := entity.WarehouseItem{
		Name:      request.Name,
		Unit:      request.Unit,
		Category:  warehouseItemCategory,
		CreatedBy: accountId,
	}

	err := w.repository.CreateWarehouseItem(&warehouseItem)
	if err != nil {
		w.log.Error("[CreateWarehouseItem] failed to create warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	return mapper.WarehouseItemToResponse(&warehouseItem), nil
}

func (w *WarehouseService) GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error) {
	if filter.StoreId > 0 && filter.WarehouseId > 0 {
		w.log.Error("[GetWarehouseItems] storeId and warehouseId cannot be used at the same time")
		return nil, errx.BadRequest("storeId and warehouseId cannot be used at the same time")
	}

	if filter.StoreId > 0 {
		storeItems, err := w.storeService.GetStoreItems(
			dto.GetStoreItemFilter{
				StoreId:  filter.StoreId,
				Category: filter.Category,
			},
		)
		if err != nil {
			w.log.Error("[GetWarehouseItems] failed to get store items", zap.Error(err))
			return nil, err
		}

		storeItemResponses := make([]dto.WarehouseItemResponse, 0, len(storeItems))
		for _, item := range storeItems {
			storeItemResponses = append(storeItemResponses, item.WarehouseItem)
		}

		return storeItemResponses, nil
	}

	if filter.WarehouseId > 0 {
		warehouseStockItems, err := w.repository.GetWarehouseStockItems(
			dto.GetWarehouseStockItemFilter{
				WarehouseId: filter.WarehouseId,
				Category:    filter.Category,
			},
		)
		if err != nil {
			w.log.Error("[GetWarehouseItems] failed to get warehouse stock items", zap.Error(err))
			return nil, err
		}

		warehouseStockItemResponses := make([]dto.WarehouseItemResponse, 0, len(warehouseStockItems))
		for _, item := range warehouseStockItems {
			warehouseStockItemResponses = append(warehouseStockItemResponses, mapper.WarehouseItemToResponse(&item.WarehouseItem))
		}

		return warehouseStockItemResponses, nil
	}

	warehouseItems, err := w.repository.GetWarehouseItems(filter)
	if err != nil {
		w.log.Error("[GetWarehouseItems] failed to get warehouse items", zap.Error(err))
		return nil, err
	}

	warehouseItemResponses := make([]dto.WarehouseItemResponse, 0, len(warehouseItems))
	for _, item := range warehouseItems {
		warehouseItemResponses = append(warehouseItemResponses, mapper.WarehouseItemToResponse(&item))
	}

	return warehouseItemResponses, nil
}

func (w *WarehouseService) UpdateWarehouseItem(warehouseItemId uint64, request dto.UpdateWarehouseItemRequest, accountId uuid.UUID) (dto.WarehouseItemResponse, error) {
	warehouseItemCategory := enum.ValueOfWarehouseItemCategory(request.Category)
	if !warehouseItemCategory.IsValid() {
		w.log.Error("[UpdateWarehouseItem] invalid warehouse item category", zap.String("category", request.Category))
		return dto.WarehouseItemResponse{}, errx.BadRequest("invalid warehouse item category")
	}

	warehouseItem, err := w.repository.GetWarehouseItemById(warehouseItemId)
	if err != nil {
		w.log.Error("[UpdateWarehouseItem] failed to get warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	warehouseItem.Name = request.Name
	warehouseItem.Unit = request.Unit
	warehouseItem.Category = warehouseItemCategory
	warehouseItem.UpdatedBy = accountId

	err = w.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		w.log.Error("[UpdateWarehouseItem] failed to update warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	return mapper.WarehouseItemToResponse(&warehouseItem), nil
}

func (w *WarehouseService) GetWarehouseItemById(id uint64) (dto.WarehouseItemResponse, error) {
	warehouseItem, err := w.repository.GetWarehouseItemById(id)
	if err != nil {
		w.log.Error("[GetWarehouseItemById] failed to get warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	warehouseItemResponse := mapper.WarehouseItemToResponse(&warehouseItem)
	return warehouseItemResponse, nil
}

func (w *WarehouseService) GetWarehouses() ([]dto.WarehouseResponse, error) {
	warehouses, err := w.repository.GetWarehouses()
	if err != nil {
		w.log.Error("[GetWarehouses] failed to get warehouses", zap.Error(err))
		return nil, err
	}

	warehouseResponses := make([]dto.WarehouseResponse, 0, len(warehouses))
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, mapper.WarehouseToResponse(&warehouse))
	}

	return warehouseResponses, nil
}

func (w *WarehouseService) CreateWarehouseStockItem(request *dto.CreateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error) {
	// Todo : create estimation run out date, based on average used per day from request item from request warhouse item.

	stockWarehouseItem := entity.WarehouseStockItem{
		WarehouseId:      request.WarehouseId,
		WarehouseItemId:  request.WarehouseItemId,
		Quantity:         request.Quantity,
		EstimationRunOut: time.Now(),
		CreatedBy:        accountId,
	}

	err := w.repository.CreateWarehouseStockItem(&stockWarehouseItem)
	if err != nil {
		w.log.Error("[CreateStockWarehouseItem] failed to create stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	stockWarehouseItem, err = w.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(
		stockWarehouseItem.WarehouseId,
		stockWarehouseItem.WarehouseItemId,
	)
	if err != nil {
		w.log.Error("[CreateStockWarehouseItem] failed to get stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	var description string
	if time.Now().Add(time.Hour * 24 * 7).After(stockWarehouseItem.EstimationRunOut) {
		description = constant.StockWarehouseItemDanger
	} else {
		description = constant.StockWarehouseItemSafe
	}

	warehouseStockItemResponse := mapper.WarehouseStockItemToResponse(&stockWarehouseItem)
	warehouseStockItemResponse.Description = description
	return warehouseStockItemResponse, nil
}

func (w *WarehouseService) GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]dto.WarehouseStockItemResponse, error) {
	w.repository.UseTx(false)

	stockWarehouseItems, err := w.repository.GetWarehouseStockItems(filter)
	if err != nil {
		w.log.Error("[GetStockWarehouseItem] failed to get stock warehouse items", zap.Error(err))
		return nil, err
	}

	stockWarehouseItemResponses := make([]dto.WarehouseStockItemResponse, 0, len(stockWarehouseItems))
	for _, item := range stockWarehouseItems {
		var description string
		if time.Now().Add(time.Hour * 24 * 7).After(item.EstimationRunOut) {
			description = constant.StockWarehouseItemDanger
		} else {
			description = constant.StockWarehouseItemSafe
		}

		warehouseStockItemResponse := mapper.WarehouseStockItemToResponse(&item)
		warehouseStockItemResponse.Description = description

		stockWarehouseItemResponses = append(stockWarehouseItemResponses, warehouseStockItemResponse)
	}

	return stockWarehouseItemResponses, nil
}

func (w *WarehouseService) GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseStockItemResponse, error) {
	w.repository.UseTx(false)

	stockWarehouseItem, err := w.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		w.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to get stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	var description string
	if time.Now().Add(time.Hour * 24 * 7).After(stockWarehouseItem.EstimationRunOut) {
		description = constant.StockWarehouseItemDanger
	} else {
		description = constant.StockWarehouseItemSafe
	}

	warehouseStockItemResponse := mapper.WarehouseStockItemToResponse(&stockWarehouseItem)
	warehouseStockItemResponse.Description = description

	return warehouseStockItemResponse, nil
}

func (w *WarehouseService) UpdateWarehouseStockItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error) {
	w.repository.UseTx(false)

	stockWarehouseItem, err := w.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		w.log.Error("[UpdateStockWarehouseItem] failed to get stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	stockWarehouseItem.Quantity = request.Quantity
	stockWarehouseItem.UpdatedBy = accountId

	err = w.repository.UpdateWarehouseStockItem(&stockWarehouseItem)
	if err != nil {
		w.log.Error("[UpdateStockWarehouseItem] failed to update stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	var description string
	if time.Now().Add(time.Hour * 24 * 7).After(stockWarehouseItem.EstimationRunOut) {
		description = constant.StockWarehouseItemDanger
	} else {
		description = constant.StockWarehouseItemSafe
	}

	warehouseStockItemResponse := mapper.WarehouseStockItemToResponse(&stockWarehouseItem)
	warehouseStockItemResponse.Description = description

	return warehouseStockItemResponse, nil
}

func (w *WarehouseService) DeleteWarehouseStockItem(warehouseId uint64, warehouseItemId uint64) error {
	w.repository.UseTx(false)

	err := w.repository.DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		w.log.Error("[DeleteStockWarehouseItem] failed to delete stock warehouse item", zap.Error(err))
		return err
	}

	return nil
}

func (w *WarehouseService) CreateWarehouseOrderItem(request dto.CreateWarehouseOrderItemRequest, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error) {
	w.repository.UseTx(false)

	warehouseOrderItem := entity.WarehouseOrderItem{
		WarehouseId:     request.WarehouseId,
		WarehouseItemId: request.WarehouseItemId,
		Quantity:        request.Quantity,
	}

	err := w.repository.CreateWarehouseOrderItem(&warehouseOrderItem)
	if err != nil {
		w.log.Error("[CreateWarehouseOrderItem] failed to create warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (w *WarehouseService) GetWarehouseOrderItemById(id uint64) (dto.WarehouseOrderItemResponse, error) {
	w.repository.UseTx(false)

	warehouseOrderItem, err := w.repository.GetWarehouseOrderItemById(id)
	if err != nil {
		w.log.Error("[GetWarehouseOrderItemById] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (w *WarehouseService) GetWarehouseOrderItems() ([]dto.WarehouseOrderItemResponse, error) {
	w.repository.UseTx(false)

	warehouseOrderItems, err := w.repository.GetWarehouseOrderItems()
	if err != nil {
		w.log.Error("[GetWarehouseOrderItems] failed to get warehouse order items", zap.Error(err))
		return nil, err
	}

	warehouseOrderItemResponses := make([]dto.WarehouseOrderItemResponse, 0, len(warehouseOrderItems))
	for _, item := range warehouseOrderItems {
		warehouseOrderItemResponses = append(warehouseOrderItemResponses, mapper.WarehouseOrderItemToResponse(&item))
	}

	return warehouseOrderItemResponses, nil
}

func (w *WarehouseService) DeleteWarehouseOrderItem(id uint64) error {
	w.repository.UseTx(false)

	err := w.repository.DeleteWarehouseOrderItem(id)
	if err != nil {
		w.log.Error("[DeleteWarehouseOrderItem] failed to delete warehouse order item", zap.Error(err))
		return err
	}

	return nil
}

func (w *WarehouseService) TakeWarehouseOrderItem(id uint64, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error) {
	w.repository.UseTx(false)

	// Todo : add stock warehouse item in warehouse
	warehouseOrderItem, err := w.repository.GetWarehouseOrderItemById(id)
	if err != nil {
		w.log.Error("[TakeWarehouseOrderItem] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	warehouseOrderItem.TakenBy = accountId
	warehouseOrderItem.TakenAt = time.Now()
	warehouseOrderItem.UpdatedBy = accountId

	err = w.repository.UpdateWarehouseOrderItem(&warehouseOrderItem)
	if err != nil {
		w.log.Error("[TakeWarehouseOrderItem] failed to update warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}
