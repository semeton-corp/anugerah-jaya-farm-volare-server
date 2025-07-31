package entity

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type ChickenProcurementDraft struct {
	Id         uint64          `gorm:"primaryKey:autoIncrement"`
	CageId     uint64          `gorm:"type:bigint;not null"`
	Cage       Cage            `gorm:"foreignKey:CageId;refereces:Id"`
	SupplierId uint64          `gorm:"type:bigint;not null"`
	Supplier   Supplier        `gorm:"foreignKey:SupplierId;refereces:Id"`
	Quantity   uint64          `gorm:"type:bigint;not null"`
	Price      decimal.Decimal `gorm:"type:decimal;not null"`
	TotalPrice decimal.Decimal `gorm:"type:decimal;not null"`
	CreatedBy  uuid.NullUUID   `gorm:"type:varchar(255)"`
	CreatedAt  time.Time       `gorm:"type:timestamp;autoCreateTime"`
	UpdatedBy  uuid.NullUUID   `gorm:"type:varchar(255)"`
	UpdatedAt  time.Time       `gorm:"type:timestamp;autoUpdateTime"`
}
