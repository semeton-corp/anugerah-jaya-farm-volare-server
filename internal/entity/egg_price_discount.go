package entity

import (
	"time"

	"github.com/google/uuid"
)

type EggPriceDiscount struct {
	Id                     uint64    `gorm:"primaryKey;autoIncrement"`
	Name                   string    `gorm:"type:varchar(255);not null;unique"`
	MinimumTransactionUser uint64    `gorm:"type:bigint;not null"`
	TotalDiscount          float64   `gorm:"type:decimal;not null"`
	CreatedAt              time.Time `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy              uuid.UUID `gorm:"type:varchar(255)"`
	UpdatedAt              time.Time `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy              uuid.UUID `gorm:"type:varchar(255)"`
}
