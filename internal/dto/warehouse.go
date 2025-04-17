package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type CreateWarehouseItemRequest struct {
	Name string `json:"name" validate:"required"`
	Unit string `json:"unit" validate:"required,unit"`
}

type GetWarehouseItemFilter struct {
	Category param.WarehouseItemCategoryParam `query:"category"`
}

type WarehouseItemResponse struct {
	Id       uint64 `json:"id"`
	Name     string `json:"name"`
	Unit     string `json:"unit"`
	Category string `json:"category"`
}

type GetWarehouseStockItemFilter struct {
	WarehouseId uint64 `query:"warehouseId"`
}

type WarehouseResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Location LocationResponse `json:"location"`
}

type CreateWarehouseStockItemRequest struct {
	WarehouseId     uint64 `json:"warehouseId" validate:"required"`
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required"`
	Quantity        uint64 `json:"quantity" validate:"required"`
}

type UpdateWarehouseStockItemRequest struct {
	Quantity uint64 `json:"quantity" validate:"required"`
}

type WarehouseStockItemResponse struct {
	Warehouse        WarehouseResponse     `json:"warehouse"`
	WarehouseItem    WarehouseItemResponse `json:"warehouseItem"`
	Quantity         uint64                `json:"quantity"`
	EstimationRunOut string                `json:"estimationRunOut"`
	Description      string                `json:"description"`
}
