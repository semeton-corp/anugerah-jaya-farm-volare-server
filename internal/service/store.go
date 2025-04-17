package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type StoreService struct {
	log        *zap.Logger
	repository repository.IStoreRepository
}

type IStoreService interface {
	GetStores() ([]dto.StoreResponse, error)

	CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]dto.StoreRequestItemResponse, error)
	UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error)
}

func NewStoreService(log *zap.Logger, repository repository.IStoreRepository) IStoreService {
	return &StoreService{
		log:        log,
		repository: repository,
	}
}

func (s *StoreService) GetStores() ([]dto.StoreResponse, error) {
	stores, err := s.repository.GetStores()
	if err != nil {
		s.log.Error("[GetStores] failed to get stores", zap.Error(err))
		return nil, err
	}

	storeResponses := make([]dto.StoreResponse, len(stores))
	for i, store := range stores {
		storeResponses[i] = dto.StoreResponse{
			Id:   store.Id,
			Name: store.Name,
			Location: dto.LocationResponse{
				Id:   store.Location.Id,
				Name: store.Location.Name,
			},
		}
	}

	return storeResponses, nil
}

func (s *StoreService) CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error) {
	// Todo : Check if warehouse is hava warehouse item

	storeRequestItem := entity.StoreRequestItem{
		WarehouseId:     request.WarehouseId,
		WarehouseItemId: request.WarehouseItemId,
		StoreId:         request.StoreId,
		Quantity:        request.Quantity,
		Status:          enum.RequestItemStatusPending,
		CreatedBy:       accountId,
	}

	err := s.repository.CreateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("[CreateStoreRequestItem] failed to create store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	storeRequestItem, err = s.repository.GetStoreRequestItemById(storeRequestItem.Id)
	if err != nil {
		s.log.Error("[CreateStoreRequestItem] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return dto.StoreRequestItemResponse{
		Id: storeRequestItem.Id,
		Warehouse: dto.WarehouseResponse{
			Id:   storeRequestItem.Warehouse.Id,
			Name: storeRequestItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Warehouse.Location.Id,
				Name: storeRequestItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeRequestItem.WarehouseItem.Id,
			Name:     storeRequestItem.WarehouseItem.Name,
			Category: storeRequestItem.WarehouseItem.Category.String(),
			Unit:     storeRequestItem.WarehouseItem.Unit,
		},
		Store: dto.StoreResponse{
			Id:   storeRequestItem.Store.Id,
			Name: storeRequestItem.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Store.Location.Id,
				Name: storeRequestItem.Store.Location.Name,
			},
		},
		Quantity: storeRequestItem.Quantity,
		Status:   storeRequestItem.Status.String(),
	}, nil
}

func (s *StoreService) GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error) {
	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("[GetStoreRequestItemById] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return dto.StoreRequestItemResponse{
		Id: storeRequestItem.Id,
		Warehouse: dto.WarehouseResponse{
			Id:   storeRequestItem.Warehouse.Id,
			Name: storeRequestItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Warehouse.Location.Id,
				Name: storeRequestItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeRequestItem.WarehouseItem.Id,
			Name:     storeRequestItem.WarehouseItem.Name,
			Category: storeRequestItem.WarehouseItem.Category.String(),
			Unit:     storeRequestItem.WarehouseItem.Unit,
		},
		Store: dto.StoreResponse{
			Id:   storeRequestItem.Store.Id,
			Name: storeRequestItem.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Store.Location.Id,
				Name: storeRequestItem.Store.Location.Name,
			},
		},
		Quantity: storeRequestItem.Quantity,
		Status:   storeRequestItem.Status.String(),
	}, nil
}

