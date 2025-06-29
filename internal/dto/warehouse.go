package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type CreateWarehouseRequest struct {
	Name       string `json:"name" validate:"required"`
	LocationId uint64 `json:"locationId" validate:"required"`
}

type UpdateWarehouseRequest struct {
	Name       string `json:"name" validate:"required"`
	LocationId uint64 `json:"locationId" validate:"required"`
}

type GetWarehouseFilter struct {
	LocationId uint64 `query:"locationId"`
}


type GetWarehouseStockItemFilter struct {
	WarehouseId uint64                           `query:"warehouseId"`
	Category    param.WarehouseItemCategoryParam `query:"category"`
}

type WarehouseResponse struct {
	Id            uint64           `json:"id"`
	Name          string           `json:"name"`
	Location      LocationResponse `json:"location"`
	TotalEmployee uint64           `json:"totalEmployee"`
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
	Warehouse        WarehouseResponse `json:"warehouse"`
	WarehouseItem    ItemResponse      `json:"warehouseItem"`
	Quantity         uint64            `json:"quantity"`
	EstimationRunOut string            `json:"estimationRunOut"`
	Description      string            `json:"description"`
}

type CreateWarehouseOrderItemRequest struct {
	WarehouseId     uint64 `json:"warehouseId" validate:"required"`
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required"`
	SupplierId      uint64 `json:"supplierId" validate:"required"`
	Quantity        uint64 `json:"quantity" validate:"required"`
}

type WarehouseOrderItemResponse struct {
	Id            uint64                               `json:"id"`
	Warehouse     WarehouseResponse                    `json:"warehouse"`
	WarehouseItem ItemResponse                         `json:"warehouseItem"`
	Supplier      SupplierWithoutWarehouseItemResponse `json:"supplier"`
	TakenBy       string                               `json:"takenBy"`
	TakenAt       string                               `json:"takenAt"`
	IsTaken       bool                                 `json:"isTaken"`
	Quantity      uint64                               `json:"quantity"`
}

type GoodEggWarehouseConvertionRequest struct {
	WarehouseId uint64 `json:"warehouseId" validate:"required,number"`
	TotalKarpet uint64 `json:"totalKarpet" validate:"required,number,min=0"`
	TotalButir  uint64 `json:"totalButir" validate:"required,number,min=0"`
	TotalIkat   uint64 `json:"totalIkat" validate:"required,number,min=0"`
}

type CrackedEggWarehouseConvertionRequest struct {
	WarehouseId uint64 `json:"warehouseId" validate:"required,number"`
	TotalButir  uint64 `json:"totalButir" validate:"required,number,min=0"`
	TotalPack   uint64 `json:"totalPack" validate:"required,number,min=0"`
}

type GetWarehouseOrderItemFilter struct {
	Date    param.DateParam `query:"date"`
	IsTaken bool
}
