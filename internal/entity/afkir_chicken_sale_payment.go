package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type AfkirChickenSalePayment struct {
	Id                 uint64             `gorm:"primaryKey;autoIncrement"`
	AfkirChickenSaleId uint64             `gorm:"type:bigint"`
	AfkirChickenSale   AfkirChickenSale   `gorm:"foreignKey:AfkirChickenSaleId;references:Id"`
	PaymentDate        time.Time          `gorm:"type:date;not null"`
	Nominal            decimal.Decimal    `gorm:"decimal;not null"`
	PaymentProof       string             `gorm:"type:text;not null"`
	PaymentMethod      enum.PaymentMethod `gorm:"type:int;not null"`
	CreatedAt          time.Time          `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy          uuid.NullUUID      `gorm:"type:varchar(255)"`
	UpdatedAt          time.Time          `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy          uuid.NullUUID      `gorm:"type:varchar(255)"`
	CreatedByUser      User               `gorm:"foreignKey:CreatedBy;references:Id;constraint:OnDelete:SET NULL"`
}
