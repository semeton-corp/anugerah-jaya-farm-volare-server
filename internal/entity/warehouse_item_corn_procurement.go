package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
	"github.com/shopspring/decimal"
)

type WarehouseItemCornProcurement struct {
	Id                        uint64                                `gorm:"primaryKey;autoIncrement"`
	WarehouseId               uint64                                `gorm:"type:bigint;not null"`
	Warehouse                 Warehouse                             `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	SupplierId                uint64                                `gorm:"type:bigint;not null"`
	Supplier                  Supplier                              `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity                  float64                               `gorm:"type:decimal;not null"`
	RecieveQuantity           sql.NullFloat64                       `gorm:"type:decimal;not null;default:0"`
	Note                      string                                `gorm:"type:text"`
	Price                     decimal.Decimal                       `gorm:"type:decimal;not null"`
	TotalPrice                decimal.Decimal                       `gorm:"type:decimal;not null"`
	IsArrived                 bool                                  `gorm:"type:boolean;default:false"`
	TakenAt                   sql.NullTime                          `gorm:"type:timestamp"`
	TakenBy                   uuid.NullUUID                         `gorm:"type:varchar(255)"`
	Status                    enum.ProcurementStatus                `gorm:"type:int;not null"`
	PaymentStatus             enum.PaymentStatus                    `gorm:"type:int;not null"`
	Payments                  []WarehouseItemCornProcurementPayment `gorm:"foreignKey:WarehouseItemCornProcurementId;references:Id"`
	OvenCondition             enum.OvenCondition                    `gorm:"type:int;not null"`
	CornWaterLevel            enum.CornWaterLevel                   `gorm:"type:int;not null"`
	IsOvenCanOperateInNearDay bool                                  `gorm:"type:bool;not null"`
	ExpiredAt                 time.Time                             `gorm:"timestamp;not null"`
	CreatedAt                 time.Time                             `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy                 uuid.NullUUID                         `gorm:"type:varchar(255)"`
	UpdatedAt                 time.Time                             `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy                 uuid.NullUUID                         `gorm:"type:varchar(255)"`
}
