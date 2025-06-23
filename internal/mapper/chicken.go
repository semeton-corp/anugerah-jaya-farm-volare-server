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
		Id: chickenMonitoring.Id,
		// ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			// Id:   chickenMonitoring.Cage.Id,
			// Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				// Id:   chickenMonitoring.Cage.Location.Id,
				// Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		// Age:               chickenMonitoring.Age,
		// TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
	}
}

func ChickenMonitoringToListResponse(chickenMonitoring *entity.ChickenMonitoring) dto.ChickenMonitoringListResponse {
	return dto.ChickenMonitoringListResponse{
		Id: chickenMonitoring.Id,
		// ChickenCategory: chickenMonitoring.ChickenCategory.String(),
		Cage: dto.CageResponse{
			// Id:   chickenMonitoring.Cage.Id,
			// Name: chickenMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				// Id:   chickenMonitoring.Cage.Location.Id,
				// Name: chickenMonitoring.Cage.Location.Name,
			},
		},
		// Age:               chickenMonitoring.Age,
		// TotalLiveChicken:  chickenMonitoring.TotalLiveChicken,
		TotalSickChicken:  chickenMonitoring.TotalSickChicken,
		TotalDeathChicken: chickenMonitoring.TotalDeathChicken,
		TotalFeed:         chickenMonitoring.TotalFeed,
		// MortalityRate:     float64((chickenMonitoring.TotalDeathChicken / (chickenMonitoring.TotalLiveChicken + chickenMonitoring.TotalSickChicken)) * 100.0), // Todo : fix this calculation
	}
}
