package service

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/repository"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/errx"
	"github.com/shopspring/decimal"
	"go.uber.org/zap"
)

type StoreService struct {
	log        *zap.Logger
	repository repository.IStoreRepository
}

type IStoreService interface {
	GetStores() ([]dto.StoreResponse, error)

	CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error)
	GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]dto.StoreRequestItemResponse, error)
	UpdateStoreRequestItem(id uint64, request dto.UpdateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error)

	GetStoreItems() ([]dto.StoreItemResponse, error)

	CreateStoreSale(request dto.CreateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)
	GetStoreSaleById(id uint64) (dto.StoreSaleResponse, error)
	GetStoreSales(filter dto.GetStoreSaleFilter) ([]dto.StoreSaleListResponse, error)
	UpdateStoreSale(id uint64, request dto.UpdateStoreSaleRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)

	CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error)
}

func NewStoreService(log *zap.Logger, repository repository.IStoreRepository) IStoreService {
	return &StoreService{
		log:        log,
		repository: repository,
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
		storeResponses[i] = dto.StoreResponse{
			Id:   store.Id,
			Name: store.Name,
			Location: dto.LocationResponse{
				Id:   store.Location.Id,
				Name: store.Location.Name,
			},
		}
	}

	return storeResponses, nil
}

func (s *StoreService) CreateStoreRequestItem(request dto.CreateStoreRequestItemRequest, accountId uuid.UUID) (dto.StoreRequestItemResponse, error) {
	// Todo : Check if warehouse is hava warehouse item

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

	return dto.StoreRequestItemResponse{
		Id: storeRequestItem.Id,
		Warehouse: dto.WarehouseResponse{
			Id:   storeRequestItem.Warehouse.Id,
			Name: storeRequestItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Warehouse.Location.Id,
				Name: storeRequestItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeRequestItem.WarehouseItem.Id,
			Name:     storeRequestItem.WarehouseItem.Name,
			Category: storeRequestItem.WarehouseItem.Category.String(),
			Unit:     storeRequestItem.WarehouseItem.Unit,
		},
		Store: dto.StoreResponse{
			Id:   storeRequestItem.Store.Id,
			Name: storeRequestItem.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Store.Location.Id,
				Name: storeRequestItem.Store.Location.Name,
			},
		},
		Quantity:    storeRequestItem.Quantity,
		Status:      storeRequestItem.Status.String(),
		RequestDate: storeRequestItem.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *StoreService) GetStoreRequestItemById(id uint64) (dto.StoreRequestItemResponse, error) {
	storeRequestItem, err := s.repository.GetStoreRequestItemById(id)
	if err != nil {
		s.log.Error("[GetStoreRequestItemById] failed to get store request item by id", zap.Error(err))
		return dto.StoreRequestItemResponse{}, err
	}

	return dto.StoreRequestItemResponse{
		Id: storeRequestItem.Id,
		Warehouse: dto.WarehouseResponse{
			Id:   storeRequestItem.Warehouse.Id,
			Name: storeRequestItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Warehouse.Location.Id,
				Name: storeRequestItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeRequestItem.WarehouseItem.Id,
			Name:     storeRequestItem.WarehouseItem.Name,
			Category: storeRequestItem.WarehouseItem.Category.String(),
			Unit:     storeRequestItem.WarehouseItem.Unit,
		},
		Store: dto.StoreResponse{
			Id:   storeRequestItem.Store.Id,
			Name: storeRequestItem.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Store.Location.Id,
				Name: storeRequestItem.Store.Location.Name,
			},
		},
		Quantity:    storeRequestItem.Quantity,
		Status:      storeRequestItem.Status.String(),
		RequestDate: storeRequestItem.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *StoreService) GetStoreRequestItems(filter dto.GetStoreRequestItemFilter) ([]dto.StoreRequestItemResponse, error) {
	storeRequestItems, err := s.repository.GetStoreRequestItems(filter)
	if err != nil {
		s.log.Error("[GetStoreRequestItems] failed to get store request items", zap.Error(err))
		return nil, err
	}

	storeRequestItemResponses := make([]dto.StoreRequestItemResponse, len(storeRequestItems))
	for i, storeRequestItem := range storeRequestItems {
		storeRequestItemResponses[i] = dto.StoreRequestItemResponse{
			Id: storeRequestItem.Id,
			Warehouse: dto.WarehouseResponse{
				Id:   storeRequestItem.Warehouse.Id,
				Name: storeRequestItem.Warehouse.Name,
				Location: dto.LocationResponse{
					Id:   storeRequestItem.Warehouse.Location.Id,
					Name: storeRequestItem.Warehouse.Location.Name,
				},
			},
			WarehouseItem: dto.WarehouseItemResponse{
				Id:       storeRequestItem.WarehouseItem.Id,
				Name:     storeRequestItem.WarehouseItem.Name,
				Category: storeRequestItem.WarehouseItem.Category.String(),
				Unit:     storeRequestItem.WarehouseItem.Unit,
			},
			Store: dto.StoreResponse{
				Id:   storeRequestItem.Store.Id,
				Name: storeRequestItem.Store.Name,
				Location: dto.LocationResponse{
					Id:   storeRequestItem.Store.Location.Id,
					Name: storeRequestItem.Store.Location.Name,
				},
			},
			Quantity:    storeRequestItem.Quantity,
			Status:      storeRequestItem.Status.String(),
			RequestDate: storeRequestItem.CreatedAt.Format("2006-01-02"),
		}
	}

	return storeRequestItemResponses, nil
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

	return dto.StoreRequestItemResponse{
		Id: storeRequestItem.Id,
		Warehouse: dto.WarehouseResponse{
			Id:   storeRequestItem.Warehouse.Id,
			Name: storeRequestItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Warehouse.Location.Id,
				Name: storeRequestItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeRequestItem.WarehouseItem.Id,
			Name:     storeRequestItem.WarehouseItem.Name,
			Category: storeRequestItem.WarehouseItem.Category.String(),
			Unit:     storeRequestItem.WarehouseItem.Unit,
		},
		Store: dto.StoreResponse{
			Id:   storeRequestItem.Store.Id,
			Name: storeRequestItem.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeRequestItem.Store.Location.Id,
				Name: storeRequestItem.Store.Location.Name,
			},
		},
		Quantity:    storeRequestItem.Quantity,
		Status:      storeRequestItem.Status.String(),
		RequestDate: storeRequestItem.CreatedAt.Format("2006-01-02"),
	}, nil
}

