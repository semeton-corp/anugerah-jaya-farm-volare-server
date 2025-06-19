package dto

type CageResponse struct {
	Id       uint64           `json:"id"`
	Name     string           `json:"name"`
	Capacity uint64           `json:"capacity"`
	Location LocationResponse `json:"location"`
}

type CreateCageRequest struct {
	Name       string `json:"name"`
	Capacity   uint64 `json:"capacity"`
	LocationId uint64 `json:"locationId"`
}

type GetCageFilter struct {
	LocationId uint64 `query:"locationId"`
}
