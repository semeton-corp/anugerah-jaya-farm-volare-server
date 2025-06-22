package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenProcurement struct {
	Id                  uint64        `gorm:"primaryKey;autoIncrement"`
	Quantity            int           `gorm:"type:int;not null"`
	SupplierId          uint64        `gorm:"type:bigint;not null"`
	Supplier            Supplier      `gorm:"foreignKey:SupplierId;refereces:Id"`
	TotalPrice          float64       `gorm:"type:float;not null"`
	StatusPayment       string        `gorm:"type:varchar(255);not null"`
	PaymentType         string        `gorm:"type:varchar(255);not null"`
	DatePayment         string        `gorm:"type:date;not null"`
	InvoiceUrl          string        `gorm:"type:text;not null"`
	TakenBy             string        `gorm:"type:varchar(26);not null"`
	EstimateArrivalDate time.Time     `gorm:"type:date;not null"`
	Status              string        `gorm:"type:varchar(255);not null"`
	CreatedAt           time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy           uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt           time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy           uuid.NullUUID `gorm:"type:varchar(255)"`
}
