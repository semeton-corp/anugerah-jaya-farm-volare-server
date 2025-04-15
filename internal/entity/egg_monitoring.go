package entity

import (
	"time"

	"github.com/google/uuid"
)

type EggMonitoring struct {
	ID              uint64    `gorm:"primaryKey;autoIncrement"`
	CageID          uint64    `gorm:"type:integer;not null"`
	Date            time.Time `gorm:"type:date;not null"`
	TotalCrackedEgg uint64    `gorm:"type:integer;not null"`
	TotalGoodEgg    uint64    `gorm:"type:integer;not null"`
	TotalBrokeEgg   uint64    `gorm:"type:integer;not null"`
	TotalRejectEgg  uint64    `gorm:"type:integer;not null"`
	CreatedBy       uuid.UUID `gorm:"type:varchar(26);not null"`
	CreatedAt       time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt       time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
