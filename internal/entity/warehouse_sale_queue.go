package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type WarehouseSaleQueue struct {
	Id                  uint64            `gorm:"primaryKey;autoIncrement"`
	CustomerId          sql.NullInt64     `gorm:"type:bigint"`
	Customer            Customer          `gorm:"foreignKey:CustomerId;references:Id"`
	CustomerName        sql.NullString    `gorm:"type:varchar(255)"`
	CustomerPhoneNumber sql.NullString    `gorm:"type:varchar(255)"`
	CustomerType        enum.CustomerType `gorm:"type:varchar(255);not null"`
	ItemId              uint64            `gorm:"type:bigint;not null"`
	Item                Item              `gorm:"foreignKey:ItemId;references:Id"`
	WarehouseId         uint64            `gorm:"type:bigint;not null"`
	Warehouse           Warehouse         `gorm:"foreignKey:WarehouseId;references:Id"`
	SaleUnit            enum.SaleUnit     `gorm:"type:varchar(255);not null"`
	CreatedAt           time.Time         `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           uuid.NullUUID     `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time         `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy           uuid.NullUUID     `gorm:"type:varchar(255)"`
}
