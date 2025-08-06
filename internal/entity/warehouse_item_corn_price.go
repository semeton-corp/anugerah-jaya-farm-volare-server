package entity

import "github.com/shopspring/decimal"

type WarehouseItemCornPrice struct {
	Id          uint64          `gorm:"primaryKey;autoIncrement"`
	UpperLimit  float64         `gorm:"decimal;not null"`
	BottomLimit float64         `gorm:"decimal;not null"`
	BasePrice   decimal.Decimal `gorm:"decimal;not null"`
	Discount    float64         `gorm:"decimal;not null"`
}
