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
		RequestDate: storeRequestItem.CreatedAt.Format("02-01-2006"),
	}
}

func StoreItemToResponse(storeItem *entity.StoreItem) dto.StoreItemResponse {
	return dto.StoreItemResponse{
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

// Note : without payments, payment payment
func StoreSaleToResponse(storeSale *entity.StoreSale) dto.StoreSaleResponse {
	return dto.StoreSaleResponse{
		Id:         storeSale.Id,
		SendDate:   storeSale.SendDate.Format("02-01-2006"),
		Customer:   storeSale.Customer,
		Phone:      storeSale.Phone,
		Price:      storeSale.Price.String(),
		TotalPrice: storeSale.TotalPrice.String(),
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
		WarehouseItem: WarehouseItemToResponse(&storeSale.WarehouseItem),
		Store:         StoreToResponse(&storeSale.Store),
		Quantity:      storeSale.Quantity,
		SaleUnit:      storeSale.SaleUnit.String(),
		PaymentType:   storeSale.PaymentType.String(),
		PaymentStatus: storeSale.PaymentStatus.String(),
		IsSend:        storeSale.IsSend,
	}
}
