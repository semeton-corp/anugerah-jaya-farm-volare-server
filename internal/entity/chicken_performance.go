package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type ChickenPerformance struct {
	Id                           uint64                   `gorm:"primaryKey:autoIncrement"`
	CageName                     string                   `gorm:"type:varchar(255);not null"`
	LocationId                   uint64                   `gorm:"type:bigint;not null"`
	ChickenCategory              enum.ChickenCategory     `gorm:"type:int;not null"`
	ChickenAge                   uint64                   `gorm:"type:int;not null"`
	TotalChicken                 uint64                   `gorm:"type:int;not null"`
	TotalGoodEgg                 uint64                   `gorm:"type:int;not null"`
	AverageConsumptionPerChicken float64                  `gorm:"type:decimal;not null"`
	AverageWeightPerGoodEgg      float64                  `gorm:"type:decimal;not null"`
	FCR                          float64                  `gorm:"type:decimal;not null"`
	HDP                          float64                  `gorm:"type:decimal;not null"`
	MortalityRate                float64                  ` gorm:"type:decimal;not null"`
	Productivity                 enum.ChickenProductivity `gorm:"type:int;not null"`
	CreatedAt                    time.Time                `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy                    uuid.NullUUID            `gorm:"type:varchar(255)"`
	UpdatedAt                    time.Time                `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy                    uuid.NullUUID            `gorm:"type:varchar(255)"`
}
