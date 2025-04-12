package entity

import "time"

type DailyWork struct {
	Id             uint64    `gorm:"primary_key;auto_increment"`
	Name           string    `gorm:"type:varchar(255);not null"`
	AssignedRoleId uint      `gorm:"type:varchar(255);not null"`
	StartTime      time.Time `gorm:"type:time;not null"`
	EndTime        time.Time `gorm:"type:time;not null"`
	CreatedAt      time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt      time.Time `gorm:"type:timestamp;auto_update_time"`
}
