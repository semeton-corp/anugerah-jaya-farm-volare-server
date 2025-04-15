package entity

import (
	"time"

	"github.com/shopspring/decimal"
)

type AdditionalWork struct {
	Id          uint64          `gorm:"primaryKey;autoIncrement"`
	Name        string          `gorm:"type:varchar(255);not null"`
	StaffId     uint            `gorm:"type:varchar(255);not null"`
	StartTime   string          `gorm:"type:timestamp;not null"`
	EndTime     string          `gorm:"type:timestamp;not null"`
	Salary      decimal.Decimal `gorm:"type:decimal;not null"`
	Description string          `gorm:"type:text;not null"`
	CreatedAt   time.Time       `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time       `gorm:"type:timestamp;autoUpdateTime"`
}