func (s *StoreService) GetStoreItems() ([]dto.StoreItemResponse, error) {
	storeItems, err := s.repository.GetStoreItems()
	if err != nil {
		s.log.Error("[GetStoreItem] failed to get store items", zap.Error(err))
		return nil, err
	}

	storeItemResponses := make([]dto.StoreItemResponse, len(storeItems))
	for i, storeItem := range storeItems {
		storeItemResponses[i] = dto.StoreItemResponse{
			Store: dto.StoreResponse{
				Id:   storeItem.Store.Id,
				Name: storeItem.Store.Name,
				Location: dto.LocationResponse{
					Id:   storeItem.Store.Location.Id,
					Name: storeItem.Store.Location.Name,
				},
			},
			WarehouseItem: dto.WarehouseItemResponse{
				Id:       storeItem.WarehouseItem.Id,
				Name:     storeItem.WarehouseItem.Name,
				Category: storeItem.WarehouseItem.Category.String(),
				Unit:     storeItem.WarehouseItem.Unit,
			},
			Quantity:    storeItem.Quantity,
			Description: constant.StoreItemDescriptionDanger, // Todo : give formula for description
		}
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

	paymentMethod := enum.ValueOfPaymentMethod(request.PaymentMethod)
	if !paymentMethod.IsValid() {
		s.log.Error("[CreateStoreSale] invalid payment method", zap.String("paymentMethod", request.PaymentMethod))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid payment method")
	}

	price, err := decimal.NewFromString(request.Price)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse price", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid price format")
	}

	totalPrice := price.Mul(decimal.NewFromInt(int64(request.Quantity)))

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
		PaymentMethod:   paymentMethod,
		CreatedBy:       accountId,
	}

	nominal, err := decimal.NewFromString(request.StoreSalePayment.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	if paymentMethod == enum.PaymentMethodPaidOff {
		if !storeSale.TotalPrice.Equal(nominal) {
			s.log.Error("[CreateStoreSale] nominal is not equal to total price", zap.Error(err))
			return dto.StoreSaleResponse{}, errx.BadRequest("nominal is not equal to total price")
		}
	}

	err = s.repository.CreateStoreSale(&storeSale)
	if err != nil {
		s.log.Error("[CreateStoreSale] failed to create store sale", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	if request.StoreSalePayment.Nominal != "" || request.StoreSalePayment.PaymentProof != "" {
		storeSalePayment := entity.StoreSalePayment{
			StoreSaleId:  storeSale.Id,
			Nominal:      nominal,
			PaymentProof: request.StoreSalePayment.PaymentProof,
			CreatedBy:    accountId,
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
		storeSalePayments[i] = dto.StoreSalePaymentResponse{
			Id:           storeSalePayment.Id,
			Nominal:      storeSalePayment.Nominal.String(),
			PaymentProof: storeSalePayment.PaymentProof,
			Date:         storeSalePayment.CreatedAt.Format("2006-01-02"),
		}

		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
	}

	return dto.StoreSaleResponse{
		Id:       storeSale.Id,
		SendDate: storeSale.SendDate.Format("2006-01-02"),
		Customer: storeSale.Customer,
		Phone:    storeSale.Phone,
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeSale.WarehouseItem.Id,
			Name:     storeSale.WarehouseItem.Name,
			Unit:     storeSale.WarehouseItem.Unit,
			Category: storeSale.WarehouseItem.Category.String(),
		},
		Store: dto.StoreResponse{
			Id:   storeSale.Store.Id,
			Name: storeSale.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeSale.Store.Location.Id,
				Name: storeSale.Store.Location.Name,
			},
		},
		Quantity:         storeSale.Quantity,
		PaymentMethod:    storeSale.PaymentMethod.String(),
		IsSend:           storeSale.IsSend,
		Payments:         storeSalePayments,
		RemainingPayment: remainingPayment.String(),
	}, nil
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
		storeSalePayments[i] = dto.StoreSalePaymentResponse{
			Id:           storeSalePayment.Id,
			Nominal:      storeSalePayment.Nominal.String(),
			PaymentProof: storeSalePayment.PaymentProof,
			Date:         storeSalePayment.CreatedAt.Format("2006-01-02"),
		}

		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
	}

	return dto.StoreSaleResponse{
		Id:       storeSale.Id,
		SendDate: storeSale.SendDate.Format("2006-01-02"),
		Customer: storeSale.Customer,
		Phone:    storeSale.Phone,
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeSale.WarehouseItem.Id,
			Name:     storeSale.WarehouseItem.Name,
			Unit:     storeSale.WarehouseItem.Unit,
			Category: storeSale.WarehouseItem.Category.String(),
		},
		Store: dto.StoreResponse{
			Id:   storeSale.Store.Id,
			Name: storeSale.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeSale.Store.Location.Id,
				Name: storeSale.Store.Location.Name,
			},
		},
		Quantity:         storeSale.Quantity,
		PaymentMethod:    storeSale.PaymentMethod.String(),
		IsSend:           storeSale.IsSend,
		Payments:         storeSalePayments,
		RemainingPayment: remainingPayment.String(),
	}, nil
}

