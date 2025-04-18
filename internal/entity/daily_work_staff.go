package entity

import (
	"time"

	"github.com/google/uuid"
)

type DailyWorkStaff struct {
	Id          uint64    `gorm:"primaryKey;autoIncrement"`
	DailyWorkId uint64    `gorm:"type:biginteger;not null"`
	DailyWork   DailyWork `gorm:"foreignKey:DailyWorkId;references:Id"`
	StaffId     uuid.UUID `gorm:"type:biginteger;not null"`
	Staff       Staff     `gorm:"foreignKey:StaffId;references:AccountId"`
	IsDone      bool      `gorm:"type:bool;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
