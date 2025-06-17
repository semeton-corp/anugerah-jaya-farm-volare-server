package service

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type LocationService struct {
	log        *zap.Logger
	repository repository.ILocationRepository
}

type ILocationService interface {
	GetLocations() ([]dto.LocationResponse, error)
}

func NewLocationService(log *zap.Logger, repository repository.ILocationRepository) ILocationService {
	return &LocationService{
		log:        log,
		repository: repository,
	}
}

func (s *LocationService) GetLocations() ([]dto.LocationResponse, error) {
	s.repository.UseTx(false)

	locations, err := s.repository.GetLocations()
	if err != nil {
		s.log.Error("[GetLocations] failed to get locations", zap.Error(err))
		return nil, err
	}

	responseLocation := make([]dto.LocationResponse, len(locations))
	for i, location := range locations {
		responseLocation[i] = mapper.LocationToResponse(&location)
	}

	return responseLocation, nil
}
