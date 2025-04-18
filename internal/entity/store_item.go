package entity

import (
	"time"

	"github.com/google/uuid"
)

type StoreItem struct {
	StoreId         uint64        `gorm:"primaryKey;not null"`
	Store           Store         `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	WarehouseItemId uint64        `gorm:"primaryKey;not null"`
	WarehouseItem   WarehouseItem `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Quantity        uint64        `gorm:"type:bigint;not null"`
	CreatedAt       time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.UUID     `gorm:"type:varchar(255);not null"`
	UpdatedAt       time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.UUID     `gorm:"type:varchar(255);not null"`
}
