package entity

import "time"

type ChickenSickMedicine struct {
	Id          uint64    `gorm:"primary_key;auto_increment"`
	Name        string    `gorm:"type:varchar(255);not null"`
	Description string    `gorm:"type:varchar(255);not null"`
	Dose        uint64    `gorm:"type:integer;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt   time.Time `gorm:"type:timestamp;auto_update_time"`
}
