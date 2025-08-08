package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WarehouseItemProcurementDraft struct {
	Id            uint64          `gorm:"primaryKey;autoIncrement"`
	WarehouseId   uint64          `gorm:"type:bigint;not null"`
	Warehouse     Warehouse       `gorm:"foreignKey:WarehouseId;refereces:Id"`
	ItemId        uint64          `gorm:"type:bigint;not null"`
	Item          Item            `gorm:"foreignKey:ItemId;references:Id"`
	SupplierId    sql.NullInt64   `gorm:"type:bigint"`
	Supplier      Supplier        `gorm:"foreignKey:SupplierId;references:Id"`
	DailySpending float64         `gorm:"type:decimal;not null"`
	DaysNeed      uint64          `gorm:"type:int;not null"`
	Price         decimal.Decimal `gorm:"type:decimal;not null"`
	CreatedAt     time.Time       `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy     uuid.NullUUID   `gorm:"type:varchar(255)"`
	UpdatedAt     time.Time       `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy     uuid.NullUUID   `gorm:"type:varchar(255)"`
}
