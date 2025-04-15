package entity

import "time"

type MaterialProcurementPayment struct {
	ID                    uint64    `gorm:"primaryKey;autoIncrement"`
	MaterialProcurementID uint64    `gorm:"not null"`
	TotalPayment          float64   `gorm:"not null"`
	CreatedAt             time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt             time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
