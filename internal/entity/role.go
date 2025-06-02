package entity

import (
	"time"

	"github.com/google/uuid"
)

type Role struct {
	Id        uint64        `gorm:"primaryKey;autoIncrement"`
	Name      string        `gorm:"type:varchar(255);unique"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.NullUUID `gorm:"varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.NullUUID `gorm:"varchar(255)"`
}
