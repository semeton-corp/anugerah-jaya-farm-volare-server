package entity

import "time"

type Cage struct {
	ID         uint64    `gorm:"primary_key;auto_increment"`
	LocationID uint64    `gorm:"type:integer;not null"`
	Name       string    `gorm:"type:varchar(255);not null"`
	Capacity   uint64    `gorm:"type:integer;not null"`
	CreatedAt  time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt  time.Time `gorm:"type:timestamp;auto_update_time"`
}
