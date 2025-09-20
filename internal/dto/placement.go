package dto

type CreateCagePlacementRequest struct {
	UserId string `json:"userId" validate:"required"`
	CageId uint64 `json:"cageId" validate:"required,number"`
}

type UpdateCagePlacementRequest struct {
	CageId uint64                           `json:"cageId" validate:"required,number"`
	Users  []UpdateCagePlacementUserRequest `json:"users" validate:"dive"`
}

type UpdateCagePlacementUserRequest struct {
	UserId string `json:"userId"`
	RoleId uint64 `json:"roleId" validate:"number"`
}

type CreateStorePlacementRequest struct {
	UserId  string `json:"userId" validate:"required"`
	StoreId uint64 `json:"storeId" validate:"required,number"`
}

type CreateWarehousePlacementRequest struct {
	UserId      string `json:"userId" validate:"required"`
	WarehouseId uint64 `json:"warehouseId" validate:"required,number"`
}

type CagePlacementResponse struct {
	User UserListResponse `json:"user"`
	Cage CageResponse     `json:"cage"`
}

type StorePlacementResponse struct {
	User  UserListResponse `json:"user"`
	Store StoreResponse    `json:"store"`
}

type WarehousePlacementResponse struct {
	User      UserListResponse  `json:"user"`
	Warehouse WarehouseResponse `json:"warehouse"`
}
