package entity

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	Id          uint64    `gorm:"primary_key;auto_increment"`
	AccountId   uuid.UUID `gorm:"type:varchar(255);not null"`
	Account     Account   `gorm:"foreignKey:AccountId;references:Id"`
	Name        string    `gorm:"type:varchar(255);not null"`
	PhoneNumber string    `gorm:"type:varchar(15);not null"`
	Description string    `gorm:"type:text;not null"`
	Address     string    `gorm:"type:text;not null"`
	Salary      float64   `gorm:"type:decimal;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt   time.Time `gorm:"type:timestamp;auto_update_time"`
}
