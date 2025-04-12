package entity

import (
	"time"

	"github.com/oklog/ulid"
)

type EggMonitoring struct {
	ID              uint64    `gorm:"primary_key;auto_increment"`
	CageID          uint64    `gorm:"type:integer;not null"`
	Date            time.Time `gorm:"type:date;not null"`
	TotalCrackedEgg uint64    `gorm:"type:integer;not null"`
	TotalGoodEgg    uint64    `gorm:"type:integer;not null"`
	TotalBrokeEgg   uint64    `gorm:"type:integer;not null"`
	TotalRejectEgg  uint64    `gorm:"type:integer;not null"`
	SubmittedBy     ulid.ULID `gorm:"type:varchar(26);not null"`
	CreatedAt       time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt       time.Time `gorm:"type:timestamp;auto_update_time"`
}
