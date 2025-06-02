package entity

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type WarehouseActivity struct {
	Id          uint64              `gorm:"primaryKey;autoIncrement"`
	WarehouseId uint64              `gorm:"type:bigint;not null"`
	Warehouse   Warehouse           `gorm:"foreignKey:WarehouseId;references:Id; constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Description string              `gorm:"type:varchar(255);not null"`
	Status      enum.ActivityStatus `gorm:"type:int;not null"`
	CreatedAt   time.Time           `gorm:"type:timestamp;autoCreateTime"`
}
