package entity

import (
	"time"

	"github.com/google/uuid"
)

type CagePlacement struct {
	UserId    uuid.UUID     `gorm:"primaryKey;type:varchar(255);not null"`
	User      User          `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	CageId    uint64        `gorm:"primaryKey;type:bigint;not null"`
	Cage      Cage          `gorm:"foreignKey:CageId;references:Id;constraint:OnDelete:CASCADE"`
	CreatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
