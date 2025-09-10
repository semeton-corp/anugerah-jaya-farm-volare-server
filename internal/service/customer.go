package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"go.uber.org/zap"
)

type CustomerService struct {
	log        *zap.Logger
	repository repository.ICustomerRepository
}

type ICustomerService interface {
	GetCustomers() ([]dto.CustomerResponse, error)
	CreateCustomer(request dto.CreateCustomerRequest, userId uuid.UUID) (dto.CustomerResponse, error)
	DeleteCustomer(id uint64) error
}

func NewCustomerService(log *zap.Logger, repository repository.ICustomerRepository) ICustomerService {
	return &CustomerService{
		log:        log,
		repository: repository,
	}
}

func (s *CustomerService) GetCustomers() ([]dto.CustomerResponse, error) {
	s.repository.UseTx(false)

	customers, err := s.repository.GetCustomers()
	if err != nil {
		s.log.Error("failed to get customers", zap.Error(err))
		return nil, err
	}

	responseCustomers := make([]dto.CustomerResponse, len(customers))
	for i, Customer := range customers {
		responseCustomers[i] = mapper.CustomerToResponse(&Customer)
	}

	return responseCustomers, nil
}

func (s *CustomerService) CreateCustomer(request dto.CreateCustomerRequest, userId uuid.UUID) (dto.CustomerResponse, error) {
	s.repository.UseTx(false)

	customer := entity.Customer{
		Name:        request.Name,
		PhoneNumber: request.PhoneNumber,
		CreatedBy:   uuid.NullUUID{UUID: userId, Valid: true},
	}

	err := s.repository.CreateCustomer(&customer)
	if err != nil {
		s.log.Error("failed to create customer", zap.Error(err))
		return dto.CustomerResponse{}, err
	}

	customer, err = s.repository.GetCustomerById(customer.Id)
	if err != nil {
		s.log.Error("failed to get customer by id", zap.Error(err))
		return dto.CustomerResponse{}, err
	}

	return mapper.CustomerToResponse(&customer), nil
}

func (s *CustomerService) DeleteCustomer(id uint64) error {
	s.repository.UseTx(false)
	err := s.repository.DeleteCustomer(id)
	if err != nil {
		s.log.Error("failed to delete customer", zap.Error(err))
		return err
	}

	return nil
}