func (s *StoreService) GetStoreSales(filter dto.GetStoreSaleFilter) ([]dto.StoreSaleListResponse, error) {
	storeSales, err := s.repository.GetStoreSales(filter)
	if err != nil {
		s.log.Error("[GetStoreSales] failed to get store sales", zap.Error(err))
		return nil, err
	}

	storeSaleResponses := make([]dto.StoreSaleListResponse, len(storeSales))
	for i, storeSale := range storeSales {
		storeSaleResponses[i] = dto.StoreSaleListResponse{
			Id:       storeSale.Id,
			SendDate: storeSale.SendDate.Format("2006-01-02"),
			Customer: storeSale.Customer,
			Phone:    storeSale.Phone,
			WarehouseItem: dto.WarehouseItemResponse{
				Id:       storeSale.WarehouseItem.Id,
				Name:     storeSale.WarehouseItem.Name,
				Unit:     storeSale.WarehouseItem.Unit,
				Category: storeSale.WarehouseItem.Category.String(),
			},
			Store: dto.StoreResponse{
				Id:   storeSale.Store.Id,
				Name: storeSale.Store.Name,
				Location: dto.LocationResponse{
					Id:   storeSale.Store.Location.Id,
					Name: storeSale.Store.Location.Name,
				},
			},
			Quantity:      storeSale.Quantity,
			PaymentMethod: storeSale.PaymentMethod.String(),
			IsSend:        storeSale.IsSend,
		}
	}

	return storeSaleResponses, nil
}

