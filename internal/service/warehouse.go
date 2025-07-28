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
	GetWarehouses(filter dto.GetWarehouseFilter) ([]dto.WarehouseResponse, error)
	CreateWarehouse(request dto.CreateWarehouseRequest, createdBy uuid.UUID) (dto.WarehouseResponse, error)
	DeleteWarehouse(id uint64) error
	UpdateWarehouse(id uint64, request dto.UpdateWarehouseRequest, createdBy uuid.UUID) (dto.WarehouseResponse, error)
	GetWarehouseDetailById(id uint64) (dto.WarehouseDetailResponse, error)
	GetWarehouseOverview(id uint64) (dto.WarehouseOverview, error)

	CreateWarehouseItem(request dto.CreateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error)
	GetWarehouseItems(filter dto.GetWarehouseItemFilter) ([]dto.WarehouseItemResponse, error)
	GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, itemId uint64) (dto.WarehouseItemResponse, error)
	UpdateWarehouseItem(warehouseId uint64, itemId uint64, request dto.UpdateWarehouseItemRequest, updatedBy uuid.UUID) (dto.WarehouseItemResponse, error)
	DeleteWarehouseItem(warehouseId uint64, itemId uint64) error
	GetEggWarehouseItemSummary(warehouseId uint64) ([]dto.EggWarehouseItemSummary, error)

	CreateWarehouseOrderItem(request dto.CreateWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItempProcurementResponse, error)
	GetWarehouseOrderItemById(id uint64) (dto.WarehouseItempProcurementResponse, error)
	GetWarehouseOrderItems(filter dto.GetWarehouseItemProcurementFilter) ([]dto.WarehouseItempProcurementResponse, error)
	DeleteWarehouseOrderItem(id uint64) error
	TakeWarehouseOrderItem(id uint64, userId uuid.UUID) (dto.WarehouseItempProcurementResponse, error)

	GetWarehouseItemHistories(filter dto.GetWarehouseItemHistoryFilter) (dto.WarehouseItemHistoryListPaginationResponse, error)
	GetWarehouseItemHistoryById(id uint64) (dto.WarehouseItemHistoryResponse, error)

	CreateWarehouseSale(request dto.CreateWarehouseSaleRequest, createdBy uuid.UUID) (dto.WarehouseSaleResponse, error)
	GetWarehouseSaleById(id uint64) (dto.WarehouseSaleResponse, error)
	GetWarehouseSales(filter dto.GetWarehouseSaleFilter) (dto.WarehouseSaleListPaginationResponse, error)
	UpdateWarehouseSale(id uint64, request dto.UpdateWarehouseSaleRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	DeleteWarehouseSale(id uint64, userId uuid.UUID) error

	CreateWarehouseSalePayment(warehouseSaleId uint64, request dto.CreateWarehouseSalePaymentRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	UpdateWarehouseSalePayment(warehouseSaleId uint64, id uint64, request dto.UpdateWarehouseSalePaymentRequest, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
	DeleteWarehouseSalePayment(warehouseSaleId uint64, id uint64, userId uuid.UUID) error

	SendWarehouseSale(id uint64, userId uuid.UUID) (dto.WarehouseSaleResponse, error)
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

func (s *WarehouseService) GetWarehouseDetailById(id uint64) (dto.WarehouseDetailResponse, error) {
	s.repository.UseTx(false)

	warehouse, err := s.repository.GetWarehouseById(id)
	if err != nil {
		s.log.Error("failed to get warehouse by id")
		return dto.WarehouseDetailResponse{}, err
	}

	warehousePlacements, err := s.placementService.GetWarehousePlacementByWarehouseId(id)
	if err != nil {
		return dto.WarehouseDetailResponse{}, err
	}

	userResponses := make([]dto.UserResponse, 0)
	for _, e := range warehousePlacements {
		userResponses = append(userResponses, e.User)
	}

	return dto.WarehouseDetailResponse{
		Id:       warehouse.Id,
		Name:     warehouse.Name,
		Location: mapper.LocationToResponse(&warehouse.Location),
		Users:    userResponses,
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

func (s *WarehouseService) CreateWarehouse(request dto.CreateWarehouseRequest, createdBy uuid.UUID) (dto.WarehouseResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	warehouse := entity.Warehouse{
		Name:       request.Name,
		LocationId: request.LocationId,
		CreatedBy:  uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateWarehouse(&warehouse)
	if err != nil {
		s.log.Error("failed to create warehouse", zap.Error(err))
		return dto.WarehouseResponse{}, err
	}

	goodEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.GoodEgg, constant.EggUnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.WarehouseResponse{}, err
	}

	crackedEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.CrackedEgg, constant.EggUnitKg, enum.ItemCategoryEgg)
	if err != nil {
		return dto.WarehouseResponse{}, err

	}

	brokenEggItem, err := s.itemService.GetItemByNameAndUnitAndType(constant.BrokenEgg, constant.EggUnitPlastik, enum.ItemCategoryEgg)
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

func (s *WarehouseService) GetWarehouses(filter dto.GetWarehouseFilter) ([]dto.WarehouseResponse, error) {
	s.repository.UseTx(false)

	warehouses, err := s.repository.GetWarehouses(filter)
	if err != nil {
		s.log.Error("failed to get warehouses", zap.Error(err))
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
		s.log.Error("failed to delete warehouse", zap.Error(err))
		return err
	}

	return nil
}

func (s *WarehouseService) CreateWarehouseItem(request dto.CreateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	// Todo : create estimation run out date, based on average used per day from request item from request warhouse item.
	// Todo : fix run out item

	stockWarehouseItem := entity.WarehouseItem{
		WarehouseId: request.WarehouseId,
		ItemId:      request.ItemId,
		Quantity:    request.Quantity,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.RunOutCountDown != nil {
		stockWarehouseItem.EstimationRunOut = sql.NullTime{
			Time:  time.Now().Add(time.Hour * 24 * time.Duration(*request.RunOutCountDown)),
			Valid: true,
		}
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

	stockWarehouseItemResponses := make([]dto.WarehouseItemResponse, 0, len(warehouseItems))
	for _, item := range warehouseItems {
		stockWarehouseItemResponses = append(stockWarehouseItemResponses, mapper.WarehouseItemToResponse(&item))
	}

	return stockWarehouseItemResponses, nil
}

func (s *WarehouseService) GetWarehouseItemByWarehouseIdAndItemId(warehouseId uint64, warehouseItemId uint64) (dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item by warehouse id and item id", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	var description string
	if stockWarehouseItem.EstimationRunOut.Valid && time.Now().Add(time.Hour*24*7).After(stockWarehouseItem.EstimationRunOut.Time) {
		description = constant.StockWarehouseItemDanger
	} else {
		description = constant.StockWarehouseItemSafe
	}

	warehouseStockItemResponse := mapper.WarehouseItemToResponse(&stockWarehouseItem)
	warehouseStockItemResponse.Description = description

	return warehouseStockItemResponse, nil
}

func (s *WarehouseService) UpdateWarehouseItem(warehouseId uint64, warehouseItemId uint64, request dto.UpdateWarehouseItemRequest, userId uuid.UUID) (dto.WarehouseItemResponse, error) {
	s.repository.UseTx(false)

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseId, warehouseItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	warehouseItem.Quantity = request.Quantity
	// Todo : Fix this run out time
	if request.RunOutCountDown != nil {
		warehouseItem.EstimationRunOut = sql.NullTime{
			Time:  time.Now().Add(time.Hour * 24 * time.Duration(*request.RunOutCountDown)),
			Valid: true,
		}
	}
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update warehouse item", zap.Error(err))
		return dto.WarehouseItemResponse{}, err
	}

	var description string
	if warehouseItem.EstimationRunOut.Valid && time.Now().Add(time.Hour*24*7).After(warehouseItem.EstimationRunOut.Time) {
		description = constant.StockWarehouseItemDanger
	} else {
		description = constant.StockWarehouseItemSafe
	}

	warehouseStockItemResponse := mapper.WarehouseItemToResponse(&warehouseItem)
	warehouseStockItemResponse.Description = description

	return warehouseStockItemResponse, nil
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

	eggWarehouseItems := make([]dto.WarehouseItemResponse, 0)
	equipmentWarehouseItems := make([]dto.WarehouseItemResponse, 0)

	totalSafeStock := 0
	totalDangerStock := 0

	for _, e := range warehouseItems {
		res := mapper.WarehouseItemToResponse(&e)
		switch res.Description {
		case constant.StockWarehouseItemDanger:
			totalDangerStock++
		case constant.StockWarehouseItemSafe:
			totalSafeStock++
		}

		if e.Item.Category == enum.ItemCategoryEgg {
			eggWarehouseItems = append(eggWarehouseItems, res)
		} else {
			equipmentWarehouseItems = append(equipmentWarehouseItems, res)
		}
	}

	requestItemCount, err := s.repository.CountStoreRequestItemByWarehouseId(id)
	if err != nil {
		s.log.Error("failed to count store request item by warehouse id", zap.Error(err))
		return dto.WarehouseOverview{}, err
	}

	return dto.WarehouseOverview{
		EggStocks:         eggWarehouseItems,
		EquipmentStocks:   equipmentWarehouseItems,
		TotalSafeStock:    uint64(totalSafeStock),
		TotalDangerStock:  uint64(totalDangerStock),
		TotalStoreRequest: uint64(requestItemCount),
	}, nil
}

func (s *WarehouseService) CreateWarehouseOrderItem(request dto.CreateWarehouseItemProcurementRequest, userId uuid.UUID) (dto.WarehouseItempProcurementResponse, error) {
	s.repository.UseTx(false)

	warehouseOrderItem := entity.WarehouseItemProcurement{
		WarehouseId: request.WarehouseId,
		SupplierId:  request.SupplierId,
		ItemId:      request.WarehouseItemId,
		Quantity:    request.Quantity,
		Status:      enum.WarehouseOrderStatusInSend,
		TakenAt:     sql.NullTime{},
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateWarehouseOrderItem(&warehouseOrderItem)
	if err != nil {
		s.log.Error("[CreateWarehouseOrderItem] failed to create warehouse order item", zap.Error(err))
		return dto.WarehouseItempProcurementResponse{}, err
	}

	warehouseOrderItem, err = s.repository.GetWarehouseOrderItemById(warehouseOrderItem.Id)
	if err != nil {
		s.log.Error("[CreateWarehouseOrderItem] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseItempProcurementResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (s *WarehouseService) GetWarehouseOrderItemById(id uint64) (dto.WarehouseItempProcurementResponse, error) {
	s.repository.UseTx(false)

	warehouseOrderItem, err := s.repository.GetWarehouseOrderItemById(id)
	if err != nil {
		s.log.Error("[GetWarehouseOrderItemById] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseItempProcurementResponse{}, err
	}

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
}

func (s *WarehouseService) GetWarehouseOrderItems(filter dto.GetWarehouseItemProcurementFilter) ([]dto.WarehouseItempProcurementResponse, error) {
	s.repository.UseTx(false)

	filter.IsTaken = true
	warehouseOrderItems, err := s.repository.GetWarehouseOrderItems(filter)
	if err != nil {
		s.log.Error("[GetWarehouseOrderItems] failed to get warehouse order items", zap.Error(err))
		return nil, err
	}

	if filter.Date.Value().Equal(time.Date(time.Now().Year(), time.Now().Month(), time.Now().Day(), 0, 0, 0, 0, time.Local)) {
		untakenWarehouseOrderItems, err := s.repository.GetWarehouseOrderItems(dto.GetWarehouseItemProcurementFilter{IsTaken: false})
		if err != nil {
			s.log.Error("[GetWarehouseOrderItems] failed to get warehouse order items", zap.Error(err))
			return nil, err
		}

		warehouseOrderItems = append(warehouseOrderItems, untakenWarehouseOrderItems...)
	}

	warehouseOrderItemResponses := make([]dto.WarehouseItempProcurementResponse, 0, len(warehouseOrderItems))
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

func (s *WarehouseService) TakeWarehouseOrderItem(id uint64, userId uuid.UUID) (dto.WarehouseItempProcurementResponse, error) {
	s.repository.UseTx(false)

	// Todo : add stock warehouse item in warehouse
	warehouseOrderItem, err := s.repository.GetWarehouseOrderItemById(id)
	if err != nil {
		s.log.Error("[TakeWarehouseOrderItem] failed to get warehouse order item", zap.Error(err))
		return dto.WarehouseItempProcurementResponse{}, err
	}

	if warehouseOrderItem.IsTaken.Bool {
		s.log.Error("[TakeWarehouseOrderItem] warehouse order item already taken", zap.Error(err))
		return dto.WarehouseItempProcurementResponse{}, errx.BadRequest("warehouse order item already taken")
	}

	warehouseOrderItem.IsTaken = sql.NullBool{Bool: true, Valid: true}
	warehouseOrderItem.TakenBy = uuid.NullUUID{UUID: userId, Valid: true}
	warehouseOrderItem.TakenAt = sql.NullTime{Time: time.Now(), Valid: true}
	warehouseOrderItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseOrderItem(&warehouseOrderItem)
	if err != nil {
		s.log.Error("[TakeWarehouseOrderItem] failed to update warehouse order item", zap.Error(err))
		return dto.WarehouseItempProcurementResponse{}, err
	}

	// s.cacheService.Publish(context.Background(), constant.TopicWarehouseActivity,
	// 	entity.WarehouseItemActivity{
	// 		WarehouseId: warehouseOrderItem.WarehouseId,
	// 		Description: "Pesanan barang dari supplier " + warehouseOrderItem.Supplier.Name + " telah diambil",
	// 		Status:      enum.ActivityStatusIn,
	// 	},
	// )

	return mapper.WarehouseOrderItemToResponse(&warehouseOrderItem), nil
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

func (s *WarehouseService) GetEggWarehouseItemSummary(warehouseId uint64) ([]dto.EggWarehouseItemSummary, error) {
	s.repository.UseTx(false)

	response := make([]dto.EggWarehouseItemSummary, 0)
	warehouseItems, err := s.repository.GetWarehouseItems(dto.GetWarehouseItemFilter{
		WarehouseId: warehouseId,
		ItemNames:   []string{constant.GoodEgg, constant.CrackedEgg},
		Units:       []string{constant.EggUnitKg},
	})
	if err != nil {
		s.log.Error("failed to get warehouse items", zap.Error(err))
		return nil, err
	}

	for _, warehouseItem := range warehouseItems {
		switch warehouseItem.Item.Name {
		case constant.GoodEgg:
			response = append(response, dto.EggWarehouseItemSummary{
				Name:     constant.GoodEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.EggUnitKg,
			})

			response = append(response, dto.EggWarehouseItemSummary{
				Name:     constant.GoodEgg,
				Quantity: warehouseItem.Quantity / float64(constant.TotalEggPerIkat),
				Unit:     constant.EggUnitIkat,
			})
		case constant.CrackedEgg:
			response = append(response, dto.EggWarehouseItemSummary{
				Name:     constant.CrackedEgg,
				Quantity: warehouseItem.Quantity,
				Unit:     constant.EggUnitKg,
			})

			response = append(response, dto.EggWarehouseItemSummary{
				Name:     constant.CrackedEgg,
				Quantity: warehouseItem.Quantity / float64(constant.TotalEggPerIkat),
				Unit:     constant.EggUnitIkat,
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

	warehouseItem.Quantity -= request.Quantity
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

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		s.log.Error("invalid sale unit", zap.String("saleUnit", request.SaleUnit))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	warehouseSale := entity.WarehouseSale{
		WarehouseId: request.WarehouseId,
		ItemId:      request.ItemId,
		Quantity:    request.Quantity,
		Price:       price,
		TotalPrice:  totalPrice,
		SendDate:    sendDate,
		IsSend:      false,
		SaleUnit:    saleUnit,
		PaymentType: paymentType,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
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

		if request.CustomerPhoneNumber[:2] != "08" {
			return dto.WarehouseSaleResponse{}, errx.BadRequest("customer phone number must be in valid format 08")
		}

		resp, err := s.customerService.CreateCustomer(customer)
		if err != nil {
			return dto.WarehouseSaleResponse{}, err
		}

		warehouseSale.CustomerId = resp.Id
	}

	nominal, err := decimal.NewFromString(request.WarehouseSalePayment.Nominal)
	if err != nil {
		s.log.Error("failed to parse nominal", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !warehouseSale.TotalPrice.Equal(nominal) {
			s.log.Error("nominal is not equal to total price", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}

		warehouseSale.PaymentStatus = enum.PaymentStatusPaid
	} else {
		warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	err = s.repository.CreateWarehouseSale(&warehouseSale)
	if err != nil {
		s.log.Error("failed to create warehouse sale", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	if request.WarehouseSalePayment.Nominal != "" &&
		request.WarehouseSalePayment.PaymentDate != "" &&
		request.WarehouseSalePayment.PaymentProof != "" &&
		request.WarehouseSalePayment.PaymentMethod != "" {
		paymentMethod := enum.ValueOfPaymentMethod(request.WarehouseSalePayment.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("invalid payment method", zap.String("paymentMethod", request.WarehouseSalePayment.PaymentMethod))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment method")
		}

		paymentDate, err := time.Parse("02-01-2006", request.WarehouseSalePayment.PaymentDate)
		if err != nil {
			s.log.Error("failed to parse payment date", zap.Error(err))
			return dto.WarehouseSaleResponse{}, errx.BadRequest("invalid payment date format")
		}

		warehouseSalePayment := entity.WarehouseSalePayment{
			PaymentDate:     paymentDate,
			WarehouseSaleId: warehouseSale.Id,
			Nominal:         nominal,
			PaymentProof:    request.WarehouseSalePayment.PaymentProof,
			PaymentMethod:   paymentMethod,
			CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
		}

		err = s.repository.CreateWarehouseSalePayment(&warehouseSalePayment)
		if err != nil {
			s.log.Error("failed to create warehouse sale payment", zap.Error(err))
			return dto.WarehouseSaleResponse{}, err
		}
	}

	warehouseSale, err = s.repository.GetWarehouseSaleById(warehouseSale.Id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
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
			PaymentMethod: filter.PaymentMethod,
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

	if warehouseSale.IsSend {
		s.log.Error("warehouse sale is already sent", zap.Uint64("id", id))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("warehouse sale is already sent")
	}

	warehouseItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(warehouseSale.WarehouseId, warehouseSale.ItemId)
	if err != nil {
		s.log.Error("failed to get store item by store id and item id", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseItem.Quantity -= warehouseSale.Quantity + request.Quantity
	warehouseItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&warehouseItem)
	if err != nil {
		s.log.Error("failed to update store item", zap.Error(err))
		return dto.WarehouseSaleResponse{}, err
	}

	warehouseSale.Quantity = request.Quantity
	totalPrice := warehouseSale.Price.Mul(decimal.NewFromFloat(request.Quantity))
	discountPrice := totalPrice.Mul(decimal.NewFromFloat(warehouseSale.Discount / 100.0))
	warehouseSale.TotalPrice = totalPrice.Sub(discountPrice)

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

	if warehouseSale.IsSend {
		s.log.Error("warehouse sale is already sent", zap.Uint64("id", warehouseSale.Id))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("warehouse sale is already sent")
	}

	if warehouseSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("warehouse sale is already paid", zap.Uint64("id", warehouseSale.Id))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("warehouse sale is already paid")
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
	} else if totalPayment.GreaterThan(warehouseSale.TotalPrice) {
		s.log.Error("total payment is greater than total price", zap.Error(err))
		return dto.WarehouseSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	} else if totalPayment.LessThan(warehouseSale.TotalPrice) {
		warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
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
	storeSale, err := s.repository.GetWarehouseSaleById(id)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return err
	}

	if storeSale.IsSend {
		s.log.Error("warehouse sale is already sent", zap.Uint64("id", id))
		return errx.BadRequest("store sale already send")
	}

	storeItem, err := s.repository.GetWarehouseItemByWarehouseIdAndItemId(storeSale.WarehouseId, storeSale.ItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item by store id and item id", zap.Error(err))
		return err
	}

	storeItem.Quantity += storeSale.Quantity
	storeItem.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateWarehouseItem(&storeItem)
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
	s.repository.UseTx(false)

	warehouseSale, err := s.repository.GetWarehouseSaleById(warehouseSaleId)
	if err != nil {
		s.log.Error("failed to get warehouse sale by id", zap.Error(err))
		return err
	}

	if warehouseSale.IsSend {
		s.log.Error("warehouse sale is already sent", zap.Uint64("id", warehouseSale.Id))
		return errx.BadRequest("warehouse sale is already sent")
	}

	if warehouseSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("warehouse sale is already paid", zap.Uint64("id", warehouseSale.Id))
		return errx.BadRequest("warehouse sale is already paid")
	}

	totalPayment := decimal.Zero
	for _, payment := range warehouseSale.Payments {
		if payment.Id != id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.LessThan(warehouseSale.TotalPrice) && totalPayment.GreaterThan(decimal.Zero) {
		warehouseSale.PaymentStatus = enum.PaymentStatusUnpaid
		warehouseSale.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}
	} else if totalPayment.LessThan(decimal.Zero) {
		s.log.Error("delete this payment make minus", zap.Error(err))
		return errx.BadRequest("delete this payment make minus")
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

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("failed to parse send date", zap.Error(err))
		return dto.WarehouseSaleQueueResponse{}, err
	}

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
		SendDate:     sendDate,
		CustomerType: customerType,
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

	err = s.repository.CreateWarehouseSaleQueue(&data)
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

func (s *WarehouseService) GetWarehouseSaleQueues() ([]dto.WarehouseSaleQueueResponse, error) {
	s.repository.UseTx(false)

	// Todo : formula for integrated planning
	WarehouseSaleQueues, err := s.repository.GetWarehouseSaleQueues()
	if err != nil {
		s.log.Error("failed get Warehouse sale queues", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehouseSaleQueueResponse, 0)
	for _, WarehouseSaleQueue := range WarehouseSaleQueues {
		response = append(response, mapper.WarehouseSaleQueueToResponse(&WarehouseSaleQueue))
	}

	return response, nil
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
