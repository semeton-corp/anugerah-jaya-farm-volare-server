package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type Cage struct {
	Id              uint64               `gorm:"primaryKey;autoIncrement"`
	LocationId      uint64               `gorm:"type:bigint;not null"`
	Location        Location             `gorm:"foreignKey:LocationId;references:Id;constraint:OnDelete:CASCADE"`
	Name            string               `gorm:"type:varchar(255);not null"`
	Capacity        uint64               `gorm:"type:bigint;not null"`
	ChickenCategory enum.ChickenCategory `gorm:"type:bigint;not null"`
	CreatedAt       time.Time            `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.NullUUID        `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time            `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.NullUUID        `gorm:"type:varchar(255)"`
}
