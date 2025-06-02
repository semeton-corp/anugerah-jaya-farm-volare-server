package service

import (
	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type EggPriceService struct {
	log        *zap.Logger
	repository repository.IEggPriceRepository
}

type IEggPriceService interface {
	CreateEggPrice(request dto.CreateEggPriceRequest, accountId uuid.UUID) (dto.EggPriceResponse, error)
	GetEggPrices() ([]dto.EggPriceResponse, error)
	GetEggPriceById(id uint64) (dto.EggPriceResponse, error)
	UpdateEggPrice(id uint64, request dto.UpdateEggPriceRequest, accountId uuid.UUID) (dto.EggPriceResponse, error)
	DeleteEggPrice(id uint64) error

	CreateEggPriceDiscount(request dto.CreateEggPriceDiscountRequest, accountId uuid.UUID) (dto.EggPriceDiscountResponse, error)
	GetEggPriceDiscounts() ([]dto.EggPriceDiscountResponse, error)
	GetEggPriceDiscountById(id uint64) (dto.EggPriceDiscountResponse, error)
	UpdateEggPriceDiscount(id uint64, request dto.UpdateEggPriceDiscountRequest, accountId uuid.UUID) (dto.EggPriceDiscountResponse, error)
	DeleteEggPriceDiscount(id uint64) error
}

func NewEggPriceService(log *zap.Logger, repository repository.IEggPriceRepository) IEggPriceService {
	return &EggPriceService{
		log:        log,
		repository: repository,
	}
}

func (s *EggPriceService) CreateEggPrice(request dto.CreateEggPriceRequest, accountId uuid.UUID) (dto.EggPriceResponse, error) {
	s.repository.UseTx(false)

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[CreateEggPrice] failed to parse price", zap.Error(err))
		return dto.EggPriceResponse{}, errx.BadRequest("invalid price format")
	}

	eggPrice := entity.EggPrice{
		Category:        request.Category,
		WarehouseItemId: request.WarehouseItemId,
		Price:           price,
		CreatedBy:       uuid.NullUUID{UUID: accountId, Valid: true},
	}

	err = s.repository.CreateEggPrice(&eggPrice)
	if err != nil {
		s.log.Error("[CreateEggPrice] failed to create egg price", zap.Error(err))
		return dto.EggPriceResponse{}, err
	}

	resp, err := s.repository.GetEggPriceById(eggPrice.Id)
	if err != nil {
		s.log.Error("[CreateEggPrice] failed to get egg price by id", zap.Error(err))
		return dto.EggPriceResponse{}, err
	}

	return mapper.EggPriceToResponse(&resp), nil
}

func (s *EggPriceService) GetEggPrices() ([]dto.EggPriceResponse, error) {
	s.repository.UseTx(false)

	eggPrices, err := s.repository.GetEggPrices()
	if err != nil {
		s.log.Error("[GetEggPrices] failed to get egg prices", zap.Error(err))
		return nil, err
	}

	eggPriceResponses := make([]dto.EggPriceResponse, len(eggPrices))
	for i, eggPrice := range eggPrices {
		eggPriceResponses[i] = mapper.EggPriceToResponse(&eggPrice)
	}

	return eggPriceResponses, nil
}

func (s *EggPriceService) GetEggPriceById(id uint64) (dto.EggPriceResponse, error) {
	s.repository.UseTx(false)

	eggPrice, err := s.repository.GetEggPriceById(id)
	if err != nil {
		s.log.Error("[GetEggPriceById] failed to get egg price by id", zap.Error(err))
		return dto.EggPriceResponse{}, err
	}

	return mapper.EggPriceToResponse(&eggPrice), nil
}

