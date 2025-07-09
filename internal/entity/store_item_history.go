package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type StoreItemHistory struct {
	Id             uint64                 `gorm:"primaryKey;autoIncrement"`
	ItemId         uint64                 `gorm:"type:bigint;not null"`
	Item           Item                   `gorm:"foreignKey:ItemId;references:Id" json:"-"`
	Source         string                 `gorm:"type:varchar(255)"`
	Destination    string                 `gorm:"type:varchar(255)"`
	QuantityBefore float64                `gorm:"type:decimal;not null"`
	QuantityAfter  float64                `gorm:"type:decimal;not null"`
	Status         enum.ItemHistoryStatus `gorm:"type:int;not null"`
	UserId         uuid.UUID              `gorm:"type:varchar(255);not null"`
	User           User                   `gorm:"foreignKey:UserId;references:Id" json:"-"`
	CreatedAt      time.Time              `gorm:"type:timestamp;autoCreateTime"`
}
