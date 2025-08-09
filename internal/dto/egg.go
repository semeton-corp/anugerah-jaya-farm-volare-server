package dto

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type CreateEggMonitoringRequest struct {
	ChickenCageId            uint64  `json:"chickenCageId" validate:"required,number"`
	WarehouseId              uint64  `json:"warehouseId" validate:"required,number"`
	TotalKarpetGoodEgg       uint64  `json:"totalKarpetGoodEgg" validate:"number,min=0"`
	TotalRemainingGoodEgg    uint64  `json:"totalRemainingGoodEgg" validate:"number,min=0"`
	TotalWeightGoodEgg       float64 `json:"totalWeightGoodEgg" validate:"number,min=0"`
	TotalKarpetCrackedEgg    uint64  `json:"totalKarpetCrackedEgg" validate:"number,min=0"`
	TotalRemainingCrackedEgg uint64  `json:"totalRemainingCrackedEgg" validate:"number,min=0"`
	TotalWeightCrackedEgg    float64 `json:"totalWeightCrackedEgg" validate:"number,min=0"`
	TotalKarpetRejectEgg     uint64  `json:"totalKarpetRejectEgg" validate:"number,min=0"`
	TotalRemainingRejectEgg  uint64  `json:"totalRemainingRejectEgg" validate:"number,min=0"`
}

type UpdateEggMonitoringRequest struct {
	ChickenCageId            uint64  `json:"chickenCageId" validate:"required,number"`
	WarehouseId              uint64  `json:"warehouseId" validate:"required,number"`
	TotalKarpetGoodEgg       uint64  `json:"totalKarpetGoodEgg" validate:"number,min=0"`
	TotalRemainingGoodEgg    uint64  `json:"totalRemainingGoodEgg" validate:"number,min=0"`
	TotalWeightGoodEgg       float64 `json:"totalWeightGoodEgg" validate:"number,min=0"`
	TotalKarpetCrackedEgg    uint64  `json:"totalKarpetCrackedEgg" validate:"number,min=0"`
	TotalRemainingCrackedEgg uint64  `json:"totalRemainingCrackedEgg" validate:"number,min=0"`
	TotalWeightCrackedEgg    float64 `json:"totalWeightCrackedEgg" validate:"number,min=0"`
	TotalKarpetRejectEgg     uint64  `json:"totalKarpetRejectEgg" validate:"number,min=0"`
	TotalRemainingRejectEgg  uint64  `json:"totalRemainingRejectEgg" validate:"number,min=0"`
}

type EggMonitoringResponse struct {
	Id                       uint64              `json:"id"`
	ChickenCage              ChickenCageResponse `json:"chickenCage"`
	Warehouse                WarehouseResponse   `json:"warehouse"`
	TotalKarpetGoodEgg       uint64              `json:"totalKarpetGoodEgg"`
	TotalRemainingGoodEgg    uint64              `json:"totalRemainingGoodEgg"`
	TotalWeightGoodEgg       float64             `json:"totalWeightGoodEgg"`
	TotalKarpetCrackedEgg    uint64              `json:"totalKarpetCrackedEgg"`
	TotalRemainingCrackedEgg uint64              `json:"totalRemainingCrackedEgg"`
	TotalWeightCrackedEgg    float64             `json:"totalWeightCrackedEgg"`
	TotalKarpetRejectEgg     uint64              `json:"totalKarpetRejectEgg"`
	TotalRemainingRejectEgg  uint64              `json:"totalRemainingRejectEgg"`
	TotalWeightAllEgg        float64             `json:"totalWeightAllEgg"`
	AverageWeight            float64             `json:"averageWeight"`
	CreatedAt                string              `json:"createdAt"`
}

type EggMonitoringListResponse struct {
	Id                 uint64              `json:"id"`
	ChickenCage        ChickenCageResponse `json:"chickenCage"`
	Warehouse          WarehouseResponse   `json:"warehouse"`
	TotalAllEgg        uint64              `json:"totalAllEgg"`
	TotalGoodEgg       uint64              `json:"totalGoodEgg"`
	TotalWeightGoodEgg float64             `json:"totalWeightGoogEgg"`
	TotalCrackedEgg    uint64              `json:"totalCrackedEgg"`
	TotalRejectEgg     uint64              `json:"totalRejectEgg"`
	AbnormalityRate    float64             `json:"abnormalityRate"`
	AverageWeight      float64             `json:"averageWeight"`
	Status             string              `json:"status"`
	CreatedAt          time.Time           `json:"-"`
}

type GetEggMonitoringFilter struct {
	Date       param.DateParam `query:"date"`
	LocationId uint64          `query:"locationId"`
	CageId     uint64
	StartDate  param.DateParam
	EndDate    param.DateParam
}

type GetEggOverviewFilter struct {
	LocationId        uint64                       `query:"locationId"`
	CageId            uint64                       `query:"cageId"`
	OverviewGraphTime param.OverviewGraphTimeParam `query:"overviewGraphTime"`
}

type EggOverviewDetailResponse struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type EggGraphResponse struct {
	Key        string `json:"key"`
	GoodEgg    uint64 `json:"goodEgg"`
	CrackedEgg uint64 `json:"crackedEgg"`
	RejectEgg  uint64 `json:"rejectEgg"`
}

type EggOverviewResponse struct {
	EggOverviewDetail []EggOverviewDetailResponse `json:"eggOverviewDetail"`
	EggGraphs         []EggGraphResponse          `json:"eggGraphs"`
}
