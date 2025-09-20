package service

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type PlacementService struct {
	log        *zap.Logger
	repository repository.IPlacementRepository
}

type IPlacementService interface {
	CreateCagePlacementForAuthentication(request []dto.CreateCagePlacementRequest, userId uuid.UUID, roleId uint64) ([]dto.CagePlacementResponse, error)
	CreateStorePlacementForAuthentication(request dto.CreateStorePlacementRequest, userId uuid.UUID) ([]dto.StorePlacementResponse, error)
	CreateWarehousePlacementForAuthentication(request dto.CreateWarehousePlacementRequest, userId uuid.UUID) ([]dto.WarehousePlacementResponse, error)

	UpdateCagePlacement(request dto.UpdateCagePlacementRequest, userId uuid.UUID) ([]dto.CagePlacementResponse, error)
	CreateStorePlacement(request dto.CreateStorePlacementRequest, userId uuid.UUID) ([]dto.StorePlacementResponse, error)
	CreateWarehousePlacement(request dto.CreateWarehousePlacementRequest, userId uuid.UUID) ([]dto.WarehousePlacementResponse, error)

	GetStorePlacementByUserId(userId uuid.UUID) ([]dto.StorePlacementResponse, error)
	GetWarehousePlacementByUserId(userId uuid.UUID) ([]dto.WarehousePlacementResponse, error)
	GetCagePlacementByUserId(userId uuid.UUID) ([]dto.CagePlacementResponse, error)

	GetStorePlacementByStoreId(storeId uint64) ([]dto.StorePlacementResponse, error)
	GetWarehousePlacementByWarehouseId(warehouseId uint64) ([]dto.WarehousePlacementResponse, error)

	DeleteStorePlacementByUserId(userId uuid.UUID) error
	DeleteWarehousePlacementByUserId(userId uuid.UUID) error
	DeleteCagePlacementByUserIdAndCageId(userId uuid.UUID, cageId uint64) error
}

func NewPlacementService(log *zap.Logger, repository repository.IPlacementRepository) IPlacementService {
	return &PlacementService{
		log:        log,
		repository: repository,
	}
}

