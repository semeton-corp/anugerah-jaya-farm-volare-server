	package entity

import (
	"time"

	"github.com/google/uuid"
)

type SupplierItem struct {
	SupplierId uint64        `gorm:"bigint;not null;primaryKey"`
	Supplier   Supplier      `gorm:"foreignKey:SupplierId;references:Id"`
	ItemId     uint64        `gorm:"bigint;not null;primaryKey"`
	Item       Item          `gorm:"foreignKey:ItemId;references:Id"`
	CreatedBy  uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt  time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy  uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt  time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
