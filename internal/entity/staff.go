package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type Staff struct {
	Id          uuid.UUID       `gorm:"primaryKey;varchar(255);not null"`
	AccountId   uuid.UUID       `gorm:"type:varchar(255);not null"`
	Account     Account         `gorm:"foreignKey:AccountId;references:Id"`
	Name        string          `gorm:"type:varchar(255);not null"`
	PhoneNumber string          `gorm:"type:varchar(15);not null"`
	Address     string          `gorm:"type:text;not null"`
	Salary      decimal.Decimal `gorm:"type:decimal;not null"`
	CreatedBy   uuid.NullUUID   `gorm:"type:varchar(255)"`
	CreatedAt   time.Time       `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy   uuid.NullUUID   `gorm:"type:varchar(255)"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;autoUpdateTime"`
}
