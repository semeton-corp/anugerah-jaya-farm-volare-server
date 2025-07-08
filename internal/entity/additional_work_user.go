package entity

import (
	"time"

	"github.com/google/uuid"
)

type AdditionalWorkUser struct {
	Id               uint64         `gorm:"primaryKey;autoIncrement"`
	UserId           uuid.UUID      `gorm:"type:varchar(255);not null"`
	User             User           `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	AdditionalWorkId uint64         `gorm:"type:bigint;not null"`
	AdditionalWork   AdditionalWork `gorm:"foreignKey:AdditionalWorkId;references:Id;constraint:OnDelete:CASCADE"`
	IsDone           bool           `gorm:"type:boolean;default:false"`
	CreatedAt        time.Time      `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy        uuid.NullUUID  `gorm:"type:varchar(255)"`
	UpdatedAt        time.Time      `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy        uuid.NullUUID  `gorm:"type:varchar(255)"`
}
