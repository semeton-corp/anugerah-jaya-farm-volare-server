package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type DailyWorkUser struct {
	Id          uint64        `gorm:"primaryKey;autoIncrement"`
	DailyWorkId uint64        `gorm:"type:biginteger;not null"`
	DailyWork   DailyWork     `gorm:"foreignKey:DailyWorkId;references:Id;constraint:OnDelete:CASCADE"`
	UserId      uuid.UUID     `gorm:"type:varchar(255);not null"`
	User        User          `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	IsDone      bool          `gorm:"type:boolean;not null"`
	Note        string        `gorm:"type:text"`
	FinishedAt  sql.NullTime  `gorm:"type:timestamp"`
	CreatedAt   time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy   uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt   time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy   uuid.NullUUID `gorm:"type:varchar(255)"`
}

// Todo : change constraint in created by user
