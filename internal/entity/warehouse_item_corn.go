package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseItemCorn struct {
	Id          uint64        `gorm:"primaryKey;autoIncrement"`
	WarehouseId uint64        `gorm:"type:bigint;not null"`
	Warehouse   Warehouse     `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	SupplierId  uint64        `gorm:"type:bigint;not null"`
	Supplier    Supplier      `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity    float64       `gorm:"decimal;not null"`
	OrderDate   time.Time     `gorm:"timestamp;not null"`
	ExpiredAt   time.Time     `gorm:"date;not null"`
	CreatedAt   time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy   uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt   time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy   uuid.NullUUID `gorm:"type:varchar(255)"`
}
