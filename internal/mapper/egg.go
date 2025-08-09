package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
)

func EggMonitoringToResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringResponse {
	totalKarpetGoodEgg := eggMonitoring.TotalGoodEgg / constant.TotalEggPerKarpet
	totalRemainingGoodEgg := eggMonitoring.TotalGoodEgg % constant.TotalEggPerKarpet

	totalKarpetCrackedEgg := eggMonitoring.TotalCrackedEgg / constant.TotalEggPerKarpet
	totalRemainingCrackedEgg := eggMonitoring.TotalCrackedEgg % constant.TotalEggPerKarpet

	totalKarpetRejectEgg := eggMonitoring.TotalRejectEgg / constant.TotalEggPerKarpet
	totalRemainingRejectEgg := eggMonitoring.TotalRejectEgg % constant.TotalEggPerKarpet

	response := dto.EggMonitoringResponse{
		Id:                       eggMonitoring.Id,
		ChickenCage:              ChickenCageToResponse(&eggMonitoring.ChickenCage),
		Warehouse:                WarehouseToResponse(&eggMonitoring.Warehouse),
		TotalKarpetGoodEgg:       totalKarpetGoodEgg,
		TotalRemainingGoodEgg:    totalRemainingGoodEgg,
		TotalKarpetCrackedEgg:    totalKarpetCrackedEgg,
		TotalRemainingCrackedEgg: totalRemainingCrackedEgg,
		TotalKarpetRejectEgg:     totalKarpetRejectEgg,
		TotalRemainingRejectEgg:  totalRemainingRejectEgg,
		TotalWeightGoodEgg:       eggMonitoring.TotalWeightGoodEgg,
		TotalWeightCrackedEgg:    eggMonitoring.TotalWeightCrackedEgg,
		TotalWeightAllEgg:        eggMonitoring.TotalWeightGoodEgg + eggMonitoring.TotalWeightCrackedEgg,
		CreatedAt:                eggMonitoring.CreatedAt.Format("02 Jan 2006"),
	}

	if eggMonitoring.TotalGoodEgg == 0 {
		response.AverageWeight = 0
	} else {
		response.AverageWeight = float64(uint64(eggMonitoring.TotalWeightGoodEgg*1000.0/float64(eggMonitoring.TotalGoodEgg)*100)) / 100.0
	}

	return response
}

func EggMonitoringToListResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringListResponse {
	response := dto.EggMonitoringListResponse{
		Id:                 eggMonitoring.Id,
		ChickenCage:        ChickenCageToResponse(&eggMonitoring.ChickenCage),
		Warehouse:          WarehouseToResponse(&eggMonitoring.Warehouse),
		TotalAllEgg:        eggMonitoring.TotalGoodEgg + eggMonitoring.TotalCrackedEgg + eggMonitoring.TotalRejectEgg,
		TotalGoodEgg:       eggMonitoring.TotalGoodEgg,
		TotalCrackedEgg:    eggMonitoring.TotalCrackedEgg,
		TotalRejectEgg:     eggMonitoring.TotalRejectEgg,
		TotalWeightGoodEgg: eggMonitoring.TotalWeightGoodEgg,
		CreatedAt:          eggMonitoring.CreatedAt,
	}

	if eggMonitoring.TotalGoodEgg == 0 {
		response.AverageWeight = 0
	} else {
		response.AverageWeight = float64(uint64(eggMonitoring.TotalWeightGoodEgg*1000.0/float64(eggMonitoring.TotalGoodEgg)*100)) / 100.0
	}

	if response.TotalAllEgg == 0 {
		response.AbnormalityRate = 0
		response.Status = constant.EggMonitoringStatusSafety
	} else {
		response.AbnormalityRate = float64(uint64(float64(eggMonitoring.TotalCrackedEgg+eggMonitoring.TotalRejectEgg)/float64(response.TotalAllEgg)*10000.0)) / 100.0

		if response.AbnormalityRate < 0.8 {
			response.Status = constant.EggMonitoringStatusSafety
		} else if response.AbnormalityRate >= 0.8 && response.AbnormalityRate < 1.2 {
			response.Status = constant.EggMonitoringStatusCheck
		} else {
			response.Status = constant.EggMonitoringStatusUrgent
		}
	}

	return response
}
