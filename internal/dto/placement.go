package dto

type CreateCagePlacementRequest struct {
	UserId  string   `json:"userId" validate:"required,number"`
	CageIds []uint64 `json:"cageIds" validate:"required,number"`
}

type CreateStorePlacementRequest struct {
	UserId   string   `json:"userId" validate:"required,number"`
	StoreIds []uint64 `json:"storeIds" validate:"required,number"`
}

type CreateWarehousePlacementRequest struct {
	UserId       string   `json:"userId" validate:"required,number"`
	WarehouseIds []uint64 `json:"warehouseIds" validate:"required,number"`
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
