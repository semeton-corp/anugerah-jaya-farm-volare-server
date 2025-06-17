package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenMedicine struct {
	Id        uint64        `gorm:"primaryKey;autoIncrement"`
	Name      string        `gorm:"type:varchar(255);not null"`
	CreatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	CreatedAt time.Time     `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy uuid.NullUUID `gorm:"type:varchar(255)"`
	UpdatedAt time.Time     `gorm:"type:timestamp;autoUpdateTime"`
}
