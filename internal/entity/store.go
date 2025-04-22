package entity

import (
	"time"

	"github.com/google/uuid"
)

type Store struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	Name       string    `gorm:"type:varchar(255);not null"`
	LocationId uint64    `gorm:"type:bigint;not null"`
	Location   Location  `gorm:"foreignKey:LocationId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy  uuid.UUID `gorm:"type:varchar(255)"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy  uuid.UUID `gorm:"type:varchar(255);not null"`
}
