package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type ChickenHealthItem struct {
	Id         uint64                     `gorm:"primaryKey;autoIncrement"`
	Name       string                     `gorm:"type:varchar(255);not null"`
	Type       enum.ChickenHealthItemType `gorm:"type:int;not null"`
	ChickenAge sql.NullInt64              `gorm:"type:int"`
	Note       string                     `gorm:"type:text"`
	CreatedAt  time.Time                  `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy  uuid.NullUUID              `gorm:"type:varchar(255)"`
	UpdatedAt  time.Time                  `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy  uuid.NullUUID              `gorm:"type:varchar(255)"`
}
