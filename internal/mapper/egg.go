package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
)

func EggMonitoringToResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringResponse {
	return dto.EggMonitoringResponse{
		Id: eggMonitoring.Id,
		Cage: dto.CageResponse{
			Id:   eggMonitoring.ChickenCage.Cage.Id,
			Name: eggMonitoring.ChickenCage.Cage.Name,
			Location: dto.LocationResponse{
				Id:   eggMonitoring.ChickenCage.Cage.Location.Id,
				Name: eggMonitoring.ChickenCage.Cage.Location.Name,
			},
		},
		Warehouse: dto.WarehouseResponse{
			Id:   eggMonitoring.Warehouse.Id,
			Name: eggMonitoring.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   eggMonitoring.Warehouse.Location.Id,
				Name: eggMonitoring.Warehouse.Location.Name,
			},
		},
		TotalGoodEgg:    eggMonitoring.TotalGoodEgg,
		TotalCrackedEgg: eggMonitoring.TotalCrackedEgg,
		TotalRejectEgg:  eggMonitoring.TotalRejectEgg,
	}
}

// Note : without AbnormalityRate and Description
func EggMonitoringToListResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringListResponse {
	return dto.EggMonitoringListResponse{
		Id: eggMonitoring.Id,
		Cage: dto.CageResponse{
			Id:   eggMonitoring.ChickenCage.Cage.Id,
			Name: eggMonitoring.ChickenCage.Cage.Name,
			Location: dto.LocationResponse{
				Id:   eggMonitoring.ChickenCage.Cage.Location.Id,
				Name: eggMonitoring.ChickenCage.Cage.Location.Name,
			},
		},
		Warehouse: dto.WarehouseResponse{
			Id:   eggMonitoring.Warehouse.Id,
			Name: eggMonitoring.Warehouse.Name,
			Location: dto.LocationResponse{
				Id:   eggMonitoring.Warehouse.Location.Id,
				Name: eggMonitoring.Warehouse.Location.Name,
			},
		},
		TotalGoodEgg:    eggMonitoring.TotalGoodEgg,
		TotalCrackedEgg: eggMonitoring.TotalCrackedEgg,
		TotalRejectEgg:  eggMonitoring.TotalRejectEgg,
	}
}
