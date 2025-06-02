package entity

import (
	"time"

	"github.com/google/uuid"
)

type AdditionalWorkStaff struct {
	Id               uint64         `gorm:"primaryKey;autoIncrement"`
	StaffId          uuid.UUID      `gorm:"type:varchar(255);not null"`
	Staff            Staff          `gorm:"foreignKey:StaffId;references:Id;constraint:OnDelete:CASCADE"`
	AdditionalWorkId uint64         `gorm:"type:bigint;not null"`
	AdditionalWork   AdditionalWork `gorm:"foreignKey:AdditionalWorkId;references:Id;constraint:OnDelete:CASCADE"`
	IsDone           bool           `gorm:"type:boolean;default:false"`
	CreatedAt        time.Time      `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy        uuid.NullUUID  `gorm:"type:varchar(255)"`
	UpdatedAt        time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy        uuid.NullUUID  `gorm:"type:varchar(255)"`
}
