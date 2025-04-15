package entity

import "time"

type Cage struct {
	Id         uint64    `gorm:"primaryKey;autoIncrement"`
	LocationId uint64    `gorm:"type:bigint;not null"`
	Location   Location  `gorm:"foreignKey:LocationId;references:Id"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Capacity   uint64    `gorm:"type:integer;not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt  time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
