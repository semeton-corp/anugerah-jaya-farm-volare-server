package dto

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"

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
	IsUsed          *bool  `json:"isUsed" validate:"required"`
}

type GetCageFilter struct {
	LocationId      uint64                     `query:"locationId"`
	ChickenCategory param.ChickenCategoryParam `query:"chickenCategory"`
	IsUsed          *bool                      `query:"isUsed"`
}

type UpdateChickenCageRequest struct {
	TotalChicken uint64 `json:"totalChicken" validate:"required"`
}

type CreateChickenCageRequest struct {
	CageId               uint64 `json:"cageId" validate:"required"`
	ChickenProcurementId uint64 `json:"chickenProcurementId" validate:"required"`
	TotalChicken         uint64 `json:"totalChicken" validate:"required"`
}

type ChickenCageResponse struct {
	Cage                 CageResponse `json:"cage"`
	Id                   uint64       `json:"id"`
	ChickenProcurementId uint64       `json:"chickenProcurementId"`
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
	CageId     uint64 `query:"cageId"`
}

type MoveChickenCageRequest struct {
	SourceCageId            uint64                          `json:"sourceCageId" validate:"required"`
	DestinationChickenCages []DestinationChickenCageRequest `json:"destinationChickenCages" validate:"required,dive"`
}

type DestinationChickenCageRequest struct {
	DestinationCageId uint64 `json:"destinationCageId"`
	TotalChicken      uint64 `json:"totalChicken"`
}

type CreateCageFeedRequest struct {
	ChickenCategory string                  `json:"chickenCategory" validate:"required,chickenCategory"`
	FeedType        string                  `json:"feedType" validate:"required,feedType"`
	TotalFeed       float64                 `json:"totalFeed" validate:"required"`
	CageFeedDetails []CageFeedDetailRequest `json:"cageFeedDetails"`
}

type UpdateCageFeedRequest struct {
	ChickenCategory string                  `json:"chickenCategory" validate:"required,chickenCategory"`
	TotalFeed       float64                 `json:"totalFeed" validate:"required"`
	FeedType        string                  `json:"feedType" validate:"required,feedType"`
	CageFeedDetails []CageFeedDetailRequest `json:"cageFeedDetails"`
}

type CageFeedDetailRequest struct {
	Id         uint64  `json:"id"`
	ItemId     uint64  `json:"itemId" validate:"required"`
	Percentage float64 `json:"percentage" validate:"required"`
}

type CageFeedResponse struct {
	Id                 uint64                   `json:"id"`
	ChickenCategory    string                   `json:"chickenCategory"`
	ChickenAgeInterval string                   `json:"chickenAgeInterval"`
	TotalFeed          float64                  `json:"totalFeed"`
	FeedType           string                   `json:"feedType"`
	CageFeedDetails    []CageFeedDetailResponse `json:"cageFeedDetails"`
}
type CageFeedDetailResponse struct {
	Id         uint64       `json:"id"`
	Item       ItemResponse `json:"item"`
	Percentage float64      `json:"percentage"`
}

type ChickenCageFeedListResponse struct {
	Cage              CageResponse `json:"cage"`
	Id                uint64       `json:"id"`
	ChickenCategory   string       `json:"chickenCategory"`
	ChickenAge        uint64       `json:"chickenAge"`
	TotalChicken      uint64       `json:"totalChicken"`
	ExpectedTotalFeed float64      `json:"expectedTotalFeed"`
	TotalFeed         float64      `json:"totalFeed"`
	IsNeedFeed        bool         `json:"isNeedFeed"`
}

type ChickenCageFeedResponse struct {
	Cage              CageResponse         `json:"cage"`
	Id                uint64               `json:"id"`
	ChickenCategory   string               `json:"chickenCategory"`
	ChickenAge        uint64               `json:"chickenAge"`
	TotalChicken      uint64               `json:"totalChicken"`
	FeedType          string               `json:"feedType"`
	ExpectedTotalFeed float64              `json:"expectedTotalFeed"`
	TotalFeed         float64              `json:"totalFeed"`
	IsNeedFeed        bool                 `json:"isNeedFeed"`
	FeedDetails       []FeedDetailResponse `json:"feedDetails"`
}

type FeedDetailResponse struct {
	Item       ItemResponse `json:"item"`
	Percentage float64      `json:"percentage"`
	Quantity   float64      `json:"quantity"`
}

type GetChickenCageFeedFilter struct {
	LocationId uint64 `query:"locationid"`
}

type ConfirmationChickenCageFeedRequest struct {
	WarehouseId uint64 `json:"warehouseId"`
}

type GetCageFeedStockFilter struct {
	CageId uint64 `query:"cageId"`
}
