package service

import (
	"context"
	"database/sql"
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
	"go.uber.org/zap"
)

type WarehouseService struct {
	log          *zap.Logger
	repository   repository.IWarehouseRepository
	storeService IStoreService
	cacheService cache.ICache
}

type IWarehouseService interface {
	GetWarehouses(filter dto.GetWarehouseFilter) ([]dto.WarehouseResponse, error)
	CreateWarehouse(request dto.CreateWarehouseRequest, createdBy uuid.UUID) (dto.WarehouseResponse, error)
	DeleteWarehouse(id uint64) error
	UpdateWarehouse(id uint64, request dto.UpdateWarehouseRequest, createdBy uuid.UUID) (dto.WarehouseResponse, error)

	CreateWarehouseStockItem(request *dto.CreateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error)
	GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]dto.WarehouseStockItemResponse, error)
	GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseStockItemResponse, error)
	UpdateWarehouseStockItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error)
	DeleteWarehouseStockItem(warehouseId uint64, warehouseItemId uint64) error

	CreateWarehouseOrderItem(request dto.CreateWarehouseOrderItemRequest, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error)
	GetWarehouseOrderItemById(id uint64) (dto.WarehouseOrderItemResponse, error)
	GetWarehouseOrderItems(filter dto.GetWarehouseOrderItemFilter) ([]dto.WarehouseOrderItemResponse, error)
	DeleteWarehouseOrderItem(id uint64) error
	TakeWarehouseOrderItem(id uint64, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error)

	GoodEggConvertionButirToIkat(request dto.GoodEggWarehouseConvertionRequest, accountId uuid.UUID) ([]dto.WarehouseStockItemResponse, error)
	GoodEggConvertionIkatToButir(request dto.GoodEggWarehouseConvertionRequest, accountId uuid.UUID) ([]dto.WarehouseStockItemResponse, error)
	CrackedEggConvertionButirToPack(request dto.CrackedEggWarehouseConvertionRequest, accountId uuid.UUID) ([]dto.WarehouseStockItemResponse, error)
}

func NewWarehouseService(log *zap.Logger, repository repository.IWarehouseRepository, storeService IStoreService, cacheService cache.ICache) IWarehouseService {
	return &WarehouseService{
		log:          log,
		repository:   repository,
		storeService: storeService,
		cacheService: cacheService,
	}
}

