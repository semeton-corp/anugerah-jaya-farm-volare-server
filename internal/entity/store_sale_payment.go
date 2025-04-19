package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type StoreSalePayment struct {
	Id           uint64          `gorm:"primaryKey;autoIncrement;not null"`
	StoreSaleId  uint64          `gorm:"type:bigint;not null"`
	PaymentDate  time.Time       `gorm:"type:date;not null"`
	Nominal      decimal.Decimal `gorm:"type:decimal;not null"`
	PaymentProof string          `gorm:"type:text;not null"`
	CreatedAt    time.Time       `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy    uuid.UUID       `gorm:"type:varchar(255);not null"`
	UpdatedAt    time.Time       `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy    uuid.UUID       `gorm:"type:varchar(255);not null"`
}
