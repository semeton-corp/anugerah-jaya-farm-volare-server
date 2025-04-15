package entity

import (
	"time"

	"github.com/google/uuid"
)

type DailyWork struct {
	Id          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:text;not null"`
	RoleId      uint64    `gorm:"type:bigint;not null"`
	StartTime   time.Time `gorm:"type:time;not null"`
	EndTime     time.Time `gorm:"type:time;not null"`
	CreatedBy   uuid.UUID `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy   uuid.UUID `gorm:"type:varchar(255);not null"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
