package entity

import "time"

type FactoryMonitoring struct {
	Id                          uint64    `gorm:"primaryKey;autoIncrement"`
	Name                        string    `gorm:"type:varchar(255);not null"`
	TotalProductionsStiringFeed int       `gorm:"type:int;not null"`
	TotalUsedDedek              int       `gorm:"type:int;not null"`
	TotalUsedKonsentrat         int       `gorm:"type:int;not null"`
	TotalUsedPremix             int       `gorm:"type:int;not null"`
	TotalUsedCorn               int       `gorm:"type:int;not null"`
	CreatedAt                   time.Time `gorm:"type:timestamp;autoCreateTime"`
	UpdatedAt                   time.Time `gorm:"type:timestamp;autoUpdateTime"`
}
