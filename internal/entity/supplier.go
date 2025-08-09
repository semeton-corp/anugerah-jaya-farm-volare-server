package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type Supplier struct {
	Id            uint64            `gorm:"primaryKey;autoIncrement"`
	SupplierItems []SupplierItem    `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE;OnUpdate:CASCADE"`
	Name          string            `gorm:"type:varchar(255);not null"`
	PhoneNumber   string            `gorm:"type:varchar(255);not null;unique"`
	Address       string            `gorm:"type:text;not null"`
	SupplierType  enum.SupplierType `gorm:"type:int;not null;default:0"`
	CreatedBy     uuid.NullUUID     `gorm:"type:varchar(255)"`
	CreatedAt     time.Time         `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy     uuid.NullUUID     `gorm:"type:varchar(255)"`
	UpdatedAt     time.Time         `gorm:"type:timestamp;autoUpdateTime"`
}
