package entity

import (
	"time"

	"github.com/google/uuid"
)

type EggMonitoring struct {
	Id              uint64    `gorm:"primaryKey;autoIncrement"`
	CageId          uint64    `gorm:"type:bigint;not null"`
	Cage            Cage      `gorm:"foreignKey:CageId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	WarehouseId     uint64    `gorm:"type:bigint;not null"`
	Warehouse       Warehouse `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	TotalCrackedEgg uint64    `gorm:"type:bigint;not null"`
	TotalGoodEgg    uint64    `gorm:"type:bigint;not null"`
	TotalBrokeEgg   uint64    `gorm:"type:bigint;not null"`
	TotalRejectEgg  uint64    `gorm:"type:bigint;not null"`
	Weight          float64   `gorm:"type:decimal;not null"`
	IsArrive        bool      `gorm:"type:boolean;default:false"`
	TakenBy         uuid.UUID `gorm:"type:varchar(255)"`
	TakenAt         time.Time `gorm:"type:timestamp"`
	CreatedBy       uuid.UUID `gorm:"type:varchar(255)"`
	CreatedAt       time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy       uuid.UUID `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
