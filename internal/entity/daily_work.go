package entity

import (
	"time"

	"github.com/google/uuid"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
)

type DailyWork struct {
	Id          uint64            `gorm:"primaryKey;autoIncrement"`
	Description string            `gorm:"type:text;not null"`
	RoleId      uint64            `gorm:"type:bigint;not null"`
	Role        Role              `gorm:"foreignKey:RoleId;references:Id"`
	StartTime   datatype.TimeOnly `gorm:"not null"`
	EndTime     datatype.TimeOnly `gorm:"not null"`
	CreatedBy   uuid.UUID         `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time         `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy   uuid.UUID         `gorm:"type:varchar(255);not null"`
	UpdatedAt   time.Time         `gorm:"type:timestamp;autoUpdateTime"`
}

type DailyWorkSummary struct {
	RoleID     uint64 `gorm:"column:role_id"`
	RoleName   string `gorm:"column:role_name"`
	TotalWork  uint64 `gorm:"column:total_work"`
	TotalStaff uint64 `gorm:"column:total_staff"`
}