func (s *WarehouseService) UpdateWarehouse(id uint64, request dto.UpdateWarehouseRequest, updateBy uuid.UUID) (dto.WarehouseResponse, error) {
	s.repository.UseTx(false)

	warehouse, err := s.repository.GetWarehouseById(id)
	if err != nil {
		s.log.Error("[UpdateWarehouse] failed to get warehouse by id", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	warehouse.Name = request.Name
	warehouse.LocationId = request.LocationId
	warehouse.UpdatedBy = uuid.NullUUID{UUID: updateBy, Valid: true}

	err = s.repository.UpdateWarehouse(&warehouse)
	if err != nil {
		s.log.Error("[UpdateWarehouse] failed to get udpate warehouse", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	warehouse, err = s.repository.GetWarehouseById(warehouse.Id)
	if err != nil {
		s.log.Error("[UpdateWarehouse] failed to get warehouse by id", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	return mapper.WarehouseToResponse(&warehouse), nil
}

func (s *WarehouseService) CreateWarehouse(request dto.CreateWarehouseRequest, createdBy uuid.UUID) (dto.WarehouseResponse, error) {
	s.repository.UseTx(false)

	warehouse := entity.Warehouse{
		Name:       request.Name,
		LocationId: request.LocationId,
		CreatedBy:  uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateWarehouse(&warehouse)
	if err != nil {
		s.log.Error("[CreateWarehouse] failed to create warehouse", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	warehouse, err = s.repository.GetWarehouseById(warehouse.Id)
	if err != nil {
		s.log.Error("[CreateWarehouse] failed to get warehouse by id", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	return mapper.WarehouseToResponse(&warehouse), nil
}

func (s *WarehouseService) GetWarehouses(filter dto.GetWarehouseFilter) ([]dto.WarehouseResponse, error) {
	s.repository.UseTx(false)

	warehouses, err := s.repository.GetWarehouses(filter)
	if err != nil {
		s.log.Error("[GetWarehouses] failed to get warehouses", zap.Error(err))
		return nil, err
	}

	warehouseResponses := make([]dto.WarehouseResponse, 0, len(warehouses))
	for _, warehouse := range warehouses {
		warehouseResponses = append(warehouseResponses, mapper.WarehouseToResponse(&warehouse))
	}

	return warehouseResponses, nil
}

func (s *WarehouseService) DeleteWarehouse(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouse(id)
	if err != nil {
		s.log.Error("[DeleteWarehouse] failed to delete warehouse", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseStockItem(request *dto.CreateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error) {
	s.repository.UseTx(false)

	// Todo : create estimation run out date, based on average used per day from request item from request warhouse item.

	stockWarehouseItem := entity.WarehouseItem{
		WarehouseId:      request.WarehouseId,
		ItemId:           request.WarehouseItemId,
		Quantity:         request.Quantity,
		EstimationRunOut: time.Now(),
		CreatedBy:        uuid.NullUUID{UUID: accountId, Valid: true},
	}

	err := s.repository.CreateWarehouseStockItem(&stockWarehouseItem)
	if err != nil {
		s.log.Error("[CreateStockWarehouseItem] failed to create stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	stockWarehouseItem, err = s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(
		stockWarehouseItem.WarehouseId,
		stockWarehouseItem.ItemId,
	)
	if err != nil {
		s.log.Error("[CreateStockWarehouseItem] failed to get stock warehouse item", zap.Error(err))
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

func (s *WarehouseService) GetWarehouseStockItems(filter dto.GetWarehouseStockItemFilter) ([]dto.WarehouseStockItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItems, err := s.repository.GetWarehouseStockItems(filter)
	if err != nil {
		s.log.Error("[GetStockWarehouseItem] failed to get stock warehouse items", zap.Error(err))
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

func (s *WarehouseService) GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseStockItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItem, err := s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("[GetStockWarehouseItemByWarehouseIdAndWarehouseItemId] failed to get stock warehouse item", zap.Error(err))
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

func (s *WarehouseService) UpdateWarehouseStockItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseStockItemRequest, accountId uuid.UUID) (dto.WarehouseStockItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItem, err := s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("[UpdateStockWarehouseItem] failed to get stock warehouse item", zap.Error(err))
		return dto.WarehouseStockItemResponse{}, err
	}

	stockWarehouseItem.Quantity = request.Quantity
	stockWarehouseItem.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateWarehouseStockItem(&stockWarehouseItem)
	if err != nil {
		s.log.Error("[UpdateStockWarehouseItem] failed to update stock warehouse item", zap.Error(err))
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

func (s *WarehouseService) DeleteWarehouseStockItem(warehouseId uint64, warehouseItemId uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouseStockItemByWarehouseIdAndWarehouseItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("[DeleteStockWarehouseItem] failed to delete stock warehouse item", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseOrderItem(request dto.CreateWarehouseOrderItemRequest, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error) {
	s.repository.UseTx(false)

	warehouseOrderItem := entity.WarehouseOrderItem{
		WarehouseId: request.WarehouseId,
		SupplierId:  request.SupplierId,
		ItemId:      request.WarehouseItemId,
		Quantity:    request.Quantity,
		Status:      enum.WarehouseOrderStatusInSend,
		TakenAt:     sql.NullTime{},
		CreatedBy:   uuid.NullUUID{UUID: accountId, Valid: true},
	}

	err := s.repository.CreateWarehouseOrderItem(&warehouseOrderItem)
	if err != nil {
		s.log.Error("[CreateWarehouseOrderItem] failed to create warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	warehouseOrderItem, err = s.repository.GetWarehouseOrderItemById(warehouseOrderItem.Id)
	if err != nil {
		s.log.Error("[CreateWarehouseOrderItem] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (s *WarehouseService) GetWarehouseOrderItemById(id uint64) (dto.WarehouseOrderItemResponse, error) {
	s.repository.UseTx(false)

	warehouseOrderItem, err := s.repository.GetWarehouseOrderItemById(id)
	if err != nil {
		s.log.Error("[GetWarehouseOrderItemById] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (s *WarehouseService) GetWarehouseOrderItems(filter dto.GetWarehouseOrderItemFilter) ([]dto.WarehouseOrderItemResponse, error) {
	s.repository.UseTx(false)

	filter.IsTaken = true
	warehouseOrderItems, err := s.repository.GetWarehouseOrderItems(filter)
	if err != nil {
		s.log.Error("[GetWarehouseOrderItems] failed to get warehouse order items", zap.Error(err))
		return nil, err
	}

	if filter.Date.Value().Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)) {
		untakenWarehouseOrderItems, err := s.repository.GetWarehouseOrderItems(dto.GetWarehouseOrderItemFilter{IsTaken: false})
		if err != nil {
			s.log.Error("[GetWarehouseOrderItems] failed to get warehouse order items", zap.Error(err))
			return nil, err
		}

		warehouseOrderItems = append(warehouseOrderItems, untakenWarehouseOrderItems...)
	}

	warehouseOrderItemResponses := make([]dto.WarehouseOrderItemResponse, 0, len(warehouseOrderItems))
	for _, item := range warehouseOrderItems {
		warehouseOrderItemResponses = append(warehouseOrderItemResponses, mapper.WarehouseOrderItemToResponse(&item))
	}

	return warehouseOrderItemResponses, nil
}

func (s *WarehouseService) DeleteWarehouseOrderItem(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteWarehouseOrderItem(id)
	if err != nil {
		s.log.Error("[DeleteWarehouseOrderItem] failed to delete warehouse order item", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) TakeWarehouseOrderItem(id uint64, accountId uuid.UUID) (dto.WarehouseOrderItemResponse, error) {
	s.repository.UseTx(false)

	// Todo : add stock warehouse item in warehouse
	warehouseOrderItem, err := s.repository.GetWarehouseOrderItemById(id)
	if err != nil {
		s.log.Error("[TakeWarehouseOrderItem] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	if warehouseOrderItem.IsTaken.Bool {
		s.log.Error("[TakeWarehouseOrderItem] warehouse order item already taken", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, errx.BadRequest("warehouse order item already taken")
	}

	warehouseOrderItem.IsTaken = sql.NullBool{Bool: true, Valid: true}
	warehouseOrderItem.TakenBy = uuid.NullUUID{UUID: accountId, Valid: true}
	warehouseOrderItem.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	warehouseOrderItem.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateWarehouseOrderItem(&warehouseOrderItem)
	if err != nil {
		s.log.Error("[TakeWarehouseOrderItem] failed to update warehouse order item", zap.Error(err))
		return dto.WarehouseOrderItemResponse{}, err
	}

	s.cacheService.Publish(context.Background(), constant.TopicWarehouseActivity,
		entity.WarehouseActivity{
			WarehouseId: warehouseOrderItem.WarehouseId,
			Description: "Pesanan barang dari supplier " + warehouseOrderItem.Supplier.Name + " telah diambil",
			Status:      enum.ActivityStatusIn,
		},
	)

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (e *WarehouseService) GoodEggConvertionButirToIkat(request dto.GoodEggWarehouseConvertionRequest, accountId uuid.UUID) ([]dto.WarehouseStockItemResponse, error) {
	e.repository.UseTx(true)
	defer e.repository.Rollback()

	goodEggButir, err := e.repository.GetWarehouseItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitButir, enum.WarehouseItemCategoryEgg)
	if err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse item", zap.Error(err))
		return nil, err
	}

	goodEggIkat, err := e.repository.GetWarehouseItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitIkat, enum.WarehouseItemCategoryEgg)
	if err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggButir, err := e.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(request.WarehouseId, goodEggButir.Id)
	if err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse stock item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggIkat, err := e.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(request.WarehouseId, goodEggIkat.Id)
	if err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse stock item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggButir.Quantity = warehouseStockItemEggButir.Quantity - request.TotalButir - (request.TotalKarpet * constant.TotalEggKarpet)

	if warehouseStockItemEggButir.Quantity < 0 {
		return nil, errx.BadRequest("stok butir tidak mencukupi")
	}

	warehouseStockItemEggIkat.Quantity = warehouseStockItemEggIkat.Quantity + request.TotalIkat

	warehouseStockItemEggButir.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}
	warehouseStockItemEggIkat.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = e.repository.UpdateWarehouseStockItem(&warehouseStockItemEggButir)
	if err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to update warehouse stock item", zap.Error(err))

		return nil, err
	}

	err = e.repository.UpdateWarehouseStockItem(&warehouseStockItemEggIkat)
	if err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to update warehouse stock item", zap.Error(err))
		return nil, err
	}

	if err := e.repository.Commit(); err != nil {
		e.log.Error("[GoodEggConverterButirToIkat] failed to commit transaction", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehouseStockItemResponse, 0)

	response = append(response, mapper.WarehouseStockItemToResponse(&warehouseStockItemEggButir))
	response = append(response, mapper.WarehouseStockItemToResponse(&warehouseStockItemEggIkat))

	return response, nil
}

func (s *WarehouseService) GoodEggConvertionIkatToButir(request dto.GoodEggWarehouseConvertionRequest, accountId uuid.UUID) ([]dto.WarehouseStockItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	// Todo : change it into GetWarehouseItemByFilter
	goodEggButir, err := s.repository.GetWarehouseItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitButir, enum.WarehouseItemCategoryEgg)
	if err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse item", zap.Error(err))
		return nil, err
	}

	goodEggIkat, err := s.repository.GetWarehouseItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitIkat, enum.WarehouseItemCategoryEgg)
	if err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggButir, err := s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(request.WarehouseId, goodEggButir.Id)
	if err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse stock item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggIkat, err := s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(request.WarehouseId, goodEggIkat.Id)
	if err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to get warehouse stock item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggIkat.Quantity = warehouseStockItemEggIkat.Quantity - request.TotalIkat

	if warehouseStockItemEggIkat.Quantity < 0 {
		return nil, errx.BadRequest("stok ikat tidak mencukupi")
	}

	warehouseStockItemEggIkat.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	warehouseStockItemEggButir.Quantity = warehouseStockItemEggButir.Quantity + request.TotalButir + (request.TotalKarpet * constant.TotalEggKarpet)
	warehouseStockItemEggButir.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateWarehouseStockItem(&warehouseStockItemEggButir)
	if err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to update warehouse stock item", zap.Error(err))

		return nil, err
	}

	err = s.repository.UpdateWarehouseStockItem(&warehouseStockItemEggIkat)
	if err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to update warehouse stock item", zap.Error(err))
		return nil, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[GoodEggConverterButirToIkat] failed to commit transaction", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehouseStockItemResponse, 0)

	response = append(response, mapper.WarehouseStockItemToResponse(&warehouseStockItemEggButir))
	response = append(response, mapper.WarehouseStockItemToResponse(&warehouseStockItemEggIkat))

	return response, nil
}

func (s *WarehouseService) CrackedEggConvertionButirToPack(request dto.CrackedEggWarehouseConvertionRequest, accountId uuid.UUID) ([]dto.WarehouseStockItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	crackedEggButir, err := s.repository.GetWarehouseItemByNameAndUnitAndType(constant.CrackedEgg, constant.EggUnitButir, enum.WarehouseItemCategoryEgg)
	if err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to get warehouse item", zap.Error(err))
		return nil, err
	}

	crackedEggPack, err := s.repository.GetWarehouseItemByNameAndUnitAndType(constant.CrackedEgg, constant.EggUnitPlastik, enum.WarehouseItemCategoryEgg)
	if err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to get warehouse item", zap.Error(err))
		return nil, err
	}

	// Todo : change it into GetWarehouseItemByFilter
	warehouseStockItemEggButir, err := s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(request.WarehouseId, crackedEggButir.Id)
	if err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to get warehouse stock item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggPack, err := s.repository.GetWarehouseStockItemByWarehouseIdAndWarehouseItemId(request.WarehouseId, crackedEggPack.Id)
	if err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to get warehouse stock item", zap.Error(err))
		return nil, err
	}

	warehouseStockItemEggButir.Quantity = warehouseStockItemEggButir.Quantity - request.TotalButir
	warehouseStockItemEggButir.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	if warehouseStockItemEggButir.Quantity < 0 {
		return nil, errx.BadRequest("stok butir tidak mencukupi")
	}

	warehouseStockItemEggPack.Quantity = warehouseStockItemEggPack.Quantity + request.TotalPack
	warehouseStockItemEggPack.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateWarehouseStockItem(&warehouseStockItemEggButir)
	if err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to update warehouse stock item", zap.Error(err))
		return nil, err
	}

	err = s.repository.UpdateWarehouseStockItem(&warehouseStockItemEggPack)
	if err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to update warehouse stock item", zap.Error(err))
		return nil, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[CrackedEggConverterButirToPacket] failed to commit transaction", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehouseStockItemResponse, 0)
	response = append(response, mapper.WarehouseStockItemToResponse(&warehouseStockItemEggButir))
	response = append(response, mapper.WarehouseStockItemToResponse(&warehouseStockItemEggPack))

	return response, nil
}
