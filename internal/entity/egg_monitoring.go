package entity

import (
	"time"

	"github.com/google/uuid"
)

type EggMonitoring struct {
	Id              uint64    `gorm:"primaryKey;autoIncrement"`
	CageId          uint64    `gorm:"type:bigint;not null"`
	Cage            Cage      `gorm:"foreignKey:CageId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	TotalCrackedEgg uint64    `gorm:"type:bigint;not null"`
	TotalGoodEgg    uint64    `gorm:"type:bigint;not null"`
	TotalBrokeEgg   uint64    `gorm:"type:bigint;not null"`
	TotalRejectEgg  uint64    `gorm:"type:bigint;not null"`
	CreatedBy       uuid.UUID `gorm:"type:varchar(255);not null"`
	CreatedAt       time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy       uuid.UUID `gorm:"type:varchar(255);not null"`
	UpdatedAt       time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
