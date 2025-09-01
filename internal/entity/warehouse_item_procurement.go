package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type WarehouseItemProcurement struct {
	Id                    uint64                            `gorm:"primaryKey;autoIncrement"`
	WarehouseId           uint64                            `gorm:"type:bigint;not null"`
	Warehouse             Warehouse                         `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	ItemId                uint64                            `gorm:"type:bigint;not null"`
	Item                  Item                              `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	SupplierId            uint64                            `gorm:"type:bigint;not null"`
	Supplier              Supplier                          `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE"`
	DailySpending         float64                           `gorm:"type:decimal;not null"`
	DaysNeed              uint64                            `gorm:"type:int;not null"`
	Quantity              float64                           `gorm:"type:decimal;not null"`
	ReceiveQuantity       sql.NullFloat64                   `gorm:"type:decimal;not null;default:0"`
	Note                  string                            `gorm:"type:text"`
	Price                 decimal.Decimal                   `gorm:"type:decimal;not null"`
	TotalPrice            decimal.Decimal                   `gorm:"type:decimal;not null"`
	EstimationArrivalDate time.Time                         `gorm:"type:date;not null"`
	IsArrived             bool                              `gorm:"type:boolean;default:false"`
	TakenAt               sql.NullTime                      `gorm:"type:timestamp"`
	TakenBy               uuid.NullUUID                     `gorm:"type:varchar(255)"`
	Status                enum.ProcurementStatus            `gorm:"type:int;not null"`
	PaymentStatus         enum.PaymentStatus                `gorm:"type:int;not null"`
	Payments              []WarehouseItemProcurementPayment `gorm:"foreignKey:WarehouseItemProcurementId;references:Id"`
	ExpiredAt             sql.NullTime                      `gorm:"type:date"`
	DeadlinePaymentDate   sql.NullTime                      `gorm:"type:date"`
	PaymentType           enum.PaymentType                  `gorm:"paymentType;not null;default:0"`
	CreatedAt             time.Time                         `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy             uuid.NullUUID                     `gorm:"type:varchar(255)"`
	UpdatedAt             time.Time                         `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy             uuid.NullUUID                     `gorm:"type:varchar(255)"`
	CreatedByUser         User                              `gorm:"foreignKey:CreatedBy;references:Id;constraint:OnDelete:SET NULL"`
}
