package dto

type CreateSupplierRequest struct {
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required"`
	Name            string `json:"name" validate:"required"`
	PhoneNumber     string `json:"phoneNumber" validate:"required"`
	Address         string `json:"address" validate:"required"`
}

type UpdateSupplierRequest struct {
	WarehouseItemId uint64 `json:"warehouseItemId" validate:"required"`
	Name            string `json:"name" validate:"required"`
	PhoneNumber     string `json:"phoneNumber" validate:"required"`
	Address         string `json:"address" validate:"required"`
}

type SupplierResponse struct {
	Id            uint64       `json:"id"`
	WarehouseItem ItemResponse `json:"warehouseItem,omitempty"`
	Name          string       `json:"name"`
	PhoneNumber   string       `json:"phoneNumber"`
	Address       string       `json:"address"`
}

type SupplierWithoutWarehouseItemResponse struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}

type SupplierListResponse struct {
	Id            uint64       `json:"id"`
	WarehouseItem ItemResponse `json:"warehouseItem"`
	Name          string       `json:"name"`
	PhoneNumber   string       `json:"phoneNumber"`
	Address       string       `json:"address"`
}
