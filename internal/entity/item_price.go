package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ItemPrice struct {
	Id        uint64          `gorm:"primaryKey;autoIncrement"`
	Category  string          `gorm:"type:varchar(255);not null"`
	ItemId    uint64          `gorm:"type:bigint;not null"`
	Item      Item            `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Price     decimal.Decimal `gorm:"type:decimal;not null"`
	CreatedAt time.Time       `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.NullUUID   `gorm:"type:varchar(255)"`
	UpdatedAt time.Time       `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.NullUUID   `gorm:"type:varchar(255)"`
}
