package service

import (
	"database/sql"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type ItemService struct {
	log        *zap.Logger
	repository repository.IItemRepository
}

type IItemService interface {
	CreateItemPrice(request dto.CreateItemPriceRequest, createdBy uuid.UUID) (dto.ItemPriceResponse, error)
	GetItemPrices() ([]dto.ItemPriceResponse, error)
	GetItemPriceById(id uint64) (dto.ItemPriceResponse, error)
	UpdateItemPrice(id uint64, request dto.UpdateItemPriceRequest, updatedBy uuid.UUID) (dto.ItemPriceResponse, error)
	DeleteItemPrice(id uint64) error

	CreateItemDiscount(request dto.CreateItemPriceDiscountRequest, createdBy uuid.UUID) (dto.ItemPriceDiscountResponse, error)
	GetItemDiscounts() ([]dto.ItemPriceDiscountResponse, error)
	GetItemDiscountById(id uint64) (dto.ItemPriceDiscountResponse, error)
	UpdateItemDiscount(id uint64, request dto.UpdateItemPriceDiscountRequest, createdBy uuid.UUID) (dto.ItemPriceDiscountResponse, error)
	DeleteItemDiscount(id uint64) error

	GetItemByNameAndUnitAndType(name string, unit string, itemType enum.ItemCategory) (dto.ItemResponse, error)
	CreateItem(request dto.CreateItemRequest, createdBy uuid.UUID) (dto.ItemResponse, error)
	GetItems(filter dto.GetItemFilter) ([]dto.ItemResponse, error)
	UpdateItem(warehouseItemId uint64, request dto.UpdateItemRequest, updatedBy uuid.UUID) (dto.ItemResponse, error)
	GetItemById(id uint64) (dto.ItemResponse, error)
	DeleteItem(id uint64) error
}

func NewItemPriceService(log *zap.Logger, repository repository.IItemRepository) IItemService {
	return &ItemService{
		log:        log,
		repository: repository,
	}
}

func (s *ItemService) CreateItemPrice(request dto.CreateItemPriceRequest, createdBy uuid.UUID) (dto.ItemPriceResponse, error) {
	s.repository.UseTx(false)

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price", zap.Error(err))
		return dto.ItemPriceResponse{}, errx.BadRequest("invalid price format")
	}

	eggPrice := entity.ItemPrice{
		Category:  request.Category,
		ItemId:    request.ItemId,
		Price:     price,
		CreatedBy: uuid.NullUUID{UUID: createdBy, Valid: true},
	}

	err = s.repository.CreateItemPrice(&eggPrice)
	if err != nil {
		s.log.Error("failed to create item price", zap.Error(err))
		return dto.ItemPriceResponse{}, err
	}

	resp, err := s.repository.GetItemPriceById(eggPrice.Id)
	if err != nil {
		s.log.Error("failed to get item price by id", zap.Error(err))
		return dto.ItemPriceResponse{}, err
	}

	return mapper.ItemPriceToResponse(&resp), nil
}

func (s *ItemService) GetItemPrices() ([]dto.ItemPriceResponse, error) {
	s.repository.UseTx(false)

	eggPrices, err := s.repository.GetItemPrices()
	if err != nil {
		s.log.Error("failed to get item prices", zap.Error(err))
		return nil, err
	}

	eggPriceResponses := make([]dto.ItemPriceResponse, len(eggPrices))
	for i, eggPrice := range eggPrices {
		eggPriceResponses[i] = mapper.ItemPriceToResponse(&eggPrice)
	}

	return eggPriceResponses, nil
}

func (s *ItemService) GetItemPriceById(id uint64) (dto.ItemPriceResponse, error) {
	s.repository.UseTx(false)

	eggPrice, err := s.repository.GetItemPriceById(id)
	if err != nil {
		s.log.Error("failed to get item price by id", zap.Error(err))
		return dto.ItemPriceResponse{}, err
	}

	return mapper.ItemPriceToResponse(&eggPrice), nil
}

