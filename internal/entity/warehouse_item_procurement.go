package entity

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/semeton-corp/anugerah-jaya-farm-volare/pkg/enum"
)

type WarehouseItemProcurement struct {
	Id               uint64                    `gorm:"primaryKey;autoIncrement"`
	WarehouseId      uint64                    `gorm:"type:bigint;not null"`
	Warehouse        Warehouse                 `gorm:"foreignKey:WarehouseId;references:Id;constraint:OnDelete:CASCADE"`
	ItemId           uint64                    `gorm:"type:bigint;not null"`
	Item             Item                      `gorm:"foreignKey:ItemId;references:Id;constraint:OnDelete:CASCADE"`
	SupplierId       uint64                    `gorm:"type:bigint;not null"`
	Supplier         Supplier                  `gorm:"foreignKey:SupplierId;references:Id;constraint:OnDelete:CASCADE"`
	Quantity         float64                   `gorm:"type:decimal;not null"`
	RecieveQuantity  float64                   `gorm:"type:decimal;not null;default:0"`
	EstimationRunOut sql.NullTime              `gorm:"type:date"`
	IsTaken          sql.NullBool              `gorm:"type:boolean;default:false"`
	TakenAt          sql.NullTime              `gorm:"type:timestamp"`
	TakenBy          uuid.NullUUID             `gorm:"type:varchar(255)"`
	Status           enum.WarehouseOrderStatus `gorm:"type:int;not null"`
	CreatedAt        time.Time                 `gorm:"type:timestamp;autoCreateTime"`
	CreatedBy        uuid.NullUUID             `gorm:"type:varchar(255)"`
	UpdatedAt        time.Time                 `gorm:"type:timestamp;autoUpdateTime"`
	UpdatedBy        uuid.NullUUID             `gorm:"type:varchar(255)"`
}
