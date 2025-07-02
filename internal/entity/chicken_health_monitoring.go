package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type ChickenHealthMonitoring struct {
	Id             uint64                     `gorm:"primaryKey;autoIncrement"`
	ChickenCageId  uint64                     `gorm:"type:bigint;not null"`
	HealthItemName string                     `gorm:"varchar(255);not null"`
	Type           enum.ChickenHealthItemType `gorm:"int;not null"`
	Dose           float64                    `gorm:"type:decimal;not null"`
	Unit           string                     `gorm:"type:varchar(255);not null"`
	Disease        sql.NullString             `gorm:"type:varchar(255)"`
	CreatedBy      uuid.NullUUID              `gorm:"type:varchar(255)"`
	CreatedAt      time.Time                  `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy      uuid.NullUUID              `gorm:"type:varchar(255)"`
	UpdatedAt      time.Time                  `gorm:"type:timestamp;autoUpdateTime"`
}
