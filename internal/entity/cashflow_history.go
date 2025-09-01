package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type CashflowHistory struct {
	Id               uint64          `gorm:"primaryKey;autoIncrement"`
	LocationId       uint64          `gorm:"type:bigint;not null"`
	Income           decimal.Decimal `gorm:"type:decimal"`
	Profit           decimal.Decimal `gorm:"type:decimal"`
	Expense          decimal.Decimal `gorm:"type:decimal"`
	Cash             decimal.Decimal `gorm:"type:decimal"`
	Receivables      decimal.Decimal `gorm:"type:decimal"`
	Debt             decimal.Decimal `gorm:"type:decimal"`
	StoreEggSale     decimal.Decimal `gorm:"type:decimal"`
	WarehouseEggSale decimal.Decimal `gorm:"type:decimal"`
	CreatedBy        uuid.NullUUID   `gorm:"type:varchar(255)"`
	CreatedAt        time.Time       `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy        uuid.NullUUID   `gorm:"type:varchar(255)"`
	UpdatedAt        time.Time       `gorm:"type:timestamp;autoUpdateTime"`
}
