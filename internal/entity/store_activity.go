package entity

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"

type StoreActivity struct {
	Id          uint64              `gorm:"primaryKey;autoIncrement"`
	Description string              `gorm:"type:varchar(255);not null"`
	Status      enum.ActivityStatus `gorm:"type:int;not null"`
	CreatedAt   string              `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy   string              `gorm:"type:varchar(255);not null"`
}
