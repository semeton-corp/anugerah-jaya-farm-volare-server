package entity

import "time"

type FactoryMonitoring struct {
	Id                          uint64    `gorm:"primary_key;auto_increment"`
	Name                        string    `gorm:"type:varchar(255);not null"`
	TotalProductionsStiringFeed int       `gorm:"type:int;not null"`
	TotalUsedDedek              int       `gorm:"type:int;not null"`
	TotalUsedKonsentrat         int       `gorm:"type:int;not null"`
	TotalUsedPremix             int       `gorm:"type:int;not null"`
	TotalUsedCorn               int       `gorm:"type:int;not null"`
	CreatedAt                   time.Time `gorm:"type:timestamp;auto_create_time"`
	UpdatedAt                   time.Time `gorm:"type:timestamp;auto_update_time"`
}
