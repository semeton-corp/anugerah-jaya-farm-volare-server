package entity

import "time"

type LoanPayment struct {
	Id        uint64    `gorm:"primary_key;auto_increment"`
	LoanId    uint64    `gorm:"not null"`
	Amount    float64   `gorm:"not null"`
	CreatedAt time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt time.Time `gorm:"type:timestamp;auto_update_time"`
}
