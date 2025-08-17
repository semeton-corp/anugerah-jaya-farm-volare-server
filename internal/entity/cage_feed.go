package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type CageFeed struct {
	Id              uint64               `gorm:"primaryKey;autoIncrement"`
	ChickenCategory enum.ChickenCategory `gorm:"type:int;not null;unique"`
	FeedType        enum.FeedType        `gorm:"type:int;not null"`
	TotalFeed       float64              `gorm:"type:decimal;not null"`
	CageFeedDetails []CageFeedDetail     `gorm:"foreignKey:CageFeedId;references:Id;constraint:OnDelete:CASCADE"`
	CreatedAt       time.Time            `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.NullUUID        `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time            `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.NullUUID        `gorm:"type:varchar(255)"`
}
