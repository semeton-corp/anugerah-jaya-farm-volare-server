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
	CreateStorePlacement(request dto.CreateStorePlacementRequest, createdBy uuid.UUID) (dto.StorePlacementResponse, error)
	CreateWarehousePlacement(request dto.CreateWarehousePlacementRequest, createdBy uuid.UUID) (dto.WarehousePlacementResponse, error)

	GetStorePlacementByUserId(userId uuid.UUID) (dto.StorePlacementResponse, error)
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
			UserId:    userId,
			CageId:    cageId,
			CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
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

func (s *PlacementService) CreateStorePlacement(request dto.CreateStorePlacementRequest, createdBy uuid.UUID) (dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	userId := uuid.MustParse(request.UserId)
	data := entity.StorePlacement{
		UserId:    userId,
		StoreId:   request.StoreId,
		CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateStorePlacement(&data)
	if err != nil {
		s.log.Error("failed to create store placement in batch", zap.Error(err))
		return dto.StorePlacementResponse{}, err
	}

	data, err = s.repository.GetStorePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to get store placement by user id", zap.Error(err))
		return dto.StorePlacementResponse{}, err
	}

	return mapper.StorePlacementToResponse(&data), nil
}

func (s *PlacementService) CreateWarehousePlacement(request dto.CreateWarehousePlacementRequest, createdBy uuid.UUID) (dto.WarehousePlacementResponse, error) {
	s.repository.UseTx(false)

	userId := uuid.MustParse(request.UserId)
	data := entity.WarehousePlacement{
		UserId:      userId,
		WarehouseId: request.WarehouseId,
		CreatedBy:   uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateWarehousePlacement(&data)
	if err != nil {
		s.log.Error("failed to create warehouse placement in batch", zap.Error(err))
		return dto.WarehousePlacementResponse{}, err
	}

	data, err = s.repository.GetWarehousePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to get warehouse placement by user id", zap.Error(err))
		return dto.WarehousePlacementResponse{}, err
	}

	return mapper.WarehousePlacementToResponse(&data), nil
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

func (s *PlacementService) GetStorePlacementByUserId(userId uuid.UUID) (dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	storePlacement, err := s.repository.GetStorePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to get store placement by user id", zap.Error(err))
		return dto.StorePlacementResponse{}, err
	}

	return mapper.StorePlacementToResponse(&storePlacement), err
}
