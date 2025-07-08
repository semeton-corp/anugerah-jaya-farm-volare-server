package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func ChickenHealthItemToResponse(chickenHealthItem *entity.ChickenHealthItem) dto.ChickenHealthItemResponse {
	response := dto.ChickenHealthItemResponse{
		Id:   chickenHealthItem.Id,
		Name: chickenHealthItem.Name,
		Type: chickenHealthItem.Type.String(),
		Note: chickenHealthItem.Note,
	}

	if chickenHealthItem.ChickenAge.Valid {
		valUint64 := uint64(chickenHealthItem.ChickenAge.Int64)
		var chickenCategory string

		if chickenHealthItem.ChickenAge.Int64 >= 0 && valUint64 <= 9 {
			chickenCategory = enum.ChickenCategoryDOC.String()
		} else if valUint64 >= 10 && valUint64 <= 15 {
			chickenCategory = enum.ChickenCategoryGrower.String()
		} else if valUint64 >= 16 && valUint64 <= 17 {
			chickenCategory = enum.ChickenCategoryPreLayer.String()
		} else if valUint64 >= 18 {
			chickenCategory = enum.ChickenCategoryPreLayer.String()
		}

		response.ChickenAge = &valUint64
		response.ChickenCategory = &chickenCategory
	} else {
		response.ChickenAge = nil
		response.ChickenCategory = nil
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

	var chickenCategory string
	if chickenHealthMonitoring.ChickenAge >= 0 && chickenHealthMonitoring.ChickenAge <= 9 {
		chickenCategory = enum.ChickenCategoryDOC.String()
	} else if chickenHealthMonitoring.ChickenAge >= 10 && chickenHealthMonitoring.ChickenAge <= 15 {
		chickenCategory = enum.ChickenCategoryGrower.String()
	} else if chickenHealthMonitoring.ChickenAge >= 16 && chickenHealthMonitoring.ChickenAge <= 17 {
		chickenCategory = enum.ChickenCategoryPreLayer.String()
	} else if chickenHealthMonitoring.ChickenAge >= 18 {
		chickenCategory = enum.ChickenCategoryPreLayer.String()
	}

	response.ChickenAge = chickenHealthMonitoring.ChickenAge
	response.ChickenCategory = chickenCategory

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
