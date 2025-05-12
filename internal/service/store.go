package service

import (
	"math"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/infra/cache"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/mapper"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type StoreService struct {
	log          *zap.Logger
	repository   repository.IStoreRepository
	cacheService cache.ICache
}

type IStoreService interface {
	GetStores() ([]dto.StoreResponse, error)

	CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) (dto.StoreRequestItemListPaginationResponse, error)
	UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error)

	GetStoreItems(filter dto.GetStoreItemFilter) ([]dto.StoreItemResponse, error)

	CreateStoreSale(request dto.CreateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)
	GetStoreSaleById(id uint64) (dto.StoreSaleResponse, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) (dto.StoreSaleListPaginationResponse, error)
	UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)

	CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)
	UpdateStoreSalePayment(id uint64, request dto.UpdateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)

	SendStoreSale(id uint64, accountId uuid.UUID) (dto.StoreSaleResponse, error)
}

func NewStoreService(log *zap.Logger, repository repository.IStoreRepository, cacheService cache.ICache) IStoreService {
	return &StoreService{
		log:          log,
		repository:   repository,
		cacheService: cacheService,
	}
}

func (s *StoreService) GetStores() ([]dto.StoreResponse, error) {
	stores, err := s.repository.GetStores()
	if err != nil {
		s.log.Error("[GetStores] failed to get stores", zap.Error(err))
		return nil, err
	}

	storeResponses := make([]dto.StoreResponse, len(stores))
	for i, store := range stores {
		storeResponses[i] = mapper.StoreToResponse(&store)
	}

	return storeResponses, nil
}

