package mapper

import (
	"fmt"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
)

func WarehouseToResponse(warehouse *entity.Warehouse) dto.WarehouseResponse {
	return dto.WarehouseResponse{
		Id:   warehouse.Id,
		Name: warehouse.Name,
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
			description = constant.StockWarehouseItemDanger
		} else {
			description = constant.StockWarehouseItemSafe
		}
	} else {
		description = constant.StockWarehouseItemSafe
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

func WarehouseOrderItemToResponse(warehouseOrderItem *entity.WarehouseOrderItem) dto.WarehouseOrderItemResponse {
	warehouseItemResponse := dto.WarehouseOrderItemResponse{
		Id:            warehouseOrderItem.Id,
		TakenBy:       warehouseOrderItem.TakenBy.UUID.String(),
		IsTaken:       warehouseOrderItem.IsTaken.Bool,
		Warehouse:     WarehouseToResponse(&warehouseOrderItem.Warehouse),
		WarehouseItem: ItemToResponse(&warehouseOrderItem.Item),
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
