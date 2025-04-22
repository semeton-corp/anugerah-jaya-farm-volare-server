package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type EggService struct {
	log        *zap.Logger
	repository repository.IEggRepository
}

type IEggService interface {
	CreateEggMonitoring(request dto.CreateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error)
	GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error)
	GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error)
	UpdateEggMonitoring(id uint64, request dto.UpdateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error)
	DeleteEggMonitoring(id uint64) error
}

func NewEggService(log *zap.Logger, repository repository.IEggRepository) IEggService {
	return &EggService{
		log:        log,
		repository: repository,
	}
}

func (e *EggService) CreateEggMonitoring(request dto.CreateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error) {
	count, err := e.repository.CountEggMonitoringByCageIdToday(request.CageId)
	if err != nil {
		e.log.Error("[CreateEggMonitoring] failed to count egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	if count > 0 {
		e.log.Error("[CreateEggMonitoring] egg monitoring already exists for today", zap.Error(errx.BadRequest("egg monitoring already exists for today")))
		return dto.EggMonitoringResponse{}, errx.BadRequest("egg monitoring already exists for today")
	}

	eggMonitoring := entity.EggMonitoring{
		CageId:          request.CageId,
		WarehouseId:     request.WarehouseId,
		TotalGoodEgg:    request.TotalGoodEgg,
		TotalCrackedEgg: request.TotalCrackedEgg,
		TotalBrokeEgg:   request.TotalBrokeEgg,
		TotalRejectEgg:  request.TotalRejectEgg,
		Weight:          request.Weight,
		CreatedBy:       accountId,
	}

	if err := e.repository.CreateEggMonitoring(&eggMonitoring); err != nil {
		e.log.Error("[CreateEggMonitoring] failed to create egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = e.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		e.log.Error("[CreateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (e *EggService) GetEggMonitoringById(id uint64) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[GetEggMonitoringById] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (e *EggService) GetEggMonitorings(filter dto.GetEggMonitoringFilter) ([]dto.EggMonitoringListResponse, error) {
	eggMonitorings, err := e.repository.GetEggMonitorings(filter)
	if err != nil {
		e.log.Error("[GetEggMonitorings] failed to get egg monitorings", zap.Error(err))
		return nil, err
	}

	eggMonitoringResponses := make([]dto.EggMonitoringListResponse, 0, len(eggMonitorings))
	for _, eggMonitoring := range eggMonitorings {
		eggMonitoringResponse := mapper.EggMonitoringToListResponse(&eggMonitoring)

		if eggMonitoringResponse.TotalAll == 0 {
			eggMonitoringResponse.AbnormalityRate = 0
			eggMonitoringResponse.Description = constant.EggMonitoringStatusSafety
		} else {
			eggMonitoringResponse.AbnormalityRate = float64(eggMonitoring.TotalCrackedEgg+eggMonitoring.TotalBrokeEgg+eggMonitoring.TotalRejectEgg) / float64(eggMonitoringResponse.TotalAll) * 100
			eggMonitoringResponse.Description = constant.EggMonitoringStatusSafety
		}

		eggMonitoringResponses = append(eggMonitoringResponses, eggMonitoringResponse)
	}

	return eggMonitoringResponses, nil
}

func (e *EggService) UpdateEggMonitoring(id uint64, request dto.UpdateEggMonitoringRequest, accountId uuid.UUID) (dto.EggMonitoringResponse, error) {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[UpdateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring.Weight = request.Weight
	eggMonitoring.CageId = request.CageId
	eggMonitoring.WarehouseId = request.WarehouseId
	eggMonitoring.TotalGoodEgg = request.TotalGoodEgg
	eggMonitoring.TotalCrackedEgg = request.TotalCrackedEgg
	eggMonitoring.TotalBrokeEgg = request.TotalBrokeEgg
	eggMonitoring.TotalRejectEgg = request.TotalRejectEgg
	eggMonitoring.UpdatedBy = accountId

	if err := e.repository.UpdateEggMonitoring(&eggMonitoring); err != nil {
		e.log.Error("[UpdateEggMonitoring] failed to update egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoring, err = e.repository.GetEggMonitoringById(eggMonitoring.Id)
	if err != nil {
		e.log.Error("[UpdateEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return dto.EggMonitoringResponse{}, err
	}

	eggMonitoringResponse := mapper.EggMonitoringToResponse(&eggMonitoring)

	return eggMonitoringResponse, nil
}

func (e *EggService) DeleteEggMonitoring(id uint64) error {
	eggMonitoring, err := e.repository.GetEggMonitoringById(id)
	if err != nil {
		e.log.Error("[DeleteEggMonitoring] failed to get egg monitoring", zap.Error(err))
		return err
	}

	if err := e.repository.DeleteEggMonitoring(eggMonitoring.Id); err != nil {
		e.log.Error("[DeleteEggMonitoring] failed to delete egg monitoring", zap.Error(err))
		return err
	}

	return nil
}
