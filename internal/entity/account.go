package entity

import (
	"time"

	"github.com/google/uuid"
)

type Account struct {
	Id           uuid.UUID     `gorm:"type:varchar(255);primaryKey"`
	Username     string        `gorm:"type:varchar(255);not null;unique"`
	Email        string        `gorm:"type:varchar(255);unique"`
	Password     string        `gorm:"type:varchar(255);not null"`
	RoleId       uint64        `gorm:"type:bigint;not null"`
	Role         Role          `gorm:"foreignKey:RoleId;references:Id;constraint:on_delete:set_null"`
	PhotoProfile string        `gorm:"type:text;default:null"`
	CreatedBy    uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt    time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy    uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt    time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
