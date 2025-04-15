package entity

import "time"

type DailyWorkStaff struct {
	DailyWorkId uint64    `gorm:"type:biginteger;not null"`
	DailyWork   DailyWork `gorm:"foreignKey:Id;references:DailyWorkId"`
	StaffId     uint64    `gorm:"type:biginteger;not null"`
	Staff       Staff     `gorm:"foreignKey:Id;references:StaffId"`
	IsDone      bool      `gorm:"type:bool;not null"`
	CreatedAt   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt   time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
