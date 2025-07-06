package mapper

import (
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

func WarehouseItemToResponse(warehouseItem *entity.WarehouseItem) dto.WarehouseItemResponse {
	var description string
	if time.Now().Add(time.Hour * 24 * 7).After(warehouseItem.EstimationRunOut) {
		description = constant.StockWarehouseItemDanger
	} else {
		description = constant.StockWarehouseItemSafe
	}

	response := dto.WarehouseItemResponse{
		Warehouse:        WarehouseToResponse(&warehouseItem.Warehouse),
		Item:             ItemToResponse(&warehouseItem.Item),
		Quantity:         warehouseItem.Quantity,
		EstimationRunOut: warehouseItem.EstimationRunOut.Format("02-Jan-2006"),
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
