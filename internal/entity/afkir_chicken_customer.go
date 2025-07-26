package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AfkirChickenCustomer struct {
	Id                uint64             `gorm:"primaryKey;autoIncrement"`
	Name              string             `gorm:"type:varchar(255);not null"`
	PhoneNumber       string             `gorm:"type:varchar(255);not null"`
	Address           string             `gorm:"type:text;not null"`
	LatestPrice       decimal.Decimal    `gorm:"type:decimal;not null"`
	AfkirChickenSales []AfkirChickenSale `gorm:"foreignKey:AfkirChickenCustomerId;references:Id"`
	CreatedAt         time.Time          `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy         uuid.NullUUID      `gorm:"type:varchar(255)"`
	UpdatedAt         time.Time          `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy         uuid.NullUUID      `gorm:"type:varchar(255)"`
}
