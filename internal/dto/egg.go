package dto

type EggMonitoringRequest struct {
	CageId          uint64 `json:"cageId" validate:"required,number"`
	TotalGoodEgg    uint64 `json:"totalGoodEggs" validate:"required,number"`
	TotalCrackedEgg uint64 `json:"totalCrackedEggs" validate:"required,number"`
	TotalBrokeEgg   uint64 `json:"totalBrokeEggs" validate:"required,number"`
	TotalRejectEgg  uint64 `json:"totalRejectEggs" validate:"required,number"`
}

type EggMonitoringResponse struct {
}
