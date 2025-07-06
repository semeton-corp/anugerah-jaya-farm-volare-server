package dto

type CreateSupplierRequest struct {
	ItemIds     []uint64 `json:"itemIds" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	PhoneNumber string   `json:"phoneNumber" validate:"required"`
	Address     string   `json:"address" validate:"required"`
}

type UpdateSupplierRequest struct {
	ItemIds     []uint64 `json:"itemIds" validate:"required"`
	Name        string   `json:"name" validate:"required"`
	PhoneNumber string   `json:"phoneNumber" validate:"required"`
	Address     string   `json:"address" validate:"required"`
}

type SupplierResponse struct {
	Id          uint64         `json:"id"`
	Items       []ItemResponse `json:"items,omitempty"`
	Name        string         `json:"name"`
	PhoneNumber string         `json:"phoneNumber"`
	Address     string         `json:"address"`
}

type SupplierListResponse struct {
	Id          uint64 `json:"id"`
	Name        string `json:"name"`
	PhoneNumber string `json:"phoneNumber"`
	Address     string `json:"address"`
}
