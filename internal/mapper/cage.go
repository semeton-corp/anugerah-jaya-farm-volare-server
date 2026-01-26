package mapper

import (
	"fmt"
	"strings"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/util"
)

func CageToResponse(cage *entity.Cage) dto.CageResponse {
	return dto.CageResponse{
		Id:              cage.Id,
		Name:            cage.Name,
		Capacity:        cage.Capacity,
		ChickenCategory: cage.ChickenCategory.String(),
		IsUsed:          cage.IsUsed,
		Location: dto.LocationResponse{
			Id:   cage.Location.Id,
			Name: cage.Location.Name,
		},
	}
}

func ChickenCageToResponse(chickenCage *entity.ChickenCage) dto.ChickenCageResponse {
	var (
		batchId string = ""
	)

	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() {
		batchId = fmt.Sprintf("%s-%d", chickenCage.ChickenProcurement.CreatedAt.Format("02012006"), chickenCage.ChickenProcurement.Id)
	}

	var chickenPic, eggPic string
	for _, cagePlacement := range chickenCage.Cage.CagePlacement {
		if strings.Contains(cagePlacement.User.Role.Name, "Kandang") {
			chickenPic = cagePlacement.User.Name
		}

		if strings.Contains(cagePlacement.User.Role.Name, "Telur") {
			eggPic = cagePlacement.User.Name
		}
	}

	response := dto.ChickenCageResponse{
		Cage:                 CageToResponse(&chickenCage.Cage),
		Id:                   chickenCage.Id,
		BatchId:              batchId,
		ChickenAge:           util.GetChickenAgeByChickenCage(chickenCage),
		ChickenCategory:      util.GetChickenCategoryByChickenCage(chickenCage).String(),
		TotalChicken:         chickenCage.TotalChicken,
		ChickenPic:           chickenPic,
		EggPic:               eggPic,
		IsNeedRoutineVaccine: chickenCage.IsNeedRoutineVaccine,
	}

	if chickenCage.ChickenProcurementId.Valid {
		response.ChickenProcurementId = uint64(chickenCage.ChickenProcurementId.Int64)
	}

	if chickenCage.LatestChickenAgeVaccineRoutine.Valid {
		response.LatestChickenAgeVaccineRoutine = &chickenCage.LatestChickenAgeVaccineRoutine.Int64
	}

	return response
}

func CageFeedToResponse(data *entity.CageFeed) dto.CageFeedResponse {
	response := dto.CageFeedResponse{
		Id:              data.Id,
		ChickenCategory: data.ChickenCategory.String(),
		FeedType:        data.FeedType.String(),
		TotalFeed:       data.TotalFeed,
	}

	switch data.ChickenCategory {
	case enum.ChickenCategoryDOC:
		response.ChickenAgeInterval = "0 - 9 Minggu"
	case enum.ChickenCategoryGrower:
		response.ChickenAgeInterval = "10 - 15 Minggu"
	case enum.ChickenCategoryPreLayer:
		response.ChickenAgeInterval = "16 - 17 Minggu"
	case enum.ChickenCategoryLayer:
		response.ChickenAgeInterval = ">= 18 Minggu"
	}

	return response
}

func CageFeedDetailToResponse(data *entity.CageFeedDetail) dto.CageFeedDetailResponse {
	return dto.CageFeedDetailResponse{
		Id:         data.Id,
		Item:       ItemToResponse(&data.Item),
		Percentage: data.Percentage,
	}
}

func ChickenCageFeedToListResponse(chickenCage *entity.ChickenCage) dto.ChickenCageFeedListResponse {
	response := dto.ChickenCageFeedListResponse{
		Id:              chickenCage.Id,
		Cage:            CageToResponse(&chickenCage.Cage),
		ChickenCategory: util.GetChickenCategoryByChickenCage(chickenCage).String(),
		ChickenAge:      util.GetChickenAgeByChickenCage(chickenCage),
		TotalChicken:    chickenCage.TotalChicken,
		IsNeedFeed:      chickenCage.IsNeedFeed,
	}

	return response
}

func ChickenCageFeedToResponse(chickenCage *entity.ChickenCage) dto.ChickenCageFeedResponse {
	response := dto.ChickenCageFeedResponse{
		Id:              chickenCage.Id,
		Cage:            CageToResponse(&chickenCage.Cage),
		ChickenCategory: util.GetChickenCategoryByChickenCage(chickenCage).String(),
		ChickenAge:      util.GetChickenAgeByChickenCage(chickenCage),
		TotalChicken:    chickenCage.TotalChicken,
		IsNeedFeed:      chickenCage.IsNeedFeed,
	}

	return response
}
