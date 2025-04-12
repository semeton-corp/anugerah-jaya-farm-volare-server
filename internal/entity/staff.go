package entity

import "time"

type Staff struct {
	Id          uint64    `gorm:"primary_key;auto_increment"`
	Account     Account   `gorm:"foreignKey:AccountId"`
	Name        string    `gorm:"type:varchar(255);not null"`
	PhoneNumber string    `gorm:"type:varchar(15);not null"`
	Description string    `gorm:"type:text;not null"`
	Address     string    `gorm:"type:text;not null"`
	Salary      float64   `gorm:"type:double;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt   time.Time `gorm:"type:timestamp;auto_update_time"`
}
