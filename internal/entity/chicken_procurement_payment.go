package entity

import "time"

type ChickenProcurementPayment struct {
	ID                   uint64    `gorm:"primaryKey;autoIncrement"`
	ChickenProcurementID uint64    `gorm:"not null"`
	TotalPayment         float64   `gorm:"not null"`
	CreatedAt            time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt            time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
