package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehousePlacement struct {
	UserId      uuid.UUID     `gorm:"primaryKey;type:varchar(255);not null"`
	User        User          `gorm:"foreignKey:UserId;references:Id"`
	WarehouseId uint64        `gorm:"primaryKey;type:bigint;not null"`
	Warehouse   Warehouse     `gorm:"foreignKey:WarehouseId;refereces:Id"`
	CreatedBy   uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt   time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy   uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt   time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
