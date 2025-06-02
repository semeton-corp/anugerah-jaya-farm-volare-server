package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type SupplierService struct {
	log        *zap.Logger
	repository repository.ISupplierRepository
}

type ISupplierService interface {
	CreateSupplier(requesst *dto.CreateSupplierRequest, accountId uuid.UUID) (dto.SupplierResponse, error)
	GetSupplierById(id uint64) (dto.SupplierResponse, error)
	GetAllSuppliers() ([]dto.SupplierResponse, error)
	UpdateSupplier(id uint64, request *dto.UpdateSupplierRequest, accountId uuid.UUID) (dto.SupplierResponse, error)
	DeleteSupplier(id uint64) error
}

func NewSupplierService(log *zap.Logger, repository repository.ISupplierRepository) ISupplierService {
	return &SupplierService{
		log:        log,
		repository: repository,
	}
}

func (s *SupplierService) CreateSupplier(requesst *dto.CreateSupplierRequest, accountId uuid.UUID) (dto.SupplierResponse, error) {
	s.repository.UseTx(false)

	supplier := entity.Supplier{
		WarehouseItemId: requesst.WarehouseItemId,
		Name:            requesst.Name,
		PhoneNumber:     requesst.PhoneNumber,
		Address:         requesst.Address,
		CreatedBy:       uuid.NullUUID{UUID: accountId, Valid: true},
	}

	err := s.repository.CreateSupplier(&supplier)
	if err != nil {
		s.log.Error("[CreateSupplier] failed to create supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	supplier, err = s.repository.GetSupplierById(supplier.Id)
	if err != nil {
		s.log.Error("[CreateSupplier] failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	return mapper.SupplierToResponse(&supplier), nil
}

func (s *SupplierService) GetSupplierById(id uint64) (dto.SupplierResponse, error) {
	s.repository.UseTx(false)

	supplier, err := s.repository.GetSupplierById(id)
	if err != nil {
		s.log.Error("[GetSupplierById] failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	return mapper.SupplierToResponse(&supplier), nil
}

func (s *SupplierService) GetAllSuppliers() ([]dto.SupplierResponse, error) {
	s.repository.UseTx(false)

	suppliers, err := s.repository.GetAllSuppliers()
	if err != nil {
		s.log.Error("[GetAllSuppliers] failed to get all suppliers", zap.Error(err))
		return nil, err
	}

	supplierResponses := make([]dto.SupplierResponse, len(suppliers))
	for i, supplier := range suppliers {
		supplierResponses[i] = mapper.SupplierToResponse(&supplier)
	}

	return supplierResponses, nil
}

func (s *SupplierService) UpdateSupplier(id uint64, request *dto.UpdateSupplierRequest, accountId uuid.UUID) (dto.SupplierResponse, error) {
	s.repository.UseTx(false)

	supplier, err := s.repository.GetSupplierById(id)
	if err != nil {
		s.log.Error("[UpdateSupplier] failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	supplier.WarehouseItemId = request.WarehouseItemId
	supplier.Name = request.Name
	supplier.PhoneNumber = request.PhoneNumber
	supplier.Address = request.Address
	supplier.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	if err := s.repository.UpdateSupplier(&supplier); err != nil {
		s.log.Error("[UpdateSupplier] failed to update supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	supplier, err = s.repository.GetSupplierById(supplier.Id)
	if err != nil {
		s.log.Error("[UpdateSupplier] failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	return mapper.SupplierToResponse(&supplier), nil
}

func (s *SupplierService) DeleteSupplier(id uint64) error {
	s.repository.UseTx(false)

	if err := s.repository.DeleteSupplier(id); err != nil {
		s.log.Error("[DeleteSupplier] failed to delete supplier", zap.Error(err))
		return err
	}

	return nil
}
