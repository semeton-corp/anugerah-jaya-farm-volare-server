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

func (s *CageService) CreateCage(request dto.CreateCageRequest, userId uuid.UUID) (dto.CageResponse, error) {
	s.repository.UseTx(false)

	chickenCategory := enum.ValueOfChickenCategory(request.ChickenCategory)
	if !chickenCategory.IsValid() {
		s.log.Warn("[CreateCage] invalid chicken category")
		return dto.CageResponse{}, errx.BadRequest("invalid chicken category")
	}

	data := entity.Cage{
		LocationId:      request.LocationId,
		Name:            request.Name,
		Capacity:        request.Capacity,
		ChickenCategory: chickenCategory,
		CreatedBy:       uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateCage(&data)
	if err != nil {
		s.log.Error("[CreateCage] failed to create cage", zap.Error(err))
		return dto.CageResponse{}, err
	}

	data, err = s.repository.GetCageById(data.Id)
	if err != nil {
		s.log.Error("[CreateCage] failed to get cage by id", zap.Error(err))
		return dto.CageResponse{}, err
	}

	return mapper.CageToResponse(&data), nil
}
