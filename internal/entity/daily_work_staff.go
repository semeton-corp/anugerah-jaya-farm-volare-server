package entity

import (
	"time"

	"github.com/google/uuid"
)

type DailyWorkStaff struct {
	Id          uint64        `gorm:"primaryKey;autoIncrement"`
	DailyWorkId uint64        `gorm:"type:biginteger;not null"`
	DailyWork   DailyWork     `gorm:"foreignKey:DailyWorkId;references:Id"`
	UserId      uuid.UUID     `gorm:"type:varchar(255);not null"`
	User        User          `gorm:"foreignKey:UserId;references:Id"`
	IsDone      bool          `gorm:"type:bool;not null"`
	CreatedAt   time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy   uuid.NullUUID `gorm:"type:varchar(255);not null"`
	UpdatedAt   time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy   uuid.NullUUID `gorm:"type:varchar(255);not null"`
}
