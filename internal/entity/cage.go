package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"gorm.io/gorm"
)

type Cage struct {
	Id              uint64               `gorm:"primaryKey;autoIncrement"`
	LocationId      uint64               `gorm:"type:bigint;not null"`
	Location        Location             `gorm:"foreignKey:LocationId;references:Id;constraint:OnDelete:CASCADE"`
	Name            string               `gorm:"type:varchar(255);not null"`
	Capacity        uint64               `gorm:"type:bigint;not null"`
	ChickenCategory enum.ChickenCategory `gorm:"type:bigint;not null"`
	CagePlacement   []CagePlacement      `gorm:"foreignKey:CageId;references:Id"`
	IsUsed          bool                 `gorm:"type:boolean;default:false"`
	CreatedAt       time.Time            `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy       uuid.NullUUID        `gorm:"type:varchar(255)"`
	UpdatedAt       time.Time            `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy       uuid.NullUUID        `gorm:"type:varchar(255)"`
	DeletedAt       gorm.DeletedAt       `gorm:"type:timestamp;index"`
}

func (c *Cage) BeforeDelete(tx *gorm.DB) (err error) {
	// Append timestamp to name to make it available for reuse
	err = tx.Model(&Cage{}).Where("id = ?", c.Id).Update("name", fmt.Sprintf("%s_deleted_%d", c.Name, time.Now().Unix())).Error
	return
}
