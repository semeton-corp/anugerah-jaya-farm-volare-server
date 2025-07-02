package entity

import (
	"time"

	"github.com/google/uuid"
)

type ChickenCage struct {
	Id                      uint64                    `gorm:"primaryKey;autoIncrement"`
	CageId                  uint64                    `gorm:"bigint;not null"`
	Cage                    Cage                      `gorm:"foreingKey:CageId;references:Id"`
	ChickenProcurementId    uint64                    `gorm:"bigint;default:null"`
	ChickenProcurement      ChickenProcurement        `gorm:"foreignKey:ChickenProcurementId;references:Id"`
	ChickenMonitoring       []ChickenMonitoring       `gorm:"foreignKey:ChickenCageId;references:Id"`
	ChickenHealthMonitoring []ChickenHealthMonitoring `gorm:"foreignKey:ChickenCageId;references:Id"`
	TotalChicken            uint64                    `gorm:"int;not null;default:0"`
	TotalDeathChicken       uint64                    `gorm:"int;not null;default:0"`
	IsNeedRoutineVaccine    bool                      `gorm:"bool;not null;default:false"`
	CreatedAt               time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy               uuid.NullUUID             `gorm:"type:varchar(255)"`
	UpdatedAt               time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy               uuid.NullUUID             `gorm:"type:varchar(255)"`
}
