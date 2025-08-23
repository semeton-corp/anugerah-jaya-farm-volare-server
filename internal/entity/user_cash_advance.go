package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type UserCashAdvance struct {
	Id                  uint64                   `gorm:"primaryKey;autoIncrement"`
	UserId              uuid.UUID                `gorm:"type:varchar(255);not null"`
	User                User                     `gorm:"foreignKey:UserId;referennces:Id"`
	Nominal             decimal.Decimal          `gorm:"type:decimal;not null"`
	DeadlinePaymentDate time.Time                `gorm:"type:date;not null"`
	PaymentStatus       enum.PaymentStatus       `gorm:"type:int;not null"`
	Payments            []UserCashAdvancePayment `gorm:"foreignKey:UserCashAdvanceId;references:Id"`
	CreatedAt           time.Time                `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           uuid.NullUUID            `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time                `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy           uuid.NullUUID            `gorm:"type:varchar(255)"`
	CreatedByUser       User                     `gorm:"foreignKey:CreatedBy;references:Id;-:migration"`
}
