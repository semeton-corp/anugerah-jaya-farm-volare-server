package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenMonitoring struct {
	Id                uint64        `gorm:"primaryKey;autoIncrement"`
	ChickenCageId     uint64        `gorm:"bigint;not null"`
	ChickenCage       ChickenCage   `gorm:"foreignKey:ChickenCageId;references:Id;constraint:OnDelete:CASCADE"`
	TotalChicken      uint64        `gorm:"bigint;not null;default:0"` // Note : this total chicken is the total chicken cage when monitoring inserted
	TotalDeathChicken uint64        `gorm:"type:bigint;not null;default:0"`
	TotalSickChicken  uint64        `gorm:"type:bigint;not null"`
	TotalFeed         float64       `gorm:"type:decimal;not null"`
	Note              string        `gorm:"type:text"`
	CreatedBy         uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt         time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdateBy          uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt         time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
