package entity

import (
	"time"

	"github.com/google/uuid"
)

type WarehouseItemProcurementDraft struct {
	Id            uint64        `gorm:"primaryKey;autoIncrement"`
	WarehouseId   uint64        `gorm:"type:bigint;not null"`
	Warehouse     Warehouse     `gorm:"foreignKey:WarehouseId;refereces:Id"`
	ItemId        uint64        `gorm:"type:bigint;not null"`
	Item          Item          `gorm:"foreignKey:ItemId;references:Id"`
	SupplierId    uint64        `gorm:"type:bigint;not null"`
	Supplier      Supplier      `gorm:"foreignKey:SupplierId;references:Id"`
	DailySpending float64       `gorm:"type:decimal;not null"`
	DaysNeed      uint64        `gorm:"type:int;not null"`
	TotalOrder    float64       `gorm:"type:decimal;not null"`
	CreatedAt     time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy     uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt     time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy     uuid.NullUUID `gorm:"type:varchar(255)"`
}
