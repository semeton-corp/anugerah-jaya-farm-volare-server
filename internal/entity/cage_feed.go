package entity

import "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"

type CageFeed struct {
	Id              uint64               `gorm:"primaryKey;autoIncrement"`
	CageId          uint64               `gorm:"type:bigint;not null"`
	Cage            Cage                 `gorm:"foreignKey:CageId;references:Id;constraint:OnDelete:CASCADE"`
	ChickenCategory enum.ChickenCategory `gorm:"type:int;not null"`
}
