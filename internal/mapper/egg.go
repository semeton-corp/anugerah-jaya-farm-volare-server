package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func EggMonitoringToResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringResponse {
	return dto.EggMonitoringResponse{
		Id: eggMonitoring.Id,
		Cage: dto.CageResponse{
			Id:   eggMonitoring.Cage.Id,
			Name: eggMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   eggMonitoring.Cage.Location.Id,
				Name: eggMonitoring.Cage.Location.Name,
			},
		},
		TotalGoodEgg:    eggMonitoring.TotalGoodEgg,
		TotalCrackedEgg: eggMonitoring.TotalCrackedEgg,
		TotalBrokeEgg:   eggMonitoring.TotalBrokeEgg,
		TotalRejectEgg:  eggMonitoring.TotalRejectEgg,
	}
}

func EggMonitoringToListResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringListResponse {
	return dto.EggMonitoringListResponse{
		Id: eggMonitoring.Id,
		Cage: dto.CageResponse{
			Id:   eggMonitoring.Cage.Id,
			Name: eggMonitoring.Cage.Name,
			Location: dto.LocationResponse{
				Id:   eggMonitoring.Cage.Location.Id,
				Name: eggMonitoring.Cage.Location.Name,
			},
		},
		TotalGoodEgg:    eggMonitoring.TotalGoodEgg,
		TotalCrackedEgg: eggMonitoring.TotalCrackedEgg,
		TotalBrokeEgg:   eggMonitoring.TotalBrokeEgg,
		TotalRejectEgg:  eggMonitoring.TotalRejectEgg,
		TotalAll:        eggMonitoring.TotalGoodEgg + eggMonitoring.TotalCrackedEgg + eggMonitoring.TotalBrokeEgg + eggMonitoring.TotalRejectEgg,
	}
}
