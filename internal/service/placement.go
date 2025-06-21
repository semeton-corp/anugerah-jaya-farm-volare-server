package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type PlacementService struct {
	log        *zap.Logger
	repository repository.IPlacementRepository
}

type IPlacementService interface {
	CreateCagePlacementBatch(request dto.CreateCagePlacementRequest, createdBy uuid.UUID) ([]dto.CagePlacementResponse, error)
	CreateStorePlacementBatch(request dto.CreateStorePlacementRequest, createdBy uuid.UUID) ([]dto.StorePlacementResponse, error)
	CreateWarehousePlacementBatch(request dto.CreateWarehousePlacementRequest, createdBy uuid.UUID) ([]dto.WarehousePlacementResponse, error)
}

func NewPlacementService(log *zap.Logger, repository repository.IPlacementRepository) IPlacementService {
	return &PlacementService{
		log:        log,
		repository: repository,
	}
}

func (s *PlacementService) CreateCagePlacementBatch(request dto.CreateCagePlacementRequest, createdBy uuid.UUID) ([]dto.CagePlacementResponse, error) {
	s.repository.UseTx(false)

	data := make([]entity.CagePlacement, 0)
	userId := uuid.MustParse(request.UserId)
	for _, cageId := range request.CageIds {
		data = append(data, entity.CagePlacement{
			UserId: userId,
			CageId: cageId,
		})
	}

	err := s.repository.CreateCagePlacementBatch(data)
	if err != nil {
		s.log.Error("[CreateCagePlacementBatch] failed to create cage placement in batch", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.CagePlacementResponse, 0)
	data, err = s.repository.GetCagePlacementByUserId(userId)
	if err != nil {
		s.log.Error("[CreateCagePlacementBatch] failed to get cage placement by user id", zap.Error(err))
		return nil, err
	}

	for _, d := range data {
		dataResponse = append(dataResponse, mapper.CagePlacementToResponse(&d))
	}

	return dataResponse, nil
}

func (s *PlacementService) CreateStorePlacementBatch(request dto.CreateStorePlacementRequest, createdBy uuid.UUID) ([]dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	data := make([]entity.StorePlacement, 0)
	userId := uuid.MustParse(request.UserId)
	for _, storeId := range request.StoreIds {
		data = append(data, entity.StorePlacement{
			UserId:    userId,
			StoreId:   storeId,
			CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
		})
	}

	err := s.repository.CreateStorePlacementBatch(data)
	if err != nil {
		s.log.Error("[CreateStorePlacementBatch] failed to create store placement in batch", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.StorePlacementResponse, 0)
	data, err = s.repository.GetStorePlacementByUserId(userId)
	if err != nil {
		s.log.Error("[CreateStorePlacementBatch] failed to get store placement by user id", zap.Error(err))
		return nil, err
	}

	for _, d := range data {
		dataResponse = append(dataResponse, mapper.StorePlacementToResponse(&d))
	}

	return dataResponse, nil
}

func (s *PlacementService) CreateWarehousePlacementBatch(request dto.CreateWarehousePlacementRequest, createdBy uuid.UUID) ([]dto.WarehousePlacementResponse, error) {
	s.repository.UseTx(false)

	data := make([]entity.WarehousePlacement, 0)
	userId := uuid.MustParse(request.UserId)
	for _, WarehouseId := range request.WarehouseIds {
		data = append(data, entity.WarehousePlacement{
			UserId:      userId,
			WarehouseId: WarehouseId,
		})
	}

	err := s.repository.CreateWarehousePlacementBatch(data)
	if err != nil {
		s.log.Error("[CreateWarehousePlacementBatch] failed to create warehouse placement in batch", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.WarehousePlacementResponse, 0)
	data, err = s.repository.GetWarehousePlacementByUserId(userId)
	if err != nil {
		s.log.Error("[CreateWarehousePlacementBatch] failed to get warehouse placement by user id", zap.Error(err))
		return nil, err
	}

	for _, d := range data {
		dataResponse = append(dataResponse, mapper.WarehousePlacementToResponse(&d))
	}

	return dataResponse, nil
}

func (s *PlacementService) DeleteCagePlacementByUserId(userId uuid.UUID) error {
	s.repository.UseTx(false)
	return s.repository.DeleteCagePlacementByUserId(userId)
}

func (s *PlacementService) DeleteStorePlacementByUserId(userId uuid.UUID) error {
	s.repository.UseTx(false)
	return s.repository.DeleteStorePlacementByUserId(userId)
}

func (s *PlacementService) DeleteWarehousePlacementByUserId(userId uuid.UUID) error {
	s.repository.UseTx(false)
	return s.repository.DeleteWarehousePlacementByUserId(userId)
}
