package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

// Note : why not the id saved? because when the item changed or updated i think it will make confussion, that's why i just saved the name instead of the id
type WarehouseItemHistory struct {
	Id             uint64                 `gorm:"primaryKey;autoIncrement"`
	ItemName       string                 `gorm:"type:varchar(255)"`
	ItemUnit       string                 `gorm:"type:varchar(255)"`
	Source         string                 `gorm:"type:varchar(255)"`
	Destination    string                 `gorm:"type:varchar(255)"`
	QuantityBefore float64                `gorm:"type:decimal;not null"`
	QuantityAfter  float64                `gorm:"type:decimal;not null"`
	Status         enum.ItemHistoryStatus `gorm:"type:int;not null"`
	UserId         uuid.UUID              `gorm:"type:varchar(255);constraint:OnDelete:SET NULL"`
	User           User                   `gorm:"foreignKey:UserId;references:Id" json:"-"`
	CreatedAt      time.Time              `gorm:"type:timestamp;autoCreateTime"`
}
