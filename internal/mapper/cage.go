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
		Location: dto.LocationResponse{
			Id:   cage.Location.Id,
			Name: cage.Location.Name,
		},
	}
}

func ChickenCageToResponse(chickenCage *entity.ChickenCage) dto.ChickenCageResponse {
	var chickenAgeInWeek uint64
	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() {
		chickenAge := time.Now().Sub(chickenCage.CreatedAt)
		chickenAgeInWeek = uint64(chickenAge.Hours() / float64((7 * 24 * time.Hour)))
	}

	var chickenCategory enum.ChickenCategory
	if chickenAgeInWeek >= 0 {
		chickenCategory = enum.ChickenCategoryAfkir
	}

	var chickenPic, eggPic string
	for _, cagePlacement := range chickenCage.Cage.CagePlacement {
		if strings.Contains(cagePlacement.User.Name, "Kandang") {
			chickenPic = cagePlacement.User.Name
		}

		if strings.Contains(cagePlacement.User.Name, "Telur") {
			chickenPic = cagePlacement.User.Name
		}
	}

	response := dto.ChickenCageResponse{
		Cage:            CageToResponse(&chickenCage.Cage),
		Id:              chickenCage.Id,
		BatchId:         fmt.Sprintf("%s%d", chickenCage.CreatedAt.Format("02012006"), chickenCage.Id),
		ChickenAge:      chickenAgeInWeek,
		ChickenCategory: chickenCategory.String(),
		TotalChicken:    chickenCage.TotalChicken - chickenCage.TotalDeathChicken,
		ChickenPic:      chickenPic,
		EggPic:          eggPic,
	}

	return response
}
