package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type WarehouseItemCornProcurementDraft struct {
	Id                        uint64              `gorm:"primaryKey;autoIncrement"`
	WarehouseId               uint64              `gorm:"type:bigint;not null"`
	Warehouse                 Warehouse           `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	SupplierId                uint64              `gorm:"type:bigint;not null"`
	Supplier                  Supplier            `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE"`
	OvenCondition             enum.OvenCondition  `gorm:"type:int;not null"`
	CornWaterLevel            enum.CornWaterLevel `gorm:"type:int;not null"`
	IsOvenCanOperateInNearDay bool                `gorm:"type:bool;not null"`
	Quantity                  float64             `gorm:"type:decimal;not null"`
	Price                     decimal.Decimal     `gorm:"type:decimal;not null"`
	CreatedAt                 time.Time           `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy                 uuid.NullUUID       `gorm:"type:varchar(255)"`
	UpdatedAt                 time.Time           `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy                 uuid.NullUUID       `gorm:"type:varchar(255)"`
}
