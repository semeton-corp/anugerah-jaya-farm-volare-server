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
	"go.uber.org/zap"
)

type WarehouseService struct {
	log        *zap.Logger
	repository repository.IWarehouseRepository
}

type IWarehouseService interface {
	CreateWarehouseItem(request *dto.CreateWarehouseItemRequest, accountId uuid.UUID) (dto.WarehouseItemResponse, error)
	GetWarehouseItem(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error)

	GetWarehouses() ([]dto.WarehouseResponse, error)

	CreateWarehouseStockItem(request *dto.CreateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error)
	GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]dto.WarehouseStockItemResponse, error)
	GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseStockItemResponse, error)
	UpdateWarehouseStockItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error)
	DeleteWarehouseStockItem(warehouseId uint64, warehouseItemId uint64) error
}

func NewWarehouseService(log *zap.Logger, repository repository.IWarehouseRepository) IWarehouseService {
	return &WarehouseService{
		log:        log,
		repository: repository,
	}
}

func (w *WarehouseService) CreateWarehouseItem(request *dto.CreateWarehouseItemRequest, accountId uuid.UUID) (dto.WarehouseItemResponse, error) {
	warehouseItem := entity.WarehouseItem{
		Name:      request.Name,
		Unit:      request.Unit,
		Category:  enum.WarehouseItemCategoryFeed,
		CreatedBy: accountId,
	}

	err := w.repository.CreateWarehouseItem(&warehouseItem)
	if err != nil {
		w.log.Error("[CreateWarehouseItem] failed to create warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	return mapper.WarehouseItemToResponse(&warehouseItem), nil
}

func (w *WarehouseService) GetWarehouseItem(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error) {
	warehouseItems, err := w.repository.GetWarehouseItem(filter)
	if err != nil {
		w.log.Error("[GetWarehouseItem] failed to get warehouse items", zap.Error(err))
		return nil, err
	}

	warehouseItemResponses := make([]dto.WarehouseItemResponse, 0, len(warehouseItems))
	for _, item := range warehouseItems {
		warehouseItemResponses = append(warehouseItemResponses, mapper.WarehouseItemToResponse(&item))
	}

	return warehouseItemResponses, nil
}

func (w *WarehouseService) GetWarehouses() ([]dto.WarehouseResponse, error) {
	warehouses, err := w.repository.GetWarehouse()
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
	err := w.repository.DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		w.log.Error("[DeleteStockWarehouseItem] failed to delete stock warehouse item", zap.Error(err))
		return err
	}

	return nil
}
