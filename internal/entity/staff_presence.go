package entity

import (
	"time"

	"github.com/google/uuid"
)

type StaffPresence struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement"`
	StaffId   uuid.UUID `gorm:"type:bigint;not null"`
	Staff     Staff     `gorm:"foreignKey:StaffId;references:Id;constraint:OnDelete:CASCADE"`
	StartTime time.Time `gorm:"type:timestamp;not null"`
	EndTime   time.Time `gorm:"type:timestamp;not null"`
	IsPresent bool      `gorm:"type:bool;not null"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.UUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.UUID `gorm:"type:varchar(255)"`
}
