package dto

type CageResponse struct {
	Id              uint64           `json:"id"`
	Name            string           `json:"name"`
	Capacity        uint64           `json:"capacity"`
	ChickenCategory string           `json:"chickenCategory"`
	IsUsed          bool             `json:"isUsed"`
	Location        LocationResponse `json:"location"`
}

type CreateCageRequest struct {
	Name            string `json:"name" validate:"required"`
	Capacity        uint64 `json:"capacity" validate:"required"`
	LocationId      uint64 `json:"locationId" validate:"required"`
	ChickenCategory string `json:"chickenCategory" validate:"required"`
}

type UpdateCageRequest struct {
	Name            string `json:"name" validate:"required"`
	Capacity        uint64 `json:"capacity" validate:"required"`
	LocationId      uint64 `json:"locationId" validate:"required"`
	ChickenCategory string `json:"chickenCategory" validate:"required"`
	IsUsed          bool   `json:"isUsed" validate:"required"`
}

type GetCageFilter struct {
	LocationId uint64 `query:"locationId"`
}

type UpdateChickenCageRequest struct {
	TotalDeatchChicken bool `json:"totalDeathChicken" validate:"required"`
}

type CreateChickenCageRequest struct {
	CageId               uint64 `json:"cageId" validate:"required"`
	ChickenProcurementId uint64 `json:"chickenProcurementId" validate:"required"`
	TotalChicken         uint64 `json:"totalChicken" validate:"required"`
}

type ChickenCageResponse struct {
	Cage                 CageResponse `json:"cage"`
	Id                   uint64       `json:"id"`
	BatchId              string       `json:"batchId"`
	ChickenCategory      string       `json:"chickenCategory"`
	ChickenAge           uint64       `json:"chickenAge"`
	TotalChicken         uint64       `json:"totalChicken"`
	ChickenPic           string       `json:"chickenPic"`
	EggPic               string       `json:"eggPic"`
	IsNeedRoutineVaccine bool         `json:"isNeedRoutineVaccine"`
}

type GetChickenCageFilter struct {
	LocationId uint64 `query:"locationid"`
}
