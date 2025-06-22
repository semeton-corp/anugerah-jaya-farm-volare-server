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
	GetChickenCageByCageId(cageId uint64) (dto.ChickenCageResponse, error)
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
		s.log.Error("[GetCages] failed to get cages", zap.Error(err))
		return nil, err
	}

	cageResponses := make([]dto.CageResponse, 0)
	for _, cage := range cages {
		cageResponses = append(cageResponses, mapper.CageToResponse(&cage))
	}

	return cageResponses, nil
}

func (s *CageService) CreateCage(request dto.CreateCageRequest, createdBy uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("[CreateCage] invalid chicken category")
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
		s.log.Error("[CreateCage] failed to create cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage, err = s.repository.GetCageById(cage.Id)
	if err != nil {
		s.log.Error("[CreateCage] failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) UpdateCage(id uint64, request dto.UpdateCageRequest, updatedBy uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("[UpdateCage] invalid chicken category")
		return dto.CageResponse{}, errx.BadRequest("invalid chicken category")
	}

	cage, err := s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("[UpdateCage] failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage.Name = request.Name
	cage.LocationId = request.LocationId
	cage.Capacity = request.Capacity
	cage.ChickenCategory = chickenCategory
	cage.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	err = s.repository.UpdateCage(&cage)
	if err != nil {
		s.log.Error("[UpdateCage] failed to update cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	cage, err = s.repository.GetCageById(id)
	if err != nil {
		s.log.Error("[UpdateCage] failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&cage), nil
}

func (s *CageService) DeleteCage(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteCage(id)
	if err != nil {
		s.log.Error("[DeleteCage] failed to delete cage", zap.Error(err))
		return err
	}

	return nil
}

func (s *CageService) GetChickenCageByCageId(cageId uint64) (dto.ChickenCageResponse, error) {
	s.repository.UseTx(false)

	chickenCage, err := s.repository.GetChickenCageByCageId(cageId)
	if err != nil {
		s.log.Error("[GetChickenCageByCageId] failed to get chicken cage by cage id", zap.Error(err))
		return dto.ChickenCageResponse{}, err
	}

	return mapper.ChickenCageToResponse(&chickenCage), nil
}
