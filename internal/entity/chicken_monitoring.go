package entity

import (
	"time"

	"github.com/oklog/ulid"
)

type ChickenMonitoring struct {
	ID                       uint64    `gorm:"primary_key;auto_increment"`
	CageID                   uint64    `gorm:"type:integer;not null"`
	TotalLifeChicken         uint64    `gorm:"type:integer;not null"`
	TotalDeathChicken        uint64    `gorm:"type:integer;not null"`
	TotalSickChicken         uint64    `gorm:"type:integer;not null"`
	ChickenSickTypeID        uint64    `gorm:"type:integer;not null"`
	ChickenSickMedicineID    uint64    `gorm:"type:integer;not null"`
	TotalMedicineConsumption float64   `gorm:"type:float;not null"`
	TotalFeedConsumption     float64   `gorm:"type:float;not null"`
	SubmittedBy              ulid.ULID `gorm:"type:varchar(26);not null"`
	CreatedAt                time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt                time.Time `gorm:"type:timestamp;auto_update_time"`
}
