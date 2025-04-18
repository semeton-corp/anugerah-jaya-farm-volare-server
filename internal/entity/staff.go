package entity

import (
	"time"

	"github.com/google/uuid"
)

type Staff struct {
	AccountId   uuid.UUID `gorm:"primaryKey;type:varchar(255);not null"`
	Account     Account   `gorm:"foreignKey:AccountId;references:Id;constraint:on_delete:cascade"`
	Name        string    `gorm:"type:varchar(255);not null"`
	PhoneNumber string    `gorm:"type:varchar(15);not null"`
	Description string    `gorm:"type:text;not null"`
	Address     string    `gorm:"type:text;not null"`
	Salary      float64   `gorm:"type:decimal;not null"`
	CreatedBy   uuid.UUID `gorm:"type:varchar(255);not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy   uuid.UUID `gorm:"type:varchar(255);not null"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
