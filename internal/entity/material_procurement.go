package entity

import "time"

type MaterialProcurement struct {
	ID                  uint64    `gorm:"primaryKey;autoIncrement"`
	DateTransaction     string    `gorm:"type:date;not null"`
	ProcurementType     string    `gorm:"type:varchar(255);not null"`
	Detail              string    `gorm:"type:text"`
	Quantity            int       `gorm:"type:int;not null"`
	Unit                string    `gorm:"type:varchar(255);not null"`
	SupplierID          uint64    `gorm:"type:int;not null"`
	Price               float64   `gorm:"type:float;not null"`
	DatePayment         string    `gorm:"type:date;not null"`
	InvoiceUrl          string    `gorm:"type:text;not null"`
	StatusPayment       string    `gorm:"type:varchar(255);not null"`
	PaymentType         string    `gorm:"type:varchar(255);not null"`
	AcceptedBy          string    `gorm:"type:varchar(26);not null"`
	EstimateArrivalDate string    `gorm:"type:date;not null"`
	Status              string    `gorm:"type:varchar(255);not null"`
	CreatedBy           string    `gorm:"type:varchar(26);not null"`
	CreatedAt           time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt           time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
