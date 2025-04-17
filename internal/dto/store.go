package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type StoreResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Location LocationResponse `json:"location"`
}

type CreateStoreRequestItemRequest struct {
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required,number"`
	WarehouseId     uint64 `json:"warehouseId" validate:"required,number"`
	StoreId         uint64 `json:"storeId" validate:"required,number"`
	Quantity        uint64 `json:"quantity" validate:"required,number"`
}

type UpdateStoreRequestItemByWarehouseRequest struct {
	Status string `json:"status" validate:"required,requestItemStatus,oneof=Dikirim Ditolak"`
}

type UpdateStoreRequestItemByStoreRequest struct {
	Status   string `json:"status" validate:"required,requestItemStatus,oneof=Diterima"`
	Quantity uint64 `json:"quantity" validate:"required,number"`
}

type UpdateStoreRequestItemRequest struct {
	Status   string `json:"status" validate:"required,requestItemStatus"`
	Quantity uint64 `json:"quantity" validate:"required,number"`
}

type StoreRequestItemResponse struct {
	Id            uint64                `json:"id"`
	Warehouse     WarehouseResponse     `json:"warehouse"`
	WarehouseItem WarehouseItemResponse `json:"warehouseItem"`
	Store         StoreResponse         `json:"store"`
	Quantity      uint64                `json:"quantity"`
	Status        string                `json:"status"`
}

type GetStoreRequestItemFilter struct {
	Date param.DateParam `query:"date"`
	Page uint64          `query:"page"`
}
