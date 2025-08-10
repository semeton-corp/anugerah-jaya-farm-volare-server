package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type ChickenProcurement struct {
	Id                    uint64                      `gorm:"primaryKey;autoIncrement"`
	CageId                uint64                      `gorm:"type:bigint;not null"`
	Cage                  Cage                        `gorm:"foreignKey:CageId;refereces:Id"`
	SupplierId            uint64                      `gorm:"type:bigint;not null"`
	Supplier              Supplier                    `gorm:"foreignKey:SupplierId;refereces:Id;constraint:OnDelete:SET NULL"`
	Quantity              uint64                      `gorm:"type:bigint;not null"`
	RecieveQuantity       sql.NullInt64               `gorm:"type:bigint"`
	Note                  string                      `gorm:"type:text"`
	Price                 decimal.Decimal             `gorm:"type:decimal;not null"`
	TotalPrice            decimal.Decimal             `gorm:"type:decimal;not null"`
	PaymentStatus         enum.PaymentStatus          `gorm:"type:int;not null"`
	Status                enum.ProcurementStatus      `gorm:"type:int;not null"`
	TakenBy               uuid.NullUUID               `gorm:"type:varchar(255)"`
	TakenAt               sql.NullTime                `gorm:"type:timestamp"`
	PaymentType           enum.PaymentType            `gorm:"paymentType;not null;default:0"`
	IsArrived             bool                        `gorm:"type:boolean;not null;default:false"`
	EstimationArrivalDate time.Time                   `gorm:"type:date;not null"`
	Payments              []ChickenProcurementPayment `gorm:"foreignKey:ChickenProcurementId;references:Id"`
	CreatedAt             time.Time                   `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy             uuid.NullUUID               `gorm:"type:varchar(255)"`
	UpdatedAt             time.Time                   `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy             uuid.NullUUID               `gorm:"type:varchar(255)"`
}
