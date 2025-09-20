package entity

import (
	"time"

	"github.com/google/uuid"
)

type CageFeedStock struct {
	Id        uint64        `gorm:"primaryKey;autoIncrement"`
	CageId    uint64        `gorm:"type:bigint;not null"`
	TotalFeed float64       `gorm:"type:decimal;not null"`
	UsedFeed  float64       `gorm:"decimal;not null"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
}