func (s *EggPriceService) UpdateEggPrice(id uint64, request dto.UpdateEggPriceRequest, accountId uuid.UUID) (dto.EggPriceResponse, error) {
	s.repository.UseTx(false)

	eggPrice, err := s.repository.GetEggPriceById(id)
	if err != nil {
		s.log.Error("[UpdateEggPrice] failed to get egg price by id", zap.Error(err))
		return dto.EggPriceResponse{}, err
	}

	eggPrice.Price, err = decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[UpdateEggPrice] failed to parse price", zap.Error(err))
		return dto.EggPriceResponse{}, errx.BadRequest("invalid price format")
	}

	eggPrice.Category = request.Category
	eggPrice.WarehouseItemId = request.WarehouseItemId
	eggPrice.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateEggPrice(&eggPrice)
	if err != nil {
		s.log.Error("[UpdateEggPrice] failed to update egg price", zap.Error(err))
		return dto.EggPriceResponse{}, err
	}

	return mapper.EggPriceToResponse(&eggPrice), nil
}

func (s *EggPriceService) DeleteEggPrice(id uint64) error {
	s.repository.UseTx(false)

	return s.repository.DeleteEggPrice(id)
}

func (s *EggPriceService) CreateEggPriceDiscount(request dto.CreateEggPriceDiscountRequest, accountId uuid.UUID) (dto.EggPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscount := entity.EggPriceDiscount{
		Name:                   request.Name,
		MinimumTransactionUser: request.MinimumTransactionUser,
		TotalDiscount:          request.TotalDiscount,
	}

	err := s.repository.CreateEggPriceDiscount(&eggPriceDiscount)
	if err != nil {
		s.log.Error("[CreateEggPriceDiscount] failed to create egg price discount", zap.Error(err))
		return dto.EggPriceDiscountResponse{}, err
	}

	resp, err := s.repository.GetEggPriceDiscountById(eggPriceDiscount.Id)
	if err != nil {
		s.log.Error("[CreateEggPriceDiscount] failed to get egg price discount by id", zap.Error(err))
		return dto.EggPriceDiscountResponse{}, err
	}

	return mapper.EggPriceDiscountToResponse(&resp), nil
}

func (s *EggPriceService) GetEggPriceDiscounts() ([]dto.EggPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscounts, err := s.repository.GetEggPriceDiscounts()

	if err != nil {
		s.log.Error("[GetEggPriceDiscounts] failed to get egg price discounts", zap.Error(err))
		return nil, err
	}

	eggPriceDiscountResponses := make([]dto.EggPriceDiscountResponse, len(eggPriceDiscounts))
	for i, eggPriceDiscount := range eggPriceDiscounts {
		eggPriceDiscountResponses[i] = mapper.EggPriceDiscountToResponse(&eggPriceDiscount)
	}

	return eggPriceDiscountResponses, nil
}

func (s *EggPriceService) GetEggPriceDiscountById(id uint64) (dto.EggPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscount, err := s.repository.GetEggPriceDiscountById(id)
	if err != nil {
		s.log.Error("[GetEggPriceDiscountById] failed to get egg price discount by id", zap.Error(err))
		return dto.EggPriceDiscountResponse{}, err
	}

	return mapper.EggPriceDiscountToResponse(&eggPriceDiscount), nil
}

func (s *EggPriceService) UpdateEggPriceDiscount(id uint64, request dto.UpdateEggPriceDiscountRequest, accountId uuid.UUID) (dto.EggPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscount, err := s.repository.GetEggPriceDiscountById(id)
	if err != nil {
		s.log.Error("[UpdateEggPriceDiscount] failed to get egg price discount by id", zap.Error(err))
		return dto.EggPriceDiscountResponse{}, err
	}

	eggPriceDiscount.Name = request.Name
	eggPriceDiscount.MinimumTransactionUser = request.MinimumTransactionUser
	eggPriceDiscount.TotalDiscount = request.TotalDiscount
	eggPriceDiscount.UpdatedBy = uuid.NullUUID{UUID: accountId, Valid: true}

	err = s.repository.UpdateEggPriceDiscount(&eggPriceDiscount)
	if err != nil {
		s.log.Error("[UpdateEggPriceDiscount] failed to update egg price discount", zap.Error(err))
		return dto.EggPriceDiscountResponse{}, err
	}

	return mapper.EggPriceDiscountToResponse(&eggPriceDiscount), nil
}

func (s *EggPriceService) DeleteEggPriceDiscount(id uint64) error {
	s.repository.UseTx(false)

	return s.repository.DeleteEggPriceDiscount(id)
}
