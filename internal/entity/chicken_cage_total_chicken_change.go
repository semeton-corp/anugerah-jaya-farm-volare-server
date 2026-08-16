package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenCageTotalChickenChange struct {
	Id                   uint64        `gorm:"primaryKey;autoIncrement"`
	ChickenCageId        uint64        `gorm:"type:bigint;not null"`
	ChickenCage          ChickenCage   `gorm:"foreignKey:ChickenCageId;references:Id;constraint:OnDelete:CASCADE"`
	PreviousTotalChicken uint64        `gorm:"type:bigint;not null"`
	NewTotalChicken      uint64        `gorm:"type:bigint;not null"`
	CreatedAt            time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy            uuid.NullUUID `gorm:"type:varchar(255)"`
}
