package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type StoreSale struct {
	Id              uint64             `gorm:"primaryKey;autoIncrement;not null"`
	Customer        string             `gorm:"type:varchar(255);not null"`
	Phone           string             `gorm:"type:varchar(255);not null"`
	WarehouseItemId uint64             `gorm:"type:bigint;not null"`
	WarehouseItem   WarehouseItem      `gorm:"foreignKey:WarehouseItemId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	StoreId         uint64             `gorm:"type:bigint;not null"`
	Store           Store              `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	Quantity        uint64             `gorm:"type:bigint;not null"`
	Price           decimal.Decimal    `gorm:"type:decimal;not null"`
	TotalPrice      decimal.Decimal    `gorm:"type:decimal;not null"`
	SendDate        time.Time          `gorm:"type:date;not null"`
	PaymentMethod   enum.PaymentMethod `gorm:"type:int;not null"`
	PaymentStatus   enum.PaymentStatus `gorm:"type:int;not null"`
	IsSend          bool               `gorm:"type:boolean;not null"`
	Payments        []StoreSalePayment `gorm:"foreignKey:StoreSaleId;references:Id;constraint:OnDelete:CASCADE;onUpdate:CASCADE"`
	CreatedAt       time.Time          `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.UUID          `gorm:"type:varchar(255);not null"`
	UpdatedAt       time.Time          `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.UUID          `gorm:"type:varchar(255);not null"`
}
