package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type CageService struct {
	log        *zap.Logger
	repository repository.ICageRepository
}

type ICageService interface {
	GetCages(filter dto.GetCageFilter) ([]dto.CageResponse, error)
	CreateCage(request dto.CreateCageRequest, userId uuid.UUID) (dto.CageResponse, error)
	UpdateCage(id uint64, request dto.UpdateCageRequest, updatedBy uuid.UUID) (dto.CageResponse, error)
	DeleteCage(id uint64) error
	GetCageById(id uint64) (dto.CageResponse, error)
	GetCagesByIds(ids []uint64) ([]dto.CageResponse, error)

	CreateChickenCage(request dto.CreateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error)
	GetChickenCages(filter dto.GetChickenCageFilter) ([]dto.ChickenCageResponse, error)
	GetChickenCageById(id uint64) (dto.ChickenCageResponse, error)
	GetChickenCagesByCageIds(cageIds []uint64) ([]dto.ChickenCageResponse, error)
	CreateChickenCageInBatch(request []dto.CreateChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error)
}

func NewCageService(log *zap.Logger, repository repository.ICageRepository) ICageService {
	return &CageService{
		log:        log,
		repository: repository,
	}
}

func (s *CageService) GetCages(filter dto.GetCageFilter) ([]dto.CageResponse, error) {
	s.repository.UseTx(false)

	cages, err := s.repository.GetCages(filter)
	if err != nil {
		s.log.Error("failed to get cages", zap.Error(err))
		return nil, err
	}

	cageResponses := make([]dto.CageResponse, 0)
	for _, cage := range cages {
		cageResponses = append(cageResponses, mapper.CageToResponse(&cage))
	}

	return cageResponses, nil
}

func (s *CageService) CreateCage(request dto.CreateCageRequest, createdBy uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("invalid chicken category")
		return dto.CageResponse{}, errx.BadRequest("invalid chicken category")
	}

	cage := entity.Cage{
		LocationId:      request.LocationId,
		Name:            request.Name,
		Capacity:        request.Capacity,
		ChickenCategory: chickenCategory,
		CreatedBy:       uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateCage(&cage)
	if err != nil {
		s.log.Error("failed to create cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	chickenCage := entity.ChickenCage{
		CageId:    cage.Id,
		CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err = s.repository.CreateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed to create chicken cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transcation", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage, err = s.repository.GetCageById(cage.Id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) UpdateCage(id uint64, request dto.UpdateCageRequest, updatedBy uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("invalid chicken category")
		return dto.CageResponse{}, errx.BadRequest("invalid chicken category")
	}

	cage, err := s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage.Name = request.Name
	cage.LocationId = request.LocationId
	cage.Capacity = request.Capacity
	cage.ChickenCategory = chickenCategory
	cage.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateCage(&cage)
	if err != nil {
		s.log.Error("failed to update cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage, err = s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) DeleteCage(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteCage(id)
	if err != nil {
		s.log.Error("failed to delete cage", zap.Error(err))
		return err
	}

	return nil
}

func (s *CageService) GetChickenCages(filter dto.GetChickenCageFilter) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCageResponses := make([]dto.ChickenCageResponse, 0)
	chickenCages, err := s.repository.GetChickenCages(filter)
	if err != nil {
		return chickenCageResponses, err
	}

	for _, chickenCage := range chickenCages {
		chickenCageResponses = append(chickenCageResponses, mapper.ChickenCageToResponse(&chickenCage))
	}

	return chickenCageResponses, nil
}

func (s *CageService) GetChickenCageById(id uint64) (dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageById(id)
	if err != nil {
		return dto.ChickenCageResponse{}, err
	}

	return mapper.ChickenCageToResponse(&chickenCage), nil
}

func (s *CageService) UpdateChickenCage(id uint64, request dto.UpdateChickenCageRequest, updatedBy uuid.UUID) (dto.ChickenCageResponse, error) {
	return dto.ChickenCageResponse{}, nil
}

func (s *CageService) GetCageById(id uint64) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	cage, err := s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) CreateChickenCage(request dto.CreateChickenCageRequest, userId uuid.UUID) (dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCage := entity.ChickenCage{
		CageId:               request.CageId,
		ChickenProcurementId: request.ChickenProcurementId,
		TotalChicken:         request.TotalChicken,
		CreatedBy:            uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateChickenCage(&chickenCage)
	if err != nil {
		s.log.Error("failed to create chicken cage", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	chickenCage, err = s.repository.GetChickenCageById(chickenCage.Id)
	if err != nil {
		s.log.Error("failed get chicken cage by id", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	return mapper.ChickenCageToResponse(&chickenCage), err
}

func (s *CageService) GetCagesByIds(ids []uint64) ([]dto.CageResponse, error) {
	s.repository.UseTx(false)

	cages, err := s.repository.GetCagesByIds(ids)
	if err != nil {
		return nil, err
	}

	cageResponses := make([]dto.CageResponse, 0)
	for _, cage := range cages {
		cageResponses = append(cageResponses, mapper.CageToResponse(&cage))
	}

	return cageResponses, nil
}

func (s *CageService) GetChickenCagesByCageIds(ids []uint64) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCages, err := s.repository.GetChickenCagesByCageIds(ids)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ChickenCageResponse, 0)
	for _, chickenCage := range chickenCages {
		response = append(response, mapper.ChickenCageToResponse(&chickenCage))
	}

	return response, nil
}

func (s *CageService) CreateChickenCageInBatch(request []dto.CreateChickenCageRequest, userId uuid.UUID) ([]dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCages := make([]entity.ChickenCage, 0)
	for _, req := range request {
		chickenCages = append(chickenCages, entity.ChickenCage{
			CageId:               req.CageId,
			ChickenProcurementId: req.ChickenProcurementId,
			TotalChicken:         req.TotalChicken,
		})
	}

	err := s.repository.CreateChickenCageInBatch(&chickenCages)
	if err != nil {
		s.log.Error("failed create chicken cage in batch", zap.Error(err))
		return nil, err
	}

	chickenCageIds := make([]uint64, 0)
	for _, chichickenCage := range chickenCages {
		chickenCageIds = append(chickenCageIds, chichickenCage.Id)
	}

	chickenCages, err = s.repository.GetChickenCageByIds(chickenCageIds)
	if err != nil {
		return nil, err
	}

	response := make([]dto.ChickenCageResponse, 0)
	for _, chickenCage := range chickenCages {
		response = append(response, mapper.ChickenCageToResponse(&chickenCage))
	}

	return response, nil
}
