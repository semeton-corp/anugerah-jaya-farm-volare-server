package util

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/internal/entity"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

func GetChickenCategory(chickenCage *entity.ChickenCage) enum.ChickenCategory {
	var (
		chickenAgeInWeek uint64
		chickenCategory  enum.ChickenCategory
	)

	if !chickenCage.ChickenProcurement.CreatedAt.IsZero() {
		chickenAge := time.Since(chickenCage.CreatedAt)
		chickenAgeInWeek = uint64(chickenAge.Hours() / float64((7 * 24)))

		if chickenAgeInWeek <= 9 {
			chickenCategory = enum.ChickenCategoryDOC
		} else if chickenAgeInWeek >= 10 && chickenAgeInWeek <= 15 {
			chickenCategory = enum.ChickenCategoryGrower
		} else if chickenAgeInWeek >= 16 && chickenAgeInWeek <= 17 {
			chickenCategory = enum.ChickenCategoryPreLayer
		} else if chickenAgeInWeek >= 18 {
			chickenCategory = enum.ChickenCategoryLayer
		}
	}

	return chickenCategory
}
