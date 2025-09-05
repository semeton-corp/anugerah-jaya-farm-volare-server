package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	datatype "github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/custom/data_type"
)

type Notification struct {
	Id           uint64                    `gorm:"primaryKey;autoIncrement"`
	UserId       uuid.NullUUID             `gorm:"type:varchar(255)"`
	StoreId      sql.NullInt64             `gorm:"type:bigint"`
	WarehouseId  sql.NullInt64             `gorm:"type:bigint"`
	CageId       sql.NullInt64             `gorm:"bigint"`
	Description  string                    `gorm:"type:text;not null"`
	LocationType datatype.NullLocationType `gorm:"type:int"`
	IsMarked     bool                      `gorm:"type:boolean;default:false"`
	CreatedAt    time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy    uuid.NullUUID             `gorm:"type:varchar(255)"`
	UpdatedAt    time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy    uuid.NullUUID             `gorm:"type:varchar(255)"`
}
