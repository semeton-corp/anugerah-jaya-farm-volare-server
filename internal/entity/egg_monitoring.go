package entity

import (
	"time"

	"github.com/google/uuid"
)

type EggMonitoring struct {
	Id                    uint64        `gorm:"primaryKey;autoIncrement"`
	ChickenCageId         uint64        `gorm:"bigint;not null"`
	ChickenCage           ChickenCage   `gorm:"foreingKey:ChickenCageId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	WarehouseId           uint64        `gorm:"type:bigint;not null"`
	Warehouse             Warehouse     `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	TotalCrackedEgg       uint64        `gorm:"type:bigint;not null"`
	TotalGoodEgg          uint64        `gorm:"type:bigint;not null"`
	TotalRejectEgg        uint64        `gorm:"type:bigint;not null"`
	TotalWeightGoodEgg    float64       `gorm:"decimal;not null"`
	TotalWeightCrackedEgg float64       `gorm:"decimal;not null"`
	CreatedBy             uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt             time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy             uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt             time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
