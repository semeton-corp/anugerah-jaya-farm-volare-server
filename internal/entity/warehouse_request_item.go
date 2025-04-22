package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseRequestItem struct {
	Id              uint64        `gorm:"primaryKey;autoIncrement"`
	WarehouseItemId uint64        `gorm:"type:bigint;not null"`
	WarehouseItem   WarehouseItem `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE"`
	Supplier        string        `gorm:"type:varchar(255);not null"`
	Status          uint8         `gorm:"type:int;not null"`
	CreatedAt       time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.UUID     `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.UUID     `gorm:"type:varchar(255)"`
}
