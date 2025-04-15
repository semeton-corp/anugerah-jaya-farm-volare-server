package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenVaccineMonitoring struct {
	Id                  uint64    `gorm:"primaryKey;autoIncrement"`
	ChickenMonitoringId uint64    `gorm:"type:bigint;not null"`
	Vaccine             string    `gorm:"type:varchar(255);not null"`
	Dose                float64   `gorm:"type:decimal;not null"`
	Unit                string    `gorm:"type:varchar(255);not null"`
	CreatedBy           uuid.UUID `gorm:"type:varchar(255);not null"`
	CreatedAt           time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy           uuid.UUID `gorm:"type:varchar(255);not null"`
	UpdatedAt           time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