func (s *ItemService) UpdateItemPrice(id uint64, request dto.UpdateItemPriceRequest, userId uuid.UUID) (dto.ItemPriceResponse, error) {
	s.repository.UseTx(false)

	eggPrice, err := s.repository.GetItemPriceById(id)
	if err != nil {
		s.log.Error("failed to get item price by id", zap.Error(err))
		return dto.ItemPriceResponse{}, err
	}

	eggPrice.Price, err = decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("failed to parse price", zap.Error(err))
		return dto.ItemPriceResponse{}, errx.BadRequest("invalid price format")
	}

	eggPrice.Category = request.Category
	eggPrice.ItemId = request.ItemId
	eggPrice.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateItemPrice(&eggPrice)
	if err != nil {
		s.log.Error("failed to update item price", zap.Error(err))
		return dto.ItemPriceResponse{}, err
	}

	return mapper.ItemPriceToResponse(&eggPrice), nil
}

func (s *ItemService) DeleteItemPrice(id uint64) error {
	s.repository.UseTx(false)

	return s.repository.DeleteItemPrice(id)
}

func (s *ItemService) CreateItemDiscount(request dto.CreateItemPriceDiscountRequest, userId uuid.UUID) (dto.ItemPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscount := entity.ItemPriceDiscount{
		Name:                   request.Name,
		MinimumTransactionUser: request.MinimumTransactionUser,
		TotalDiscount:          request.TotalDiscount,
	}

	err := s.repository.CreateItemPriceDiscount(&eggPriceDiscount)
	if err != nil {
		s.log.Error("failed to create item price discount", zap.Error(err))
		return dto.ItemPriceDiscountResponse{}, err
	}

	resp, err := s.repository.GetItemPriceDiscountById(eggPriceDiscount.Id)
	if err != nil {
		s.log.Error("failed to get item price discount by id", zap.Error(err))
		return dto.ItemPriceDiscountResponse{}, err
	}

	return mapper.ItemPriceDiscountToResponse(&resp), nil
}

func (s *ItemService) GetItemDiscounts() ([]dto.ItemPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscounts, err := s.repository.GetItemPriceDiscounts()

	if err != nil {
		s.log.Error("failed to get item price discounts", zap.Error(err))
		return nil, err
	}

	eggPriceDiscountResponses := make([]dto.ItemPriceDiscountResponse, len(eggPriceDiscounts))
	for i, eggPriceDiscount := range eggPriceDiscounts {
		eggPriceDiscountResponses[i] = mapper.ItemPriceDiscountToResponse(&eggPriceDiscount)
	}

	return eggPriceDiscountResponses, nil
}

func (s *ItemService) GetItemDiscountById(id uint64) (dto.ItemPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscount, err := s.repository.GetItemPriceDiscountById(id)
	if err != nil {
		s.log.Error("failed to get item price discount by id", zap.Error(err))
		return dto.ItemPriceDiscountResponse{}, err
	}

	return mapper.ItemPriceDiscountToResponse(&eggPriceDiscount), nil
}

func (s *ItemService) UpdateItemDiscount(id uint64, request dto.UpdateItemPriceDiscountRequest, userId uuid.UUID) (dto.ItemPriceDiscountResponse, error) {
	s.repository.UseTx(false)

	eggPriceDiscount, err := s.repository.GetItemPriceDiscountById(id)
	if err != nil {
		s.log.Error("failed to get item price discount by id", zap.Error(err))
		return dto.ItemPriceDiscountResponse{}, err
	}

	eggPriceDiscount.Name = request.Name
	eggPriceDiscount.MinimumTransactionUser = request.MinimumTransactionUser
	eggPriceDiscount.TotalDiscount = request.TotalDiscount
	eggPriceDiscount.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	err = s.repository.UpdateItemPriceDiscount(&eggPriceDiscount)
	if err != nil {
		s.log.Error("failed to update item price discount", zap.Error(err))
		return dto.ItemPriceDiscountResponse{}, err
	}

	return mapper.ItemPriceDiscountToResponse(&eggPriceDiscount), nil
}

func (s *ItemService) DeleteItemDiscount(id uint64) error {
	s.repository.UseTx(false)

	err := s.repository.DeleteItemPriceDiscount(id)
	if err != nil {
		s.log.Error("failed to delete item price discount")
		return err
	}

	return nil
}

