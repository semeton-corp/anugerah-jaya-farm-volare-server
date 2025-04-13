package entity

import "time"

type ChickenVaccineMonitoring struct {
	Id                  uint64    `gorm:"primary_key;auto_increment"`
	ChickenMonitoringId uint64    `gorm:"type:bigint;not null"`
	Vaccine             string    `gorm:"type:varchar(255);not null"`
	Dose                uint64    `gorm:"type:integer;not null"`
	Unit                string    `gorm:"type:varchar(255);not null"`
	CreatedAt           time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt           time.Time `gorm:"type:timestamp;auto_update_time"`
}
