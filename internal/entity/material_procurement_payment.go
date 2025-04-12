package entity

import "time"

type MaterialProcurementPayment struct {
	ID                    uint64    `gorm:"primary_key;auto_increment"`
	MaterialProcurementID uint64    `gorm:"not null"`
	TotalPayment          float64   `gorm:"not null"`
	CreatedAt             time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt             time.Time `gorm:"type:timestamp;auto_update_time"`
}
