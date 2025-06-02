package entity

import (
	"time"

	"github.com/google/uuid"
)

type Supplier struct {
	Id              uint64        `gorm:"primaryKey;autoIncrement"`
	WarehouseItemId uint64        `gorm:"bigint;not null"`
	WarehouseItem   WarehouseItem `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Name            string        `gorm:"type:varchar(255);not null"`
	PhoneNumber     string        `gorm:"type:varchar(255);not null"`
	Address         string        `gorm:"type:text;not null"`
	CreatedBy       uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt       time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy       uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
