package entity

import (
	"time"

	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type StoreActivity struct {
	Id          uint64              `gorm:"primaryKey;autoIncrement"`
	Description string              `gorm:"type:varchar(255);not null"`
	Status      enum.ActivityStatus `gorm:"type:int;not null"`
	CreatedAt   time.Time           `gorm:"type:timestamp;autoCreateTime"`
}
