package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type StoreRequestItem struct {
	Id                   uint64                 `gorm:"primaryKey;autoIncrement"`
	WarehouseId          uint64                 `gorm:"type:bigint;not null"`
	Warehouse            Warehouse              `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	ItemId               uint64                 `gorm:"type:bigint;not null"`
	Item                 Item                   `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	StoreId              sql.NullInt64          `gorm:"type:bigint"`
	Store                Store                  `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity             float64                `gorm:"type:bigint;not null"`
	RecieveQuantity      float64                `gorm:"type:bigint;not null;default:0"`
	WarehouseFulfillment float64                `gorm:"type:bigint;not null;default:-1"`
	Status               enum.RequestItemStatus `gorm:"type:int;not null"`
	WarehouseNote        string                 `gorm:"type:text"`
	StoreNote            string                 `gorm:"type:text"`
	IsSorted             bool                   `gorm:"type:boolean;default:false;not null"`
	RecieveDate          sql.NullTime           `gorm:"type:timestamp"`
	CreatedAt            time.Time              `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy            uuid.NullUUID          `gorm:"type:varchar(255)"`
	UpdatedAt            time.Time              `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy            uuid.NullUUID          `gorm:"type:varchar(255)"`
}
