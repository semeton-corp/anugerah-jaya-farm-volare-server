package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenVaccineMonitoring struct {
	Id                  uint64        `gorm:"primaryKey;autoIncrement"`
	ChickenMonitoringId uint64        `gorm:"type:bigint;not null"`
	Vaccine             string        `gorm:"type:varchar(255);not null"`
	Dose                float64       `gorm:"type:decimal;not null"`
	Unit                string        `gorm:"type:varchar(255);not null"`
	CreatedBy           uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt           time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy           uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
