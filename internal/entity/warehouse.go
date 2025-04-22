package entity

import "github.com/google/uuid"

type Warehouse struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	LocationId uint64    `gorm:"type:bigint;not null"`
	Location   Location  `gorm:"foreignKey:LocationId;references:Id;constraint:OnDelete:CASCADE"`
	Name       string    `gorm:"type:varchar(255);not null"`
	CreatedAt  string    `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy  uuid.UUID `gorm:"type:varchar(255)"`
	UpdatedAt  string    `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy  uuid.UUID `gorm:"type:varchar(255)"`
}