func (s *PlacementService) CreateCagePlacementForAuthentication(requests []dto.CreateCagePlacementRequest, userId uuid.UUID, roleId uint64) ([]dto.CagePlacementResponse, error) {
	s.repository.UseTx(false)

	for _, r := range requests {
		cagePlacements, err := s.repository.GetCagePlacementByCageId(r.CageId)
		if err != nil {
			s.log.Error("failed get cage placement by cage id")
			return nil, err
		}

		for _, cagePlacement := range cagePlacements {
			if cagePlacement.User.RoleId == roleId {
				return nil, errx.BadRequest(fmt.Sprintf("user with this role already exist in cage %s", cagePlacement.Cage.Name))
			}
		}
	}

	data := make([]entity.CagePlacement, 0)
	for _, request := range requests {
		data = append(data, entity.CagePlacement{
			UserId:    uuid.MustParse(request.UserId),
			CageId:    request.CageId,
			CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
		})
	}

	err := s.repository.CreateCagePlacementBatch(data)
	if err != nil {
		s.log.Error("failed to create cage placement in batch", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.CagePlacementResponse, 0)
	data, err = s.repository.GetCagePlacementByUserId(uuid.MustParse(requests[0].UserId))
	if err != nil {
		s.log.Error("failed to get cage placement by user id", zap.Error(err))
		return nil, err
	}

	for _, d := range data {
		dataResponse = append(dataResponse, mapper.CagePlacementToResponse(&d))
	}

	return dataResponse, nil
}

func (s *PlacementService) CreateStorePlacementForAuthentication(request dto.CreateStorePlacementRequest, userId uuid.UUID) ([]dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	userIdRequest := uuid.MustParse(request.UserId)
	data := entity.StorePlacement{
		UserId:    userIdRequest,
		StoreId:   request.StoreId,
		CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateStorePlacement(&data)
	if err != nil {
		s.log.Error("failed to create store placement in batch", zap.Error(err))
		return nil, err
	}

	placements, err := s.repository.GetStorePlacementByUserId(userIdRequest)
	if err != nil {
		s.log.Error("failed to get store placement by user id", zap.Error(err))
		return nil, err
	}

	response := make([]dto.StorePlacementResponse, 0)
	for _, e := range placements {
		response = append(response, mapper.StorePlacementToResponse(&e))
	}

	return response, nil
}

func (s *PlacementService) CreateWarehousePlacementForAuthentication(request dto.CreateWarehousePlacementRequest, userId uuid.UUID) ([]dto.WarehousePlacementResponse, error) {
	s.repository.UseTx(false)

	userIdRequest := uuid.MustParse(request.UserId)
	data := entity.WarehousePlacement{
		UserId:      userIdRequest,
		WarehouseId: request.WarehouseId,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateWarehousePlacement(&data)
	if err != nil {
		s.log.Error("failed to create warehouse placement in batch", zap.Error(err))
		return nil, err
	}

	placements, err := s.repository.GetWarehousePlacementByUserId(userIdRequest)
	if err != nil {
		s.log.Error("failed to get warehouse placement by user id", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehousePlacementResponse, 0)
	for _, e := range placements {
		response = append(response, mapper.WarehousePlacementToResponse(&e))
	}

	return response, nil
}

func (s *PlacementService) DeleteCagePlacementByUserIdAndCageId(userId uuid.UUID, cageId uint64) error {
	s.repository.UseTx(false)
	err := s.repository.DeleteCagePlacementByUserIdAndCageId(userId, cageId)
	if err != nil {
		s.log.Error("failed to delete cage placement by user id", zap.Error(err))
		return err
	}
	return nil
}

func (s *PlacementService) DeleteStorePlacementByUserId(userId uuid.UUID) error {
	s.repository.UseTx(false)
	err := s.repository.DeleteStorePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to delete store placement by user id")
		return err
	}
	return nil
}

func (s *PlacementService) DeleteWarehousePlacementByUserId(userId uuid.UUID) error {
	s.repository.UseTx(false)
	err := s.repository.DeleteWarehousePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to delete warehouse placement by user id")
		return err
	}
	return nil
}

func (s *PlacementService) GetStorePlacementByUserId(userId uuid.UUID) ([]dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	storePlacement, err := s.repository.GetStorePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to get store placement by user id", zap.Error(err))
		return nil, err
	}

	response := make([]dto.StorePlacementResponse, 0)
	for _, e := range storePlacement {
		response = append(response, mapper.StorePlacementToResponse(&e))
	}

	return response, nil
}

func (s *PlacementService) GetWarehousePlacementByUserId(userId uuid.UUID) ([]dto.WarehousePlacementResponse, error) {
	s.repository.UseTx(false)

	warehousePlacement, err := s.repository.GetWarehousePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to get warehouse placement by user id", zap.Error(err))
		return nil, err
	}

	response := make([]dto.WarehousePlacementResponse, 0)
	for _, e := range warehousePlacement {
		response = append(response, mapper.WarehousePlacementToResponse(&e))
	}

	return response, nil
}

func (s *PlacementService) GetCagePlacementByUserId(userId uuid.UUID) ([]dto.CagePlacementResponse, error) {
	s.repository.UseTx(false)

	dataResponse := make([]dto.CagePlacementResponse, 0)
	data, err := s.repository.GetCagePlacementByUserId(userId)
	if err != nil {
		s.log.Error("failed to get cage placement by user id", zap.Error(err))
		return nil, err
	}

	for _, d := range data {
		dataResponse = append(dataResponse, mapper.CagePlacementToResponse(&d))
	}

	return dataResponse, nil
}

func (s *PlacementService) UpdateCagePlacement(request dto.UpdateCagePlacementRequest, userId uuid.UUID) ([]dto.CagePlacementResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	err := s.repository.DeleteCagePlacementByCageId(request.CageId)
	if err != nil {
		s.log.Error("failed delete cage placement by cage id", zap.Error(err))
		return nil, err
	}

	if len(request.Users) > 0 {
		data := make([]entity.CagePlacement, 0)
		for _, r := range request.Users {
			data = append(data, entity.CagePlacement{
				UserId:    uuid.MustParse(r.UserId),
				CageId:    request.CageId,
				CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
			})
		}

		err := s.repository.CreateCagePlacementBatch(data)
		if err != nil {
			s.log.Error("failed to create cage placement in batch", zap.Error(err))
			return nil, err
		}
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
	}

	dataResponse := make([]dto.CagePlacementResponse, 0)
	data, err := s.repository.GetCagePlacementByCageId(request.CageId)
	if err != nil {
		s.log.Error("failed to get cage placement by cage id", zap.Error(err))
		return nil, err
	}

	for _, d := range data {
		dataResponse = append(dataResponse, mapper.CagePlacementToResponse(&d))
	}

	return dataResponse, nil
}

