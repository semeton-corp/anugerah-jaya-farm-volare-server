package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type Item struct {
	Id            uint64            `gorm:"primaryKey;autoIncrement"`
	Name          string            `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_category_unit"`
	Category      enum.ItemCategory `gorm:"type:int;not null;uniqueIndex:idx_name_category_unit"`
	Unit          string            `gorm:"type:varchar(255);not null;uniqueIndex:idx_name_category_unit"`
	DailySpending sql.NullFloat64   `gorm:"type:decimal"`
	CreatedAt     time.Time         `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy     uuid.NullUUID     `gorm:"type:varchar(255)"`
	UpdatedAt     time.Time         `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy     uuid.NullUUID     `gorm:"type:varchar(255)"`
}