func (s *StoreService) GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]dto.StoreRequestItemResponse, error) {
	storeRequestItems, err := s.repository.GetStoreRequestItems(filter)
	if err != nil {
		s.log.Error("[GetStoreRequestItems] failed to get store request items", zap.Error(err))
		return nil, err
	}

	storeRequestItemResponses := make([]dto.StoreRequestItemResponse, len(storeRequestItems))
	for i, storeRequestItem := range storeRequestItems {
		storeRequestItemResponses[i] = dto.StoreRequestItemResponse{
			Id: storeRequestItem.Id,
			Warehouse: dto.WarehouseResponse{
				Id:   storeRequestItem.Warehouse.Id,
				Name: storeRequestItem.Warehouse.Name,
				Location: dto.LocationResponse{
					Id:   storeRequestItem.Warehouse.Location.Id,
					Name: storeRequestItem.Warehouse.Location.Name,
				},
			},
			WarehouseItem: dto.WarehouseItemResponse{
				Id:       storeRequestItem.WarehouseItem.Id,
				Name:     storeRequestItem.WarehouseItem.Name,
				Category: storeRequestItem.WarehouseItem.Category.String(),
				Unit:     storeRequestItem.WarehouseItem.Unit,
			},
			Store: dto.StoreResponse{
				Id:   storeRequestItem.Store.Id,
				Name: storeRequestItem.Store.Name,
				Location: dto.LocationResponse{
					Id:   storeRequestItem.Store.Location.Id,
					Name: storeRequestItem.Store.Location.Name,
				},
			},
			Quantity: storeRequestItem.Quantity,
			Status:   storeRequestItem.Status.String(),
		}
	}

	return storeRequestItemResponses, nil
}

func (s *StoreService) UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	status := enum.ValueOfRequestItemStatus(request.Status)
	if !status.IsValid() {
		s.log.Error("[UpdateStoreRequestItem] invalid status", zap.String("status", request.Status))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("invalid status")
	}

	if storeRequestItem.Status == enum.RequestItemStatusAccepted || storeRequestItem.Status == enum.RequestItemStatusRejected {
		s.log.Error("[UpdateStoreRequestItem] store request item is already accepted or rejected", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is already accepted or rejected")
	}

	if status == enum.RequestItemStatusAccepted && storeRequestItem.Status != enum.RequestItemStatusSentOff {
		s.log.Error("[UpdateStoreRequestItem] store request item is not sent off", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is not sent off")
	}

	if status == enum.RequestItemStatusSentOff && storeRequestItem.Status != enum.RequestItemStatusPending {
		s.log.Error("[UpdateStoreRequestItem] store request item is not pending", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is not pending")
	}

	if status == enum.RequestItemStatusPending && storeRequestItem.Status != enum.RequestItemStatusPending {
		s.log.Error("[UpdateStoreRequestItem] store request item is not pending", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is not pending")
	}

	if status == enum.RequestItemStatusAccepted {
		// Todo update into stock store_item
	}

	if storeRequestItem.Status == enum.RequestItemStatusPending {
		storeRequestItem.Quantity = request.Quantity
	} else {
		s.log.Error("[UpdateStoreRequestItem] can't update quantity when status is not pending", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("can't update quantity when status is not pending")
	}

	storeRequestItem.Status = status
	storeRequestItem.UpdatedBy = accountId

	err = s.repository.UpdateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to update store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	storeRequestItem, err = s.repository.GetStoreRequestItemById(storeRequestItem.Id)
	if err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to commit transaction", zap.Error(err))
	}

	return dto.StoreRequestItemResponse{
		Id: storeRequestItem.Id,
		Warehouse: dto.WarehouseResponse{
			Id:   storeRequestItem.Warehouse.Id,
			Name: storeRequestItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Warehouse.Location.Id,
				Name: storeRequestItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeRequestItem.WarehouseItem.Id,
			Name:     storeRequestItem.WarehouseItem.Name,
			Category: storeRequestItem.WarehouseItem.Category.String(),
			Unit:     storeRequestItem.WarehouseItem.Unit,
		},
		Store: dto.StoreResponse{
			Id:   storeRequestItem.Store.Id,
			Name: storeRequestItem.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Store.Location.Id,
				Name: storeRequestItem.Store.Location.Name,
			},
		},
		Quantity: storeRequestItem.Quantity,
		Status:   storeRequestItem.Status.String(),
	}, nil
}
