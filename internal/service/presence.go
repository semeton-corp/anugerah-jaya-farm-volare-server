package service

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type PresenceService struct {
	log        *zap.Logger
	repository repository.IPresenceRepository
}

type IPresenceService interface {
}

func NewPresenceService(log *zap.Logger, repository repository.IPresenceRepository) IPresenceService {
	return &PresenceService{
		log:        log,
		repository: repository,
	}
}

func (p *PresenceService) GetOwnPresence() {

}
