package entity

import (
	"time"

	"github.com/google/uuid"
)

type Location struct {
	Id        uint64        `gorm:"primaryKey;autoIncrement"`
	Name      string        `gorm:"type:varchar(255);not null;unique"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
}
