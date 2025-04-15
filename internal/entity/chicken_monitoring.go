package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenMonitoring struct {
	Id                       uint64                     `gorm:"primaryKey;autoIncrement"`
	CageId                   uint64                     `gorm:"type:bigint;not null"`
	Cage                     Cage                       `gorm:"foreignKey:CageId;references:Id"`
	Age                      uint64                     `gorm:"type:bigint;not null"`
	TotalLiveChicken         uint64                     `gorm:"type:bigint;not null"`
	TotalDeathChicken        uint64                     `gorm:"type:bigint;not null"`
	TotalSickChicken         uint64                     `gorm:"type:bigint;not null"`
	TotalFeed                float64                    `gorm:"type:decimal;not null"`
	ChickenDiseaseMonitoring []ChickenDiseaseMonitoring `gorm:"foreignKey:ChickenMonitoringId;references:Id;constraint:OnDelete:CASCADE"`
	ChickenVaccineMonitoring []ChickenVaccineMonitoring `gorm:"foreignKey:ChickenMonitoringId;references:Id;constraint:OnDelete:CASCADE"`
	CreatedBy                uuid.UUID                  `gorm:"type:varchar(255);not null"`
	CreatedAt                time.Time                  `gorm:"type:timestamp;autoCreateTime"`
	UpdateBy                 uuid.UUID                  `gorm:"type:varchar(255);not null"`
	UpdatedAt                time.Time                  `gorm:"type:timestamp;autoUpdateTime"`
}
