package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehousePlacement struct {
	UserId     uuid.UUID     `gorm:"type:varchar(255);not null"`
	User       User          `gorm:"foreignKey:UserId;references:Id"`
	WarehousId uint64        `gorm:"type:bigint;not null"`
	Warehouse  Warehouse     `gorm:"foreignKey:WarehouseId;refereces:Id"`
	CreatedBy  uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt  time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy  uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt  time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
