package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
)

func StoreToResponse(store *entity.Store) dto.StoreResponse {
	return dto.StoreResponse{
		Id:   store.Id,
		Name: store.Name,
		Location: dto.LocationResponse{
			Id:   store.Location.Id,
			Name: store.Location.Name,
		},
		TotalEmployee: uint64(len(store.StorePlacement)),
	}
}

func StoreRequestItemToResponse(storeRequestItem *entity.StoreRequestItem) dto.StoreRequestItemResponse {
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
		WarehouseItem: dto.ItemResponse{
			Id:       storeRequestItem.Item.Id,
			Name:     storeRequestItem.Item.Name,
			Category: storeRequestItem.Item.Category.String(),
			Unit:     storeRequestItem.Item.Unit,
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
		RequestDate: storeRequestItem.CreatedAt.Format("02-01-2006"),
	}
}

func StoreItemToResponse(storeItem *entity.StoreItem) dto.StoreItemResponse {
	response := dto.StoreItemResponse{
		Store:    StoreToResponse(&storeItem.Store),
		Item:     ItemToResponse(&storeItem.Item),
		Quantity: storeItem.Quantity,
	}

	if storeItem.Quantity/float64(constant.TotalEggPerIkat) >= 20.0 {
		response.Description = constant.StoreItemDescriptionSafety
	} else {
		response.Description = constant.StoreItemDescriptionDanger
	}

	return response
}

// Note : without payments, payment payment
func StoreSaleToResponse(storeSale *entity.StoreSale) dto.StoreSaleResponse {
	return dto.StoreSaleResponse{
		Id:         storeSale.Id,
		SendDate:   storeSale.SendDate.Format("02-01-2006"),
		Customer:   storeSale.Customer,
		Phone:      storeSale.Phone,
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
}

// Note : without remaining payment
func StoreSalePaymentToResponse(storeSalePayment *entity.StoreSalePayment) dto.StoreSalePaymentResponse {
	return dto.StoreSalePaymentResponse{
		Id:            storeSalePayment.Id,
		Nominal:       storeSalePayment.Nominal.String(),
		PaymentProof:  storeSalePayment.PaymentProof,
		PaymentMethod: storeSalePayment.PaymentMethod.String(),
		Date:          storeSalePayment.PaymentDate.Format("02-01-2006"),
	}
}

func StoreSaleToListResponse(storeSale *entity.StoreSale) dto.StoreSaleListResponse {
	return dto.StoreSaleListResponse{
		Id:            storeSale.Id,
		SendDate:      storeSale.SendDate.Format("02-01-2006"),
		Customer:      storeSale.Customer,
		Phone:         storeSale.Phone,
		WarehouseItem: ItemToResponse(&storeSale.Item),
		Store:         StoreToResponse(&storeSale.Store),
		Quantity:      storeSale.Quantity,
		SaleUnit:      storeSale.SaleUnit.String(),
		PaymentType:   storeSale.PaymentType.String(),
		PaymentStatus: storeSale.PaymentStatus.String(),
		IsSend:        storeSale.IsSend,
	}
}
