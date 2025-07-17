package dto

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type CreateChickenMonitoringRequest struct {
	ChickenCageId     uint64  `json:"chickenCageId" validate:"required"`
	TotalSickChicken  uint64  `json:"totalSickChicken" validate:"number,min=0"`
	TotalDeathChicken uint64  `json:"totalDeathChicken" validate:"number,min=0"`
	TotalFeed         float64 `json:"totalFeed" validate:"number,min=0"`
	Note              string  `json:"note"`
}

type UpdateChickenMonitoringRequest struct {
	ChickenCageId     uint64  `json:"chickenCageId" validate:"required"`
	TotalSickChicken  uint64  `json:"totalSickChicken" validate:"number,min=0"`
	TotalDeathChicken uint64  `json:"totalDeathChicken" validate:"number,min=0"`
	TotalFeed         float64 `json:"totalFeed" validate:"number,min=0"`
	Note              string  `json:"note"`
}

type ChickenMonitoringResponse struct {
	Id                 uint64              `json:"id"`
	ChickenCage        ChickenCageResponse `json:"chickenCage"`
	TotalLiveChicken   uint64              `json:"totalLiveChicken"`
	TotalSickChicken   uint64              `json:"totalSickChicken"`
	TotalDeatchChicken uint64              `json:"totalDeathChicken"`
	TotalFeed          float64             `json:"totalFeed"`
	Note               string              `json:"note"`
}

type ChickenMonitoringListResponse struct {
	Id                uint64              `json:"id"`
	ChickenCage       ChickenCageResponse `json:"chickenCage"`
	TotalLiveChicken  uint64              `json:"totalLiveChicken"`
	TotalSickChicken  uint64              `json:"totalSickChicken"`
	TotalDeathChicken uint64              `json:"totalDeathChicken"`
	TotalFeed         float64             `json:"totalFeed"`
	MortalityRate     float64             `json:"mortalityRate"`
}

type CreateChickenHealthItemRequest struct {
	Name       string  `json:"name" validate:"required"`
	Type       string  `json:"type" validate:"required,chickenHealthItemType"`
	ChickenAge *uint64 `json:"chickenAge"`
	Note       string  `json:"note"`
}

type UpdateChickenHealthItemRequest struct {
	Name       string  `json:"name" validate:"required"`
	Type       string  `json:"type" validate:"required,chickenHealthItemType"`
	ChickenAge *uint64 `json:"chickenAge"`
	Note       string  `json:"note"`
}

type ChickenHealthItemResponse struct {
	Id              uint64  `json:"id"`
	Name            string  `json:"name"`
	Type            string  `json:"type"`
	ChickenAge      *uint64 `json:"chickenAge"`
	ChickenCategory *string `json:"chickenCategory"`
	Note            string  `json:"note"`
}

type GetChickenHealthItemFilter struct {
	Type param.ChickenHealthItemTypeParam `query:"type"`
}

type CreateChickenHealthMonitoringRequest struct {
	ChickenCageId  uint64  `json:"chickenCageId" validate:"required"`
	HealthItemName string  `json:"healthItemName" validate:"required"`
	Type           string  `json:"type" validate:"required,chickenHealthItemType"`
	Dose           float64 `json:"dose" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	Disease        *string `json:"disease"`
}

type UpdateChickenHealthMonitoringRequest struct {
	ChickenCageId  uint64  `json:"chickenCageId" validate:"required"`
	HealthItemName string  `json:"healthItemName" validate:"required"`
	Type           string  `json:"type" validate:"required,chickenHealthItemType"`
	Dose           float64 `json:"dose" validate:"required"`
	Unit           string  `json:"unit" validate:"required"`
	Disease        *string `json:"disease"`
}

type ChickenHealthMonitoringResponse struct {
	Id              uint64  `json:"id"`
	HealthItemName  string  `json:"healthItemName"`
	Type            string  `json:"type"`
	Dose            float64 `json:"dose"`
	Unit            string  `json:"unit"`
	Disease         string  `json:"disease"`
	Date            string  `json:"date"`
	ChickenAge      uint64  `json:"chickenAge"`
	ChickenCategory string  `json:"chickenCategory"`
	CreatedAt       string  `json:"createdAt"`
}

type ChickenHealthMonitoringDetailResponse struct {
	ChickenCage              ChickenCageResponse               `json:"chickenCage"`
	ChickenHealthMonitorings []ChickenHealthMonitoringResponse `json:"chickenHealthMonitorings"`
}

type GetChickenMonitoringFilter struct {
	Date       param.DateParam `query:"date"`
	LocationId uint64          `query:"locationId"`
	CageId     uint64
	StartDate  param.DateParam
	EndDate    param.DateParam
}

type GetChickenOverviewFilter struct {
	LocationId        uint64                       `query:"locationId"`
	CageId            uint64                       `query:"cageId"`
	OverviewGraphTime param.OverviewGraphTimeParam `query:"overviewGraphTime"`
}

type ChickenDetailOverview struct {
	TotalLiveChicken    uint64  `json:"totalLiveChicken"`
	TotalSickChicken    uint64  `json:"totalSickChicken"`
	TotalDeathChicken   uint64  `json:"totalDeathChicken"`
	TotalKPIPerformance float64 `json:"totalKPIPerformance"`
}

type ChickenGraphResponse struct {
	Key          string `json:"key"`
	SickChicken  uint64 `json:"sickChicken"`
	DeathChicken uint64 `json:"deathChicken"`
}

type ChickenBarChartResponse struct {
	ChickenDOC       float64 `json:"chickenDOC"`
	ChickenGrower    float64 `json:"chickenGrower"`
	ChickentPreLayer float64 `json:"chickentPreLayer"`
	ChickenLayer     float64 `json:"chickenLayer"`
	ChickenAfkir     float64 `json:"chickenAfkir"`
}

type ChickenOverviewResponse struct {
	ChickenDetail ChickenDetailOverview   `json:"chickenDetail"`
	ChickenGraphs []ChickenGraphResponse  `json:"chickenGraphs"`
	ChickenPie    ChickenBarChartResponse `json:"chickenPie"`
}

type ChickenProcurementRequest struct {
}
