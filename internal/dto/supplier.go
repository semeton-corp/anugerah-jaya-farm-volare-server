package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

type CreateSupplierRequest struct {
	ItemIds      []uint64 `json:"itemIds"`
	Name         string   `json:"name" validate:"required"`
	PhoneNumber  string   `json:"phoneNumber" validate:"required,phoneNumber"`
	Address      string   `json:"address" validate:"required"`
	SupplierType string   `json:"supplierType" validate:"required,supplierType"`
}

type UpdateSupplierRequest struct {
	ItemIds      []uint64 `json:"itemIds" validate:"required"`
	Name         string   `json:"name" validate:"required"`
	PhoneNumber  string   `json:"phoneNumber" validate:"required"`
	Address      string   `json:"address" validate:"required"`
	SupplierType string   `json:"supplierType" validate:"required,supplierType"`
}

type SupplierResponse struct {
	Id           uint64         `json:"id"`
	Items        []ItemResponse `json:"items,omitempty"`
	Name         string         `json:"name"`
	PhoneNumber  string         `json:"phoneNumber"`
	Address      string         `json:"address"`
	SupplierType string         `json:"supplierType"`
}

type SupplierListResponse struct {
	Id           uint64 `json:"id"`
	Name         string `json:"name"`
	PhoneNumber  string `json:"phoneNumber"`
	Address      string `json:"address"`
	SupplierType string `json:"supplierType"`
}

type GetSupplierFilter struct {
	SupplierType param.SupplierTypeParam `query:"supplierType"`
	ItemId       uint64                  `query:"itemId"`
}
