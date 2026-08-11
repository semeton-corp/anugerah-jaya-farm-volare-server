package util

import (
	"math"
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func GetChickenCategoryByChickenCage(chickenCage *entity.ChickenCage) enum.ChickenCategory {
	return GetChickenCategoryByChickenCageAt(chickenCage, time.Now())
}

func GetChickenCategoryByChickenCageAt(chickenCage *entity.ChickenCage, asOf time.Time) enum.ChickenCategory {
	chickenAgeInWeek := GetChickenAgeByChickenCageAt(chickenCage, asOf)

	return GetChickenCategoryByAge(chickenAgeInWeek)
}

func GetChickenCategoryByAge(chickenAgeInWeek uint64) enum.ChickenCategory {
	if chickenAgeInWeek <= 9 {
		return enum.ChickenCategoryDOC
	} else if chickenAgeInWeek >= 10 && chickenAgeInWeek <= 15 {
		return enum.ChickenCategoryGrower
	} else if chickenAgeInWeek >= 16 && chickenAgeInWeek <= 17 {
		return enum.ChickenCategoryPreLayer
	} else if chickenAgeInWeek >= 18 && chickenAgeInWeek < 90 {
		return enum.ChickenCategoryLayer
	} else if chickenAgeInWeek >= 90 {
		return enum.ChickenCategoryAfkir
	}

	return enum.ChickenCategoryUnknown
}

func GetChickenAgeByChickenCage(chickenCage *entity.ChickenCage) uint64 {
	return GetChickenAgeByChickenCageAt(chickenCage, time.Now())
}

func GetChickenAgeByChickenCageAt(chickenCage *entity.ChickenCage, asOf time.Time) uint64 {
	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() && !asOf.IsZero() {
		chickenAge := asOf.UTC().Add(time.Hour * 7).Sub(chickenCage.ChickenProcurement.CreatedAt)
		chickenAgeInWeek := uint64(math.Floor(math.Abs(chickenAge.Hours()) / float64(7*24)))
		return chickenAgeInWeek
	}
	return 0
}
