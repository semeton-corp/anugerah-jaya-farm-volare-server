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
	GetRoleById(id uint64) (dto.RoleResponse, error)
	GetRoleByName(name string) (dto.RoleResponse, error)
}

func NewRoleService(log *zap.Logger, repository repository.IRoleRepository) IRoleService {
	return &RoleService{
		log:        log,
		repository: repository,
	}
}

func (s *RoleService) GetRoles() ([]dto.RoleResponse, error) {
	s.repository.UseTx(false)

	roles, err := s.repository.GetRoles()
	if err != nil {
		s.log.Error("failed to get roles", zap.Error(err))
		return nil, err
	}

	roleResponses := make([]dto.RoleResponse, 0)
	for _, role := range roles {
		roleResponses = append(roleResponses, mapper.RoleToResponse(&role))
	}

	return roleResponses, nil
}

func (s *RoleService) GetRoleById(id uint64) (dto.RoleResponse, error) {
	s.repository.UseTx(false)

	role, err := s.repository.GetRoleById(id)
	if err != nil {
		s.log.Error("[GetRoleById] failed to get role by id", zap.Error(err))
		return dto.RoleResponse{}, err
	}

	return mapper.RoleToResponse(&role), nil
}

func (s *RoleService) GetRoleByName(name string) (dto.RoleResponse, error) {
	s.repository.UseTx(false)

	role, err := s.repository.GetRoleByName(name)
	if err != nil {
		s.log.Error("failed get role by name", zap.Error(err))
		return dto.RoleResponse{}, err
	}

	return mapper.RoleToResponse(&role), nil
}