func (s *StoreService) CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error) {
	// Todo : Check if warehouse is have warehouse item

	storeRequestItem := entity.StoreRequestItem{
		WarehouseId:     request.WarehouseId,
		WarehouseItemId: request.WarehouseItemId,
		StoreId:         request.StoreId,
		Quantity:        request.Quantity,
		Status:          enum.RequestItemStatusPending,
		CreatedBy:       accountId,
	}

	err := s.repository.CreateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("[CreateStoreRequestItem] failed to create store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	storeRequestItem, err = s.repository.GetStoreRequestItemById(storeRequestItem.Id)
	if err != nil {
		s.log.Error("[CreateStoreRequestItem] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return mapper.StoreRequestItemToResponse(&storeRequestItem), nil
}

func (s *StoreService) GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error) {
	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("[GetStoreRequestItemById] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return mapper.StoreRequestItemToResponse(&storeRequestItem), nil
}

func (s *StoreService) GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) (dto.StoreRequestItemListPaginationResponse, error) {
	storeRequestItems, err := s.repository.GetStoreRequestItems(filter)
	if err != nil {
		s.log.Error("[GetStoreRequestItems] failed to get store request items", zap.Error(err))
		return dto.StoreRequestItemListPaginationResponse{}, err
	}

	storeRequestItemResponses := make([]dto.StoreRequestItemResponse, len(storeRequestItems))
	for i, storeRequestItem := range storeRequestItems {
		storeRequestItemResponses[i] = mapper.StoreRequestItemToResponse(&storeRequestItem)
	}

	totalData, err := s.repository.CountTotalStoreRequestItem(dto.GetStoreRequestItemFilter{
		Date: filter.Date,
	})
	if err != nil {
		s.log.Error("[GetStoreRequestItems] failed to count request items", zap.Error(err))
		return dto.StoreRequestItemListPaginationResponse{}, err
	}

	resp := dto.StoreRequestItemListPaginationResponse{
		TotalPage:         uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData:         totalData,
		StoreRequestItems: storeRequestItemResponses,
	}

	return resp, nil
}

func (s *StoreService) UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	status := enum.ValueOfRequestItemStatus(request.Status)
	if !status.IsValid() {
		s.log.Error("[UpdateStoreRequestItem] invalid status", zap.String("status", request.Status))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("invalid status")
	}

	if storeRequestItem.Status == enum.RequestItemStatusAccepted || storeRequestItem.Status == enum.RequestItemStatusRejected {
		s.log.Error("[UpdateStoreRequestItem] store request item is already accepted or rejected", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is already accepted or rejected")
	}

	if status == enum.RequestItemStatusAccepted && storeRequestItem.Status != enum.RequestItemStatusSentOff {
		s.log.Error("[UpdateStoreRequestItem] store request item is not sent off", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is not sent off")
	}

	if status == enum.RequestItemStatusSentOff && storeRequestItem.Status != enum.RequestItemStatusPending {
		s.log.Error("[UpdateStoreRequestItem] store request item is not pending", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is not pending")
	}

	if status == enum.RequestItemStatusPending && storeRequestItem.Status != enum.RequestItemStatusPending {
		s.log.Error("[UpdateStoreRequestItem] store request item is not pending", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("store request item is not pending")
	}

	if status == enum.RequestItemStatusAccepted {
		// Todo : Create store activity
		storeItem := entity.StoreItem{
			WarehouseItemId: storeRequestItem.WarehouseItemId,
			StoreId:         storeRequestItem.StoreId,
		}

		err = s.repository.FirstOrCreateStoreItem(&storeItem)
		if err != nil {
			s.log.Error("[UpdateStoreRequestItem] failed to first or create store item", zap.Error(err))
			return dto.StoreRequestItemResponse{}, err
		}

		storeItem.Quantity += storeRequestItem.Quantity
		err = s.repository.UpdateStoreItem(&storeItem)
		if err != nil {
			s.log.Error("[UpdateStoreRequestItem] failed to update store item", zap.Error(err))
			return dto.StoreRequestItemResponse{}, err
		}
	}

	if storeRequestItem.Status == enum.RequestItemStatusPending {
		storeRequestItem.Quantity = request.Quantity
	} else if status != enum.RequestItemStatusAccepted {
		s.log.Error("[UpdateStoreRequestItem] can't update quantity when status is not pending", zap.Uint64("id", id))
		return dto.StoreRequestItemResponse{}, errx.BadRequest("can't update quantity when status is not pending")
	}

	storeRequestItem.Status = status
	storeRequestItem.UpdatedBy = accountId

	err = s.repository.UpdateStoreRequestItem(&storeRequestItem)
	if err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to update store request item", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	storeRequestItem, err = s.repository.GetStoreRequestItemById(storeRequestItem.Id)
	if err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStoreRequestItem] failed to commit transaction", zap.Error(err))
	}

	return mapper.StoreRequestItemToResponse(&storeRequestItem), nil
}

func (s *StoreService) GetStoreItems(filter dto.GetStoreItemFilter) ([]dto.StoreItemResponse, error) {
	storeItems, err := s.repository.GetStoreItems(filter)
	if err != nil {
		s.log.Error("[GetStoreItem] failed to get store items", zap.Error(err))
		return nil, err
	}

	storeItemResponses := make([]dto.StoreItemResponse, len(storeItems))
	for i, storeItem := range storeItems {
		storeItemResponses[i] = mapper.StoreItemToResponse(&storeItem)
	}

	return storeItemResponses, nil
}

func (s *StoreService) CreateStoreSale(request dto.CreateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	// Todo : reduce stock in warehouse item id in store id

	s.repository.UseTx(true)
	defer s.repository.Rollback()

	sendDate, err := time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse sent date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sent date format")
	}

	paymentType := enum.ValueOfPaymentType(request.PaymentType)
	if !paymentType.IsValid() {
		s.log.Error("[CreateStoreSale] invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid price format")
	}

	totalPrice := price.Mul(decimal.NewFromInt(int64(request.Quantity)))

	saleUnit := enum.ValueOfSaleUnit(request.SaleUnit)
	if !saleUnit.IsValid() {
		s.log.Error("[CreateStoreSale] invalid sale unit", zap.String("saleUnit", request.SaleUnit))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	storeSale := entity.StoreSale{
		Customer:        request.Customer,
		Phone:           request.Phone,
		StoreId:         request.StoreId,
		WarehouseItemId: request.WarehouseItemId,
		Quantity:        request.Quantity,
		Price:           price,
		TotalPrice:      totalPrice,
		SendDate:        sendDate,
		IsSend:          false,
		SaleUnit:        saleUnit,
		PaymentType:     paymentType,
		CreatedBy:       accountId,
	}

	nominal, err := decimal.NewFromString(request.StoreSalePayment.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	if paymentType == enum.PaymentTypePaidOff {
		if !storeSale.TotalPrice.Equal(nominal) {
			s.log.Error("[CreateStoreSale] nominal is not equal to total price", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}

		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	err = s.repository.CreateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to create store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if request.StoreSalePayment.Nominal != "" &&
		request.StoreSalePayment.PaymentDate != "" &&
		request.StoreSalePayment.PaymentProof != "" &&
		request.StoreSalePayment.PaymentMethod != "" {
		paymentMethod := enum.ValueOfPaymentMethod(request.StoreSalePayment.PaymentMethod)
		if !paymentMethod.IsValid() {
			s.log.Error("[CreateStoreSale] invalid payment method", zap.String("paymentMethod", request.StoreSalePayment.PaymentMethod))
			return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
		}

		paymentDate, err := time.Parse("02-01-2006", request.StoreSalePayment.PaymentDate)
		if err != nil {
			s.log.Error("[CreateStoreSale] failed to parse payment date", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
		}

		storeSalePayment := entity.StoreSalePayment{
			PaymentDate:   paymentDate,
			StoreSaleId:   storeSale.Id,
			Nominal:       nominal,
			PaymentProof:  request.StoreSalePayment.PaymentProof,
			PaymentMethod: paymentMethod,
			CreatedBy:     accountId,
		}

		err = s.repository.CreateStoreSalePayment(&storeSalePayment)
		if err != nil {
			s.log.Error("[CreateStoreSale] failed to create store sale payment", zap.Error(err))
			return dto.StoreSaleResponse{}, err
		}
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.Commit()
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))
	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) GetStoreSaleById(id uint64) (dto.StoreSaleResponse, error) {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) GetStoreSales(filter dto.GetStoreSaleFilter) (dto.StoreSaleListPaginationResponse, error) {
	storeSales, err := s.repository.GetStoreSales(filter)
	if err != nil {
		s.log.Error("[GetStoreSales] failed to get store sales", zap.Error(err))
		return dto.StoreSaleListPaginationResponse{}, err
	}

	storeSaleResponses := make([]dto.StoreSaleListResponse, len(storeSales))
	for i, storeSale := range storeSales {
		storeSaleResponses[i] = mapper.StoreSaleToListResponse(&storeSale)
	}

	totalData, err := s.repository.CountTotalStoreSale(
		dto.GetStoreSaleFilter{
			Date:          filter.Date,
			PaymentMethod: filter.PaymentMethod,
		},
	)
	if err != nil {
		s.log.Error("[GetStoreSales] failed to get store sales", zap.Error(err))
		return dto.StoreSaleListPaginationResponse{}, err
	}

	resp := dto.StoreSaleListPaginationResponse{
		TotalPage:  uint64(math.Ceil(float64(totalData) / float64(constant.PaginationDefaultLimit))),
		TotalData:  totalData,
		StoreSales: storeSaleResponses,
	}

	return resp, nil
}

func (s *StoreService) CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("[CreateStoreSalePayment] invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse payment date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	storeSalePayment := entity.StoreSalePayment{
		StoreSaleId:   storeSaleId,
		PaymentDate:   paymentDate,
		PaymentMethod: paymentMethod,
		Nominal:       nominal,
		PaymentProof:  request.PaymentProof,
		CreatedBy:     accountId,
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSaleId)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("[CreateStoreSalePayment] store sale is already paid", zap.Uint64("id", storeSaleId))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already paid")
	}

	totalPayment := nominal
	for _, payment := range storeSale.Payments {
		totalPayment = totalPayment.Add(payment.Nominal)
	}

	if totalPayment.Equal(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(storeSale.TotalPrice) {
		s.log.Error("[CreateStoreSalePayment] total payment is greater than total price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	}

	err = s.repository.CreateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to create store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale.Payments = append(storeSale.Payments, storeSalePayment)
	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("[UpdateStoreSale] store sale is already sent", zap.Uint64("id", id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already sent")
	}

	storeSale.Customer = request.Customer
	storeSale.Phone = request.Phone
	storeSale.StoreId = request.StoreId
	storeSale.WarehouseItemId = request.WarehouseItemId
	storeSale.Quantity = request.Quantity
	storeSale.Price, err = decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to parse price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid price format")
	}

	storeSale.TotalPrice = storeSale.Price.Mul(decimal.NewFromInt(int64(request.Quantity)))

	storeSale.SaleUnit = enum.ValueOfSaleUnit(request.SaleUnit)
	if !storeSale.SaleUnit.IsValid() {
		s.log.Error("[UpdateStoreSale] invalid sale unit", zap.String("saleUnit", request.SaleUnit))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid sale unit")
	}

	storeSale.PaymentType = enum.ValueOfPaymentType(request.PaymentType)
	if !storeSale.PaymentType.IsValid() {
		s.log.Error("[UpdateStoreSale] invalid payment type", zap.String("paymentType", request.PaymentType))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment type")
	}

	storeSale.SendDate, err = time.Parse("02-01-2006", request.SendDate)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to parse send date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid send date format")
	}

	storeSale.UpdatedBy = accountId

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[UpdateStoreSale] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) UpdateStoreSalePayment(id uint64, request dto.UpdateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	s.repository.UseTx(true)
	defer s.repository.Rollback()

	storeSalePayment, err := s.repository.GetStoreSalePaymentById(id)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to get store sale payment by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSalePayment.StoreSaleId)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("[UpdateStoreSalePayment] store sale is already sent", zap.Uint64("id", storeSale.Id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already sent")
	}

	if storeSale.PaymentStatus == enum.PaymentStatusPaid {
		s.log.Error("[CreateStoreSalePayment] store sale is already paid", zap.Uint64("id", storeSale.Id))
		return dto.StoreSaleResponse{}, errx.BadRequest("store sale is already paid")
	}

	paymentDate, err := time.Parse("02-01-2006", request.PaymentDate)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse payment date", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment date format")
	}

	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	totalPayment := nominal
	for _, payment := range storeSale.Payments {
		if payment.Id != storeSalePayment.Id {
			totalPayment = totalPayment.Add(payment.Nominal)
		}
	}

	if totalPayment.Equal(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusPaid
	} else if totalPayment.GreaterThan(storeSale.TotalPrice) {
		s.log.Error("[CreateStoreSalePayment] total payment is greater than total price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("total payment is greater than total price")
	} else if totalPayment.LessThan(storeSale.TotalPrice) {
		storeSale.PaymentStatus = enum.PaymentStatusUnpaid
	}

	storeSalePayment.Nominal = nominal
	storeSalePayment.PaymentProof = request.PaymentProof
	storeSalePayment.PaymentDate = paymentDate
	storeSalePayment.UpdatedBy = accountId

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	err = s.repository.UpdateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to update store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if err := s.repository.Commit(); err != nil {
		s.log.Error("[UpdateStoreSalePayment] failed to commit transaction", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, payment := range storeSale.Payments {
		if payment.Id == id {
			storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
			remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
			storeSalePayments[i].Remaining = remainingPayment.String()
		} else {
			storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&payment)
			remainingPayment = remainingPayment.Sub(payment.Nominal)
			storeSalePayments[i].Remaining = remainingPayment.String()
		}
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}

func (s *StoreService) SendStoreSale(id uint64, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	storeSale, err := s.repository.GetStoreSaleById(id)
	if err != nil {
		s.log.Error("[SendStoreSale] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if storeSale.IsSend {
		s.log.Error("[SendStoreSale] store sale is already sent", zap.Uint64("id", id))
		return dto.StoreSaleResponse{}, err
	}

	storeSale.IsSend = true
	storeSale.UpdatedBy = accountId

	err = s.repository.UpdateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[SendStoreSale] failed to update store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err = s.repository.GetStoreSaleById(storeSale.Id)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = mapper.StoreSalePaymentToResponse(&storeSalePayment)
		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
		storeSalePayments[i].Remaining = remainingPayment.String()
	}

	storeSaleResponse := mapper.StoreSaleToResponse(&storeSale)
	storeSaleResponse.Payments = storeSalePayments
	storeSaleResponse.RemainingPayment = remainingPayment.String()

	return storeSaleResponse, nil
}
