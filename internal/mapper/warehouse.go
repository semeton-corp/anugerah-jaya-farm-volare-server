package mapper

import (
	"fmt"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

func WarehouseToResponse(warehouse *entity.Warehouse) dto.WarehouseResponse {
	return dto.WarehouseResponse{
		Id:           warehouse.Id,
		Name:         warehouse.Name,
		CornCapacity: warehouse.CornCapacity,
		Location: dto.LocationResponse{
			Id:   warehouse.Location.Id,
			Name: warehouse.Location.Name,
		},
		TotalEmployee: uint64(len(warehouse.WarehousePlacement)),
	}
}

// Todo : fix this!!
func WarehouseItemToResponse(warehouseItem *entity.WarehouseItem) dto.WarehouseItemResponse {
	var description string
	var estimationRunOutStr string

	if warehouseItem.EstimationRunOut.Valid {
		now := time.Now()
		runOutTime := warehouseItem.EstimationRunOut.Time
		daysLeft := int(runOutTime.Sub(now).Hours() / 24)
		if daysLeft < 0 {
			daysLeft = 0
		}
		estimationRunOutStr = fmt.Sprintf("%d hari lagi", daysLeft)

		if now.Add(time.Hour * 24 * 7).After(runOutTime) {
			description = constant.WarehouseItemDescriptionDanger
		} else {
			description = constant.WarehouseItemDescriptionSafe
		}
	} else {
		description = constant.WarehouseItemDescriptionSafe
		estimationRunOutStr = ""
	}

	response := dto.WarehouseItemResponse{
		Warehouse:        WarehouseToResponse(&warehouseItem.Warehouse),
		Item:             ItemToResponse(&warehouseItem.Item),
		Quantity:         warehouseItem.Quantity,
		EstimationRunOut: estimationRunOutStr,
		Description:      description,
	}

	return response
}

func WarehouseOrderItemToResponse(warehouseOrderItem *entity.WarehouseItemProcurement) dto.WarehouseItemProcurementListResponse {
	warehouseItemResponse := dto.WarehouseItemProcurementListResponse{
		Id:        warehouseOrderItem.Id,
		TakenBy:   warehouseOrderItem.TakenBy.UUID.String(),
		IsTaken:   warehouseOrderItem.IsArrived,
		Warehouse: WarehouseToResponse(&warehouseOrderItem.Warehouse),
		Item:      ItemToResponse(&warehouseOrderItem.Item),
		Supplier: dto.SupplierListResponse{
			Id:          warehouseOrderItem.Supplier.Id,
			Name:        warehouseOrderItem.Supplier.Name,
			PhoneNumber: warehouseOrderItem.Supplier.PhoneNumber,
			Address:     warehouseOrderItem.Supplier.Address,
		},
		Quantity: warehouseOrderItem.Quantity,
	}

	if warehouseOrderItem.TakenAt.Valid {
		warehouseItemResponse.TakenAt = warehouseOrderItem.TakenAt.Time.Format("02-Jan-2006")
	}

	return warehouseItemResponse
}

func WarehouseItemHistoryToResponse(warehouseItemHistory *entity.WarehouseItemHistory) dto.WarehouseItemHistoryResponse {
	return dto.WarehouseItemHistoryResponse{
		Id:             warehouseItemHistory.Id,
		Item:           ItemToResponse(&warehouseItemHistory.Item),
		Source:         warehouseItemHistory.Source,
		Destination:    warehouseItemHistory.Destination,
		QuantityBefore: warehouseItemHistory.QuantityBefore,
		QuantityAfter:  warehouseItemHistory.QuantityAfter,
		Status:         warehouseItemHistory.Status.String(),
		UpdatedBy:      warehouseItemHistory.User.Name,
		Date:           warehouseItemHistory.CreatedAt.Format("02-Jan-2006"),
		Time:           warehouseItemHistory.CreatedAt.Format("15:04"),
	}
}

func WarehouseItemHistoryToListResponse(warehouseItemHistory *entity.WarehouseItemHistory) dto.WarehouseItemHistoryListResponse {
	return dto.WarehouseItemHistoryListResponse{
		Id:          warehouseItemHistory.Id,
		Item:        ItemToResponse(&warehouseItemHistory.Item),
		Source:      warehouseItemHistory.Source,
		Destination: warehouseItemHistory.Destination,
		Status:      warehouseItemHistory.Status.String(),
		Quantity:    warehouseItemHistory.QuantityAfter - warehouseItemHistory.QuantityBefore,
		Time:        warehouseItemHistory.CreatedAt.Format("15:04"),
	}
}

// Note : without payments, payment payment
func WarehouseSaleToResponse(warehouseSale *entity.WarehouseSale) dto.WarehouseSaleResponse {
	return dto.WarehouseSaleResponse{
		Id:         warehouseSale.Id,
		SendDate:   warehouseSale.SendDate.Format("02 Jan 2006"),
		Customer:   CustomerToResponse(&warehouseSale.Customer),
		Price:      warehouseSale.Price.String(),
		TotalPrice: warehouseSale.TotalPrice.String(),
		WarehouseItem: dto.ItemResponse{
			Id:       warehouseSale.Item.Id,
			Name:     warehouseSale.Item.Name,
			Unit:     warehouseSale.Item.Unit,
			Category: warehouseSale.Item.Category.String(),
		},
		Warehouse:     WarehouseToResponse(&warehouseSale.Warehouse),
		Quantity:      warehouseSale.Quantity,
		SaleUnit:      warehouseSale.SaleUnit.String(),
		PaymentType:   warehouseSale.PaymentType.String(),
		PaymentStatus: warehouseSale.PaymentStatus.String(),
		IsSend:        warehouseSale.IsSend,
	}
}

// Note : without remaining payment
func WarehouseSalePaymentToResponse(warehouseSalePayment *entity.WarehouseSalePayment) dto.WarehouseSalePaymentResponse {
	return dto.WarehouseSalePaymentResponse{
		Id:            warehouseSalePayment.Id,
		Nominal:       warehouseSalePayment.Nominal.String(),
		PaymentProof:  warehouseSalePayment.PaymentProof,
		PaymentMethod: warehouseSalePayment.PaymentMethod.String(),
		Date:          warehouseSalePayment.PaymentDate.Format("02 Jan 2006"),
	}
}

func WarehouseSaleToListResponse(warehouseSale *entity.WarehouseSale) dto.WarehouseSaleListResponse {
	return dto.WarehouseSaleListResponse{
		Id:            warehouseSale.Id,
		OrderDate:     warehouseSale.CreatedAt.Format("02 Jan 2006"),
		SendDate:      warehouseSale.SendDate.Format("02 Jan 2006"),
		Customer:      CustomerToResponse(&warehouseSale.Customer),
		Item:          ItemToResponse(&warehouseSale.Item),
		Warehouse:     WarehouseToResponse(&warehouseSale.Warehouse),
		Quantity:      warehouseSale.Quantity,
		SaleUnit:      warehouseSale.SaleUnit.String(),
		PaymentStatus: warehouseSale.PaymentStatus.String(),
		IsSend:        warehouseSale.IsSend,
	}
}

func WarehouseSaleQueueToResponse(storeSaleQueue *entity.WarehouseSaleQueue) dto.WarehouseSaleQueueResponse {
	response := dto.WarehouseSaleQueueResponse{
		Id:        storeSaleQueue.Id,
		Item:      ItemToResponse(&storeSaleQueue.Item),
		Warehouse: WarehouseToResponse(&storeSaleQueue.Warehouse),
		Quantity:  storeSaleQueue.Quantity,
		SaleUnit:  storeSaleQueue.SaleUnit.String(),
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

func WarehouseItemProcurementDraftToResponse(data *entity.WarehouseItemProcurementDraft) dto.WarehouseItemProcurementDraftResponse {
	return dto.WarehouseItemProcurementDraftResponse{
		Warehouse:     WarehouseToResponse(&data.Warehouse),
		Item:          ItemToResponse(&data.Item),
		Supplier:      SupplierToListResponse(&data.Supplier),
		DailySpending: data.DailySpending,
		DaysNeed:      data.DaysNeed,
		TotalOrder:    data.DailySpending * float64(data.DaysNeed),
	}
}

func WarehouseItemProcurementToResponse(data *entity.WarehouseItemProcurement) dto.WarehouseItemProcurementResponse {
	return dto.WarehouseItemProcurementResponse{}
}

func WarehouseItemProcurementToListResponse(data *entity.WarehouseItemProcurement) dto.WarehouseItemProcurementListResponse {
	return dto.WarehouseItemProcurementListResponse{}
}

// Note : without remaining payment
func WarehouseItemProcurementPaymentToResponse(storeSalePayment *entity.WarehouseItemProcurementPayment) dto.WarehouseItemProcurementPaymentResponse {
	return dto.WarehouseItemProcurementPaymentResponse{
		Id:            storeSalePayment.Id,
		Nominal:       storeSalePayment.Nominal.String(),
		PaymentProof:  storeSalePayment.PaymentProof,
		PaymentMethod: storeSalePayment.PaymentMethod.String(),
		Date:          storeSalePayment.PaymentDate.Format("02 Jan 2006"),
	}
}

func WarehouseItemCornProcurementDraftToResponse(warehouseItemCornProcurementDraft *entity.WarehouseItemCornProcurementDraft, cornItem dto.ItemResponse) dto.WarehouseItemCornProcurementDraftResponse {
	response := dto.WarehouseItemCornProcurementDraftResponse{
		Warehouse:                 WarehouseToResponse(&warehouseItemCornProcurementDraft.Warehouse),
		Supplier:                  SupplierToListResponse(&warehouseItemCornProcurementDraft.Supplier),
		Item:                      cornItem,
		OvenCondition:             warehouseItemCornProcurementDraft.OvenCondition.String(),
		CornWaterLevel:            warehouseItemCornProcurementDraft.CornWaterLevel.String(),
		IsOvenCanOperateInNearDay: &warehouseItemCornProcurementDraft.IsOvenCanOperateInNearDay,
		Quantity:                  warehouseItemCornProcurementDraft.Quantity,
		TotalPrice:                warehouseItemCornProcurementDraft.Price.Mul(decimal.NewFromFloat(warehouseItemCornProcurementDraft.Quantity)).String(),
	}

	if warehouseItemCornProcurementDraft.Discount.Valid {
		discountPrice := warehouseItemCornProcurementDraft.Price.Mul(decimal.NewFromFloat(warehouseItemCornProcurementDraft.Discount.Float64))
		response.TotalPrice = warehouseItemCornProcurementDraft.Price.Sub(discountPrice).Mul(decimal.NewFromFloat(response.Quantity)).String()
		response.Discount = &warehouseItemCornProcurementDraft.Discount.Float64
	}

	return response
}

// Note : without remaining payment
func WarehouseItemCornProcurementPaymentToResponse(storeSalePayment *entity.WarehouseItemCornProcurementPayment) dto.WarehouseItemCornProcurementPaymentResponse {
	return dto.WarehouseItemCornProcurementPaymentResponse{
		Id:            storeSalePayment.Id,
		Nominal:       storeSalePayment.Nominal.String(),
		PaymentProof:  storeSalePayment.PaymentProof,
		PaymentMethod: storeSalePayment.PaymentMethod.String(),
		Date:          storeSalePayment.PaymentDate.Format("02 Jan 2006"),
	}
}

func WarehouseItemCornProcurementToResponse(data *entity.WarehouseItemCornProcurement) dto.WarehouseItemCornProcurementResponse {
	return dto.WarehouseItemCornProcurementResponse{}
}

func WarehouseItemCornProcurementToListResponse(data *entity.WarehouseItemCornProcurement) dto.WarehouseItemCornProcurementListResponse {
	return dto.WarehouseItemCornProcurementListResponse{}
}

func WarehouseItemCornPriceToResponse(data *entity.WarehouseItemCornPrice) dto.WarehouseItemCornPriceResponse {
	return dto.WarehouseItemCornPriceResponse{
		Id:          data.Id,
		BottomLimit: data.BottomLimit,
		UpperLimit:  data.UpperLimit,
		Discount:    data.Discount,
	}
}
