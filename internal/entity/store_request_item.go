package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

// Note : Quantity field is in KG
type StoreRequestItem struct {
	Id                   uint64                 `gorm:"primaryKey;autoIncrement"`
	WarehouseId          uint64                 `gorm:"type:bigint;not null"`
	Warehouse            Warehouse              `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	ItemId               uint64                 `gorm:"type:bigint;not null"`
	Item                 Item                   `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	StoreId              sql.NullInt64          `gorm:"type:bigint"`
	Store                Store                  `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity             float64                `gorm:"type:decimal;not null"`
	ReceiveQuantity      sql.NullFloat64        `gorm:"type:decimal"`
	WarehouseFulfillment sql.NullFloat64        `gorm:"type:decimal"`
	Status               enum.RequestItemStatus `gorm:"type:int;not null"`
	WarehouseNote        string                 `gorm:"type:text"`
	StoreNote            string                 `gorm:"type:text"`
	IsSorted             bool                   `gorm:"type:boolean;default:false;not null"`
	ReceiveDate          sql.NullTime           `gorm:"type:timestamp"`
	CreatedAt            time.Time              `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy            uuid.NullUUID          `gorm:"type:varchar(255)"`
	UpdatedAt            time.Time              `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy            uuid.NullUUID          `gorm:"type:varchar(255)"`
	CreatedByUser        User                   `gorm:"foreignKey:CreatedBy;references:Id;constraint:OnDelete:SET NULL"`
}
