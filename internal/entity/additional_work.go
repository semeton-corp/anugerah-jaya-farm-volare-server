package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type AdditionalWork struct {
	Id                  uint64                `gorm:"primaryKey;autoIncrement"`
	Name                string                `gorm:"type:varchar(255);not null"`
	LocationId          uint64                `gorm:"type:bigint;not null"`
	Location            Location              `gorm:"foreignKey:LocationId;references:Id"`
	WarehouseId         sql.NullInt64         `gorm:"type:bigint"`
	Warehouse           Warehouse             `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	StoreId             sql.NullInt64         `gorm:"type:bigint"`
	Store               Store                 `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE"`
	CageId              sql.NullInt64         `gorm:"type:bigint"`
	Cage                Cage                  `gorm:"foreignKey:CageId;refereces:Id;constraint:OnDelete:CASCADE"`
	Description         string                `gorm:"type:text;not null"`
	Slot                uint64                `gorm:"type:bigint;not null"`
	WorkDate            time.Time             `gorm:"type:timestamp;not null"`
	Salary              decimal.Decimal       `gorm:"decimal;not null;default:0"`
	LocationType        enum.LocationWorkType `gorm:"int;not null"`
	AdditionalWorkUsers []AdditionalWorkUser  `gorm:"foreignKey:AdditionalWorkId;references:Id;constraint:OnDelete:CASCADE"`
	CreatedBy           uuid.NullUUID         `gorm:"type:varchar(255)"`
	CreatedAt           time.Time             `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy           uuid.NullUUID         `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time             `gorm:"type:timestamp;autoUpdateTime"`
	DeletedAt           gorm.DeletedAt        `gorm:"type:timestamp;index"` // soft delete
}
