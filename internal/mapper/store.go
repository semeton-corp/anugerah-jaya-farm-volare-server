package mapper

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func StoreToResponse(store *entity.Store) dto.StoreResponse {
	return dto.StoreResponse{
		Id:            store.Id,
		Name:          store.Name,
		Location:      LocationToResponse(&store.Location),
		TotalEmployee: uint64(len(store.StorePlacement)),
	}
}

func StoreRequestItemToResponse(storeRequestItem *entity.StoreRequestItem) dto.StoreRequestItemResponse {
	response := dto.StoreRequestItemResponse{
		Id:                   storeRequestItem.Id,
		Warehouse:            WarehouseToResponse(&storeRequestItem.Warehouse),
		Store:                StoreToResponse(&storeRequestItem.Store),
		Item:                 ItemToResponse(&storeRequestItem.Item),
		Quantity:             storeRequestItem.Quantity,
		Status:               storeRequestItem.Status.String(),
		RequestDate:          storeRequestItem.CreatedAt.Format("15:04, 02 Jan 2006"),
		IsSorted:             storeRequestItem.IsSorted,
		WarehouseFulFillment: storeRequestItem.WarehouseFulfillment,
	}

	if storeRequestItem.ReceiveDate.Valid {
		response.ReceiveDate = storeRequestItem.ReceiveDate.Time.Format("15:04, 02 Jan 2006")
	} else {
		response.ReceiveDate = "-"
	}

	return response
}

func StoreItemToResponse(storeItem *entity.StoreItem) dto.StoreItemResponse {
	response := dto.StoreItemResponse{
		Store:    StoreToResponse(&storeItem.Store),
		Item:     ItemToResponse(&storeItem.Item),
		Quantity: storeItem.Quantity,
	}

	if storeItem.Quantity/float64(constant.TotalEggPerIkat) >= 20.0 {
		response.Description = constant.StoreItemDescriptionSafe
	} else {
		response.Description = constant.StoreItemDescriptionDanger
	}

	return response
}

// Note : without payments
func StoreSaleToResponse(storeSale *entity.StoreSale) dto.StoreSaleResponse {
	response := dto.StoreSaleResponse{
		Id:         storeSale.Id,
		SendDate:   storeSale.SendDate.Format("02 Jan 2006"),
		Customer:   CustomerToResponse(&storeSale.Customer),
		Price:      storeSale.Price.String(),
		TotalPrice: storeSale.TotalPrice.String(),
		WarehouseItem: dto.ItemResponse{
			Id:       storeSale.Item.Id,
			Name:     storeSale.Item.Name,
			Unit:     storeSale.Item.Unit,
			Category: storeSale.Item.Category.String(),
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
		SaleUnit:      storeSale.SaleUnit.String(),
		PaymentType:   storeSale.PaymentType.String(),
		PaymentStatus: storeSale.PaymentStatus.String(),
		IsSend:        storeSale.IsSend,
	}

	if storeSale.DeadlinePaymentDate.Valid {
		response.DeadlinePaymentDate = storeSale.DeadlinePaymentDate.Time.Format("02 Jan 2006")
		if time.Now().After(storeSale.DeadlinePaymentDate.Time) {
			response.IsMoreThanDeadlinePaymentDate = true
		} else {
			response.IsMoreThanDeadlinePaymentDate = false
		}
	} else {
		response.DeadlinePaymentDate = "-"
		response.IsMoreThanDeadlinePaymentDate = false
	}

	return response
}

// Note : without remaining payment
func StoreSalePaymentToResponse(storeSalePayment *entity.StoreSalePayment) dto.StoreSalePaymentResponse {
	return dto.StoreSalePaymentResponse{
		Id:            storeSalePayment.Id,
		Nominal:       storeSalePayment.Nominal.String(),
		PaymentProof:  storeSalePayment.PaymentProof,
		PaymentMethod: storeSalePayment.PaymentMethod.String(),
		Date:          storeSalePayment.PaymentDate.Format("02 Jan 2006"),
	}
}

func StoreSaleToListResponse(storeSale *entity.StoreSale) dto.StoreSaleListResponse {
	response := dto.StoreSaleListResponse{
		Id:            storeSale.Id,
		OrderDate:     storeSale.CreatedAt.Format("02 Jan 2006"),
		SendDate:      storeSale.SendDate.Format("02 Jan 2006"),
		Customer:      CustomerToResponse(&storeSale.Customer),
		Item:          ItemToResponse(&storeSale.Item),
		Store:         StoreToResponse(&storeSale.Store),
		Quantity:      storeSale.Quantity,
		SaleUnit:      storeSale.SaleUnit.String(),
		PaymentStatus: storeSale.PaymentStatus.String(),
		IsSend:        storeSale.IsSend,
	}

	if storeSale.DeadlinePaymentDate.Valid {
		response.DeadlinePaymentDate = storeSale.DeadlinePaymentDate.Time.Format("02 Jan 2006")
		if time.Now().After(storeSale.DeadlinePaymentDate.Time) {
			response.IsMoreThanDeadlinePaymentDate = true
		} else {
			response.IsMoreThanDeadlinePaymentDate = false
		}
	} else {
		response.DeadlinePaymentDate = "-"
		response.IsMoreThanDeadlinePaymentDate = false
	}

	return response
}

func StoreItemHistoryToResponse(storeItemHistory *entity.StoreItemHistory) dto.StoreItemHistoryResponse {
	return dto.StoreItemHistoryResponse{
		Id:             storeItemHistory.Id,
		Item:           ItemToResponse(&storeItemHistory.Item),
		Source:         storeItemHistory.Source,
		Destination:    storeItemHistory.Destination,
		QuantityBefore: storeItemHistory.QuantityBefore,
		QuantityAfter:  storeItemHistory.QuantityAfter,
		Status:         storeItemHistory.Status.String(),
		UpdatedBy:      storeItemHistory.User.Name,
		Date:           storeItemHistory.CreatedAt.Format("02-Jan-2006"),
		Time:           storeItemHistory.CreatedAt.Format("15:04"),
	}
}

func StoreItemHistoryToListResponse(storeItemHistory *entity.StoreItemHistory) dto.StoreItemHistoryListResponse {
	return dto.StoreItemHistoryListResponse{
		Id:          storeItemHistory.Id,
		Item:        ItemToResponse(&storeItemHistory.Item),
		Source:      storeItemHistory.Source,
		Destination: storeItemHistory.Destination,
		Status:      storeItemHistory.Status.String(),
		Quantity:    storeItemHistory.QuantityAfter - storeItemHistory.QuantityBefore,
		Time:        storeItemHistory.CreatedAt.Format("15:04"),
	}
}

func StoreSaleQueueToResponse(storeSaleQueue *entity.StoreSaleQueue) dto.StoreSaleQueueResponse {
	response := dto.StoreSaleQueueResponse{
		Id:           storeSaleQueue.Id,
		Item:         ItemToResponse(&storeSaleQueue.Item),
		Store:        StoreToResponse(&storeSaleQueue.Store),
		Quantity:     storeSaleQueue.Quantity,
		CustomerType: storeSaleQueue.CustomerType.String(),
		SaleUnit:     storeSaleQueue.SaleUnit.String(),
	}

	if storeSaleQueue.CustomerType == enum.CustomerTypeNew {
		response.Customer = dto.CustomerResponse{
			Name:        storeSaleQueue.CustomerName.String,
			PhoneNumber: storeSaleQueue.CustomerPhoneNumber.String,
		}
	} else {
		response.Customer = CustomerToResponse(&storeSaleQueue.Customer)
	}

	return response
}
