package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

// Note : created every month
type UserSalaryPayment struct {
	Id                   uint64             `gorm:"primaryKey;autoIncrement"`
	UserId               uuid.UUID          `gorm:"type:varchar(255);not null"`
	User                 User               `gorm:"foreignKey:UserId;referennces:Id"`
	BaseSalary           decimal.Decimal    `gorm:"type:decimal;not null"`
	BonusSalary          decimal.Decimal    `gorm:"type:decimal;not null"`
	CompentationSalary   decimal.Decimal    `gorm:"type:decimal;not null"`
	AdditionalWorkSalary decimal.Decimal    `gorm:"type:decimal;not null"`
	PaymentProof         string             `gorm:"type:text;not null"`
	PaymentMethod        enum.PaymentMethod `gorm:"type:int;not null;default:0"`
	IsPaid               bool               `gorm:"type:bool;not null;default:false"`
	CreatedAt            time.Time          `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy            uuid.NullUUID      `gorm:"type:varchar(255)"`
	UpdatedAt            time.Time          `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy            uuid.NullUUID      `gorm:"type:varchar(255)"`
	CreatedByUser        User               `gorm:"foreignKey:CreatedBy;references:Id"`
}
