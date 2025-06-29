package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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

// Note : without description
func WarehouseStockItemToResponse(warehouseStockItem *entity.WarehouseItem) dto.WarehouseStockItemResponse {
	return dto.WarehouseStockItemResponse{
		Warehouse:        WarehouseToResponse(&warehouseStockItem.Warehouse),
		WarehouseItem:    ItemToResponse(&warehouseStockItem.Item),
		Quantity:         warehouseStockItem.Quantity,
		EstimationRunOut: warehouseStockItem.EstimationRunOut.Format("02-Jan-2006"),
	}
}

func WarehouseOrderItemToResponse(warehouseOrderItem *entity.WarehouseOrderItem) dto.WarehouseOrderItemResponse {
	warehouseItemResponse := dto.WarehouseOrderItemResponse{
		Id:            warehouseOrderItem.Id,
		TakenBy:       warehouseOrderItem.TakenBy.UUID.String(),
		IsTaken:       warehouseOrderItem.IsTaken.Bool,
		Warehouse:     WarehouseToResponse(&warehouseOrderItem.Warehouse),
		WarehouseItem: ItemToResponse(&warehouseOrderItem.Item),
		Supplier: dto.SupplierWithoutWarehouseItemResponse{
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
