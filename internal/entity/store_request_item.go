package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type StoreRequestItem struct {
	Id              uint64                 `gorm:"primaryKey;autoIncrement"`
	WarehouseId     uint64                 `gorm:"bigint;not null"`
	Warehouse       Warehouse              `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	WarehouseItemId uint64                 `gorm:"bigint;not null"`
	WarehouseItem   WarehouseItem          `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	StoreId         uint64                 `gorm:"bigint;not null"`
	Store           Store                  `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Quantity        uint64                 `gorm:"bigint;not null"`
	Status          enum.RequestItemStatus `gorm:"int;not null"`
	CreatedAt       time.Time              `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.UUID              `gorm:"type:varchar(255);not null"`
	UpdatedAt       time.Time              `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.UUID              `gorm:"type:varchar(255);not null"`
}
