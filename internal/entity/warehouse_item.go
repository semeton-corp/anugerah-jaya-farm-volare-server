package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type WarehouseItem struct {
	Id        uint64                     `gorm:"primaryKey;autoIncrement"`
	Name      string                     `gorm:"type:varchar(255);not null"` // Telur Ok, Jagung
	Category  enum.WarehouseItemCategory `gorm:"type:int;not null"`          // Telur, Pakan
	Unit      string                     `gorm:"type:varchar(255);not null"` // Kg, Ltr
	CreatedAt time.Time                  `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy uuid.UUID                  `gorm:"type:varchar(255);not null"`
	UpdatedAt time.Time                  `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy uuid.UUID                  `gorm:"type:varchar(255);not null"`
}
