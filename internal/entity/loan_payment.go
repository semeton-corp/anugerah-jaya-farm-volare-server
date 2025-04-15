package entity

import "time"

type LoanPayment struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement"`
	LoanId    uint64    `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
