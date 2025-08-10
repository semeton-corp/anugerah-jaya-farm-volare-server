package mapper

import (
	"fmt"
	"strings"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
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
		chickenAgeInWeek uint64
		chickenCategory  enum.ChickenCategory
		batchId          string = ""
	)

	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() {
		batchId = fmt.Sprintf("%s%d", chickenCage.ChickenProcurement.CreatedAt.Format("02012006"), chickenCage.Id)
		chickenAge := time.Since(chickenCage.CreatedAt)
		chickenAgeInWeek = uint64(chickenAge.Hours() / float64((7 * 24)))

		if chickenAgeInWeek >= 0 && chickenAgeInWeek <= 9 {
			chickenCategory = enum.ChickenCategoryDOC
		} else if chickenAgeInWeek >= 10 && chickenAgeInWeek <= 15 {
			chickenCategory = enum.ChickenCategoryGrower
		} else if chickenAgeInWeek >= 16 && chickenAgeInWeek <= 17 {
			chickenCategory = enum.ChickenCategoryPreLayer
		} else if chickenAgeInWeek >= 18 {
			chickenCategory = enum.ChickenCategoryPreLayer
		}
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
		ChickenAge:           chickenAgeInWeek,
		ChickenCategory:      chickenCategory.String(),
		TotalChicken:         chickenCage.TotalChicken - chickenCage.TotalDeathChicken,
		ChickenPic:           chickenPic,
		EggPic:               eggPic,
		TotalDeathChicken:    chickenCage.TotalDeathChicken,
		IsNeedRoutineVaccine: chickenCage.IsNeedRoutineVaccine,
		ChickenProcurementId: chickenCage.ChickenProcurementId,
	}

	return response
}

func CageFeedToResponse(data *entity.CageFeed) dto.CageFeedResponse {
	return dto.CageFeedResponse{
		Id:              data.Id,
		ChickenCategory: data.ChickenCategory.String(),
		FeedType:        data.FeedType.String(),
		TotalFeed:       data.TotalFeed,
	}
}

func CageFeedDetailToResponse(data *entity.CageFeedDetail) dto.CageFeedDetailResponse {
	return dto.CageFeedDetailResponse{
		Id:         data.Id,
		Item:       ItemToResponse(&data.Item),
		Percentage: data.Percentage,
	}
}

// Note : without totalFeed and yesterdayTotalFeed
func ChickenCageFeedToListResponse(chickenCage *entity.ChickenCage) dto.ChickenCageFeedListResponse {
	var (
		chickenAgeInWeek uint64
		chickenCategory  enum.ChickenCategory
	)

	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() {
		chickenAge := time.Since(chickenCage.CreatedAt)
		chickenAgeInWeek = uint64(chickenAge.Hours() / float64((7 * 24)))

		if chickenAgeInWeek >= 0 && chickenAgeInWeek <= 9 {
			chickenCategory = enum.ChickenCategoryDOC
		} else if chickenAgeInWeek >= 10 && chickenAgeInWeek <= 15 {
			chickenCategory = enum.ChickenCategoryGrower
		} else if chickenAgeInWeek >= 16 && chickenAgeInWeek <= 17 {
			chickenCategory = enum.ChickenCategoryPreLayer
		} else if chickenAgeInWeek >= 18 {
			chickenCategory = enum.ChickenCategoryPreLayer
		}
	}

	response := dto.ChickenCageFeedListResponse{
		Id:                chickenCage.Id,
		Cage:              CageToResponse(&chickenCage.Cage),
		ChickenCategory:   chickenCategory.String(),
		ChickenAge:        chickenAgeInWeek,
		TotalChicken:      chickenCage.TotalChicken - chickenCage.TotalDeathChicken,
		IsNeedFeed:        chickenCage.IsNeedFeed,
		ExpectedTotalFeed: chickenCage.Cage.CageFeed.TotalFeed,
	}

	return response
}

// Note : without totalFeed and yesterdayTotalFeed and feedDetails
func ChickenCageFeedToResponse(chickenCage *entity.ChickenCage) dto.ChickenCageFeedResponse {
	var (
		chickenAgeInWeek uint64
		chickenCategory  enum.ChickenCategory
	)

	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() {
		chickenAge := time.Since(chickenCage.CreatedAt)
		chickenAgeInWeek = uint64(chickenAge.Hours() / float64((7 * 24)))

		if chickenAgeInWeek >= 0 && chickenAgeInWeek <= 9 {
			chickenCategory = enum.ChickenCategoryDOC
		} else if chickenAgeInWeek >= 10 && chickenAgeInWeek <= 15 {
			chickenCategory = enum.ChickenCategoryGrower
		} else if chickenAgeInWeek >= 16 && chickenAgeInWeek <= 17 {
			chickenCategory = enum.ChickenCategoryPreLayer
		} else if chickenAgeInWeek >= 18 {
			chickenCategory = enum.ChickenCategoryPreLayer
		}
	}

	response := dto.ChickenCageFeedResponse{
		Id:                chickenCage.Id,
		Cage:              CageToResponse(&chickenCage.Cage),
		ChickenCategory:   chickenCategory.String(),
		ChickenAge:        chickenAgeInWeek,
		TotalChicken:      chickenCage.TotalChicken - chickenCage.TotalDeathChicken,
		IsNeedFeed:        chickenCage.IsNeedFeed,
		ExpectedTotalFeed: chickenCage.Cage.CageFeed.TotalFeed,
	}

	return response
}
