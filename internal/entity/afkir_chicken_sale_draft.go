package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type AfkirChickenSaleDraft struct {
	Id                     uint64               `gorm:"primaryKey;autoIncrement"`
	AfkirChickenCustomerId uint64               `gorm:"type:bigint"`
	AfkirChickenCustomer   AfkirChickenCustomer `gorm:"foreignKey:AfkirChickenCustomerId;references:Id;constraint:OnDelete:CASCADE"`
	ChickenCageId          uint64               `gorm:"type:bigint;not null"`
	ChickenCage            ChickenCage          `gorm:"foreignKey:ChickenCageId;references:Id;constraint:OnDelete:CASCADE"`
	TotalSellChicken       uint64               `gorm:"type:bigint;not null"`
	PricePerChicken        decimal.Decimal      `gorm:"type:decimal;not null"`
	TotalPrice             decimal.Decimal      `gorm:"type:decimal;not null"`
	CreatedAt              time.Time            `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy              uuid.NullUUID        `gorm:"type:varchar(255)"`
	UpdatedAt              time.Time            `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy              uuid.NullUUID        `gorm:"type:varchar(255)"`
}
