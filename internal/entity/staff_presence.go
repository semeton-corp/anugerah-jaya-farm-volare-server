package entity

import (
	"time"

	"github.com/google/uuid"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
)

type StaffPresence struct {
	Id        uint64            `gorm:"primaryKey;autoIncrement"`
	UserId    uuid.UUID         `gorm:"type:bigint;not null"`
	User      User              `gorm:"foreignKey:UserId;references:Id;constraint:OnDelete:CASCADE"`
	StartTime datatype.TimeOnly `gorm:"type:timestamp"`
	EndTime   datatype.TimeOnly `gorm:"type:timestamp"`
	IsPresent bool              `gorm:"type:bool;not null"`
	CreatedAt time.Time         `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.NullUUID     `gorm:"type:varchar(255)"`
	UpdatedAt time.Time         `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.NullUUID     `gorm:"type:varchar(255)"`
}
