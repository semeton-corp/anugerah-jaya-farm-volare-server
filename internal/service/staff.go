package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type StaffService struct {
	log        *zap.Logger
	repository repository.IStaffRepository
}

type IStaffService interface {
	GetStaffById(id uuid.UUID) (dto.StaffResponse, error)
	UpdateStaff(id uuid.UUID, request dto.UpdateStaffRequest, accountId uuid.UUID) (dto.StaffResponse, error)
}

func NewStaffService(log *zap.Logger, repository repository.IStaffRepository) IStaffService {
	return &StaffService{
		log:        log,
		repository: repository,
	}
}

func (s *StaffService) GetStaffById(id uuid.UUID) (dto.StaffResponse, error) {
	s.repository.UseTx(false)

	staff, err := s.repository.GetStaffById(id)
	if err != nil {
		s.log.Error("[GetStaffById] failed to get staff by id", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	return mapper.StaffToResponse(&staff), nil
}

func (s *StaffService) UpdateStaff(id uuid.UUID, request dto.UpdateStaffRequest, accountId uuid.UUID) (dto.StaffResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	staff, err := s.repository.GetStaffById(id)
	if err != nil {
		s.log.Error("[UpdateStaff] failed to get staff by id", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	staff.Name = request.Name

	if err := s.repository.UpdateStaff(&staff); err != nil {
		s.log.Error("[UpdateStaff] failed to update staff", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStaff] failed to commit transaction", zap.Error(err))
		return dto.StaffResponse{}, err
	}

	return mapper.StaffToResponse(&staff), nil
}
