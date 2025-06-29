package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/constant"
)

func EggMonitoringToResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringResponse {
	totalKarpetGoodEgg := eggMonitoring.TotalGoodEgg / constant.TotalEggKarpet
	totalRemainingGoodEgg := eggMonitoring.TotalGoodEgg % constant.TotalEggKarpet

	totalKarpetCrackedEgg := eggMonitoring.TotalCrackedEgg / constant.TotalEggKarpet
	totalRemainingCrackedEgg := eggMonitoring.TotalCrackedEgg % constant.TotalEggKarpet

	totalKarpetRejectEgg := eggMonitoring.TotalRejectEgg / constant.TotalEggKarpet
	totalRemainingRejectEgg := eggMonitoring.TotalRejectEgg % constant.TotalEggKarpet

	response := dto.EggMonitoringResponse{
		Id:                       eggMonitoring.Id,
		Cage:                     CageToResponse(&eggMonitoring.ChickenCage.Cage),
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
		AverageWeight:            float64(uint64(eggMonitoring.TotalWeightGoodEgg*1000.0/float64(eggMonitoring.TotalGoodEgg)*100.0)) / 100.0,
		IsTaken:                  eggMonitoring.IsTaken,
	}

	return response
}

func EggMonitoringToListResponse(eggMonitoring *entity.EggMonitoring) dto.EggMonitoringListResponse {
	response := dto.EggMonitoringListResponse{
		Id:              eggMonitoring.Id,
		Cage:            CageToResponse(&eggMonitoring.ChickenCage.Cage),
		Warehouse:       WarehouseToResponse(&eggMonitoring.Warehouse),
		TotalAllEgg:     eggMonitoring.TotalGoodEgg + eggMonitoring.TotalCrackedEgg + eggMonitoring.TotalRejectEgg,
		TotalGoodEgg:    eggMonitoring.TotalGoodEgg,
		TotalCrackedEgg: eggMonitoring.TotalCrackedEgg,
		TotalRejectEgg:  eggMonitoring.TotalRejectEgg,
		AverageWeight:   float64(uint64(eggMonitoring.TotalWeightGoodEgg*1000.0/float64(eggMonitoring.TotalGoodEgg)*100)) / 100.0,
		IsTaken:         eggMonitoring.IsTaken,
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
