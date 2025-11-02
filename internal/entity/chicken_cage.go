package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type ChickenCage struct {
	Id                             uint64                    `gorm:"primaryKey;autoIncrement"`
	CageId                         uint64                    `gorm:"bigint;not null"`
	Cage                           Cage                      `gorm:"foreignKey:CageId;references:Id;constraint:OnDelete:CASCADE"`
	ChickenProcurementId           sql.NullInt64             `gorm:"bigint;default:null"`
	ChickenProcurement             ChickenProcurement        `gorm:"foreignKey:ChickenProcurementId;references:Id;constraint:OnDelete:CASCADE"`
	ChickenMonitoring              []ChickenMonitoring       `gorm:"foreignKey:ChickenCageId;references:Id"`
	ChickenHealthMonitoring        []ChickenHealthMonitoring `gorm:"foreignKey:ChickenCageId;references:Id"`
	TotalChicken                   uint64                    `gorm:"type:int;not null;default:0"`
	LatestChickenAgeVaccineRoutine sql.NullInt64             `gorm:"type:int"`
	IsNeedRoutineVaccine           bool                      `gorm:"bool;not null;default:false"`
	IsNeedFeed                     bool                      `gorm:"bool;not null;default:true"`
	CreatedAt                      time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy                      uuid.NullUUID             `gorm:"type:varchar(255)"`
	UpdatedAt                      time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy                      uuid.NullUUID             `gorm:"type:varchar(255)"`
}