func (s *PlacementService) CreateStorePlacement(request dto.CreateStorePlacementRequest, userId uuid.UUID) ([]dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	userIdRequest := uuid.MustParse(request.UserId)
	userStorePlacement, err := s.repository.GetStorePlacementByUserId(userIdRequest)
	if err != nil {
		s.log.Error("failed get store placement by user id", zap.Error(err))
		return nil, err
	}

	if len(userStorePlacement) > 0 {
		return nil, errx.BadRequest(fmt.Sprintf("user has been exits in store %s", userStorePlacement[0].Store.Name))
	}

	data := entity.StorePlacement{
		UserId:    userIdRequest,
		StoreId:   request.StoreId,
		CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
	}

	if err := s.repository.CreateStorePlacement(&data); err != nil {
		s.log.Error("failed to create store placement in batch", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.StorePlacementResponse, 0)
	storePlacements, err := s.repository.GetStorePlacementByStoreId(data.StoreId)
	if err != nil {
		s.log.Error("failed to get store placement by user id", zap.Error(err))
		return nil, err
	}

	for _, e := range storePlacements {
		dataResponse = append(dataResponse, mapper.StorePlacementToResponse(&e))
	}

	return dataResponse, nil
}

func (s *PlacementService) CreateWarehousePlacement(request dto.CreateWarehousePlacementRequest, userId uuid.UUID) ([]dto.WarehousePlacementResponse, error) {
	s.repository.UseTx(false)

	userIdRequest := uuid.MustParse(request.UserId)
	userWarehousePlacement, err := s.repository.GetWarehousePlacementByUserId(userIdRequest)
	if err != nil {
		s.log.Error("failed get warehouse placement by user id", zap.Error(err))
		return nil, err
	}

	if len(userWarehousePlacement) > 0 {
		return nil, errx.BadRequest(fmt.Sprintf("user has been exits in warehouse %s", userWarehousePlacement[0].Warehouse.Name))
	}

	data := entity.WarehousePlacement{
		UserId:      userIdRequest,
		WarehouseId: request.WarehouseId,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	err = s.repository.CreateWarehousePlacement(&data)
	if err != nil {
		s.log.Error("failed to create warehouse placement in batch", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.WarehousePlacementResponse, 0)
	warehousePlacements, err := s.repository.GetWarehousePlacementByWarehouseId(data.WarehouseId)
	if err != nil {
		s.log.Error("failed to get warehouse placement by user id", zap.Error(err))
		return nil, err
	}

	for _, e := range warehousePlacements {
		dataResponse = append(dataResponse, mapper.WarehousePlacementToResponse(&e))
	}

	return dataResponse, nil
}

func (s *PlacementService) GetWarehousePlacementByWarehouseId(warehouseId uint64) ([]dto.WarehousePlacementResponse, error) {
	s.repository.UseTx(false)

	warehousePlacements, err := s.repository.GetWarehousePlacementByWarehouseId(warehouseId)
	if err != nil {
		s.log.Error("failed to get warehouse placement by warehouse id", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.WarehousePlacementResponse, 0)
	for _, e := range warehousePlacements {
		dataResponse = append(dataResponse, mapper.WarehousePlacementToResponse(&e))
	}

	return dataResponse, nil
}

func (s *PlacementService) GetStorePlacementByStoreId(storeId uint64) ([]dto.StorePlacementResponse, error) {
	s.repository.UseTx(false)

	storePlacements, err := s.repository.GetStorePlacementByStoreId(storeId)
	if err != nil {
		s.log.Error("failed to get warehouse placement by store id", zap.Error(err))
		return nil, err
	}

	dataResponse := make([]dto.StorePlacementResponse, 0)
	for _, e := range storePlacements {
		dataResponse = append(dataResponse, mapper.StorePlacementToResponse(&e))
	}

	return dataResponse, nil
}
