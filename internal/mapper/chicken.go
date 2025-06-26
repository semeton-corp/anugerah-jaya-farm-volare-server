package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func ChickenDiseaseMonitoringToResponse(chickenDisease *entity.ChickenDiseaseMonitoring) dto.ChickenDiseaseMonitoringResponse {
	return dto.ChickenDiseaseMonitoringResponse{
		Id:       chickenDisease.Id,
		Disease:  chickenDisease.Disease,
		Medicine: chickenDisease.Medicine,
		Dose:     chickenDisease.Dose,
		Unit:     chickenDisease.Unit,
	}
}

func ChickenVaccineMonitoringToResponse(chickenVaccine *entity.ChickenVaccineMonitoring) dto.ChickenVaccineMonitoringResponse {
	return dto.ChickenVaccineMonitoringResponse{
		Id:      chickenVaccine.Id,
		Vaccine: chickenVaccine.Vaccine,
		Dose:    chickenVaccine.Dose,
		Unit:    chickenVaccine.Unit,
	}
}

// Note : without chickenDiseasesResponse and chickenVaccinesResponse
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
