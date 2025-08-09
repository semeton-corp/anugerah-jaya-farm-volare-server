package service

import (
	"slices"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"go.uber.org/zap"
)

type SupplierService struct {
	log        *zap.Logger
	repository repository.ISupplierRepository
}

type ISupplierService interface {
	CreateSupplier(requesst *dto.CreateSupplierRequest, createdBy uuid.UUID) (dto.SupplierResponse, error)
	GetSupplierById(id uint64) (dto.SupplierResponse, error)
	GetSuppliers(filter dto.GetSupplierFilter) ([]dto.SupplierListResponse, error)
	UpdateSupplier(id uint64, request *dto.UpdateSupplierRequest, updatedBy uuid.UUID) (dto.SupplierResponse, error)
	DeleteSupplier(id uint64) error
}

func NewSupplierService(log *zap.Logger, repository repository.ISupplierRepository) ISupplierService {
	return &SupplierService{
		log:        log,
		repository: repository,
	}
}

func (s *SupplierService) CreateSupplier(request *dto.CreateSupplierRequest, createdBy uuid.UUID) (dto.SupplierResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	supplierType := enum.ValueOfSupplierType(request.SupplierType)
	if !supplierType.IsValid() {
		return dto.SupplierResponse{}, errx.BadRequest("invalid supplier type")
	}

	supplier := entity.Supplier{
		Name:         request.Name,
		PhoneNumber:  request.PhoneNumber,
		Address:      request.Address,
		SupplierType: supplierType,
		CreatedBy:    uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err := s.repository.CreateSupplier(&supplier)
	if err != nil {
		s.log.Error("failed to create supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	supplierItems := make([]entity.SupplierItem, 0)
	for _, e := range request.ItemIds {
		supplierItems = append(supplierItems, entity.SupplierItem{
			SupplierId: supplier.Id,
			ItemId:     e,
			CreatedBy:  uuid.NullUUID{UUID: createdBy, Valid: true},
		})
	}

	err = s.repository.CreateSupplierItemInBatch(&supplierItems)
	if err != nil {
		s.log.Error("failed to create supplier items in batch", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transaction", zap.Error(err))
		return dto.SupplierResponse{}, nil
	}

	supplier, err = s.repository.GetSupplierById(supplier.Id)
	if err != nil {
		s.log.Error("failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	return mapper.SupplierToResponse(&supplier), nil
}

func (s *SupplierService) GetSupplierById(id uint64) (dto.SupplierResponse, error) {
	s.repository.UseTx(false)

	supplier, err := s.repository.GetSupplierById(id)
	if err != nil {
		s.log.Error("failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	return mapper.SupplierToResponse(&supplier), nil
}

func (s *SupplierService) GetSuppliers(filter dto.GetSupplierFilter) ([]dto.SupplierListResponse, error) {
	s.repository.UseTx(false)

	suppliers, err := s.repository.GetSuppliers(filter)
	if err != nil {
		s.log.Error("failed to get suppliers", zap.Error(err))
		return nil, err
	}

	supplierResponses := make([]dto.SupplierListResponse, len(suppliers))
	for i, supplier := range suppliers {
		supplierResponses[i] = mapper.SupplierToListResponse(&supplier)
	}

	return supplierResponses, nil
}

func (s *SupplierService) UpdateSupplier(id uint64, request *dto.UpdateSupplierRequest, updatedBy uuid.UUID) (dto.SupplierResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	supplier, err := s.repository.GetSupplierById(id)
	if err != nil {
		s.log.Error("failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	currentItemIds := make([]uint64, 0)
	for _, e := range supplier.SupplierItems {
		currentItemIds = append(currentItemIds, e.ItemId)
	}

	deletedItemIds := make([]uint64, 0)
	for _, e := range currentItemIds {
		if !slices.Contains(request.ItemIds, e) {
			deletedItemIds = append(deletedItemIds, e)
		}
	}

	newItemIds := make([]uint64, 0)
	for _, e := range request.ItemIds {
		if !slices.Contains(currentItemIds, e) {
			newItemIds = append(newItemIds, e)
		}
	}

	supplierType := enum.ValueOfSupplierType(request.SupplierType)
	if !supplierType.IsValid() {
		return dto.SupplierResponse{}, errx.BadRequest("invalid supplier type")
	}

	supplier.Name = request.Name
	supplier.PhoneNumber = request.PhoneNumber
	supplier.Address = request.Address
	supplier.SupplierType = supplierType
	supplier.UpdatedBy = uuid.NullUUID{UUID: updatedBy, Valid: true}

	if err := s.repository.UpdateSupplier(&supplier); err != nil {
		s.log.Error("failed to update supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	if deletedItemIds != nil {
		err := s.repository.DeleteSupplierItemInBatch(deletedItemIds, supplier.Id)
		if err != nil {
			s.log.Error("failed to delete supplier item in batch", zap.Error(err))
			return dto.SupplierResponse{}, err
		}
	}

	if newItemIds != nil {
		supplierItems := make([]entity.SupplierItem, 0)
		for _, e := range newItemIds {
			supplierItems = append(supplierItems, entity.SupplierItem{
				SupplierId: supplier.Id,
				ItemId:     e,
				CreatedBy:  uuid.NullUUID{UUID: updatedBy, Valid: true},
			})
		}
		err := s.repository.CreateSupplierItemInBatch(&supplierItems)
		if err != nil {
			s.log.Error("failed to create supplier item in batch", zap.Error(err))
			return dto.SupplierResponse{}, err
		}
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("failed to commit transcation", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	supplier, err = s.repository.GetSupplierById(supplier.Id)
	if err != nil {
		s.log.Error("failed to get supplier", zap.Error(err))
		return dto.SupplierResponse{}, err
	}

	return mapper.SupplierToResponse(&supplier), nil
}

func (s *SupplierService) DeleteSupplier(id uint64) error {
	s.repository.UseTx(false)

	if err := s.repository.DeleteSupplier(id); err != nil {
		s.log.Error("failed to delete supplier", zap.Error(err))
		return err
	}

	return nil
}
