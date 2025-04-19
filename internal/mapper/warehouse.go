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
		Warehouse: dto.WarehouseResponse{
			Id:   warehouseStockItem.Warehouse.Id,
			Name: warehouseStockItem.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   warehouseStockItem.Warehouse.Location.Id,
				Name: warehouseStockItem.Warehouse.Location.Name,
			},
		},
		WarehouseItem: dto.WarehouseItemResponse{
			Id:       warehouseStockItem.WarehouseItem.Id,
			Name:     warehouseStockItem.WarehouseItem.Name,
			Unit:     warehouseStockItem.WarehouseItem.Unit,
			Category: warehouseStockItem.WarehouseItem.Category.String(),
		},
		Quantity:         warehouseStockItem.Quantity,
		EstimationRunOut: warehouseStockItem.EstimationRunOut.Format("02-Jan-2006"),
	}
}
