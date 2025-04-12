package entity

import "time"

type ChickenProcurement struct {
	ID                  uint64    `gorm:"primary_key;auto_increment"`
	ChickenType         string    `gorm:"type:varchar(255);not null"`
	Age                 int       `gorm:"type:int;not null"`
	Quantity            int       `gorm:"type:int;not null"`
	SupplierID          uint64    `gorm:"type:bigint;not null"`
	CreatedBy           string    `gorm:"type:varchar(26);not null"`
	TotalPrice          float64   `gorm:"type:float;not null"`
	StatusPayment       string    `gorm:"type:varchar(255);not null"`
	PaymentType         string    `gorm:"type:varchar(255);not null"`
	DatePayment         string    `gorm:"type:date;not null"`
	InvoiceUrl          string    `gorm:"type:text;not null"`
	AcceptedBy          string    `gorm:"type:varchar(26);not null"`
	EstimateArrivalDate string    `gorm:"type:date;not null"`
	Status              string    `gorm:"type:varchar(255);not null"`
	CreatedAt           time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt           time.Time `gorm:"type:timestamp;auto_update_time"`
}
