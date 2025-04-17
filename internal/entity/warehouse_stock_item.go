package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseStockItem struct {
	WarehouseItemId  uint64        `gorm:"primaryKey;type:bigint;not null"`
	WarehouseItem    WarehouseItem `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE"`
	WarehouseId      uint64        `gorm:"primaryKey;type:bigint;not null"`
	Warehouse        Warehouse     `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity         uint64        `gorm:"type:bigint;not null"`
	EstimationRunOut time.Time     `gorm:"type:date;not null"`
	CreatedAt        time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy        uuid.UUID     `gorm:"type:varchar(255);not null"`
	UpdatedAt        time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy        uuid.UUID     `gorm:"type:varchar(255);not null"`
}
