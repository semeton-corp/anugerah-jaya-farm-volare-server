package entity

import "time"

type Supplier struct {
	ID          uint64    `gorm:"primaryKey;autoIncrement"`
	Name        string    `gorm:"type:varchar(255);not null"`
	PhoneNumber string    `gorm:"type:varchar(15);not null"`
	Address     string    `gorm:"type:text;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
