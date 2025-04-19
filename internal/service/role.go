package service

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type RoleService struct {
	log        *zap.Logger
	repository repository.IRoleRepository
}

type IRoleService interface {
	GetRoles() ([]dto.RoleResponse, error)
}

func NewRoleService(log *zap.Logger, repository repository.IRoleRepository) IRoleService {
	return &RoleService{
		log:        log,
		repository: repository,
	}
}

func (r *RoleService) GetRoles() ([]dto.RoleResponse, error) {
	r.repository.UseTx(false)

	roles, err := r.repository.GetRoles()
	if err != nil {
		r.log.Error("[GetRoles] failed to get roles", zap.Error(err))
		return nil, err
	}

	var roleResponses []dto.RoleResponse
	for _, role := range roles {
		roleResponses = append(roleResponses, mapper.RoleToResponse(&role))
	}

	return roleResponses, nil
}
