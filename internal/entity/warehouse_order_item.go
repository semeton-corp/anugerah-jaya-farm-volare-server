package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type WarehouseOrderItem struct {
	Id              uint64                    `gorm:"primaryKey;autoIncrement"`
	WarehouseId     uint64                    `gorm:"type:bigint;not null"`
	Warehouse       Warehouse                 `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	WarehouseItemId uint64                    `gorm:"type:bigint;not null"`
	WarehouseItem   WarehouseItem             `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE"`
	SupplierId      uint64                    `gorm:"type:bigint;not null"`
	Supplier        Supplier                  `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity        uint64                    `gorm:"type:bigint;not null"`
	IsTaken         bool                      `gorm:"type:boolean;default:false"`
	TakenAt         time.Time                 `gorm:"type:timestamp"`
	TakenBy         uuid.UUID                 `gorm:"type:timestamp"`
	Status          enum.WarehouseOrderStatus `gorm:"type:int;not null"`
	CreatedAt       time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.UUID                 `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.UUID                 `gorm:"type:varchar(255)"`
}
