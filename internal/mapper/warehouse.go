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
	}
}

func WarehouseItemToResponse(warehouseItem *entity.WarehouseItem) dto.WarehouseItemResponse {
	return dto.WarehouseItemResponse{
		Id:       warehouseItem.Id,
		Name:     warehouseItem.Name,
		Category: warehouseItem.Category.String(),
		Unit:     warehouseItem.Unit,
	}
}

// Note : without description
func WarehouseStockItemToResponse(warehouseStockItem *entity.WarehouseStockItem) dto.WarehouseStockItemResponse {
	return dto.WarehouseStockItemResponse{
		Warehouse:        WarehouseToResponse(&warehouseStockItem.Warehouse),
		WarehouseItem:    WarehouseItemToResponse(&warehouseStockItem.WarehouseItem),
		Quantity:         warehouseStockItem.Quantity,
		EstimationRunOut: warehouseStockItem.EstimationRunOut.Format("02-Jan-2006"),
	}
}

func WarehouseOrderItemToResponse(warehouseOrderItem *entity.WarehouseOrderItem) dto.WarehouseOrderItemResponse {
	return dto.WarehouseOrderItemResponse{
		Warehouse:     WarehouseToResponse(&warehouseOrderItem.Warehouse),
		WarehouseItem: WarehouseItemToResponse(&warehouseOrderItem.WarehouseItem),
		Supplier:      SupplierToResponse(&warehouseOrderItem.Supplier),
		Quantity:      warehouseOrderItem.Quantity,
	}
}
