package entity

import (
	"time"
)

type Location struct {
	Id        uint64    `gorm:"primary_key;auto_increment"`
	Name      string    `gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt time.Time `gorm:"type:timestamp;auto_update_time"`
}
