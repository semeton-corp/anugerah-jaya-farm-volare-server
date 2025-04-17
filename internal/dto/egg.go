package dto

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type EggMonitoringRequest struct {
	CageId          uint64 `json:"cageId" validate:"required,number"`
	TotalGoodEgg    uint64 `json:"totalGoodEgg" validate:"required,number"`
	TotalCrackedEgg uint64 `json:"totalCrackedEgg" validate:"required,number"`
	TotalBrokeEgg   uint64 `json:"totalBrokeEgg" validate:"required,number"`
	TotalRejectEgg  uint64 `json:"totalRejectEgg" validate:"required,number"`
}

type EggMonitoringResponse struct {
	Id              uint64       `json:"id"`
	Cage            CageResponse `json:"cage"`
	TotalGoodEgg    uint64       `json:"totalGoodEggs"`
	TotalCrackedEgg uint64       `json:"totalCrackedEggs"`
	TotalBrokeEgg   uint64       `json:"totalBrokeEggs"`
	TotalRejectEgg  uint64       `json:"totalRejectEggs"`
}

type EggMonitoringListResponse struct {
	Id              uint64       `json:"id"`
	Cage            CageResponse `json:"cage"`
	TotalAll        uint64       `json:"totalAllEgg"`
	TotalGoodEgg    uint64       `json:"totalGoodEgg"`
	TotalCrackedEgg uint64       `json:"totalCrackedEgg"`
	TotalBrokeEgg   uint64       `json:"totalBrokeEgg"`
	TotalRejectEgg  uint64       `json:"totalRejectEgg"`
	AbnormalityRate float64      `json:"abnormalityRate"`
	Description     string       `json:"description"`
}

type GetEggMonitoringFilter struct {
	Date param.DateParam `query:"date"`
}