func (s *StoreService) CreateStoreSalePayment(storeSaleId uint64, request dto.CreateStoreSalePaymentRequest, accountId uuid.UUID) (dto.StoreSaleResponse, error) {
	nominal, err := decimal.NewFromString(request.Nominal)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to parse nominal", zap.Error(err))
		return dto.StoreSaleResponse{}, errx.BadRequest("invalid nominal format")
	}

	storeSalePayment := entity.StoreSalePayment{
		StoreSaleId:  storeSaleId,
		Nominal:      nominal,
		PaymentProof: request.PaymentProof,
		CreatedBy:    accountId,
	}

	err = s.repository.CreateStoreSalePayment(&storeSalePayment)
	if err != nil {
		s.log.Error("[CreateStoreSalePayment] failed to create store sale payment", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSale, err := s.repository.GetStoreSaleById(storeSaleId)
	if err != nil {
		s.log.Error("[GetStoreSaleById] failed to get store sale by id", zap.Error(err))
		return dto.StoreSaleResponse{}, err
	}

	storeSalePayments := make([]dto.StoreSalePaymentResponse, len(storeSale.Payments))

	remainingPayment := storeSale.TotalPrice
	for i, storeSalePayment := range storeSale.Payments {
		storeSalePayments[i] = dto.StoreSalePaymentResponse{
			Id:           storeSalePayment.Id,
			Nominal:      storeSalePayment.Nominal.String(),
			PaymentProof: storeSalePayment.PaymentProof,
			Date:         storeSalePayment.CreatedAt.Format("2006-01-02"),
		}

		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
	}

	return dto.StoreSaleResponse{
		Id:       storeSale.Id,
		SendDate: storeSale.SendDate.Format("2006-01-02"),
		Customer: storeSale.Customer,
		Phone:    storeSale.Phone,
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeSale.WarehouseItem.Id,
			Name:     storeSale.WarehouseItem.Name,
			Unit:     storeSale.WarehouseItem.Unit,
			Category: storeSale.WarehouseItem.Category.String(),
		},
		Store: dto.StoreResponse{
			Id:   storeSale.Store.Id,
			Name: storeSale.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeSale.Store.Location.Id,
				Name: storeSale.Store.Location.Name,
			},
		},
		Quantity:         storeSale.Quantity,
		PaymentMethod:    storeSale.PaymentMethod.String(),
		IsSend:           storeSale.IsSend,
		Payments:         storeSalePayments,
		RemainingPayment: remainingPayment.String(),
	}, nil
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

	storeSale.IsSend = request.IsSend
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
		storeSalePayments[i] = dto.StoreSalePaymentResponse{
			Id:           storeSalePayment.Id,
			Nominal:      storeSalePayment.Nominal.String(),
			PaymentProof: storeSalePayment.PaymentProof,
			Date:         storeSalePayment.CreatedAt.Format("2006-01-02"),
		}

		remainingPayment = remainingPayment.Sub(storeSalePayment.Nominal)
	}

	return dto.StoreSaleResponse{
		Id:       storeSale.Id,
		SendDate: storeSale.SendDate.Format("2006-01-02"),
		Customer: storeSale.Customer,
		Phone:    storeSale.Phone,
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       storeSale.WarehouseItem.Id,
			Name:     storeSale.WarehouseItem.Name,
			Unit:     storeSale.WarehouseItem.Unit,
			Category: storeSale.WarehouseItem.Category.String(),
		},
		Store: dto.StoreResponse{
			Id:   storeSale.Store.Id,
			Name: storeSale.Store.Name,
			Location: dto.LocationResponse{
				Id:   storeSale.Store.Location.Id,
				Name: storeSale.Store.Location.Name,
			},
		},
		Quantity:         storeSale.Quantity,
		IsSend:           storeSale.IsSend,
		Payments:         storeSalePayments,
		RemainingPayment: remainingPayment.String(),
	}, nil
}
