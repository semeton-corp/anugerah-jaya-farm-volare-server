package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenHealthMonitoring struct {
	Id                     uint64            `gorm:"primaryKey;autoIncrement"`
	ChickenCageId          uint64            `gorm:"type:bigint;not null"`
	ChickenHealthProductId uint64            `gorm:"type:int;not null"`
	ChickenHealthProduct   ChickenHealthItem `gorm:"foreingKey:ChickenHealthProductId;references:Id"`
	Dose                   float64           `gorm:"type:decimal;not null"`
	Unit                   string            `gorm:"type:varchar(255);not null"`
	Disease                string            `gorm:"type:varchar(255);default:'-'"`
	CreatedBy              uuid.NullUUID     `gorm:"type:varchar(255)"`
	CreatedAt              time.Time         `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy              uuid.NullUUID     `gorm:"type:varchar(255)"`
	UpdatedAt              time.Time         `gorm:"type:timestamp;autoUpdateTime"`
}
