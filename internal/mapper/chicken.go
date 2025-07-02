package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func ChickenHealthItemToResponse(chickenHealthItem *entity.ChickenHealthItem) dto.ChickenHealthItemResponse {
	response := dto.ChickenHealthItemResponse{
		Id:   chickenHealthItem.Id,
		Name: chickenHealthItem.Name,
		Type: chickenHealthItem.Type.String(),
	}

	if chickenHealthItem.ChickenAge.Valid {
		valUint64 := uint64(chickenHealthItem.ChickenAge.Int64)
		response.ChickenAge = &valUint64
	} else {
		response.ChickenAge = nil
	}

	return response
}

func ChickenHealthMonitoringToResponse(chickenHealthMonitoring *entity.ChickenHealthMonitoring) dto.ChickenHealthMonitoringResponse {
	response := dto.ChickenHealthMonitoringResponse{
		Id:             chickenHealthMonitoring.Id,
		HealthItemName: chickenHealthMonitoring.HealthItemName,
		Type:           chickenHealthMonitoring.Type.String(),
		Dose:           chickenHealthMonitoring.Dose,
		Unit:           chickenHealthMonitoring.Unit,
	}

	if chickenHealthMonitoring.Disease.Valid {
		response.Disease = chickenHealthMonitoring.Disease.String
	} else {
		response.Disease = "-"
	}

	return response
}

func ChickenMonitoringToResponse(chickenMonitoring *entity.ChickenMonitoring) dto.ChickenMonitoringResponse {
	return dto.ChickenMonitoringResponse{
		Id:                 chickenMonitoring.Id,
		ChickenCage:        ChickenCageToResponse(&chickenMonitoring.ChickenCage),
		TotalLiveChicken:   chickenMonitoring.ChickenCage.TotalChicken - chickenMonitoring.ChickenCage.TotalDeathChicken,
		TotalSickChicken:   chickenMonitoring.TotalSickChicken,
		TotalDeatchChicken: chickenMonitoring.TotalSickChicken,
		TotalFeed:          chickenMonitoring.TotalFeed,
		Note:               chickenMonitoring.Note,
	}
}

func ChickenMonitoringToListResponse(chickenMonitoring *entity.ChickenMonitoring) dto.ChickenMonitoringListResponse {
	return dto.ChickenMonitoringListResponse{
		Id:                chickenMonitoring.Id,
		ChickenCage:       ChickenCageToResponse(&chickenMonitoring.ChickenCage),
		TotalLiveChicken:  chickenMonitoring.ChickenCage.TotalChicken - chickenMonitoring.ChickenCage.TotalDeathChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		MortalityRate:     float64((chickenMonitoring.TotalDeathChicken / (chickenMonitoring.ChickenCage.TotalChicken + chickenMonitoring.ChickenCage.TotalDeathChicken + chickenMonitoring.TotalDeathChicken)) * 100.0),
	}
}
