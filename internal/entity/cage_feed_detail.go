package entity

import (
	"time"

	"github.com/google/uuid"
)

type CageFeedDetail struct {
	Id         uint64        `gorm:"primaryKey;autoIncrement"`
	CageFeedId uint64        `gorm:"type:cageFeedId;not null"`
	ItemId     uint64        `gorm:"type:bigint;not null"`
	Item       Item          `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	Percentage float64       `gorm:"decimal;not null"`
	CreatedAt  time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy  uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt  time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy  uuid.NullUUID `gorm:"type:varchar(255)"`
}
