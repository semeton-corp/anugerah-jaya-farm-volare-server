package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AdditionalWork struct {
	Id                  uint64                    `gorm:"primaryKey;autoIncrement"`
	Description         string                    `gorm:"type:text;not null"`
	Slot                uint64                    `gorm:"type:bigint;not null"`
	AdditionalWorkStaff []AdditionalWorkStaff     `gorm:"foreignKey:AdditionalWorkId;references:Id"`
	Location            enum.LocationAddionalWork `gorm:"int;not null"`
	Salary              decimal.Decimal           `gorm:"decimal;not null;default:0"`
	CreatedBy           uuid.NullUUID             `gorm:"type:varchar(255)"`
	CreatedAt           time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy           uuid.NullUUID             `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt           gorm.DeletedAt            `gorm:"type:timestamp;index"` // soft delete
}
