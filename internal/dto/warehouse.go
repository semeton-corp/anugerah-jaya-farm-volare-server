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

type GetWarehouseItemFilter struct {
	WarehouseId uint64                  `query:"warehouseId"`
	Category    param.ItemCategoryParam `query:"category"`
	ItemNames   []string                `query:"itemNames"`
	Units       []string                `query:"units"`
}

type WarehouseResponse struct {
	Id            uint64           `json:"id"`
	Name          string           `json:"name"`
	Location      LocationResponse `json:"location"`
	TotalEmployee uint64           `json:"totalEmployee"`
}

type WarehouseDetailResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Location LocationResponse `json:"location"`
	Users    []UserResponse   `json:"users"`
}

type CreateWarehouseItemRequest struct {
	WarehouseId     uint64  `json:"warehouseId" validate:"required"`
	ItemId          uint64  `json:"itemId" validate:"required"`
	Quantity        float64 `json:"quantity" validate:"required"`
	RunOutCountDown *uint64 `json:"runOutCountDown"`
}

type UpdateWarehouseItemRequest struct {
	Quantity        float64 `json:"quantity" validate:"required"`
	RunOutCountDown *uint64 `json:"runOutCountDown"`
}

type WarehouseItemResponse struct {
	Warehouse        WarehouseResponse `json:"warehouse"`
	Item             ItemResponse      `json:"item"`
	Quantity         float64           `json:"quantity"`
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
	Id            uint64               `json:"id"`
	Warehouse     WarehouseResponse    `json:"warehouse"`
	WarehouseItem ItemResponse         `json:"warehouseItem"`
	Supplier      SupplierListResponse `json:"supplier"`
	TakenBy       string               `json:"takenBy"`
	TakenAt       string               `json:"takenAt"`
	IsTaken       bool                 `json:"isTaken"`
	Quantity      uint64               `json:"quantity"`
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

type WarehouseOverview struct {
	TotalSafeStock    uint64                  `json:"totalSafeStock"`
	TotalDangerStock  uint64                  `json:"totalDangerStock"`
	TotalStoreRequest uint64                  `json:"totalStoreRequest"`
	EggStocks         []WarehouseItemResponse `json:"eggStocks"`
	EquipmentStocks   []WarehouseItemResponse `json:"equipmentStocks"`
}

type WarehouseItemHistoryListResponse struct {
	Id          uint64       `json:"id"`
	Item        ItemResponse `json:"item"`
	Source      string       `json:"source"`
	Destination string       `json:"destination"`
	Quantity    float64      `json:"quantity"`
	Status      string       `json:"status"`
	Time        string       `json:"time"`
}

type WarehouseItemHistoryListPaginationResponse struct {
	TotalPage              uint64                             `json:"totalPage,omitempty"`
	TotalData              uint64                             `json:"totalData,omitempty"`
	WarehouseItemHistories []WarehouseItemHistoryListResponse `json:"warehouseItemHistories"`
}

type WarehouseItemHistoryResponse struct {
	Id             uint64       `json:"id"`
	Item           ItemResponse `json:"item"`
	Source         string       `json:"source"`
	Destination    string       `json:"destination"`
	QuantityBefore float64      `json:"quantityBefore"`
	QuantityAfter  float64      `json:"quantityAfter"`
	Status         string       `json:"status"`
	UpdatedBy      string       `json:"updatedBy"`
	Time           string       `json:"time"`
	Date           string       `json:"date"`
}

type GetWarehouseItemHistoryFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}

type EggWarehouseItemSummary struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type GetEggWarehouseItemSummary struct {
	WarehouseId uint64 `query:"warehouseId"`
}
