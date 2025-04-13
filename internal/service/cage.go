package service

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type CageService struct {
	log        *zap.Logger
	repository repository.ICageRepository
}

type ICageService interface {
	GetCages() ([]dto.CageResponse, error)
}

func NewCageService(log *zap.Logger, repository repository.ICageRepository) ICageService {
	return &CageService{
		log:        log,
		repository: repository,
	}
}

func (c *CageService) GetCages() ([]dto.CageResponse, error) {
	c.repository.UseTx(false)

	cages, err := c.repository.GetCages()
	if err != nil {
		c.log.Error("[GetCages] failed to get cages", zap.Error(err))
		return nil, err
	}

	var cageResponses []dto.CageResponse
	for _, cage := range cages {
		cageResponses = append(cageResponses, dto.CageResponse{
			Id:   cage.Id,
			Name: cage.Name,
			Location: dto.LocationResponse{
				Id:   cage.Location.Id,
				Name: cage.Location.Name,
			},
		})
	}

	return cageResponses, nil
}
