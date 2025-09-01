package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type StoreSale struct {
	Id                  uint64             `gorm:"primaryKey;autoIncrement;not null"`
	CustomerId          uint64             `gorm:"type:varchar(255);not null"`
	Customer            Customer           `gorm:"foreignKey:CustomerId;references:Id"`
	ItemId              uint64             `gorm:"type:bigint;not null"`
	Item                Item               `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	StoreId             uint64             `gorm:"type:bigint;not null"`
	Store               Store              `gorm:"foreignKey:StoreId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity            float64            `gorm:"type:bigint;not null"`
	SaleUnit            enum.SaleUnit      `gorm:"type:int;not null"`
	Price               decimal.Decimal    `gorm:"type:decimal;not null"`
	TotalPrice          decimal.Decimal    `gorm:"type:decimal;not null"`
	Discount            float64            `gorm:"type:decimal;not null"`
	SendDate            time.Time          `gorm:"type:date;not null"`
	PaymentType         enum.PaymentType   `gorm:"type:int;not null"`
	PaymentStatus       enum.PaymentStatus `gorm:"type:int;not null"`
	IsSend              bool               `gorm:"type:boolean;not null"`
	Payments            []StoreSalePayment `gorm:"foreignKey:StoreSaleId;references:Id;constraint:OnDelete:CASCADE"`
	DeadlinePaymentDate sql.NullTime       `gorm:"date"`
	CreatedAt           time.Time          `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           uuid.NullUUID      `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time          `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy           uuid.NullUUID      `gorm:"type:varchar(255)"`
	CreatedByUser       User               `gorm:"foreignKey:CreatedBy;references:Id;constraint:OnDelete:SET NULL"`
}
