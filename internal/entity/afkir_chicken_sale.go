package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type AfkirChickenSale struct {
	Id                     uint64                    `gorm:"primaryKey;autoIncrement"`
	AfkirChickenCustomerId uint64                    `gorm:"type:bigint;not null"`
	AfkirChickenCustomer   AfkirChickenCustomer      `gorm:"foreignKey:AfkirChickenCustomerId;references:Id;constraint:OnDelete:CASCADE"`
	ChickenCageId          uint64                    `gorm:"type:bigint;not null"`
	ChickenCage            ChickenCage               `gorm:"foreignKey:ChickenCageId;references:Id;constraint:OnDelete:CASCADE"`
	TotalSellChicken       uint64                    `gorm:"type:bigint;not null"`
	PricePerChicken        decimal.Decimal           `gorm:"decimal;not null"`
	TotalPrice             decimal.Decimal           `gorm:"decimal;not null"`
	ChickenAge             uint64                    `gorm:"type:int;not null"`
	Payments               []AfkirChickenSalePayment `gorm:"foreignKey:AfkirChickenSaleId;references:Id"`
	PaymentType            enum.PaymentType          `gorm:"type:int;not null"`
	PaymentStatus          enum.PaymentStatus        `gorm:"type:int;not null"`
	CreatedAt              time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy              uuid.NullUUID             `gorm:"type:varchar(255)"`
	UpdatedAt              time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy              uuid.NullUUID             `gorm:"type:varchar(255)"`
}
