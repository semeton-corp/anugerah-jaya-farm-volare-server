package entity

import (
	"time"
)

type Role struct {
	Id        uint64    `gorm:"primaryKey;autoIncrement"`
	Name      string    `gorm:"type:varchar(255);unique"`
	CreatedAt time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
