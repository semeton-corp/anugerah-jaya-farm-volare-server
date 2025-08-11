package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type WarehouseItem struct {
	ItemId           uint64        `gorm:"primaryKey;type:bigint;not null"`
	Item             Item          `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	WarehouseId      uint64        `gorm:"primaryKey;type:bigint;not null"`
	Warehouse        Warehouse     `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity         float64       `gorm:"type:decimal;not null;default:0"`
	EstimationRunOut sql.NullTime  `gorm:"type:date"`
	ExpiredAt        sql.NullTime  `gorm:"type:date"`
	CreatedAt        time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy        uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt        time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy        uuid.NullUUID `gorm:"type:varchar(255)"`
}