func (s *ItemService) GetItemByNameAndUnitAndType(name string, unit string, itemType enum.ItemCategory) (dto.ItemResponse, error) {
	s.repository.UseTx(false)

	stockWarehouseItem, err := s.repository.GetItemByNameAndUnitAndType(name, unit, itemType)
	if err != nil {
		s.log.Error("failed to get item by name and unit and type", zap.Error(err))
		return dto.ItemResponse{}, err
	}

	warehouseStockItemResponse := mapper.ItemToResponse(&stockWarehouseItem)
	return warehouseStockItemResponse, nil
}

func (s *ItemService) CreateItem(request dto.CreateItemRequest, userId uuid.UUID) (dto.ItemResponse, error) {
	s.repository.UseTx(false)

	itemCategory := enum.ValueOfItemCategory(request.Category)
	if !itemCategory.IsValid() {
		s.log.Error("invalid warehouse item category", zap.String("category", request.Category))
		return dto.ItemResponse{}, errx.BadRequest("invalid warehouse item category")
	}

	if (itemCategory != enum.ItemCategoryEgg && itemCategory != enum.ItemCategoryChicken) && request.DailySpending == nil {
		return dto.ItemResponse{}, errx.BadRequest("daily spending is required")
	}

	item := entity.Item{
		Name:      request.Name,
		Unit:      request.Unit,
		Category:  itemCategory,
		CreatedBy: uuid.NullUUID{UUID: userId, Valid: true},
	}

	if request.DailySpending != nil {
		item.DailySpending = sql.NullFloat64{Float64: *request.DailySpending, Valid: true}
	}

	err := s.repository.CreateItem(&item)
	if err != nil {
		s.log.Error("failed to create warehouse item", zap.Error(err))
		return dto.ItemResponse{}, err
	}

	return mapper.ItemToResponse(&item), nil
}

func (s *ItemService) GetItems(filter dto.GetItemFilter) ([]dto.ItemResponse, error) {
	s.repository.UseTx(false)

	items, err := s.repository.GetItems(filter)
	if err != nil {
		s.log.Error("failed to get warehouse items", zap.Error(err))
		return nil, err
	}

	warehouseItemResponses := make([]dto.ItemResponse, 0, len(items))
	for _, item := range items {
		warehouseItemResponses = append(warehouseItemResponses, mapper.ItemToResponse(&item))
	}

	return warehouseItemResponses, nil
}

func (s *ItemService) UpdateItem(warehouseItemId uint64, request dto.UpdateItemRequest, userId uuid.UUID) (dto.ItemResponse, error) {
	s.repository.UseTx(false)

	itemCategory := enum.ValueOfItemCategory(request.Category)
	if !itemCategory.IsValid() {
		s.log.Error("invalid warehouse item category", zap.String("category", request.Category))
		return dto.ItemResponse{}, errx.BadRequest("invalid warehouse item category")
	}

	if (itemCategory != enum.ItemCategoryEgg && itemCategory != enum.ItemCategoryChicken) && request.DailySpending == nil {
		return dto.ItemResponse{}, errx.BadRequest("daily spending is required")
	}

	item, err := s.repository.GetItemById(warehouseItemId)
	if err != nil {
		s.log.Error("failed to get warehouse item", zap.Error(err))
		return dto.ItemResponse{}, err
	}

	item.Name = request.Name
	item.Unit = request.Unit
	item.Category = itemCategory

	item.UpdatedBy = uuid.NullUUID{UUID: userId, Valid: true}

	if request.DailySpending != nil {
		item.DailySpending = sql.NullFloat64{Float64: *request.DailySpending, Valid: true}
	}

	err = s.repository.UpdateItem(&item)
	if err != nil {
		s.log.Error("failed to update warehouse item", zap.Error(err))
		return dto.ItemResponse{}, err
	}

	return mapper.ItemToResponse(&item), nil
}

func (s *ItemService) GetItemById(id uint64) (dto.ItemResponse, error) {
	s.repository.UseTx(false)

	item, err := s.repository.GetItemById(id)
	if err != nil {
		s.log.Error("failed to get warehouse item", zap.Error(err))
		return dto.ItemResponse{}, err
	}

	warehouseItemResponse := mapper.ItemToResponse(&item)
	return warehouseItemResponse, nil
}

func (s *ItemService) DeleteItem(id uint64) error {
	s.repository.UseTx(false)
	err := s.repository.DeleteItem(id)
	if err != nil {
		s.log.Error("failed delete item", zap.Error(err))
		return err
	}

	return nil
}
