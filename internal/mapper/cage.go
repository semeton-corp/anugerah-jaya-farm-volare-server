package mapper

import (
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/dto"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
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
	return dto.ChickenCageResponse{
		Id:           chickenCage.Id,
		TotalChicken: chickenCage.TotalChicken,
	}
}
