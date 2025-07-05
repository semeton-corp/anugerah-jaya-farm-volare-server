package dto

type CreateCagePlacementRequest struct {
	UserId  string   `json:"userId" validate:"required,number"`
	CageIds []uint64 `json:"cageIds" validate:"required,number"`
}

type CreateStorePlacementRequest struct {
	UserId  string `json:"userId" validate:"required,number"`
	StoreId uint64 `json:"storeId" validate:"required,number"`
}

type CreateWarehousePlacementRequest struct {
	UserId      string `json:"userId" validate:"required,number"`
	WarehouseId uint64 `json:"warehouseId" validate:"required,number"`
}

type CagePlacementResponse struct {
	User UserResponse `json:"user"`
	Cage CageResponse `json:"cage"`
}

type StorePlacementResponse struct {
	User  UserResponse  `json:"user"`
	Store StoreResponse `json:"store"`
}

type WarehousePlacementResponse struct {
	User      UserResponse      `json:"user"`
	Warehouse WarehouseResponse `json:"warehouse"`
}
