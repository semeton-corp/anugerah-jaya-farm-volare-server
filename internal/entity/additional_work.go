package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type AdditionalWork struct {
	Id                  uint64                    `gorm:"primaryKey;autoIncrement"`
	Description         string                    `gorm:"type:text;not null"`
	Slot                uint64                    `gorm:"type:bigint;not null"`
	AdditionalWorkStaff []AdditionalWorkStaff     `gorm:"foreignKey:AdditionalWorkId;references:Id"`
	Location            enum.LocationAddionalWork `gorm:"int:text;not null"`
	CreatedBy           uuid.UUID                 `gorm:"type:varchar(255)"`
	CreatedAt           time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy           uuid.UUID                 `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
}
