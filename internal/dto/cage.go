package dto

type CageResponse struct {
	Id              uint64           `json:"id"`
	Name            string           `json:"name"`
	Capacity        uint64           `json:"capacity"`
	ChickenCategory string           `json:"chickenCategory"`
	Location        LocationResponse `json:"location"`
}

type CreateCageRequest struct {
	Name            string `json:"name" validate:"required"`
	Capacity        uint64 `json:"capacity" validate:"required"`
	LocationId      uint64 `json:"locationId" validate:"required"`
	ChickenCategory string `json:"chickenCategory" validate:"required"`
}

type GetCageFilter struct {
	LocationId uint64 `query:"locationId"`
}
