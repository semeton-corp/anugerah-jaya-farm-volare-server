package service

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
	"go.uber.org/zap"
)

type PresenceService struct {
	log        *zap.Logger
	repository repository.IPresenceRepository
}

type IPresenceService interface {
	GetCurrentStaffPresence(staffId uuid.UUID) (dto.PresenceResponse, error)
	GetAllStaffPresences(staffId uuid.UUID, filter dto.GetPresenceFilter) ([]dto.PresenceListResponse, error)
	ArrivalPresence(id uint64, acountId uuid.UUID) (dto.PresenceResponse, error)
	DeparturePresence(id uint64, accountId uuid.UUID) (dto.PresenceResponse, error)
}

func NewPresenceService(log *zap.Logger, repository repository.IPresenceRepository) IPresenceService {
	return &PresenceService{
		log:        log,
		repository: repository,
	}
}

func (p *PresenceService) GetCurrentStaffPresence(staffId uuid.UUID) (dto.PresenceResponse, error) {
	p.repository.UseTx(false)

	staffPresence, err := p.repository.GetStaffPresenceTodayByStaffId(staffId)
	if err != nil {
		p.log.Error("[GetCurrentStaffPresence] failed to get staff presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	return mapper.PresenceToResponse(&staffPresence), nil
}

func (p *PresenceService) GetAllStaffPresences(staffId uuid.UUID, filter dto.GetPresenceFilter) ([]dto.PresenceListResponse, error) {
	p.repository.UseTx(false)

	staffPresence, err := p.repository.GetStaffPresenceByStaffId(staffId, filter)
	if err != nil {
		p.log.Error("[GetAllStaffPresences] failed to get staff presence", zap.Error(err))
		return nil, err
	}

	presenceResponses := make([]dto.PresenceListResponse, len(staffPresence))
	for i, presence := range staffPresence {
		presenceResponses[i] = mapper.PresenceToResponseList(&presence)
		extraTime := presence.EndTime.Sub(time.Date(0, 0, 0, 5, 0, 0, 0, time.Local))
		if extraTime > 0 {
			presenceResponses[i].ExtraTime = fmt.Sprintf("%02d Jam, %02d Menit", int(extraTime.Hours()), int(extraTime.Minutes())%60)
		} else {
			presenceResponses[i].ExtraTime = ""
		}
	}

	return presenceResponses, nil
}

func (p *PresenceService) ArrivalPresence(id uint64, accountId uuid.UUID) (dto.PresenceResponse, error) {
	p.repository.UseTx(false)

	staffPresence, err := p.repository.GetStaffPresenceById(id)
	if err != nil {
		p.log.Error("[ArrivalPresence] failed to get staff presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	staffPresence.IsPresent = true
	staffPresence.StartTime = datatype.TimeOnly{time.Now()}

	err = p.repository.UpdateStaffPresence(&staffPresence)
	if err != nil {
		p.log.Error("[ArrivalPresence] failed to update staff presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	return mapper.PresenceToResponse(&staffPresence), nil
}

func (p *PresenceService) DeparturePresence(id uint64, acountId uuid.UUID) (dto.PresenceResponse, error) {
	p.repository.UseTx(false)

	staffPresence, err := p.repository.GetStaffPresenceById(id)
	if err != nil {
		p.log.Error("[DeparturePresence] failed to get staff presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	staffPresence.IsPresent = true
	staffPresence.EndTime = datatype.TimeOnly{time.Now()}

	err = p.repository.UpdateStaffPresence(&staffPresence)
	if err != nil {
		p.log.Error("[DeparturePresence] failed to update staff presence", zap.Error(err))
		return dto.PresenceResponse{}, err
	}

	return mapper.PresenceToResponse(&staffPresence), nil
}
