package dto

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/param"
)

type CreateEggMonitoringRequest struct {
	CageId          uint64  `json:"cageId" validate:"required,number"`
	WarehouseId     uint64  `json:"warehouseId" validate:"required,number"`
	TotalGoodEgg    uint64  `json:"totalGoodEgg" validate:"number,min=0"`
	TotalCrackedEgg uint64  `json:"totalCrackedEgg" validate:"number,min=0"`
	TotalBrokeEgg   uint64  `json:"totalBrokeEgg" validate:"number,min=0"`
	TotalRejectEgg  uint64  `json:"totalRejectEgg" validate:"number,min=0"`
	Weight          float64 `json:"weight" validate:"number,min=0"`
}

type UpdateEggMonitoringRequest struct {
	CageId          uint64  `json:"cageId" validate:"required,number"`
	WarehouseId     uint64  `json:"warehouseId" validate:"required,number"`
	TotalGoodEgg    uint64  `json:"totalGoodEgg" validate:"required,min=0"`
	TotalCrackedEgg uint64  `json:"totalCrackedEgg" validate:"required,min=0"`
	TotalBrokeEgg   uint64  `json:"totalBrokeEgg" validate:"required,min=0"`
	TotalRejectEgg  uint64  `json:"totalRejectEgg" validate:"required,min=0"`
	Weight          float64 `json:"weight" validate:"required,min=0"`
}

type EggMonitoringResponse struct {
	Id              uint64            `json:"id"`
	Cage            CageResponse      `json:"cage"`
	Warehouse       WarehouseResponse `json:"warehouse"`
	TotalGoodEgg    uint64            `json:"totalGoodEggs"`
	Weight          float64           `json:"weight"`
	TotalCrackedEgg uint64            `json:"totalCrackedEggs"`
	TotalBrokeEgg   uint64            `json:"totalBrokeEggs"`
	TotalRejectEgg  uint64            `json:"totalRejectEggs"`
}

type EggMonitoringListResponse struct {
	Id              uint64            `json:"id"`
	Cage            CageResponse      `json:"cage"`
	Warehouse       WarehouseResponse `json:"warehouse"`
	TotalAll        uint64            `json:"totalAllEgg"`
	TotalGoodEgg    uint64            `json:"totalGoodEgg"`
	TotalCrackedEgg uint64            `json:"totalCrackedEgg"`
	TotalBrokeEgg   uint64            `json:"totalBrokeEgg"`
	TotalRejectEgg  uint64            `json:"totalRejectEgg"`
	AbnormalityRate float64           `json:"abnormalityRate"`
	Weight          float64           `json:"weight"`
	Description     string            `json:"description"`
}

type GetEggMonitoringFilter struct {
	Date  param.DateParam `query:"date"`
	Limit uint64
}
