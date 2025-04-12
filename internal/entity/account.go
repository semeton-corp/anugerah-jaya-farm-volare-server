package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id           uuid.UUID `gorm:"type:varchar(255);primary_key"`
	Name         string    `gorm:"type:varchar(255);not null"`
	Email        string    `gorm:"type:varchar(255);unique"`
	Password     string    `gorm:"type:varchar(255);not null"`
	RoleId       uint64    `gorm:"type:bigint;not null"`
	Role         Role      `gorm:"foreignKey:RoleId;references:Id"`
	PhotoProfile string    `gorm:"type:text;default:null"`
	CreatedAt    time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt    time.Time `gorm:"type:timestamp;auto_update_time"`
}
