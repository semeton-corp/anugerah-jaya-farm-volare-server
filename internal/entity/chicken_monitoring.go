package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenMonitoring struct {
	Id                       uint64                     `gorm:"primary_key;auto_increment"`
	CageId                   uint64                     `gorm:"type:bigint;not null"`
	TotalLiveChicken         uint64                     `gorm:"type:integer;not null"`
	TotalDeathChicken        uint64                     `gorm:"type:integer;not null"`
	TotalSickChicken         uint64                     `gorm:"type:integer;not null"`
	TotalFeed                float64                    `gorm:"type:decimal;not null"`
	ChickenDisease           []ChickenDiseaseMonitoring `gorm:"foreignKey:ChickenMonitoringId;references:Id"`
	ChickenVaccineMonitoring []ChickenVaccineMonitoring `gorm:"foreignKey:ChickenMonitoringId;references:Id"`
	CreatedBy                uuid.UUID                  `gorm:"type:varchar(255);not null"`
	CreatedAt                time.Time                  `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt                time.Time                  `gorm:"type:timestamp;auto_update_time"`
}
