package entity

import (
	"time"

	"github.com/google/uuid"
)

type StoreItem struct {
	StoreId   uint64        `gorm:"primaryKey;type:bigint;not null"`
	Store     Store         `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	ItemId    uint64        `gorm:"primaryKey;type:bigint;not null"`
	Item      Item          `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Quantity  float64       `gorm:"type:decimal;not null"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
}
