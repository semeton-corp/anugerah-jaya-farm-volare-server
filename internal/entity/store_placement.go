package entity

import (
	"time"

	"github.com/google/uuid"
)

type StorePlacement struct {
	UserId    uuid.UUID     `gorm:"primaryKey;type:varchar(255);not null"`
	User      User          `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	StoreId   uint64        `gorm:"type:bigint;not null"`
	Store     Store         `gorm:"foreignKey:StoreId;refereces:Id;constraint:OnDelete:CASCADE"`
	CreatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
